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

	pipelinesv1alpha1 "github.com/davidlynch-sd/lugh/api/v1alpha1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// TaskReconciler reconciles a Task object
type TaskReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=pipelines.lugh.dev,resources=tasks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=pipelines.lugh.dev,resources=tasks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=pipelines.lugh.dev,resources=tasks/finalizers,verbs=update

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

	// TODO(user): your logic here
	task := &pipelinesv1alpha1.Task{}
	err := r.Get(ctx, req.NamespacedName, task)

	if err != nil {
		return ctrl.Result{}, err
	}
	fmt.Println("Image: ", task.Spec.Image, "\nCommand: ", task.Spec.Command)

	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test",
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    task.Spec.Image,
					Image:   task.Spec.Image,
					Command: task.Spec.Command,
				},
			},
		},
	}
	fmt.Println(pod)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TaskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipelinesv1alpha1.Task{}).
		Complete(r)
}
