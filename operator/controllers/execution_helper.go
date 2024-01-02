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

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

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

func execute_pipeline_using_dfs(matrix [][]int, pipeline *pipelinesv1alpha1.Pipeline) {

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