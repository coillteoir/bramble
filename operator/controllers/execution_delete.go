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

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func teardownExecution(ctx context.Context, reconciler *ExecutionReconciler, execution *pipelinesv1alpha1.Execution) error {
	exePods := &corev1.PodList{}

	podSelector := metav1.LabelSelector{
		MatchLabels: map[string]string{
			"bramble-execution": execution.ObjectMeta.Name,
		},
	}

	podListOptions := &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(
			labels.Set(podSelector.MatchLabels),
		),
	}

	err := reconciler.Client.List(ctx, exePods, podListOptions)
	if err != nil {
		return err
	}

	for i := range exePods.Items {
		err = reconciler.Client.Delete(ctx, &exePods.Items[i])
		if err != nil {
			return err
		}
	}

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
	err = reconciler.Client.List(ctx, pvcList, listOptions)
	if err != nil {
		return err
	}

	for i := range pvcList.Items {
		err = reconciler.Client.Delete(ctx, &pvcList.Items[i])
		if err != nil {
			return err
		}
	}

	pvList := &corev1.PersistentVolumeList{}
	err = reconciler.Client.List(ctx, pvList)
	if err != nil {
		return err
	}

	for i, pv := range pvList.Items {
		if pv.ObjectMeta.Name == execution.ObjectMeta.Name+"-pv" {
			err = reconciler.Client.Delete(ctx, &pvList.Items[i])
			if err != nil {
				return err
			}

			break
		}
	}

	if err != nil {
		return err
	}

	controllerutil.RemoveFinalizer(execution, executionFinalizer)
	err = reconciler.Update(ctx, execution)
	if err != nil {
		return err
	}

	return nil
}
