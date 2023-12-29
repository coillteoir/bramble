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

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// TaskReconciler reconciles a Task object
type TaskReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pipelines.bramble.dev,resources=tasks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pipelines.bramble.dev,resources=tasks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pipelines.bramble.dev,resources=tasks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Task object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile

func (r *TaskReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	task := &pipelinesv1alpha1.Task{}
	err := r.Get(ctx, req.NamespacedName, task)

	if err != nil {
		return ctrl.Result{}, err
	}
	log.Log.WithName("task_logs").Info(fmt.Sprintf("Image: %v Command: %v", task.Spec.Image, task.Spec.Command))

	if err != nil {
		return ctrl.Result{}, err
	}

	for _, dep := range task.Spec.Dependencies {
		log.Log.WithName("task_logs_dependencies").Info(dep)
		log.Log.WithName("request_name").Info(fmt.Sprintf("%v", req.NamespacedName))

		depTask := &pipelinesv1alpha1.Task{}
		err = r.Get(
			ctx,
			types.NamespacedName{
				Namespace: task.ObjectMeta.Namespace,
				Name:      dep,
			},
			depTask)
		if err != nil {
			log.Log.WithName("ERROR").Error(err, "Dependency not found")
			return ctrl.Result{}, err
		}
		log.Log.WithName("dependency_log").Info(fmt.Sprintf("%v", depTask.ObjectMeta.Name))
	}
	return ctrl.Result{ /*RequeueAfter: time.Duration(30 * time.Second)*/ }, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TaskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipelinesv1alpha1.Task{}).
		Complete(r)
}
