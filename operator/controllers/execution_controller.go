/*
Copyright 2023 David Lynch.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// ExecutionReconciler reconciles a Execution object
type ExecutionReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pipelines.bramble.dev,resources=executions,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pipelines.bramble.dev,resources=executions/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pipelines.bramble.dev,resources=executions/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=create;delete;list;get
//+kubebuilder:rbac:groups="batch",resources=jobs,verbs=create;delete;list;get;watch

func toContinue(execution *pipelinesv1alpha1.Execution) bool {
	if execution.Status.Phase == pipelinesv1alpha1.ExecutionError {
		return false
	}
	if execution.Status.Phase == pipelinesv1alpha1.ExecutionCompleted {
		return false
	}
	if execution.GetDeletionTimestamp() != nil {
		return false
	}
	return true
}

func (reconciler *ExecutionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	execution := &pipelinesv1alpha1.Execution{}

	err := reconciler.Get(ctx, req.NamespacedName, execution)
	if err != nil {
		return ctrl.Result{}, err
	}

	logger := log.Log.WithName(fmt.Sprintf("Execution: %v", execution.ObjectMeta.Name))

	if !toContinue(execution) {
		return ctrl.Result{}, nil
	}

	pipeline, err := loadPipeline(execution, reconciler, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	pvc := &corev1.PersistentVolumeClaim{}
	listOptions := generateListOptions(execution)
	pvcList, err := getExecutionPvcs(listOptions, reconciler, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	pvcIndex := slices.IndexFunc(pvcList.Items, func(p corev1.PersistentVolumeClaim) bool {
		return p.ObjectMeta.Labels["bramble-execution"] == execution.ObjectMeta.Name
	})
	if pvcIndex != -1 {
		pvc = &pvcList.Items[pvcIndex]

		execution.Status.VolumeProvisioned = pvc != nil
		err = reconciler.Status().Update(ctx, execution)
		if err != nil {
			return ctrl.Result{}, err
		}
		err = ctrl.SetControllerReference(execution, pvc, reconciler.Scheme)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else {
		logger.Info("no pvc found after repo cloned")
	}

	exeJobs, err := getExecutionJobs(listOptions, reconciler, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	if verifyClone(exeJobs, execution) {
		execution.Status.RepoCloned = true
		logger.Info(fmt.Sprintf("Execution: %v repo cloned", execution.ObjectMeta.Name))
		err = reconciler.Status().Update(ctx, execution)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if !execution.Status.VolumeProvisioned {
		err = initExecution(ctx, reconciler, execution, pvc)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	if execution.Status.VolumeProvisioned &&
		execution.Status.RepoCloned {
		err = executePipeline(reconciler, ctx, execution, pipeline, exeJobs, pvc)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	err = reconciler.Update(ctx, execution)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = reconciler.Status().Update(ctx, execution)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: 15 * time.Duration(time.Second)}, nil
}

func executePipeline(
	reconciler *ExecutionReconciler,
	ctx context.Context,
	execution *pipelinesv1alpha1.Execution,
	pipeline *pipelinesv1alpha1.Pipeline,
	exeJobs *batchv1.JobList,
	pvc *corev1.PersistentVolumeClaim,
) error {
	matrix := generateAssociationMatrix(pipeline)

	visited := make([]bool, len(matrix))
	jobsToExecute, err := executeUsingDfs(
		matrix,
		0,
		visited,
		pipeline,
		execution,
		exeJobs,
		pvc,
	)
	if err != nil {
		execution.Status.Phase = pipelinesv1alpha1.ExecutionError
		job_err := reconciler.Status().Update(ctx, execution)
		if err != nil {
			return err
		}
		return job_err
	}
	err = reconciler.Status().Update(ctx, execution)
	if err != nil {
		return err
	}

	for i := range jobsToExecute.Items {
		err = ctrl.SetControllerReference(execution, &jobsToExecute.Items[i], reconciler.Scheme)
		if err != nil {
			return err
		}
		err = reconciler.Client.Create(ctx, &jobsToExecute.Items[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func verifyClone(exeJobs *batchv1.JobList, execution *pipelinesv1alpha1.Execution) bool {
	jobIndex := slices.IndexFunc(exeJobs.Items, func(job batchv1.Job) bool {
		return (job.ObjectMeta.Labels["bramble-task"] == "cloner" && job.ObjectMeta.Labels["bramble-execution"] == execution.ObjectMeta.Name)
	})
	if jobIndex != -1 {
		return exeJobs.Items[jobIndex].Status.Succeeded == 1
	}
	return false
}

func loadPipeline(
	execution *pipelinesv1alpha1.Execution,
	reconciler *ExecutionReconciler,
	ctx context.Context,
) (*pipelinesv1alpha1.Pipeline, error) {
	pipeline := &pipelinesv1alpha1.Pipeline{}
	err := reconciler.Get(ctx, types.NamespacedName{
		Name:      execution.Spec.Pipeline,
		Namespace: execution.ObjectMeta.Namespace,
	}, pipeline)
	if err != nil {
		return nil, err
	}

	if !pipeline.Status.ValidDeps {
		return nil, errors.New("invalid pipeline")
	}
	return pipeline, nil
}

func generateListOptions(execution *pipelinesv1alpha1.Execution) *client.ListOptions {
	selector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"bramble-execution": execution.ObjectMeta.Name,
		},
	}

	return &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			labels.Set(selector.MatchLabels),
		),
	}
}

func getExecutionJobs(
	listOptions *client.ListOptions,
	reconciler *ExecutionReconciler,
	ctx context.Context,
) (*batchv1.JobList, error) {
	jobs := &batchv1.JobList{}

	err := reconciler.Client.List(ctx, jobs, listOptions)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func getExecutionPvcs(
	listOptions *client.ListOptions,
	reconciler *ExecutionReconciler,
	ctx context.Context,
) (*corev1.PersistentVolumeClaimList, error) {
	pvcList := &corev1.PersistentVolumeClaimList{}

	err := reconciler.Client.List(ctx, pvcList, listOptions)
	if err != nil {
		return nil, err
	}
	return pvcList, nil
}

// SetupWithManager sets up the controller with the Manager.
func (reconciler *ExecutionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipelinesv1alpha1.Execution{}).
		Complete(reconciler)
}
