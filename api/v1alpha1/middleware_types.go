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

package v1alpha1

import (
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MiddlewareSpec defines the desired state of Middleware.
type MiddlewareSpec struct {
	// Type indicates the middleware type.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=basic;jwt;oauth;rateLimit;access;accessPolicy;addPrefix;redirectRegex;rewriteRegex;forwardAuth;httpCache;redirectScheme;bodyLimit;responseHeaders;errorInterceptor;ldap;userAgentBlock
	Type string `json:"type"`

	// Paths lists the route paths this middleware protects.
	// +optional
	Paths []string `json:"paths,omitempty"`

	// Rule contains the type-specific middleware configuration.
	// The structure depends on the middleware Type.
	// +optional
	// +kubebuilder:pruning:PreserveUnknownFields
	Rule *apiextensionsv1.JSON `json:"rule,omitempty"`
}

// MiddlewareStatus defines the observed state of Middleware.
type MiddlewareStatus struct {
	// Ready indicates whether the middleware configuration is valid.
	Ready bool `json:"ready,omitempty"`

	// ReferencedBy lists the names of Routes that use this middleware.
	// +optional
	ReferencedBy []string `json:"referencedBy,omitempty"`

	// Conditions represent the latest available observations of the Middleware's state.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.type"
// +kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

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
