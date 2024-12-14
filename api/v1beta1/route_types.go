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

// RouteSpec defines the desired state of a route.
type RouteSpec struct {
	// Gateway specifies the name of the Gateway resource associated with this route.
	Gateway string `json:"gateway"`
	// Path specifies the route path.
	Path string `json:"path" yaml:"path"`
	// Hosts defines a list of domains or hosts for host-based request routing.
	Hosts []string `json:"hosts,omitempty" yaml:"hosts"`
	// Rewrite specifies the new path to rewrite the incoming route path to.
	Rewrite string `json:"rewrite,omitempty" yaml:"rewrite"`
	// Methods specifies the HTTP methods allowed for this route (e.g., GET, POST, PUT).
	Methods []string `json:"methods,omitempty" yaml:"methods"`
	// Destination defines the backend URL to which requests will be proxied.
	Destination string `json:"destination,omitempty" yaml:"destination"`
	// Backends specifies a list of backend URLs for load balancing.
	Backends []string `json:"backends,omitempty" yaml:"backends"`
	// InsecureSkipVerify allows skipping TLS certificate verification for backend connections.
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty" yaml:"insecureSkipVerify"`
	// HealthCheck defines the settings for backend health checks.
	HealthCheck RouteHealthCheck `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty"`
	// Cors specifies the CORS (Cross-Origin Resource Sharing) configuration for the route.
	Cors Cors `json:"cors,omitempty" yaml:"cors"`
	// RateLimit defines the maximum number of requests allowed per minute for this route.
	RateLimit int `json:"rateLimit,omitempty" yaml:"rateLimit"`
	// DisableHostFording disables host forwarding for this route.
	// Deprecated: Use DisableHostForwarding instead.
	DisableHostFording bool `json:"disableHostFording,omitempty" yaml:"disableHostFording"`
	// InterceptErrors specifies a list of HTTP status codes for which backend errors should be intercepted.
	// Deprecated: Use ErrorInterceptor for advanced error handling.
	InterceptErrors []int `json:"interceptErrors,omitempty" yaml:"interceptErrors"`
	// DisableHostForwarding disables forwarding the host header to the backend.
	DisableHostForwarding bool `json:"disableHostForwarding,omitempty" yaml:"disableHostForwarding"`
	// ErrorInterceptor defines the configuration for handling backend error interception.
	ErrorInterceptor RouteErrorInterceptor `yaml:"errorInterceptor,omitempty" json:"errorInterceptor,omitempty"`
	// BlockCommonExploits enables or disables blocking common exploits, such as basic SQL injection or XSS attempts.
	BlockCommonExploits bool `json:"blockCommonExploits,omitempty" yaml:"blockCommonExploits"`
	// Middlewares specifies a list of middleware names to apply to this route.
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
