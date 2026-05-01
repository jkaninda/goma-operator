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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GatewaySpec defines the desired state of Gateway.
type GatewaySpec struct {
	// Image is the full container image reference for goma-gateway.
	// +kubebuilder:default="jkaninda/goma-gateway:latest"
	Image string `json:"image,omitempty"`

	// Replicas is the desired number of gateway pods.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=0
	Replicas *int32 `json:"replicas,omitempty"`

	// ImagePullSecrets is a list of references to secrets for pulling the gateway image.
	// +optional
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`

	// Resources defines the compute resource requirements for the gateway container.
	// +optional
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// Affinity defines scheduling constraints for gateway pods.
	// +optional
	Affinity *corev1.Affinity `json:"affinity,omitempty"`

	// AutoScaling defines the horizontal pod autoscaling configuration.
	// +optional
	AutoScaling *AutoScalingSpec `json:"autoScaling,omitempty"`

	// Server defines the gateway server configuration.
	// +optional
	Server ServerSpec `json:"server,omitempty"`

	// Providers defines dynamic configuration providers (file/k8s sidecar, http, git).
	// +optional
	Providers *ProvidersSpec `json:"providers,omitempty"`

	// CertManager configures Goma Gateway's built-in ACME / Let's Encrypt
	// certificate manager. When enabled, the gateway will automatically
	// obtain and renew certificates for the domains listed in Route hosts.
	// +optional
	CertManager *CertManagerSpec `json:"certManager,omitempty"`

	// Service configures the Kubernetes Service created for this gateway.
	// Use this to expose the gateway externally (LoadBalancer / NodePort)
	// and/or remap the Service ports (e.g. to the standard 80 / 443 for
	// Ingress-style exposure). Container ports are unchanged (8080 / 8443).
	// +optional
	Service *ServiceSpec `json:"service,omitempty"`
}

// ServiceSpec configures the Kubernetes Service created for the gateway.
type ServiceSpec struct {
	// Type is the Service type: ClusterIP, NodePort, or LoadBalancer.
	// +kubebuilder:default=ClusterIP
	// +kubebuilder:validation:Enum=ClusterIP;NodePort;LoadBalancer
	// +optional
	Type corev1.ServiceType `json:"type,omitempty"`

	// HTTPPort is the external port for HTTP traffic (the Service port).
	// Defaults to 8080 to match the gateway's container port. Set to 80
	// for Ingress-style exposure.
	// +kubebuilder:default=8080
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +optional
	HTTPPort int32 `json:"httpPort,omitempty"`

	// HTTPSPort is the external port for HTTPS traffic (the Service port).
	// Defaults to 8443 to match the gateway's container port. Set to 443
	// for Ingress-style exposure.
	// +kubebuilder:default=8443
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +optional
	HTTPSPort int32 `json:"httpsPort,omitempty"`

	// HTTPNodePort is the NodePort for HTTP (only honored when type=NodePort).
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65535
	// +optional
	HTTPNodePort int32 `json:"httpNodePort,omitempty"`

	// HTTPSNodePort is the NodePort for HTTPS (only honored when type=NodePort).
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=65535
	// +optional
	HTTPSNodePort int32 `json:"httpsNodePort,omitempty"`

	// LoadBalancerIP requests a specific static IP (cloud-dependent).
	// +optional
	LoadBalancerIP string `json:"loadBalancerIP,omitempty"`

	// LoadBalancerSourceRanges restricts LoadBalancer access to specific CIDRs.
	// +optional
	LoadBalancerSourceRanges []string `json:"loadBalancerSourceRanges,omitempty"`

	// LoadBalancerClass selects a specific LB implementation
	// (e.g. "service.k8s.aws/nlb").
	// +optional
	LoadBalancerClass *string `json:"loadBalancerClass,omitempty"`

	// ExternalTrafficPolicy: Cluster (default) or Local (preserves source IP,
	// only routes to nodes with pods).
	// +kubebuilder:validation:Enum=Cluster;Local
	// +optional
	ExternalTrafficPolicy corev1.ServiceExternalTrafficPolicy `json:"externalTrafficPolicy,omitempty"`

	// SessionAffinity: None (default) or ClientIP.
	// +kubebuilder:validation:Enum=None;ClientIP
	// +optional
	SessionAffinity corev1.ServiceAffinity `json:"sessionAffinity,omitempty"`

	// Annotations are merged onto the Service (for provider-specific config
	// like AWS NLB, GCP, MetalLB, etc.).
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels are merged onto the Service.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// IPFamilyPolicy: SingleStack, PreferDualStack, or RequireDualStack.
	// +optional
	IPFamilyPolicy *corev1.IPFamilyPolicy `json:"ipFamilyPolicy,omitempty"`

	// IPFamilies is the list of IP families (IPv4, IPv6).
	// +optional
	IPFamilies []corev1.IPFamily `json:"ipFamilies,omitempty"`
}

// ServiceType returns the configured Service type, defaulting to ClusterIP.
func (s *GatewaySpec) ServiceType() corev1.ServiceType {
	if s.Service == nil || s.Service.Type == "" {
		return corev1.ServiceTypeClusterIP
	}
	return s.Service.Type
}

// HTTPPort returns the Service-level HTTP port, defaulting to the
// container's port (8080).
func (s *GatewaySpec) HTTPPort() int32 {
	if s.Service != nil && s.Service.HTTPPort > 0 {
		return s.Service.HTTPPort
	}
	return 8080
}

// HTTPSPort returns the Service-level HTTPS port, defaulting to the
// container's port (8443).
func (s *GatewaySpec) HTTPSPort() int32 {
	if s.Service != nil && s.Service.HTTPSPort > 0 {
		return s.Service.HTTPSPort
	}
	return 8443
}

// CertManagerSpec configures the gateway's ACME certificate manager.
type CertManagerSpec struct {
	// Provider selects the cert manager backend. Currently only "acme"
	// is supported.
	// +kubebuilder:default="acme"
	// +kubebuilder:validation:Enum=acme
	Provider string `json:"provider,omitempty"`

	// ACME configures Let's Encrypt / ACME certificate issuance.
	// +optional
	ACME *ACMESpec `json:"acme,omitempty"`
}

// ACMESpec configures ACME / Let's Encrypt certificate issuance.
type ACMESpec struct {
	// Email is the contact email for the ACME account. Required for
	// Let's Encrypt expiration warnings and account recovery.
	// +kubebuilder:validation:Required
	Email string `json:"email"`

	// DirectoryURL overrides the ACME directory endpoint (defaults to
	// Let's Encrypt production). Use the staging URL for testing:
	// https://acme-staging-v02.api.letsencrypt.org/directory
	// +optional
	DirectoryURL string `json:"directoryUrl,omitempty"`

	// TermsAccepted indicates acceptance of the ACME provider's Terms
	// of Service. Must be true for Let's Encrypt.
	// +kubebuilder:default=true
	TermsAccepted bool `json:"termsAccepted,omitempty"`

	// ChallengeType is the ACME challenge solver: "http-01" or "dns-01".
	// Use "dns-01" for wildcard certificates.
	// +kubebuilder:default="http-01"
	// +kubebuilder:validation:Enum=http-01;dns-01
	ChallengeType string `json:"challengeType,omitempty"`

	// DNSProvider is the DNS-01 challenge provider (e.g., "cloudflare",
	// "route53", "digitalocean"). Required when ChallengeType is "dns-01".
	// +optional
	DNSProvider string `json:"dnsProvider,omitempty"`

	// CredentialsSecret references a Kubernetes Secret containing DNS
	// provider credentials. The Secret must have an "apiToken" key (and
	// any other provider-specific keys mounted as environment variables).
	// Required when ChallengeType is "dns-01".
	// +optional
	CredentialsSecret string `json:"credentialsSecret,omitempty"`
}

// AutoScalingSpec defines HPA configuration for the gateway.
type AutoScalingSpec struct {
	// Enabled controls whether HPA is created.
	Enabled bool `json:"enabled"`

	// MinReplicas is the lower limit for the number of replicas.
	// +kubebuilder:default=1
	// +kubebuilder:validation:Minimum=1
	MinReplicas *int32 `json:"minReplicas,omitempty"`

	// MaxReplicas is the upper limit for the number of replicas.
	// +kubebuilder:default=10
	// +kubebuilder:validation:Minimum=1
	MaxReplicas int32 `json:"maxReplicas"`

	// TargetCPUUtilizationPercentage is the target average CPU utilization.
	// +optional
	TargetCPUUtilizationPercentage *int32 `json:"targetCPUUtilizationPercentage,omitempty"`

	// TargetMemoryUtilizationPercentage is the target average memory utilization.
	// +optional
	TargetMemoryUtilizationPercentage *int32 `json:"targetMemoryUtilizationPercentage,omitempty"`
}

// ServerSpec defines the gateway server configuration.
type ServerSpec struct {
	// Timeouts defines server timeout settings in seconds.
	// +optional
	Timeouts TimeoutSpec `json:"timeouts,omitempty"`

	// LogLevel sets the logging level (info, debug, trace, off).
	// +kubebuilder:default="info"
	// +kubebuilder:validation:Enum=info;debug;trace;off
	LogLevel string `json:"logLevel,omitempty"`

	// Redis contains the configuration for a Redis backend.
	// +optional
	Redis *RedisSpec `json:"redis,omitempty"`

	// Monitoring defines metrics and health check configuration.
	// +optional
	Monitoring *MonitoringSpec `json:"monitoring,omitempty"`

	// TLS is a list of TLS certificate configurations via Kubernetes Secrets.
	// +optional
	TLS []TLSSpec `json:"tls,omitempty"`

	// Networking defines transport and DNS cache settings.
	// +optional
	Networking *NetworkingSpec `json:"networking,omitempty"`
}

// TimeoutSpec defines server timeout values in seconds.
type TimeoutSpec struct {
	// Read timeout in seconds.
	// +kubebuilder:default=30
	Read int `json:"read,omitempty"`

	// Write timeout in seconds.
	// +kubebuilder:default=60
	Write int `json:"write,omitempty"`

	// Idle timeout in seconds.
	// +kubebuilder:default=90
	Idle int `json:"idle,omitempty"`
}

// RedisSpec defines Redis connection settings.
type RedisSpec struct {
	// Addr is the Redis server address (host:port).
	Addr string `json:"addr"`

	// Password is the Redis password.
	// +optional
	Password string `json:"password,omitempty"`
}

// MonitoringSpec defines observability configuration.
//
// /readyz and /healthz are always enabled — the operator's Deployment
// uses them for Kubernetes readiness/liveness probes.
type MonitoringSpec struct {
	// EnableMetrics enables Prometheus metrics collection.
	// +kubebuilder:default=false
	EnableMetrics bool `json:"enableMetrics,omitempty"`

	// MetricsPath sets the path for the metrics endpoint.
	// +kubebuilder:default="/metrics"
	MetricsPath string `json:"metricsPath,omitempty"`

	// Host restricts access to the monitoring endpoints (/metrics,
	// /readyz, /healthz, /healthz/routes) to this hostname. When set,
	// requests with a different Host header are rejected.
	// +optional
	Host string `json:"host,omitempty"`

	// Middleware attaches Goma Middleware CRs to the monitoring endpoints
	// for authentication / access control (e.g. protect /metrics with a
	// basic-auth or JWT middleware).
	// +optional
	Middleware *MonitoringMiddlewareSpec `json:"middleware,omitempty"`
}

// MonitoringMiddlewareSpec attaches middlewares to specific monitoring
// endpoints. Each entry is the name of a Middleware CR in the same namespace
// as the Gateway.
type MonitoringMiddlewareSpec struct {
	// Metrics is the list of Middleware CR names applied to the /metrics
	// endpoint. Useful for restricting Prometheus scraping to authenticated
	// clients.
	// +optional
	Metrics []string `json:"metrics,omitempty"`
}

// TLSSpec references a Kubernetes TLS Secret.
type TLSSpec struct {
	// SecretName is the name of a Kubernetes Secret of type kubernetes.io/tls.
	SecretName string `json:"secretName"`
}

// NetworkingSpec defines transport and DNS settings.
type NetworkingSpec struct {
	// DNSCache configures DNS caching.
	// +optional
	DNSCache *DNSCacheSpec `json:"dnsCache,omitempty"`

	// Transport configures HTTP transport settings.
	// +optional
	Transport *TransportSpec `json:"transport,omitempty"`
}

// DNSCacheSpec defines DNS cache settings.
type DNSCacheSpec struct {
	// TTL is the cache TTL in seconds.
	TTL int `json:"ttl,omitempty"`
}

// TransportSpec defines HTTP transport configuration.
type TransportSpec struct {
	// MaxIdleConns controls the maximum number of idle connections.
	// +kubebuilder:default=512
	MaxIdleConns int `json:"maxIdleConns,omitempty"`

	// MaxIdleConnsPerHost controls the maximum idle connections per host.
	// +kubebuilder:default=256
	MaxIdleConnsPerHost int `json:"maxIdleConnsPerHost,omitempty"`

	// MaxConnsPerHost limits total connections per host.
	// +kubebuilder:default=256
	MaxConnsPerHost int `json:"maxConnsPerHost,omitempty"`
}

// ProvidersSpec defines all dynamic configuration providers.
type ProvidersSpec struct {
	// Kubernetes configures the goma-k8s-provider sidecar which watches
	// Route/Middleware CRDs and writes config to /etc/goma/providers/k8s.
	// +optional
	Kubernetes *KubernetesProviderSpec `json:"kubernetes,omitempty"`

	// HTTP configures a remote HTTP provider that periodically fetches
	// configuration from a URL.
	// +optional
	HTTP *HTTPProviderSpec `json:"http,omitempty"`

	// Git configures a Git-based provider that pulls configuration from
	// a repository.
	// +optional
	Git *GitProviderSpec `json:"git,omitempty"`
}

// KubernetesProviderSpec defines the K8s provider sidecar configuration.
//
// The sidecar is enabled by default. To disable it, set `enabled: false`.
type KubernetesProviderSpec struct {
	// Enabled controls whether the goma-k8s-provider sidecar is injected
	// into the gateway pod. Defaults to true when the spec is omitted.
	// Set to false to disable the sidecar and deliver routes/middlewares
	// only through the static ConfigMap (restart on every change).
	// +kubebuilder:default=true
	// +optional
	Enabled *bool `json:"enabled,omitempty"`

	// Image is the sidecar container image.
	// +kubebuilder:default="jkaninda/goma-k8s-provider:latest"
	Image string `json:"image,omitempty"`
}

// KubernetesProviderEnabled returns true when the goma-k8s-provider sidecar
// should be injected for this Gateway. Defaults to true unless explicitly
// disabled via spec.providers.kubernetes.enabled=false.
func (s *GatewaySpec) KubernetesProviderEnabled() bool {
	if s.Providers == nil || s.Providers.Kubernetes == nil {
		return true
	}
	if s.Providers.Kubernetes.Enabled == nil {
		return true
	}
	return *s.Providers.Kubernetes.Enabled
}

// KubernetesProviderImage returns the sidecar image, falling back to the
// default when unset.
func (s *GatewaySpec) KubernetesProviderImage() string {
	if s.Providers != nil && s.Providers.Kubernetes != nil && s.Providers.Kubernetes.Image != "" {
		return s.Providers.Kubernetes.Image
	}
	return "jkaninda/goma-k8s-provider:latest"
}

// HTTPProviderSpec defines the HTTP provider configuration.
type HTTPProviderSpec struct {
	// Enabled controls whether the HTTP provider is active.
	Enabled bool `json:"enabled"`

	// Endpoint is the URL to fetch configuration from.
	// +kubebuilder:validation:Required
	Endpoint string `json:"endpoint"`

	// Interval between fetches (e.g., "60s", "5m").
	// +kubebuilder:default="60s"
	Interval string `json:"interval,omitempty"`

	// Timeout for each fetch request (e.g., "10s").
	// +kubebuilder:default="10s"
	Timeout string `json:"timeout,omitempty"`

	// Headers are additional HTTP headers sent with each request.
	// For secrets, prefer HeadersSecret.
	// +optional
	Headers map[string]string `json:"headers,omitempty"`

	// HeadersSecret references a K8s Secret whose keys/values are mounted
	// as environment variables and can be referenced via ${VAR} in Headers.
	// +optional
	HeadersSecret string `json:"headersSecret,omitempty"`

	// InsecureSkipVerify disables TLS verification.
	// +optional
	InsecureSkipVerify bool `json:"insecureSkipVerify,omitempty"`

	// RetryAttempts is the number of retries on failure.
	// +kubebuilder:default=3
	RetryAttempts int `json:"retryAttempts,omitempty"`

	// RetryDelay between retries (e.g., "2s").
	// +kubebuilder:default="2s"
	RetryDelay string `json:"retryDelay,omitempty"`

	// CacheDir is where responses are cached.
	// +optional
	CacheDir string `json:"cacheDir,omitempty"`
}

// GitProviderSpec defines the Git provider configuration.
type GitProviderSpec struct {
	// Enabled controls whether the Git provider is active.
	Enabled bool `json:"enabled"`

	// URL is the Git repository URL.
	// +kubebuilder:validation:Required
	URL string `json:"url"`

	// Branch to check out.
	// +kubebuilder:default="main"
	Branch string `json:"branch,omitempty"`

	// Path is the subdirectory inside the repo containing config files.
	// +optional
	Path string `json:"path,omitempty"`

	// Interval between pulls (e.g., "60s").
	// +kubebuilder:default="60s"
	Interval string `json:"interval,omitempty"`

	// Auth configures Git authentication.
	// +optional
	Auth *GitAuthSpec `json:"auth,omitempty"`

	// CloneDir is the local directory where the repo is cloned.
	// +optional
	CloneDir string `json:"cloneDir,omitempty"`
}

// GitAuthSpec defines Git authentication settings.
// Secret values (token, password, SSH key data) must come from the referenced
// Kubernetes Secret.
type GitAuthSpec struct {
	// Type is the auth type: "token", "basic", or "ssh".
	// +kubebuilder:validation:Enum=token;basic;ssh
	Type string `json:"type"`

	// SecretName references a K8s Secret containing auth material:
	//   - type=token: key "token"
	//   - type=basic: keys "username" and "password"
	//   - type=ssh:   key "ssh-privatekey"
	// +kubebuilder:validation:Required
	SecretName string `json:"secretName"`
}

// GatewayStatus defines the observed state of Gateway.
type GatewayStatus struct {
	// Replicas is the current number of gateway pods.
	Replicas int32 `json:"replicas,omitempty"`

	// ReadyReplicas is the number of ready gateway pods.
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`

	// Routes is the number of Route CRs associated with this gateway.
	Routes int32 `json:"routes,omitempty"`

	// Middlewares is the number of Middleware CRs referenced by the routes.
	Middlewares int32 `json:"middlewares,omitempty"`

	// ConfigChecksum is the SHA256 checksum of the generated config.
	ConfigChecksum string `json:"configChecksum,omitempty"`

	// Addresses is the list of addresses at which the gateway is reachable.
	// Populated from the Service status (LoadBalancer ingress, external IPs,
	// node IPs with NodePort, or cluster IP depending on Service type).
	// +optional
	Addresses []GatewayAddress `json:"addresses,omitempty"`

	// Conditions represent the latest available observations of the Gateway's state.
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// GatewayAddress is a single reachable address for the gateway.
type GatewayAddress struct {
	// Type is the address type: "IPAddress" or "Hostname".
	// +kubebuilder:validation:Enum=IPAddress;Hostname
	Type string `json:"type"`

	// Value is the address itself (IP literal or DNS name).
	Value string `json:"value"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Type",type="string",JSONPath=".spec.service.type"
// +kubebuilder:printcolumn:name="Address",type="string",JSONPath=".status.addresses[0].value"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".status.replicas"
// +kubebuilder:printcolumn:name="Ready",type="integer",JSONPath=".status.readyReplicas"
// +kubebuilder:printcolumn:name="Routes",type="integer",JSONPath=".status.routes"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Gateway is the Schema for the gateways API.
type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GatewaySpec   `json:"spec,omitempty"`
	Status GatewayStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GatewayList contains a list of Gateway.
type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gateway `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Gateway{}, &GatewayList{})
}
