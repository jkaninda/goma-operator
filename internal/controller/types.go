package controller

import gomaprojv1beta1 "github.com/jkaninda/goma-operator/api/v1beta1"

// Gateway contains Goma Proxy Gateway's configs
type Gateway struct {
	// TlsCertFile  SSL Certificate file
	TlsCertFile string `yaml:"tlsCertFile"`
	// TlsKeyFile SSL Private key  file
	TlsKeyFile string `yaml:"tlsKeyFile"`
	// Redis contains redis database details
	Redis gomaprojv1beta1.Redis `yaml:"redis"`
	// WriteTimeout defines proxy write timeout
	WriteTimeout int `yaml:"writeTimeout"`
	// ReadTimeout defines proxy read timeout
	ReadTimeout int `yaml:"readTimeout"`
	// IdleTimeout defines proxy idle timeout
	IdleTimeout int                  `yaml:"idleTimeout"`
	LogLevel    string               `yaml:"logLevel"`
	Cors        gomaprojv1beta1.Cors `yaml:"cors"`
	// DisableHealthCheckStatus enable and disable routes health check
	DisableHealthCheckStatus bool `yaml:"disableHealthCheckStatus"`
	// DisableRouteHealthCheckError allows enabling and disabling backend healthcheck errors
	DisableRouteHealthCheckError bool `yaml:"disableRouteHealthCheckError"`
	// Disable allows enabling and disabling displaying routes on start
	DisableDisplayRouteOnStart bool `yaml:"disableDisplayRouteOnStart"`
	// DisableKeepAlive allows enabling and disabling KeepALive server
	DisableKeepAlive bool `yaml:"disableKeepAlive"`
	EnableMetrics    bool `yaml:"enableMetrics"`
	// InterceptErrors holds the status codes to intercept the error from backend
	InterceptErrors       []int `yaml:"interceptErrors,omitempty"` // Deprecated
	DisableHostForwarding bool  `yaml:"disableHostForwarding"`
	//  ErrorInterceptor handles backend error interceptor
	ErrorInterceptor gomaprojv1beta1.RouteErrorInterceptor `yaml:"errorInterceptor,omitempty" json:"errorInterceptor,omitempty"`
	Routes           []Route                               `json:"routes,omitempty" yaml:"routes,omitempty"`
}
type Route struct {
	// Path defines route path
	Path string `json:"path" yaml:"path"`
	// Name defines route name
	Name string `json:"name" yaml:"name"`
	// Hosts Domains/hosts based request routing
	Hosts []string `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	// Rewrite rewrites route path to desired path
	Rewrite string `json:"rewrite,omitempty" yaml:"rewrite,omitempty"`
	// Methods allowed method
	Methods []string `json:"methods,omitempty" yaml:"methods,omitempty"`
	// Destination Defines backend URL
	Destination        string   `json:"destination,omitempty" yaml:"destination,omitempty"`
	Backends           []string `json:"backends,omitempty" yaml:"backends,omitempty"`
	InsecureSkipVerify bool     `json:"insecureSkipVerify,omitempty" yaml:"insecureSkipVerify,omitempty"`
	// HealthCheck Defines the backend is health
	HealthCheck gomaprojv1beta1.RouteHealthCheck `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty"`
	// Cors contains the route cors headers
	Cors                  gomaprojv1beta1.Cors `json:"cors,omitempty" yaml:"cors"`
	RateLimit             int                  `json:"rateLimit,omitempty" yaml:"rateLimit,omitempty"`
	DisableHostForwarding bool                 `json:"disableHostForwarding,omitempty" yaml:"disableHostForwarding,omitempty"`
	//  ErrorInterceptor handles backend error interceptor
	ErrorInterceptor gomaprojv1beta1.RouteErrorInterceptor `yaml:"errorInterceptor,omitempty" json:"errorInterceptor,omitempty"`
	// BlockCommonExploits enable, disable block common exploits
	BlockCommonExploits bool `json:"blockCommonExploits,omitempty" yaml:"blockCommonExploits,omitempty"`
	EnableBotDetection  bool `json:"enableBotDetection,omitempty" yaml:"enableBotDetection,omitempty"`
	// Middlewares Defines route middleware
	Middlewares []string `json:"middlewares,omitempty" yaml:"middlewares"`
}
type Redis struct {
	// Addr redis hostname and port number :
	Addr string `yaml:"addr"`
	// Password redis password
	Password string `yaml:"password"`
}
type Middleware struct {
	// Path contains the name of middlewares and must be unique
	Name string `json:"name" yaml:"name"`
	// Type contains authentication types
	//
	// basic, jwt, oauth, rateLimit, access
	Type  string   `json:"type" yaml:"type"`   // Middleware type [basic, jwt, oauth, rateLimit, access]
	Paths []string `json:"paths" yaml:"paths"` // Protected paths
	// Rule contains route middleware rule
	Rule interface{} `json:"rule" yaml:"rule"`
}

type Middlewares struct {
	Middlewares []Middleware `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
}

type GatewayConfig struct {
	Version     string       `json:"version" yaml:"version"`
	Gateway     Gateway      `json:"gateway" yaml:"gateway"`
	Middlewares []Middleware `json:"middlewares,omitempty" yaml:"middlewares,omitempty"`
}
type BasicRuleMiddleware struct {
	Realm    string   `yaml:"realm,omitempty" json:"realm,omitempty"`
	Users    []string `yaml:"users" json:"users"`
	Username string   `yaml:"username,omitempty" json:"username"` // Deprecated
	Password string   `yaml:"password,omitempty" json:"password"` // Deprecated
}
type AddPrefixRuleMiddleware struct {
	Prefix string `yaml:"prefix" json:"prefix"`
}
type RedirectRegexRuleMiddleware struct {
	Pattern     string `yaml:"pattern" json:"pattern"`
	Replacement string `yaml:"replacement" json:"replacement"`
}
type JWTRuleMiddleware struct {
	URL             string            `yaml:"url" json:"url"`
	RequiredHeaders []string          `yaml:"requiredHeaders" json:"requiredHeaders"`
	Headers         map[string]string `yaml:"headers" json:"headers"`
	Params          map[string]string `yaml:"params" json:"params"`
}

// AccessPolicyRuleMiddleware access policy
type AccessPolicyRuleMiddleware struct {
	Action       string   `yaml:"action,omitempty" json:"action,omitempty"` // action, ALLOW or DENY
	SourceRanges []string `yaml:"sourceRanges" json:"sourceRanges"`         //  list of Ips
}

type RateLimitRuleMiddleware struct {
	Unit            string `yaml:"unit" json:"unit"`
	RequestsPerUnit int    `yaml:"requestsPerUnit" json:"requestsPerUnit"`
}
type OauthRulerMiddleware struct {
	// ClientID is the application's ID.
	ClientID string `yaml:"clientId"`

	// ClientSecret is the application's secret.
	ClientSecret string `yaml:"clientSecret"`
	// oauth provider google, gitlab, github, amazon, facebook, custom
	Provider string `yaml:"provider"`
	// Endpoint contains the resource server's token endpoint
	Endpoint OauthEndpoint `yaml:"endpoint"`

	// RedirectURL is the URL to redirect users going through
	// the OAuth flow, after the resource owner's URLs.
	RedirectURL string `yaml:"redirectUrl"`
	// RedirectPath is the PATH to redirect users after authentication, e.g: /my-protected-path/dashboard
	RedirectPath string `yaml:"redirectPath"`
	// CookiePath e.g: /my-protected-path or / || by default is applied on a route path
	CookiePath string `yaml:"cookiePath"`

	// Scope specifies optional requested permissions.
	Scopes []string `yaml:"scopes"`
	// contains filtered or unexported fields
	State     string `yaml:"state"`
	JWTSecret string `yaml:"jwtSecret"`
}
type OauthEndpoint struct {
	AuthURL     string `yaml:"authUrl"`
	TokenURL    string `yaml:"tokenUrl"`
	UserInfoURL string `yaml:"userInfoUrl"`
}
