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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type TaskRef struct {
	Name         string   `json:"name"`
	Dependencies []string `json:"dependencies,omitempty"`
}
type TaskSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Docker image which will be used.
	Image string `json:"image"`

	// Command executed by the container,
	// can be used to determine the behaviour of a CLI app.
	Command []string `json:"command"`

	Dependencies []string `json:"dependencies,omitempty"`
}

type PLTask struct {
	// Name of task to be ran
	Name string `json:"name"`
	// Spec of given task
	Spec TaskSpec `json:"spec"`
}

// PipelineSpec defines the desired state of Pipeline
type PipelineSpec struct {
	// Allows developers to create a list of tasks
	Tasks []PLTask `json:"tasks,omitempty"`
	// Allows developers to use pre applied tasks in the same namespace
	TaskRefs []TaskRef `json:"taskRefs,omitempty"`
}

// PipelineStatus defines the observed state of Pipeline
type PipelineStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ValidDeps    bool `json:"validdeps,omitempty" default:"false"`
	TasksCreated bool `json:"taskscreated,omitempty" default:"false"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Pipeline is the Schema for the pipelines API
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec,omitempty"`
	Status PipelineStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PipelineList contains a list of Pipeline
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}
