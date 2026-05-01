// Package config defines the internal types that represent the generated goma.yml
// configuration file. These types mirror the structure expected by Goma Gateway's
// top-level GatewayConfig (see goma/internal/types.go).
//
// The YAML structure is:
//
//	gateway:
//	  timeouts: ...
//	  tls: ...
//	  redis: ...
//	  monitoring: ...
//	  log: ...
//	  networking: ...
//	  providers: ...
//	  routes: [...]       # routes are nested under gateway
//	middlewares: [...]     # middlewares are at the top level
//	certManager: ...
package config

// GatewayConfig represents the complete goma.yml configuration file.
// This matches goma/internal/types.go:GatewayConfig.
type GatewayConfig struct {
	Gateway     Gateway      `yaml:"gateway"`
	Middlewares []Middleware `yaml:"middlewares,omitempty"`
	CertManager *CertManager `yaml:"certManager,omitempty"`
}

// Gateway represents the nested "gateway" section of the config file.
// This matches goma/internal/gateway.go:Gateway.
type Gateway struct {
	TLS        *TLS        `yaml:"tls,omitempty"`
	Redis      *Redis      `yaml:"redis,omitempty"`
	Timeouts   *Timeouts   `yaml:"timeouts,omitempty"`
	Monitoring *Monitoring `yaml:"monitoring,omitempty"`
	Networking *Networking `yaml:"networking,omitempty"`
	Log        *Log        `yaml:"log,omitempty"`
	Providers  *Providers  `yaml:"providers,omitempty"`

	Routes []Route `yaml:"routes"`
}

// TLS defines TLS certificate configuration (matches TlsCertificates).
type TLS struct {
	CertsDir     string        `yaml:"certsDir,omitempty"`
	Certificates []Certificate `yaml:"certificates,omitempty"`
}

// Certificate defines a single TLS certificate pair.
type Certificate struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

// Redis defines Redis connection settings.
type Redis struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password,omitempty"`
}

// Timeouts defines server timeout settings in seconds.
type Timeouts struct {
	Write int `yaml:"write,omitempty"`
	Read  int `yaml:"read,omitempty"`
	Idle  int `yaml:"idle,omitempty"`
}

// Monitoring defines observability configuration. Mirrors
// goma/internal/gateway.go:Monitoring.
type Monitoring struct {
	EnableMetrics   bool                  `yaml:"enableMetrics,omitempty"`
	Host            string                `yaml:"host,omitempty"`
	MetricsPath     string                `yaml:"metricsPath,omitempty"`
	EnableReadiness bool                  `yaml:"enableReadiness,omitempty"`
	EnableLiveness  bool                  `yaml:"enableLiveness,omitempty"`
	Middleware      *MonitoringMiddleware `yaml:"middleware,omitempty"`
}

// MonitoringMiddleware attaches middleware names to specific monitoring endpoints.
type MonitoringMiddleware struct {
	Metrics []string `yaml:"metrics,omitempty"`
}

// Networking defines transport and DNS settings.
type Networking struct {
	DNSCache  *DNSCache  `yaml:"dnsCache,omitempty"`
	Transport *Transport `yaml:"transport,omitempty"`
}

// DNSCache defines DNS cache settings.
type DNSCache struct {
	TTL int `yaml:"ttl,omitempty"`
}

// Transport defines HTTP transport configuration.
type Transport struct {
	MaxIdleConns        int `yaml:"maxIdleConns,omitempty"`
	MaxIdleConnsPerHost int `yaml:"maxIdleConnsPerHost,omitempty"`
	MaxConnsPerHost     int `yaml:"maxConnsPerHost,omitempty"`
}

// Log defines logging configuration.
type Log struct {
	Level string `yaml:"level,omitempty"`
}

// Providers defines the configuration for dynamic providers.
type Providers struct {
	File *FileProvider `yaml:"file,omitempty"`
	HTTP *HTTPProvider `yaml:"http,omitempty"`
	Git  *GitProvider  `yaml:"git,omitempty"`
}

// FileProvider defines the file provider configuration.
type FileProvider struct {
	Enabled   bool   `yaml:"enabled"`
	Directory string `yaml:"directory"`
	Watch     bool   `yaml:"watch"`
}

// HTTPProvider defines the HTTP provider configuration.
type HTTPProvider struct {
	Enabled            bool              `yaml:"enabled"`
	Endpoint           string            `yaml:"endpoint"`
	Interval           string            `yaml:"interval,omitempty"`
	Timeout            string            `yaml:"timeout,omitempty"`
	Headers            map[string]string `yaml:"headers,omitempty"`
	InsecureSkipVerify bool              `yaml:"insecureSkipVerify,omitempty"`
	RetryAttempts      int               `yaml:"retryAttempts,omitempty"`
	RetryDelay         string            `yaml:"retryDelay,omitempty"`
	CacheDir           string            `yaml:"cacheDir,omitempty"`
}

// GitProvider defines the Git provider configuration.
type GitProvider struct {
	Enabled  bool     `yaml:"enabled"`
	URL      string   `yaml:"url"`
	Branch   string   `yaml:"branch,omitempty"`
	Path     string   `yaml:"path,omitempty"`
	Interval string   `yaml:"interval,omitempty"`
	Auth     *GitAuth `yaml:"auth,omitempty"`
	CloneDir string   `yaml:"cloneDir,omitempty"`
}

// GitAuth defines Git authentication settings.
// Secret values are injected at runtime via env var substitution (${VAR}).
type GitAuth struct {
	Type       string `yaml:"type"`
	Token      string `yaml:"token,omitempty"`
	Username   string `yaml:"username,omitempty"`
	Password   string `yaml:"password,omitempty"`
	SSHKeyPath string `yaml:"sshKeyPath,omitempty"`
}

// CertManager defines ACME / certificate manager configuration.
// This matches goma/pkg/certmanager/types.go:Config.
type CertManager struct {
	Provider string `yaml:"provider,omitempty"`
	ACME     *ACME  `yaml:"acme,omitempty"`
}

// ACME defines Let's Encrypt / ACME configuration.
// This matches goma/pkg/certmanager/types.go:Acme.
type ACME struct {
	Email         string           `yaml:"email"`
	DirectoryURL  string           `yaml:"directoryUrl,omitempty"`
	StorageFile   string           `yaml:"storageFile,omitempty"`
	TermsAccepted bool             `yaml:"termsAccepted,omitempty"`
	ChallengeType string           `yaml:"challengeType,omitempty"`
	DNSProvider   string           `yaml:"dnsProvider,omitempty"`
	Credentials   *ACMECredentials `yaml:"credentials,omitempty"`
}

// ACMECredentials holds DNS provider API credentials.
// The apiToken value comes from env var substitution (${VAR}) — the operator
// mounts the referenced K8s Secret as env vars on the pod.
type ACMECredentials struct {
	APIToken string `yaml:"apiToken,omitempty"`
}

// Route represents a proxy route in the generated config.
type Route struct {
	Name           string        `yaml:"name"`
	Path           string        `yaml:"path"`
	Rewrite        string        `yaml:"rewrite,omitempty"`
	Target         string        `yaml:"target,omitempty"`
	Methods        []string      `yaml:"methods,omitempty"`
	Hosts          []string      `yaml:"hosts,omitempty"`
	Backends       []Backend     `yaml:"backends,omitempty"`
	Priority       int           `yaml:"priority,omitempty"`
	Enabled        bool          `yaml:"enabled"`
	Middlewares    []string      `yaml:"middlewares,omitempty"`
	HealthCheck    *HealthCheck  `yaml:"healthCheck,omitempty"`
	Security       *Security     `yaml:"security,omitempty"`
	TLS            *RouteTLSCert `yaml:"tls,omitempty"`
	Maintenance    *Maintenance  `yaml:"maintenance,omitempty"`
	DisableMetrics bool          `yaml:"disableMetrics,omitempty"`
}

// RouteTLSCert is the per-route serving certificate emitted in goma config.
// Matches the native TlsCertificate struct: tls.certificate.cert / tls.certificate.key.
type RouteTLSCert struct {
	Certificate RouteTLSCertPair `yaml:"certificate"`
}

// RouteTLSCertPair holds cert/key file paths.
type RouteTLSCertPair struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

// Backend defines a backend server.
type Backend struct {
	Endpoint  string         `yaml:"endpoint"`
	Weight    int            `yaml:"weight,omitempty"`
	Match     []BackendMatch `yaml:"match,omitempty"`
	Exclusive bool           `yaml:"exclusive,omitempty"`
}

// BackendMatch is a request condition pinning traffic to a backend.
type BackendMatch struct {
	Source   string `yaml:"source"`
	Name     string `yaml:"name,omitempty"`
	Operator string `yaml:"operator"`
	Value    string `yaml:"value"`
}

// HealthCheck defines health check configuration.
type HealthCheck struct {
	Path            string `yaml:"path,omitempty"`
	Interval        string `yaml:"interval,omitempty"`
	Timeout         string `yaml:"timeout,omitempty"`
	HealthyStatuses []int  `yaml:"healthyStatuses,omitempty"`
}

// Security defines route security settings.
type Security struct {
	ForwardHostHeaders      bool         `yaml:"forwardHostHeaders"`
	EnableExploitProtection bool         `yaml:"enableExploitProtection,omitempty"`
	TLS                     *SecurityTLS `yaml:"tls,omitempty"`
}

// SecurityTLS defines TLS settings for backend connections.
type SecurityTLS struct {
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify,omitempty"`
	RootCAs            string `yaml:"rootCAs,omitempty"`
	ClientCert         string `yaml:"clientCert,omitempty"`
	ClientKey          string `yaml:"clientKey,omitempty"`
}

// Middleware represents a middleware in the generated config.
type Middleware struct {
	Name  string      `yaml:"name"`
	Type  string      `yaml:"type"`
	Paths []string    `yaml:"paths,omitempty"`
	Rule  interface{} `yaml:"rule,omitempty"`
}

// Maintenance defines maintenance mode settings.
type Maintenance struct {
	Enabled bool   `yaml:"enabled"`
	Body    string `yaml:"body,omitempty"`
	Status  int    `yaml:"status,omitempty"`
}
