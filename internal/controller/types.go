package controller

import gomaprojv1beta1 "github.com/jkaninda/goma-operator/api/v1beta1"

// Gateway contains Goma Proxy Gateway's configs
type Gateway struct {
	// TlsCertFile SSL Certificate file
	// Deprecated use TLS
	TlsCertFile string `yaml:"tlsCertFile,omitempty"`
	// TlsKeyFile SSL Private key file
	// Deprecated use TLS
	TlsKeyFile string `yaml:"tlsKeyFile,omitempty"`
	TLS        TLS    `yaml:"tls,omitempty"`
	// Redis contains redis database details
	Redis gomaprojv1beta1.Redis `yaml:"redis,omitempty"`
	// WriteTimeout defines proxy write timeout
	WriteTimeout int `yaml:"writeTimeout"`
	// ReadTimeout defines proxy read timeout
	ReadTimeout int `yaml:"readTimeout"`
	// IdleTimeout defines proxy idle timeout
	IdleTimeout int                  `yaml:"idleTimeout"`
	LogLevel    string               `yaml:"logLevel"`
	Cors        gomaprojv1beta1.Cors `yaml:"cors"`
	// DisableHealthCheckStatus enable and disable routes health check
	DisableHealthCheckStatus bool `yaml:"disableHealthCheckStatus,omitempty"`
	// DisableRouteHealthCheckError allows enabling and disabling backend healthcheck errors
	DisableRouteHealthCheckError bool `yaml:"disableRouteHealthCheckError,omitempty"`
	// Disable allows enabling and disabling displaying routes on start
	DisableDisplayRouteOnStart bool `yaml:"disableDisplayRouteOnStart,omitempty"`
	// DisableKeepAlive allows enabling and disabling KeepALive server
	DisableKeepAlive bool `yaml:"disableKeepAlive,omitempty"`
	EnableMetrics    bool `yaml:"enableMetrics,omitempty"`
	// InterceptErrors holds the status codes to intercept the error from backend
	InterceptErrors   []int `yaml:"interceptErrors,omitempty"` // Deprecated
	EnableStrictSlash bool  `json:"enableStrictSlash,omitempty" yaml:"enableStrictSlash,omitempty"`
	//  ErrorInterceptor handles backend error interceptor
	ErrorInterceptor gomaprojv1beta1.RouteErrorInterceptor `yaml:"errorInterceptor,omitempty" json:"errorInterceptor,omitempty"`
	Routes           []Route                               `json:"routes,omitempty" yaml:"routes"`
}
type Route struct {
	// Path defines route path
	Path string `json:"path" yaml:"path"`
	// Name defines route name
	Name     string `json:"name" yaml:"name"`
	Disabled bool   `json:"disabled,omitempty" yaml:"disabled"`
	TLS      TLS    `yaml:"tls,omitempty"`
	// Route order priority
	Priority int `yaml:"priority,omitempty" json:"priority,omitempty"`
	// Hosts Domains/hosts based request routing
	Hosts []string `json:"hosts,omitempty" yaml:"hosts,omitempty"`
	// Rewrite rewrites route path to desired path
	Rewrite string `json:"rewrite,omitempty" yaml:"rewrite,omitempty"`
	// Methods allowed method
	Methods []string `json:"methods,omitempty" yaml:"methods,omitempty"`
	// Destination Defines backend URL
	Destination string `json:"destination,omitempty" yaml:"destination,omitempty"`
	// Backends specifies a list of backend URLs for load balancing.
	Backends           gomaprojv1beta1.Backends `yaml:"backends,omitempty" json:"backends,omitempty"`
	InsecureSkipVerify bool                     `json:"insecureSkipVerify,omitempty" yaml:"insecureSkipVerify,omitempty"`
	// HealthCheck Defines the backend is health
	HealthCheck gomaprojv1beta1.RouteHealthCheck `json:"healthCheck,omitempty" yaml:"healthCheck,omitempty"`
	// Cors contains the route cors headers
	Cors                  gomaprojv1beta1.Cors `json:"cors,omitempty" yaml:"cors,omitempty"`
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
	Type  string   `json:"type" yaml:"type"`                       // Middleware type [basic, jwt, oauth, rateLimit, access]
	Paths []string `json:"paths,omitempty" yaml:"paths,omitempty"` // Protected paths
	// Rule contains route middleware rule
	Rule interface{} `json:"rule,omitempty" yaml:"rule,omitempty"`
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
type RewriteRegexRuleMiddleware struct {
	Pattern     string `yaml:"pattern" json:"pattern"`
	Replacement string `yaml:"replacement" json:"replacement"`
}
type JWTRuleMiddleware struct {
	Alg                  string `yaml:"alg,omitempty" json:"alg,omitempty"`
	Secret               string `yaml:"secret,omitempty" json:"secret,omitempty"`
	PublicKey            string `yaml:"publicKey,omitempty" json:"publicKey,omitempty"`
	JwksUrl              string `yaml:"jwksUrl,omitempty" json:"jwksUrl,omitempty"`
	ForwardAuthorization bool   `yaml:"forwardAuthorization,omitempty" json:"forwardAuthorization,omitempty"`
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
type ForwardAuthRuleMiddleware struct {
	AuthURL                     string   `yaml:"authUrl" json:"authUrl"`
	AuthSignIn                  string   `yaml:"authSignIn,omitempty" json:"authSignIn,omitempty"`
	EnableHostForwarding        bool     `yaml:"enableHostForwarding,omitempty" json:"enableHostForwarding,omitempty"`
	SkipInsecureVerify          bool     `yaml:"skipInsecureVerify,omitempty" json:"skipInsecureVerify,omitempty"`
	AuthRequestHeaders          []string `yaml:"authRequestHeaders,omitempty" json:"authRequestHeaders,omitempty"`
	AddAuthCookiesToResponse    []string `yaml:"addAuthCookiesToResponse,omitempty" json:"addAuthCookiesToResponse,omitempty"`
	AuthResponseHeaders         []string `yaml:"authResponseHeaders,omitempty" json:"authResponseHeaders,omitempty"`
	AuthResponseHeadersAsParams []string `yaml:"authResponseHeadersAsParams,omitempty" json:"authResponseHeadersAsParams,omitempty"`
}
type httpCacheRule struct {
	MaxTtl                   int64    `yaml:"maxTtl" json:"maxTtl"`
	MaxStale                 int64    `yaml:"maxStale,omitempty" json:"maxStale,omitempty"`
	DisableCacheStatusHeader bool     `yaml:"disableCacheStatusHeader,omitempty" json:"disableCacheStatusHeader,omitempty"`
	ExcludedResponseCodes    []string `yaml:"excludedResponseCodes,omitempty" json:"excludedResponseCodes,omitempty"`
	MemoryLimit              string   `yaml:"memoryLimit,omitempty" json:"memoryLimit,omitempty"`
}
type RedirectSchemeRuleMiddleware struct {
	Scheme    string `yaml:"scheme" json:"scheme"`
	Port      int64  `yaml:"port,omitempty" json:"port,omitempty"`
	Permanent bool   `yaml:"permanent,omitempty" json:"permanent,omitempty"`
}
type AccessRuleMiddleware struct {
	StatusCode int `yaml:"statusCode,omitempty"` // HTTP Response code
}

type TLS struct {
	// Keys contains the list of TLS keys
	Keys []Key `yaml:"keys,omitempty" json:"keys,omitempty"`
}
type Key struct {
	Cert string `yaml:"cert" json:"cert"`
	Key  string `yaml:"key" json:"key"`
}
type BodyLimitRuleMiddleware struct {
	Limit string `yaml:"limit" json:"limit"`
}
