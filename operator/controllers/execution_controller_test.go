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
	"reflect"
	"testing"

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestExecutionReconciler_Reconcile(t *testing.T) {
	type args struct {
		ctx context.Context
		req ctrl.Request
	}
	tests := []struct {
		name       string
		reconciler *ExecutionReconciler
		args       args
		want       ctrl.Result
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.reconciler.Reconcile(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecutionReconciler.Reconcile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecutionReconciler.Reconcile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loadPipeline(t *testing.T) {
	type args struct {
		execution  *pipelinesv1alpha1.Execution
		reconciler *ExecutionReconciler
		ctx        context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *pipelinesv1alpha1.Pipeline
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadPipeline(tt.args.execution, tt.args.reconciler, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("loadPipeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loadPipeline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_generateListOptions(t *testing.T) {
	type args struct {
		execution *pipelinesv1alpha1.Execution
	}
	tests := []struct {
		name string
		args args
		want *client.ListOptions
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateListOptions(tt.args.execution); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateListOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getExecutionPods(t *testing.T) {
	type args struct {
		listOptions *client.ListOptions
		reconciler  *ExecutionReconciler
		ctx         context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *corev1.PodList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getExecutionPods(tt.args.listOptions, tt.args.reconciler, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getExecutionPods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getExecutionPods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getExecutionPvcs(t *testing.T) {
	type args struct {
		listOptions *client.ListOptions
		reconciler  *ExecutionReconciler
		ctx         context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *corev1.PersistentVolumeClaimList
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getExecutionPvcs(tt.args.listOptions, tt.args.reconciler, tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("getExecutionPvcs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getExecutionPvcs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExecutionReconciler_SetupWithManager(t *testing.T) {
	type args struct {
		mgr ctrl.Manager
	}
	tests := []struct {
		name       string
		reconciler *ExecutionReconciler
		args       args
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.reconciler.SetupWithManager(tt.args.mgr); (err != nil) != tt.wantErr {
				t.Errorf("ExecutionReconciler.SetupWithManager() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
