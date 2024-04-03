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
	"testing"

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func Test_initExecution(t *testing.T) {
	type args struct {
		ctx       context.Context
		r         *ExecutionReconciler
		execution *pipelinesv1alpha1.Execution
		pvc       *corev1.PersistentVolumeClaim
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := initExecution(tt.args.ctx, tt.args.r, tt.args.execution, tt.args.pvc); (err != nil) != tt.wantErr {
				t.Errorf("initExecution() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
