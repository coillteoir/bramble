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
	"fmt"
	"slices"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
)

// PipelineReconciler reconciles a Pipeline object
type PipelineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pipelines.bramble.dev,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pipelines.bramble.dev,resources=pipelines/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pipelines.bramble.dev,resources=pipelines/finalizers,verbs=update

func (r *PipelineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	pipeline := &pipelinesv1alpha1.Pipeline{}
	err := r.Get(ctx, req.NamespacedName, pipeline)

	if err != nil {
		return ctrl.Result{}, err
	}

	log.Log.WithName("pipeline_logs").
		Info(fmt.Sprintf("Name: %v", pipeline.ObjectMeta.Name))
	if !pipeline.Status.TasksCreated {
		for _, task := range pipeline.Spec.Tasks {
			err = r.Create(ctx, &pipelinesv1alpha1.Task{
				ObjectMeta: metav1.ObjectMeta{
					Name:      (pipeline.ObjectMeta.Name + "-" + task.Name),
					Namespace: pipeline.ObjectMeta.Namespace,
				},
				Spec: task.Spec,
			})
			if err != nil {
				return ctrl.Result{}, err
			}

		}
	}
	pipeline.Status.TasksCreated = true
	err = validateDependencies(pipeline)
	if err != nil {
		log.Log.WithName("pipeline logs").Error(err, "invalid dependencies")
		return ctrl.Result{}, err
	}
	err = r.Status().Update(ctx, pipeline)
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func validateDependencies(pipeline *pipelinesv1alpha1.Pipeline) error {
	deps := make([]string, 0)
	for _, task := range pipeline.Spec.Tasks {
		if slices.Contains(task.Spec.Dependencies, task.Name) {
			return fmt.Errorf("%v cannot contain itself as a dependency", task.Name)
		}
		deps = append(deps, task.Spec.Dependencies...)
	}
	// Check all dependencies are valid tasks
	for _, dep := range deps {
		depflag := false
		for _, task := range pipeline.Spec.Tasks {
			if dep == task.Name {
				depflag = true
				break
			}
		}
		if !depflag {
			return fmt.Errorf("Invalid dependency: %v. Not referenced in the pipeline. Please apply the task to the cluster, or describe it within the pipeline", dep)
		}
	}

	pipeline.Status.ValidDeps = true

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipelinesv1alpha1.Pipeline{}).
		Complete(r)
}
