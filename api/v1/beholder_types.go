/*
Copyright 2023.

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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BeholderSpec defines the desired state of Beholder
type BeholderSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:default:=1
	ReplicaCount int32 `json:"replicaCount"`
	// +kubebuilder:validation:Required
	Resources Resources `json:"resources"`
	// +kubebuilder:validation:Required
	Image Image `json:"image"`
	// +kubebuilder:validation:Required
	UI UI `json:"ui"`
	// +kubebuilder:default:true
	Redis Redis `json:"redis"`
}

// BeholderStatus defines the observed state of Beholder
type BeholderStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	AssociatedResources []corev1.ObjectReference `json:"associatedResources"`
	Conditions          []metav1.Condition       `json:"conditions"`
	// ActiveReplicas is the number of ready pods
	AvailableReplicas int32 `json:"availableReplicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:singular=beholder,scope=Namespaced

// Beholder is the Schema for the beholders API
type Beholder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BeholderSpec   `json:"spec,omitempty"`
	Status BeholderStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// BeholderList contains a list of Beholder
type BeholderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Beholder `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Beholder{}, &BeholderList{})
}

type Resources struct {
	MemoryLimit   resource.Quantity `json:"memoryLimit"`
	MemoryRequest resource.Quantity `json:"memoryRequest"`
	CPULimit      resource.Quantity `json:"cpuLimit"`
	CPURequest    resource.Quantity `json:"cpuRequest"`
}

type Image struct {
	Registry string `json:"registry"`
	Tag      string `json:"tag"`
}

type UI struct {
	Color   string `json:"color"`
	Message string `json:"message"`
}

type Redis struct {
	Enabled bool `json:"enabled"`
}
