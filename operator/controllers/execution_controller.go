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
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"
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

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Execution object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ExecutionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	execution := &pipelinesv1alpha1.Execution{}
	err := r.Get(ctx, req.NamespacedName, execution)
	if err != nil {
		return ctrl.Result{}, err
	}

	pipeline := &pipelinesv1alpha1.Pipeline{}
	err = r.Get(ctx, types.NamespacedName{
		Name:      execution.Spec.Pipeline,
		Namespace: execution.ObjectMeta.Namespace,
	}, pipeline)

	if err != nil {
		return ctrl.Result{}, err
	}

	log.Log.WithName("execution logs").
		Info(fmt.Sprintf("Name: %v", execution.ObjectMeta.Name))

	var pv *corev1.PersistentVolume
	var pvc *corev1.PersistentVolumeClaim
	var clonePod *corev1.Pod

	err = initExecution(ctx, r, execution, pv, pvc, clonePod)

	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.Status().Update(ctx, execution)
	if err != nil {
		log.Log.WithName("execution_logs").Error(err, "Couldn't update execution")
		return ctrl.Result{}, err
	}

	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

func initExecution(ctx context.Context, r *ExecutionReconciler, execution *pipelinesv1alpha1.Execution, pv *corev1.PersistentVolume, pvc *corev1.PersistentVolumeClaim, clonePod *corev1.Pod) error {

	if !execution.Status.VolumeProvisioned {
		log.Log.WithName("execution logs").Info("Provisioning PV")
		pv = &corev1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name: execution.ObjectMeta.Name + "-pv",
			},
			Spec: corev1.PersistentVolumeSpec{
				Capacity: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("1Gi"),
				},
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				StorageClassName: "standard",
				PersistentVolumeSource: corev1.PersistentVolumeSource{
					HostPath: &corev1.HostPathVolumeSource{
						Path: "/src",
					},
				},
			},
		}

		err := r.Create(ctx, pv)
		if err != nil {
			return err
		}
		log.Log.WithName("execution logs").
			Info(fmt.Sprintf("PV %v created", execution.ObjectMeta.Name))

		log.Log.WithName("execution logs").Info("Provisioning PVC")
		pvc = &corev1.PersistentVolumeClaim{

			ObjectMeta: metav1.ObjectMeta{
				Name:      pv.ObjectMeta.Name + "-pvc",
				Namespace: execution.ObjectMeta.Namespace,
			},
			Spec: corev1.PersistentVolumeClaimSpec{
				AccessModes: []corev1.PersistentVolumeAccessMode{
					corev1.ReadWriteOnce,
				},
				Resources: corev1.ResourceRequirements{
					Requests: corev1.ResourceList{
						corev1.ResourceStorage: resource.MustParse("1Gi"),
					},
				},
			},
		}
		err = r.Create(ctx, pvc)
		if err != nil {
			return err
		}
		execution.Status.VolumeProvisioned = true
		log.Log.WithName("execution_logs").Info("Provisioning Cloner Pod")
	}

	if !execution.Status.RepoCloned {
		clonePod = &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      execution.ObjectMeta.Name + "-cloner",
				Namespace: execution.ObjectMeta.Namespace,
			}, Spec: corev1.PodSpec{
				RestartPolicy: corev1.RestartPolicyOnFailure,
				Containers: []corev1.Container{
					{
						Name:  "cloner",
						Image: "alpine/git",
						Command: []string{
							"git",
							"clone",
							execution.Spec.Repo,
							execution.Spec.CloneDir,
							"--branch=" + execution.Spec.Branch,
						},
						VolumeMounts: []corev1.VolumeMount{
							{
								Name:      "cloner-volume",
								MountPath: "/src/",
							},
						},
					},
				},
				Volumes: []corev1.Volume{
					{
						Name: "cloner-volume",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: pvc.ObjectMeta.Name,
							},
						},
					},
				},
			},
		}

		err := r.Create(ctx, clonePod)
		if err != nil {
			return err
		}
		execution.Status.RepoCloned = true
		log.Log.WithName("execution_logs").Info("Cloner Pod provisioned")
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExecutionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&pipelinesv1alpha1.Execution{}).
		Complete(r)
}
