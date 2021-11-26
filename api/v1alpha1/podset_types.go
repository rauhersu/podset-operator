/*
Copyright 2021.

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

// PodsetSpec defines the desired state of Podset
type PodsetSpec struct {
	// Replicas is the desired number of pods for the PodSet
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	Replicas int32 `json:"replicas,omitempty"`
}

// PodSetStatus defines the current status of PodSet
type PodsetStatus struct {
	PodNames        []string        `json:"podNames"`
	AvailableReplicas    int32    `json:"availableReplicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Podset is the Schema for the podsets API
// +kubebuilder:printcolumn:JSONPath=".spec.replicas",name=Desired,type=string
// +kubebuilder:printcolumn:JSONPath=".status.availableReplicas",name=Available,type=string
type Podset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodsetSpec   `json:"spec,omitempty"`
	Status PodsetStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodsetList contains a list of Podset
type PodsetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Podset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Podset{}, &PodsetList{})
}
