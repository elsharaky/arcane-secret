/*
Copyright 2025 elsharaky.

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

// KeyPairSpec defines the desired state of KeyPair.
type KeyPairSpec struct {
	// Algorithm specifies the cryptographic algorithm for the key pair.
	// Valid values are RSA.
	// +kubebuilder:validation:Enum=RSA
	Algorithm string `json:"algorithm"`

	// Size defines the key size in bits.
	// - RSA: 1024, 2048, 3072, 4096, 8192 (default: 2048)
	// +kubebuilder:validation:Enum=1024;2048;3072;4096;8192
	// +optional
	Size *int `json:"size,omitempty"`

	// SSHFormat enables SSH-compatible key generation for RSA.
	// Only applicable when algorithm=RSA.
	// +kubebuilder:default=false
	// +optional
	SSHFormat bool `json:"sshFormat,omitempty"`
}

// KeyPairStatus defines the observed state of KeyPair.
type KeyPairStatus struct {
	Conditions  []metav1.Condition `json:"conditions,omitempty"`
	Fingerprint string             `json:"fingerprint,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KeyPair is the Schema for the keypairs API.
type KeyPair struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KeyPairSpec   `json:"spec,omitempty"`
	Status KeyPairStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KeyPairList contains a list of KeyPair.
type KeyPairList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeyPair `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeyPair{}, &KeyPairList{})
}
