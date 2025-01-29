package v1beta1

type GatewayConfig struct {
	GatewayConf GatewaySpec `yaml:"gateway,omitempty"`
}

type Cors struct {
	// Cors contains Allowed origins,
	Origins []string `json:"origins,omitempty" yaml:"origins,omitempty"`
	// Headers contains custom headers
	Headers map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
}

type Tls struct {
	CredentialName string `json:"credentialName"` // CredentialName fetches certs from Kubernetes secret
}

type AutoScaling struct {
	Enabled                           bool  `json:"enabled,omitempty"`
	MinReplicas                       int32 `json:"minReplicas,omitempty"`
	MaxReplicas                       int32 `json:"maxReplicas,omitempty"`
	TargetCPUUtilizationPercentage    int32 `json:"targetCPUUtilizationPercentage,omitempty"`
	TargetMemoryUtilizationPercentage int32 `json:"targetMemoryUtilizationPercentage,omitempty"`
}

type Middlewares struct {
	Middlewares []MiddlewareSpec `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
}

type Server struct {
	// WriteTimeout specifies the proxy's write timeout in seconds.
	WriteTimeout int `json:"writeTimeout,omitempty" yaml:"writeTimeout,omitempty"`
	// ReadTimeout specifies the proxy's read timeout in seconds.
	ReadTimeout int `json:"readTimeout,omitempty" yaml:"readTimeout,omitempty"`
	// IdleTimeout defines the proxy's idle timeout in seconds.
	IdleTimeout int `json:"idleTimeout,omitempty" yaml:"idleTimeout,omitempty"`
	// LogLevel specifies the logging level for the proxy. Accepted values: "info", "debug", "trace", "off".
	LogLevel string `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	// TlsSecretName specifies the name of the secret containing the TLS certificate and key.
	// Deprecated: Use TLS instead.
	TlsSecretName string `json:"tlsSecretName,omitempty" yaml:"tlsSecretName,omitempty"`

	// Redis contains the configuration details for connecting to a Redis database.
	Redis Redis `json:"redis,omitempty" yaml:"redis,omitempty"`
	// TLS contains the TLS configuration for the proxy.
	TLS *TLS `json:"tls,omitempty" yaml:"tls"`
	// Cors holds the global CORS (Cross-Origin Resource Sharing) configuration for the proxy.
	Cors Cors `json:"cors,omitempty" yaml:"cors,omitempty"`
	// ErrorInterceptor defines the configuration for intercepting and handling backend errors.
	ErrorInterceptor RouteErrorInterceptor `json:"errorInterceptor,omitempty" yaml:"errorInterceptor,omitempty"`
	// DisableHealthCheckStatus enables or disables health checks for routes.
	DisableHealthCheckStatus bool `json:"disableHealthCheckStatus,omitempty" yaml:"disableHealthCheckStatus"`
	// DisableKeepAlive enables or disables the server's KeepAlive connections.
	DisableKeepAlive bool `json:"disableKeepAlive,omitempty" yaml:"disableKeepAlive"`
	// EnableMetrics toggles the collection and exposure of server metrics.
	EnableMetrics bool `json:"enableMetrics,omitempty" yaml:"enableMetrics"`
	// EnableStrictSlash enables or disables strict routing and trailing slashes.
	//
	// When enabled, the router will match the path with or without a trailing slash.
	EnableStrictSlash bool `json:"enableStrictSlash,omitempty" yaml:"enableStrictSlash,omitempty"`
}

type RouteHealthCheck struct {
	Path            string `json:"path,omitempty" yaml:"path"`
	Interval        string `json:"interval,omitempty" yaml:"interval"`
	Timeout         string `json:"timeout,omitempty" yaml:"timeout"`
	HealthyStatuses []int  `json:"healthyStatuses,omitempty" yaml:"healthyStatuses"`
}
type Redis struct {
	// Addr redis hostname and port number :
	Addr     string `json:"addr,omitempty" yaml:"addr,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}
type RouteErrorInterceptor struct {
	// Enabled, enable and disable backend errors interceptor
	Enabled     bool   `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	ContentType string `yaml:"contentType,omitempty,omitempty" json:"contentType,omitempty"`
	// Errors provides configuration for handling backend errors.
	Errors []RouteError `yaml:"errors,omitempty" json:"errors,omitempty"`
}
type RouteError struct {
	Code int `yaml:"code,omitempty" json:"code,omitempty"` // Deprecated
	// Status contains the status code to intercept
	Status int `yaml:"status,omitempty" json:"status,omitempty"`
	// Body, contains error response custom body
	Body string `yaml:"body,omitempty,omitempty" json:"body,omitempty"`
}

type TLS struct {
	// Keys contains the list of TLS keys
	Keys []Key `yaml:"keys,omitempty" json:"keys,omitempty"`
}
type Key struct {
	TlsSecretName string `yaml:"tlsSecretName" json:"tlsSecretName"`
}

// Backend defines backend server to route traffic to
type Backend struct {
	// Endpoint defines the endpoint of the backend
	Endpoint string `yaml:"endpoint,omitempty" json:"endpoint"`
	// Weight defines Weight for weighted algorithm, it optional
	Weight int `yaml:"weight,omitempty" json:"weight,omitempty"`
}

// Backends defines List of backend servers to route traffic to
type Backends []Backend
