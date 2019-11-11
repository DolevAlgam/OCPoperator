package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DolevOpSpec defines the desired state of DolevOp
// +k8s:openapi-gen=true
type DolevOpSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Name string `json:"name"`
	Message string `json:"message"`
}

// DolevOpStatus defines the observed state of DolevOp
// +k8s:openapi-gen=true
type DolevOpStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Message string `json:"message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DolevOp is the Schema for the dolevops API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type DolevOp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DolevOpSpec   `json:"spec,omitempty"`
	Status DolevOpStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DolevOpList contains a list of DolevOp
type DolevOpList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DolevOp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DolevOp{}, &DolevOpList{})
}
