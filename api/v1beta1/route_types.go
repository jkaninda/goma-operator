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
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RouteSpec defines the desired state of Route.
type RouteSpec struct {
	// Gateway defines the name of the Gateway resource
	Gateway string `json:"gateway"`
	// Path defines route path
	Path string `json:"path" yaml:"path"`
	// Hosts Domains/hosts based request routing
	Hosts []string `json:"hosts,omitempty" yaml:"hosts"`
	// Rewrite rewrites route path to desired path
	Rewrite string `json:"rewrite,omitempty" yaml:"rewrite"`
	// Methods allowed method
	Methods []string `json:"methods,omitempty" yaml:"methods"`
	// Destination Defines backend URL
	Destination        string   `json:"destination,omitempty" yaml:"destination"`
	Backends           []string `json:"backends,omitempty" yaml:"backends"`
	InsecureSkipVerify bool     `json:"insecureSkipVerify,omitempty" yaml:"insecureSkipVerify"`
	// HealthCheck Defines the backend is health
	HealthCheck RouteHealthCheck `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty"`
	// Cors contains the route cors headers
	Cors      Cors `json:"cors,omitempty" yaml:"cors"`
	RateLimit int  `json:"rateLimit,omitempty" yaml:"rateLimit"`
	// DisableHostFording Disables host forwarding.
	DisableHostFording bool `json:"disableHostFording,omitempty" yaml:"disableHostFording"`
	// InterceptErrors intercepts backend errors based on the status codes
	InterceptErrors []int `json:"interceptErrors,omitempty" yaml:"interceptErrors"`
	// BlockCommonExploits enable, disable block common exploits
	BlockCommonExploits bool `json:"blockCommonExploits,omitempty" yaml:"blockCommonExploits"`
	// Middlewares Defines route middleware
	Middlewares []string `json:"middlewares,omitempty" yaml:"middlewares"`
}

// RouteStatus defines the observed state of Route.
type RouteStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Route is the Schema for the routes API.
type Route struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RouteSpec   `json:"spec,omitempty"`
	Status RouteStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RouteList contains a list of Route.
type RouteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Route `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Route{}, &RouteList{})
}
