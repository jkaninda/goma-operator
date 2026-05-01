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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RouteSpec defines the desired state of Route.
type RouteSpec struct {
	// Gateways is the list of Gateway CR names this route is attached to
	// (same namespace). A route may be consumed by one or more gateways.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Gateways []string `json:"gateways"`

	// Path is the URL path for this route.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	Path string `json:"path"`

	// Rewrite rewrites the incoming request path.
	// +optional
	Rewrite string `json:"rewrite,omitempty"`

	// Target is the backend URL for this route.
	// +optional
	Target string `json:"target,omitempty"`

	// Methods lists the HTTP methods allowed for this route.
	// +optional
	Methods []string `json:"methods,omitempty"`

	// Enabled controls whether the route is active.
	// +kubebuilder:default=true
	Enabled bool `json:"enabled,omitempty"`

	// Priority determines route matching order (higher = matched first).
	// +optional
	Priority int `json:"priority,omitempty"`

	// Hosts lists domains for host-based routing.
	// +optional
	Hosts []string `json:"hosts,omitempty"`

	// Backends defines multiple backend servers for load balancing.
	// +optional
	Backends []BackendSpec `json:"backends,omitempty"`

	// HealthCheck configures backend health monitoring.
	// +optional
	HealthCheck *HealthCheckSpec `json:"healthCheck,omitempty"`

	// Security defines route-level security settings.
	// +optional
	Security *RouteSecuritySpec `json:"security,omitempty"`

	// Middlewares lists Middleware CR names to apply to this route.
	// +optional
	Middlewares []string `json:"middlewares,omitempty"`

	// DisableMetrics disables Prometheus metrics collection for this route.
	// Useful for silencing high-volume or low-value endpoints (health probes,
	// static assets) from per-route metric series.
	// +optional
	DisableMetrics bool `json:"disableMetrics,omitempty"`

	// TLS references a Kubernetes TLS Secret (type kubernetes.io/tls) whose
	// certificate is served for this route's hosts. When the sidecar is
	// enabled, the cert/key are written to disk and hot-reloaded into the
	// gateway without a pod restart. This is the per-route serving
	// certificate — for backend (upstream) TLS settings see Security.TLS.
	// +optional
	TLS *RouteTLSCertificateSpec `json:"tls,omitempty"`

	// Maintenance puts the route into maintenance mode.
	// +optional
	Maintenance *MaintenanceSpec `json:"maintenance,omitempty"`
}

// BackendSpec defines a backend server for load balancing.
type BackendSpec struct {
	// Endpoint is the backend URL.
	Endpoint string `json:"endpoint"`

	// Weight is the load balancing weight.
	// +optional
	Weight int `json:"weight,omitempty"`

	// Match defines request conditions that pin traffic to this backend
	// (e.g. route a header/cookie/query/IP to a specific backend). All
	// entries must match (AND).
	// +optional
	Match []BackendMatchSpec `json:"match,omitempty"`

	// Exclusive, when true, restricts this backend to traffic that
	// satisfies its Match rules — unmatched requests are never sent here,
	// even during load balancing. When false (default), matching requests
	// are pinned but the backend still participates in the general pool.
	// +optional
	Exclusive bool `json:"exclusive,omitempty"`
}

// BackendMatchSpec is a single request-matching condition for a Backend.
type BackendMatchSpec struct {
	// Source is the part of the request to inspect.
	// +kubebuilder:validation:Enum=header;cookie;query;ip
	Source string `json:"source"`

	// Name of the header/cookie/query parameter (ignored for source=ip).
	// +optional
	Name string `json:"name,omitempty"`

	// Operator is the comparison applied between the source value and Value.
	// +kubebuilder:validation:Enum=equals;not_equals;contains;not_contains;starts_with;ends_with;regex;in
	Operator string `json:"operator"`

	// Value to compare against. For operator=in, use a comma-separated list.
	Value string `json:"value"`
}

// HealthCheckSpec defines health check settings.
type HealthCheckSpec struct {
	// Path is the health check endpoint path.
	Path string `json:"path"`

	// Interval is the check interval (e.g., "30s").
	// +kubebuilder:default="30s"
	Interval string `json:"interval,omitempty"`

	// Timeout is the check timeout (e.g., "5s").
	// +kubebuilder:default="5s"
	Timeout string `json:"timeout,omitempty"`

	// HealthyStatuses is the list of HTTP status codes considered healthy.
	// +optional
	HealthyStatuses []int `json:"healthyStatuses,omitempty"`
}

// RouteSecuritySpec defines route-level security configuration.
type RouteSecuritySpec struct {
	// ForwardHostHeaders controls forwarding of X-Forwarded-Host and related headers.
	// +kubebuilder:default=true
	ForwardHostHeaders bool `json:"forwardHostHeaders,omitempty"`

	// EnableExploitProtection enables SQL injection and XSS protection.
	// +optional
	EnableExploitProtection bool `json:"enableExploitProtection,omitempty"`

	// TLS defines route-level TLS settings for backend connections.
	// +optional
	TLS *RouteTLSSpec `json:"tls,omitempty"`
}

// RouteTLSSpec defines TLS settings for backend connections.
type RouteTLSSpec struct {
	// InsecureSkipVerify disables TLS verification for the backend.
	// +optional
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`

	// RootCAsSecret is the name of a K8s Secret containing CA certificates.
	// +optional
	RootCAsSecret string `json:"rootCAsSecret,omitempty"`

	// ClientCertSecret is the name of a K8s TLS Secret for mTLS client certificates.
	// +optional
	ClientCertSecret string `json:"clientCertSecret,omitempty"`
}

// RouteTLSCertificateSpec references a Kubernetes TLS Secret whose cert/key
// pair is presented to clients connecting to this route's hosts. The Secret
// must be of type kubernetes.io/tls with keys "tls.crt" and "tls.key".
type RouteTLSCertificateSpec struct {
	// SecretName is the name of a kubernetes.io/tls Secret in the same
	// namespace as the Route. The gateway will serve this certificate for
	// the route's hosts.
	// +kubebuilder:validation:Required
	SecretName string `json:"secretName"`
}

// MaintenanceSpec defines maintenance mode settings.
type MaintenanceSpec struct {
	// Enabled controls whether maintenance mode is active.
	Enabled bool `json:"enabled"`

	// Body is the response body returned during maintenance.
	// +optional
	Body string `json:"body,omitempty"`

	// Status is the HTTP status code returned during maintenance.
	// +kubebuilder:default=503
	Status int `json:"status,omitempty"`
}

// RouteStatus defines the observed state of Route.
type RouteStatus struct {
	// Ready indicates whether the route is valid and synced.
	Ready bool `json:"ready,omitempty"`

	// Synced indicates whether the route config has been written to the Gateway's ConfigMap.
	Synced bool `json:"synced,omitempty"`

	// GatewayGeneration tracks which Gateway config generation includes this route.
	GatewayGeneration int64 `json:"gatewayGeneration,omitempty"`

	// Conditions represent the latest available observations of the Route's state.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Gateways",type="string",JSONPath=".spec.gateways"
// +kubebuilder:printcolumn:name="Path",type="string",JSONPath=".spec.path"
// +kubebuilder:printcolumn:name="Target",type="string",JSONPath=".spec.target"
// +kubebuilder:printcolumn:name="Ready",type="boolean",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

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
