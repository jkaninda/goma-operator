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
type RouteConfig struct {
	// Path defines route path
	Path string `json:"path" yaml:"path"`
	// Name defines route name
	Name string `json:"name" yaml:"name"`
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
	// DisableHostFording Disable host forwarding.
	DisableHostFording bool `json:"disableHostFording,omitempty" yaml:"disableHostFording"`
	// InterceptErrors intercepts backend errors based on the status codes
	InterceptErrors []int `json:"interceptErrors,omitempty" yaml:"interceptErrors"`
	// BlockCommonExploits enable, disable block common exploits
	BlockCommonExploits bool `json:"blockCommonExploits,omitempty" yaml:"blockCommonExploits"`
	// Middlewares Defines route middleware
	Middlewares []string `json:"middlewares,omitempty" yaml:"middlewares"`
}
type RoutesConfig struct {
	Routes []RouteConfig `json:"routes" yaml:"routes"`
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
