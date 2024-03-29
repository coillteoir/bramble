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
	"fmt"
	"slices"

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
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

func validateTask(
	start int,
	pipeline *pipelinesv1alpha1.Pipeline,
	execution *pipelinesv1alpha1.Execution,
	podList *corev1.PodList,
) (bool, error) {
	logger := log.Log.WithName(fmt.Sprintf("Execution: %v", execution.ObjectMeta.Name))

	task := pipeline.Spec.Tasks[start]

	for _, pod := range podList.Items {
		if pod.ObjectMeta.Labels["bramble-task"] == task.Name {
			baseStr := fmt.Sprintf("Execution: %v Task: %v", execution.ObjectMeta.Name, task.Name)

			logger.Info(fmt.Sprintf("%v %v", baseStr, pod.Status.Phase))

			switch pod.Status.Phase {
			case corev1.PodRunning, corev1.PodPending:
				return false, nil

			case corev1.PodSucceeded:

				if start == 0 {
					execution.Status.Completed = true
					return false, nil
				}
				return false, nil

			case corev1.PodFailed:
				return false, fmt.Errorf("pod: %v failed", pod.ObjectMeta.Name)

			default:
				logger.Info(fmt.Sprintf("%v %v Pod created but Phase undefined", baseStr, pod.Status.Phase))
				return false, nil
			}
		}
	}

	logger.Info(fmt.Sprintf("Checking dependencies of task: %v", task.Name))

	// Check dependencies have completed before running task.
	count := len(task.Spec.Dependencies)

	for _, dep := range task.Spec.Dependencies {
		for _, pod := range podList.Items {
			if pod.ObjectMeta.Labels["bramble-execution"] == execution.ObjectMeta.Name &&
				pod.ObjectMeta.Labels["bramble-task"] == dep &&
				pod.Status.Phase == corev1.PodSucceeded {
				count--
			}
		}

		logger.Info(
			fmt.Sprintf(
				"Task: %v, Dependency: %v, Count: %v",
				task.Name,
				dep,
				count,
			),
		)
	}

	return count == 0, nil
}

// I think creating a custom struct to handle execution task logic would be the move
// This function will have 10 args by the time it's finished.
func executeUsingDfs(
	matrix [][]int,
	start int,
	visited []bool,
	pipeline *pipelinesv1alpha1.Pipeline,
	execution *pipelinesv1alpha1.Execution,
	podList *corev1.PodList,
	pvc *corev1.PersistentVolumeClaim,
) ([]*corev1.Pod, error) {
	logger := log.Log.WithName(fmt.Sprintf("Execution: %v", execution.ObjectMeta.Name))

	visited[start] = true

	task := pipeline.Spec.Tasks[start]

	podsToExecute := make([]*corev1.Pod, 0)

	toRun, err := validateTask(
		start,
		pipeline,
		execution,
		podList,
	)
	if err != nil {
		return nil, err
	}
	// BUG: When running controller-manger in cluster, pods are created multiple times and executions do not work
	if toRun {
		logger.Info(fmt.Sprintf("Executing task: %v", task.Name))

		pod, err := runTask(execution, &task, pvc)
		if err != nil {
			return nil, err
		}

		podsToExecute = append(podsToExecute, pod)

		return podsToExecute, nil
	}

	for i, node := range matrix[start] {
		if node == 1 && !visited[i] {
			logger.Info(
				fmt.Sprintf(
					"recursing to dependency %v of task %v",
					pipeline.Spec.Tasks[i].Name,
					task.Name,
				),
			)

			downstreamPods, err := executeUsingDfs(
				matrix,
				i,
				visited,
				pipeline,
				execution,
				podList,
				pvc,
			)
			if err != nil {
				return nil, err
			}

			podsToExecute = append(podsToExecute, downstreamPods...)
		}
	}
	return podsToExecute, nil
}

func runTask(
	execution *pipelinesv1alpha1.Execution,
	task *pipelinesv1alpha1.PLTask,
	pvc *corev1.PersistentVolumeClaim,
) (*corev1.Pod, error) {
	if pvc == nil {
		return nil, fmt.Errorf("NO PVC for pod")
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
							MountPath: sourceRoot,
						},
					},
					WorkingDir: sourceRoot + execution.ObjectMeta.Name,
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

	return pod, nil
}
