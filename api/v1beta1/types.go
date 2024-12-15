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
	TlsSecretName string `json:"tlsSecretName,omitempty" yaml:"tlsSecretName,omitempty"`
	// Redis contains the configuration details for connecting to a Redis database.
	Redis Redis `json:"redis,omitempty" yaml:"redis,omitempty"`
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
	Enabled     bool         `yaml:"enabled" json:"enabled"`
	ContentType string       `yaml:"contentType,omitempty,omitempty" json:"contentType,omitempty"`
	Errors      []RouteError `yaml:"errors,omitempty" json:"errors,omitempty"`
}
type RouteError struct {
	Code int    `yaml:"code" json:"code"`
	Body string `yaml:"body,omitempty,omitempty" json:"body,omitempty"`
}
