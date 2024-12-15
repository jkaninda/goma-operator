/*
Copyright 2024.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MiddlewareSpec defines the desired configuration for middleware in the system.
type MiddlewareSpec struct {
	// Type specifies the type of middleware to be applied.
	// Available values:
	// - "basic": Basic authentication.
	// - "jwt": JSON Web Token (JWT) authentication.
	// - "auth": Authentication using Auth0 service.
	// - "rateLimit": Middleware for rate-limiting requests.
	// - "access": General access control middleware.
	// - "accessPolicy": Middleware for IP-based access policies.
	Type string `json:"type" yaml:"type"` // Type of middleware to apply [basic, jwt, auth0, rateLimit, access, accessPolicy]

	// Paths defines the list of paths to which the middleware will be applied.
	// These paths will be protected by the middleware specified in the 'Type' field.
	Paths []string `json:"paths,omitempty" yaml:"paths,omitempty"` // List of paths to protect with the middleware

	// Rule contains the specific rule or configuration for the middleware.
	// This field allows for flexible rule configurations, such as access control or rate limiting.
	// It is represented as a RawExtension to accommodate different formats.
	// The content of this field depends on the middleware type and is optional.
	Rule runtime.RawExtension `json:"rule,omitempty" yaml:"rule,omitempty"` // Specific middleware rule or configuration
}

// MiddlewareStatus defines the observed state of Middleware.
type MiddlewareStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Middleware is the Schema for the middlewares API.
type Middleware struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MiddlewareSpec   `json:"spec,omitempty"`
	Status MiddlewareStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MiddlewareList contains a list of Middleware.
type MiddlewareList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Middleware `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Middleware{}, &MiddlewareList{})
}
