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

	pipelinesv1alpha1 "github.com/davidlynch-sd/bramble/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

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
