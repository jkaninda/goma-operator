package controller

import gomaprojv1beta1 "github.com/jkaninda/goma-operator/api/v1beta1"

// Gateway contains Goma Proxy Gateway's configs
type Gateway struct {
	// SSLCertFile  SSL Certificate file
	SSLCertFile string `yaml:"sslCertFile"`
	// SSLKeyFile SSL Private key  file
	SSLKeyFile string `yaml:"sslKeyFile"`
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
	InterceptErrors []int                         `yaml:"interceptErrors,omitempty"`
	Routes          []gomaprojv1beta1.RouteConfig `json:"routes,omitempty" yaml:"routes,omitempty"`
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
	// basic, jwt, auth0, rateLimit, access
	Type  string   `json:"type" yaml:"type"`   // Middleware type [basic, jwt, auth0, rateLimit, access]
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
	Username string `yaml:"username" json:"username"`
	Password string `yaml:"password" json:"password"`
}
type JWTRuleMiddleware struct {
	URL             string            `yaml:"url" json:"url"`
	RequiredHeaders []string          `yaml:"requiredHeaders" json:"requiredHeaders"`
	Headers         map[string]string `yaml:"headers" json:"headers"`
	Params          map[string]string `yaml:"params" json:"params"`
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
