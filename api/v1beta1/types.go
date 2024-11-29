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
	// WriteTimeout defines proxy write timeout
	WriteTimeout int `json:"writeTimeout,omitempty" yaml:"writeTimeout,omitempty"`
	// ReadTimeout defines proxy read timeout
	ReadTimeout int `json:"readTimeout,omitempty" yaml:"readTimeout,omitempty"`
	// IdleTimeout defines proxy idle timeout
	IdleTimeout int `json:"idleTimeout,omitempty" yaml:"idleTimeout,omitempty"`
	// LogLevel log level, info, debug, trace, off
	LogLevel string `json:"logLevel,omitempty" yaml:"logLevel,omitempty"`
	// tls secret name
	TlsSecretName string `json:"tlsSecretName,omitempty" yaml:"tlsSecretName,omitempty"`
	// Redis contains redis database details
	Redis Redis `json:"redis,omitempty" yaml:"redis,omitempty"`
	// Cors holds proxy global cors
	Cors Cors `json:"cors,omitempty" yaml:"cors,omitempty,omitempty"`
	// InterceptErrors holds the status codes to intercept the error from backend
	InterceptErrors []int `json:"interceptErrors,omitempty" yaml:"interceptErrors,omitempty"`
	// DisableHealthCheckStatus enable and disable routes health check
	DisableHealthCheckStatus bool `json:"disableHealthCheckStatus,omitempty" yaml:"disableHealthCheckStatus"`
	// DisableKeepAlive allows enabling and disabling KeepALive server
	DisableKeepAlive bool `json:"disableKeepAlive,omitempty" yaml:"disableKeepAlive"`
	// EnableMetrics enable and disable server metrics
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
