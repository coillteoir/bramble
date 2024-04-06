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
	"reflect"
	"testing"

	//	"testing/quick"

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Test_generateAssociationMatrix(t *testing.T) {
	type args struct {
		pipeline *pipelinesv1alpha1.Pipeline
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateAssociationMatrix(tt.args.pipeline); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateAssociationMatrix() = %v, want %v", got, tt.want)
			}
		})
	}

	// TestMatrixSize := func(
	// 	pipeline *pipelinesv1alpha1.Pipeline,
	// ) bool {
	// 	matrix := generateAssociationMatrix(pipeline)

	// 	return len(matrix) == len(pipeline.Spec.Tasks)
	// }

	// if err := quick.Check(TestMatrixSize, nil); err != nil {
	// 	t.Error(err)
	// }
}

func Test_validateTask(t *testing.T) {
	type args struct {
		start     int
		pipeline  *pipelinesv1alpha1.Pipeline
		execution *pipelinesv1alpha1.Execution
		podList   *corev1.PodList
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := validateTask(tt.args.start, tt.args.pipeline, tt.args.execution, tt.args.podList)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("validateTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_executeUsingDfs(t *testing.T) {
	type args struct {
		matrix    [][]int
		start     int
		visited   []bool
		pipeline  *pipelinesv1alpha1.Pipeline
		execution *pipelinesv1alpha1.Execution
		podList   *corev1.PodList
		pvc       *corev1.PersistentVolumeClaim
	}
	tests := []struct {
		name    string
		args    args
		want    []*corev1.Pod
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := executeUsingDfs(tt.args.matrix, tt.args.start, tt.args.visited, tt.args.pipeline, tt.args.execution, tt.args.podList, tt.args.pvc)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeUsingDfs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("executeUsingDfs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateTaskPod(t *testing.T) {
	type args struct {
		execution *pipelinesv1alpha1.Execution
		task      *pipelinesv1alpha1.PLTask
		pvc       *corev1.PersistentVolumeClaim
	}
	tests := []struct {
		name    string
		args    args
		want    *corev1.Pod
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateTaskPod(tt.args.execution, tt.args.task, tt.args.pvc)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateTaskPod() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateTaskPod() = %v, want %v", got, tt.want)
			}
		})
	}

	//	TestPodGeneration := func(
	//		execution *pipelinesv1alpha1.Execution,
	//		task *pipelinesv1alpha1.PLTask,
	//		pvc *corev1.PersistentVolumeClaim,
	//	) bool {
	//		return true
	//	}
	//	if err := quick.Check(TestPodGeneration, nil); err != nil {
	//		t.Error(err)
	//	}
}
