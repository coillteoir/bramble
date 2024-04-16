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
	//"path/filepath"

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	// batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	sourceRoot = "/src/"
)

func initExecution(
	ctx context.Context,
	reconciler *ExecutionReconciler,
	execution *pipelinesv1alpha1.Execution,
	pvc *corev1.PersistentVolumeClaim,
) error {
	logger := log.Log.WithName(fmt.Sprintf("Execution: %v", execution.ObjectMeta.Name))

	if !execution.Status.VolumeProvisioned {
		logger.Info("Provisioning PVC")

		pvc = &corev1.PersistentVolumeClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%v-pvc", execution.ObjectMeta.Name),
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

		err := ctrl.SetControllerReference(execution, pvc, reconciler.Scheme)
		if err != nil {
			return err
		}
		err = reconciler.Create(ctx, pvc)
		if err != nil {
			return err
		}

		logger.Info("Provisioning Cloner Pod")
	}

	if !execution.Status.RepoCloned {
		cloneTask := &pipelinesv1alpha1.PLTask{
			Name: "cloner",
			Spec: pipelinesv1alpha1.TaskSpec{
				Image: "alpine/git",
				Command: []string{
					"sh", "-c",
					fmt.Sprintf("rm -rf %v && git clone %v --depth=1 --branch=%v %v",
						execution.ObjectMeta.Name,
						execution.Spec.Repo,
						execution.Spec.Branch,
						execution.ObjectMeta.Name,
					),
				},
			},
		}
		cloneJob, err := generateTaskJob(execution, cloneTask, pvc)
		if err != nil {
			return err
		}

		err = ctrl.SetControllerReference(execution, cloneJob, reconciler.Scheme)
		if err != nil {
			return err
		}
		err = reconciler.Create(ctx, cloneJob)
		if err != nil {
			return err
		}

		logger.Info("Cloner Pod provisioned")
	}

	return nil
}
