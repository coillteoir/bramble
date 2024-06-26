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
	"path/filepath"
	"slices"

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	sourceVolumeName = "source-volume"
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
	jobList *batchv1.JobList,
) (bool, error) {
	logger := log.Log.WithName(fmt.Sprintf("Execution: %v", execution.ObjectMeta.Name))

	task := pipeline.Spec.Tasks[start]

	jobIndex := slices.IndexFunc(jobList.Items, func(job batchv1.Job) bool {
		return (job.ObjectMeta.Labels["bramble-task"] == task.Name &&
			job.ObjectMeta.Labels["bramble-execution"] == execution.ObjectMeta.Name)
	})
	if jobIndex != -1 {
		job := jobList.Items[jobIndex]
		baseStr := fmt.Sprintf("Execution: %v Task: %v", execution.ObjectMeta.Name, task.Name)
		logger.Info(fmt.Sprintf("%v %v", baseStr, job.Status.Succeeded))
		if job.Status.Succeeded != 0 {
			if start == 0 {
				execution.Status.Phase = pipelinesv1alpha1.ExecutionCompleted
				return false, nil
			}
			return false, nil
		}
		if job.Status.Active != 0 {
			return false, nil
		}
		if job.Status.Failed != 0 {
			return false, fmt.Errorf("job: %v failed", job.ObjectMeta.Name)
		}
	}
	logger.Info(fmt.Sprintf("Checking dependencies of task: %v", task.Name))

	// Check dependencies have completed before running task.
	count := len(task.Spec.Dependencies)

	for _, dep := range task.Spec.Dependencies {
		jobIndex := slices.IndexFunc(jobList.Items, func(job batchv1.Job) bool {
			return (job.ObjectMeta.Labels["bramble-execution"] == execution.ObjectMeta.Name &&
				job.ObjectMeta.Labels["bramble-task"] == dep &&
				job.Status.Succeeded == 1)
		})
		if jobIndex != -1 {
			count--
		}
		logger.Info(fmt.Sprintf("Task: %v, Dependency: %v, Count: %v",
			task.Name,
			dep,
			count,
		))
	}

	return count == 0, nil
}

func executeUsingDfs(
	matrix [][]int,
	start int,
	visited []bool,
	pipeline *pipelinesv1alpha1.Pipeline,
	execution *pipelinesv1alpha1.Execution,
	jobList *batchv1.JobList,
	pvc *corev1.PersistentVolumeClaim,
) (*batchv1.JobList, error) {
	logger := log.Log.WithName(fmt.Sprintf("Execution: %v", execution.ObjectMeta.Name))

	visited[start] = true

	task := pipeline.Spec.Tasks[start]

	jobs := &batchv1.JobList{}

	toRun, err := validateTask(
		start,
		pipeline,
		execution,
		jobList,
	)
	if err != nil {
		return nil, err
	}
	// BUG: When running controller-manger in cluster,
	// jobs are created multiple times and executions do not work
	if toRun {
		logger.Info(fmt.Sprintf("Executing task: %v", task.Name))

		job, err := generateTaskJob(execution, &task, pvc)
		if err != nil {
			return nil, err
		}
		jobs.Items = append(jobs.Items, *job)
		execution.Status.Phase = pipelinesv1alpha1.ExecutionRunning
		execution.Status.Running = append(execution.Status.Running, task.Name)
		return jobs, nil
	}

	for i, node := range matrix[start] {
		if node == 1 && !visited[i] {
			downstreamJobs, err := executeUsingDfs(
				matrix,
				i,
				visited,
				pipeline,
				execution,
				jobList,
				pvc,
			)
			if err != nil {
				return nil, err
			}

			jobs.Items = append(jobs.Items, downstreamJobs.Items...)
		}
	}
	return jobs, nil
}

func generateTaskJob(
	execution *pipelinesv1alpha1.Execution,
	task *pipelinesv1alpha1.PLTask,
	pvc *corev1.PersistentVolumeClaim,
) (*batchv1.Job, error) {
	if pvc == nil {
		return nil, fmt.Errorf("NO PVC for pod")
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: execution.ObjectMeta.Name + "-" + task.Name + "-",
			Namespace:    execution.ObjectMeta.Namespace,
			Labels: map[string]string{
				"bramble-execution": execution.ObjectMeta.Name,
				"bramble-task":      task.Name,
			},
		}, Spec: batchv1.JobSpec{
			Completions:  (func(num int32) *int32 { return &num }(1)),
			BackoffLimit: (func(num int32) *int32 { return &num }(1)),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						{
							Name:    task.Name,
							Image:   task.Spec.Image,
							Command: task.Spec.Command,
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      sourceVolumeName,
									MountPath: sourceRoot,
								},
							},
							WorkingDir:      filepath.Join(sourceRoot, execution.ObjectMeta.Name, task.Spec.Workdir),
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: sourceVolumeName,
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: pvc.ObjectMeta.Name,
								},
							},
						},
					},
				},
			},
		},
	}

	return job, nil
}
