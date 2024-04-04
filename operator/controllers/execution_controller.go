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
	"time"

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
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
//+kubebuilder:rbac:groups="",resources=pods;persistentvolumes;persistentvolumeclaims,verbs=create;delete;list;get

const (
	executionFinalizer = "executions.pipelines.bramble.dev/finalizer"
	pvSuffix           = "-pv"
)

func (reconciler *ExecutionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	execution := &pipelinesv1alpha1.Execution{}

	err := reconciler.Get(ctx, req.NamespacedName, execution)
	if err != nil {
		return ctrl.Result{}, err
	}

	logger := log.Log.WithName(fmt.Sprintf("Execution: %v", execution.ObjectMeta.Name))
	if execution.Status.Phase == pipelinesv1alpha1.ExecutionError {
		logger.Error(err, "Execution in failed state")

		return ctrl.Result{}, errors.New("execution in failed state")
	}
	if execution.Status.Phase == pipelinesv1alpha1.ExecutionCompleted {
		logger.Info("Completed!")

		return ctrl.Result{}, nil
	}

	if execution.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(execution, executionFinalizer) {
			err = teardownExecution(ctx, reconciler, execution)
			if err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
	}

	pipeline, err := loadPipeline(execution, reconciler, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if volume exists
	pvList := &corev1.PersistentVolumeList{}
	err = reconciler.Client.List(ctx, pvList)
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, pv := range pvList.Items {
		if pv.ObjectMeta.Name == execution.ObjectMeta.Name+pvSuffix {
			execution.Status.VolumeProvisioned = (pv.Status.Phase == corev1.VolumeBound || pv.Status.Phase == corev1.VolumeAvailable)
		}
	}

	var pvc *corev1.PersistentVolumeClaim
	if !execution.Status.VolumeProvisioned {
		err = initExecution(ctx, reconciler, execution, pvc)

		if err != nil {
			return ctrl.Result{}, err
		}
	}

	listOptions := generateListOptions(execution)
	pvcList, err := getExecutionPvcs(listOptions, reconciler, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	for i, p := range pvcList.Items {
		if p.ObjectMeta.Labels["bramble-execution"] == execution.ObjectMeta.Name {
			pvc = &pvcList.Items[i]
		}
	}

	exePods, err := getExecutionPods(listOptions, reconciler, ctx)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Check if the repo is cloned before proceeding
	if !execution.Status.RepoCloned {
		for _, pod := range exePods.Items {
			if pod.ObjectMeta.Name == execution.ObjectMeta.Name+"-cloner" {
				execution.Status.RepoCloned = pod.Status.Phase == corev1.PodSucceeded
				if execution.Status.RepoCloned {
					logger.Info(fmt.Sprintf("Execution: %v repo cloned", execution.ObjectMeta.Name))
				}
			}
		}
	}

	// NOTE This algorithm assumes tasks are in their topological order
	// Needs to be reworked to handle unsorted matricies

	if execution.Status.VolumeProvisioned &&
		execution.Status.RepoCloned && !(execution.Status.Phase == "completed") {

		matrix := generateAssociationMatrix(pipeline)
		logger.Info(fmt.Sprintf("MATRIX: %v", matrix))

		visited := make([]bool, len(matrix))
		podsToExecute, err := executeUsingDfs(
			matrix,
			0,
			visited,
			pipeline,
			execution,
			exePods,
			pvc,
		)
		if err != nil {
			execution.Status.Phase = "error"
			return ctrl.Result{}, err
		}

		for _, pod := range podsToExecute {
			err = reconciler.Client.Create(ctx, pod)
			if err != nil {
				return ctrl.Result{}, err
			}
		}

	}

	if !controllerutil.ContainsFinalizer(execution, executionFinalizer) {
		controllerutil.AddFinalizer(execution, executionFinalizer)
	}

	err = reconciler.Update(ctx, execution)
	if err != nil {
		return ctrl.Result{}, err
	}
	err = reconciler.Status().Update(ctx, execution)
	if err != nil {
		return ctrl.Result{}, err
	} else {
		logger.Info(fmt.Sprintf("Updated status with: %v", execution.Status))
	}

	return ctrl.Result{RequeueAfter: 15 * time.Duration(time.Second)}, nil
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
		// logger.Info("Invalid dependency tree in pipeline")

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

func getExecutionPods(
	listOptions *client.ListOptions,
	reconciler *ExecutionReconciler,
	ctx context.Context,
) (*corev1.PodList, error) {
	exePods := &corev1.PodList{}

	err := reconciler.Client.List(ctx, exePods, listOptions)
	if err != nil {
		return nil, err
	}
	return exePods, nil
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
