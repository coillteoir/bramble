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

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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

func (r *ExecutionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	logger := log.Log.WithName("execution_logs")

	execution := &pipelinesv1alpha1.Execution{}
	err := r.Get(ctx, req.NamespacedName, execution)
	if err != nil {
		return ctrl.Result{}, err
	}

	if execution.Status.Error {
		logger.Error(err, "Execution in failed state")
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

	if !pipeline.Status.ValidDeps {
		logger.Info("Invalid dependency tree in pipeline")
		return ctrl.Result{}, errors.New("invalid pipeline")
	}

	//matrix := generateAssociationMatrix(pipeline)

	logger.Info(fmt.Sprintf("Name: %v", execution.ObjectMeta.Name))

	var pv *corev1.PersistentVolume
	var pvc *corev1.PersistentVolumeClaim

	// Check if volume exists
	pvList := &corev1.PersistentVolumeList{}
	err = r.Client.List(ctx, pvList)

	if err != nil {
		return ctrl.Result{}, err
	}

	if len(pvList.Items) > 0 {
		for _, p := range pvList.Items {
			if p.ObjectMeta.Name == execution.ObjectMeta.Name+"-pv" {
				if p.Status.Phase == corev1.VolumeBound || p.Status.Phase == corev1.VolumeAvailable {
					execution.Status.VolumeProvisioned = true
					err = r.Status().Update(ctx, execution)
					if err != nil {
						return ctrl.Result{}, err
					}
				}
			}
		}
	}

	if !execution.Status.VolumeProvisioned {
		err = initExecution(ctx, r, execution, pv, pvc)

		if err != nil {
			return ctrl.Result{}, err
		}
	} else {
		pvcSelector := metav1.LabelSelector{
			MatchLabels: map[string]string{
				"bramble-execution": execution.ObjectMeta.Name,
			},
		}
		listOptions := &client.ListOptions{
			LabelSelector: labels.SelectorFromSet(
				labels.Set(pvcSelector.MatchLabels),
			),
		}
		pvcList := &corev1.PersistentVolumeClaimList{}
		err = r.Client.List(ctx, pvcList, listOptions)
		if err != nil {
			return ctrl.Result{}, err
		}
		for _, p := range pvcList.Items {
			if p.ObjectMeta.Labels["bramble-execution"] == execution.ObjectMeta.Name {
				pvc = &p
			}
		}
	}
	if err != nil {
		logger.Error(err, "Couldn't update execution")
		return ctrl.Result{}, err
	}

	exePods := &corev1.PodList{}

	podSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"bramble-execution": execution.ObjectMeta.Name,
		},
	}

	listOptions := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			labels.Set(podSelector.MatchLabels),
		),
	}

	err = r.Client.List(ctx, exePods, listOptions)

	if err != nil {
		return ctrl.Result{}, err
	}


    // Check if the repo is cloned before proceeding
	for _, pod := range exePods.Items {
		if pod.ObjectMeta.Name == execution.ObjectMeta.Name+"-cloner" {
			execution.Status.RepoCloned = pod.Status.Phase == corev1.PodSucceeded
            logger.Info(fmt.Sprintf("Execution: %v repo cloned", execution.ObjectMeta.Name))
		}
	}
	err = r.Status().Update(ctx, execution)
    if err != nil {
        return ctrl.Result{}, err
    }
	// NOTE This algorithm assumes tasks are in their topological order
	// Needs to be reworked to handle unsorted matricies

	//tasks := pipeline.Spec.Tasks

	// this code is not good, too indented
	// Recursion will provide a better solution to this problem
	if execution.Status.VolumeProvisioned && execution.Status.RepoCloned {
	}
	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

func generateAssociationMatrix(pipeline *pipelinesv1alpha1.Pipeline) [][]int {

	matrix := make([][]int, len(pipeline.Spec.Tasks))

	for i, task := range pipeline.Spec.Tasks {
		for _, tasktask := range pipeline.Spec.Tasks {
			if slices.Contains(task.Spec.Dependencies, tasktask.Name) {
				matrix[i] = append(matrix[i], 1)
			} else {
				matrix[i] = append(matrix[i], 0)
			}
		}
	}
	return matrix
}

func runTask(ctx context.Context, r *ExecutionReconciler, execution *pipelinesv1alpha1.Execution, task *pipelinesv1alpha1.PLTask, pvc *corev1.PersistentVolumeClaim) error {

	if pvc == nil {
		return fmt.Errorf("NO PVC for pod")
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: execution.ObjectMeta.Name + "-" + task.Name + "-",
			Namespace:    execution.ObjectMeta.Namespace,
			Labels: map[string]string{
				"bramble-execution": execution.ObjectMeta.Name,
				"bramble-task":      task.Name,
			},
		}, Spec: corev1.PodSpec{
			RestartPolicy: corev1.RestartPolicyNever,
			Containers: []corev1.Container{
				{
					Name:    task.Name,
					Image:   task.Spec.Image,
					Command: task.Spec.Command,
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

	err := r.Create(ctx, pod)
	if err != nil {
		return err
	}

	return nil
}

func initExecution(ctx context.Context, r *ExecutionReconciler, execution *pipelinesv1alpha1.Execution, pv *corev1.PersistentVolume, pvc *corev1.PersistentVolumeClaim) error {
	if !execution.Status.VolumeProvisioned {
		log.Log.WithName("execution logs").Info("Provisioning PV")
		pv = &corev1.PersistentVolume{
			ObjectMeta: metav1.ObjectMeta{
				Name:   execution.ObjectMeta.Name + "-pv",
				Labels: map[string]string{"bramble-execution": execution.ObjectMeta.Name},
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
				Labels:    map[string]string{"bramble-execution": execution.ObjectMeta.Name},
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
		clonePod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      execution.ObjectMeta.Name + "-cloner",
				Namespace: execution.ObjectMeta.Namespace,
				Labels:    map[string]string{"bramble-execution": execution.ObjectMeta.Name},
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
							"/src/",
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
