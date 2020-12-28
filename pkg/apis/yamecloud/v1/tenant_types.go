package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BaseTenant defines the desired state of BaseUser
type BaseTenantSpec struct {
	Name        *string  `json:"name,omitempty"`
	Departments []string `json:"departments,omitempty"`
}

// BaseTenantStatus defines the observed state of BaseUser
type BaseTenantStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BaseTenant is the Schema for the baseusers API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=basetenant,scope=Namespaced
type BaseTenant struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BaseTenantSpec   `json:"spec,omitempty"`
	Status BaseTenantStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BaseUserList contains a list of BaseUser
type BaseTenantList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BaseTenant `json:"items"`
}

func init() {
	register(&BaseTenant{}, &BaseUserList{})
}
