// Package converter provides functions to convert CRD API types to internal
// config types used for generating goma.yml.
package converter

import (
	gatewayv1alpha1 "github.com/jkaninda/goma-operator/api/v1alpha1"
	"github.com/jkaninda/goma-operator/internal/config"
)

const (
	certsBasePath = "/etc/goma/certs"
	certsDirPath  = "/etc/goma/certs" // used as certsDir when any TLS secret is mounted

	// K8sProviderDir is the directory where the goma-k8s-provider sidecar
	// writes its generated config. The gateway's FileProvider watches it.
	K8sProviderDir = "/etc/goma/providers/k8s"

	// GitCloneDir is the default clone directory for the Git provider.
	GitCloneDir = "/etc/goma/providers/git"

	// HTTPCacheDir is the default cache directory for the HTTP provider.
	HTTPCacheDir = "/etc/goma/providers/http/cache.json"
)

// GatewayConfigFromCRs builds a complete GatewayConfig from a Gateway CR,
// its associated Routes, and the referenced Middlewares.
//
// The output YAML structure is:
//
//	gateway:
//	  timeouts, tls, redis, monitoring, log, networking, providers, routes
//	middlewares: [...]
//
// nolint:gocyclo
func GatewayConfigFromCRs(gw *gatewayv1alpha1.Gateway, routes []gatewayv1alpha1.Route, middlewares []gatewayv1alpha1.Middleware) config.GatewayConfig {
	cfg := config.GatewayConfig{}

	server := gw.Spec.Server
	gwCfg := &cfg.Gateway

	// Timeouts
	if server.Timeouts.Write != 0 || server.Timeouts.Read != 0 || server.Timeouts.Idle != 0 {
		gwCfg.Timeouts = &config.Timeouts{
			Write: server.Timeouts.Write,
			Read:  server.Timeouts.Read,
			Idle:  server.Timeouts.Idle,
		}
	}

	// Log level
	if server.LogLevel != "" {
		gwCfg.Log = &config.Log{Level: server.LogLevel}
	}

	// Redis
	if server.Redis != nil {
		gwCfg.Redis = &config.Redis{
			Addr:     server.Redis.Addr,
			Password: server.Redis.Password,
		}
	}

	// TLS — when secrets are mounted at /etc/goma/certs/<secretName>/, we set
	// certsDir so Goma Gateway auto-discovers all cert/key pairs in the tree.
	if len(server.TLS) > 0 {
		gwCfg.TLS = &config.TLS{
			CertsDir: certsDirPath,
		}
	}

	// Monitoring — /readyz and /healthz are always enabled so the
	// operator-managed Deployment's readiness/liveness probes succeed.
	mon := &config.Monitoring{
		EnableReadiness: true,
		EnableLiveness:  true,
	}
	if server.Monitoring != nil {
		mon.EnableMetrics = server.Monitoring.EnableMetrics
		mon.MetricsPath = server.Monitoring.MetricsPath
		mon.Host = server.Monitoring.Host
		if server.Monitoring.Middleware != nil && len(server.Monitoring.Middleware.Metrics) > 0 {
			mon.Middleware = &config.MonitoringMiddleware{
				Metrics: server.Monitoring.Middleware.Metrics,
			}
		}
	}
	gwCfg.Monitoring = mon

	// Networking
	if server.Networking != nil {
		net := &config.Networking{}
		if server.Networking.DNSCache != nil {
			net.DNSCache = &config.DNSCache{TTL: server.Networking.DNSCache.TTL}
		}
		if server.Networking.Transport != nil {
			net.Transport = &config.Transport{
				MaxIdleConns:        server.Networking.Transport.MaxIdleConns,
				MaxIdleConnsPerHost: server.Networking.Transport.MaxIdleConnsPerHost,
				MaxConnsPerHost:     server.Networking.Transport.MaxConnsPerHost,
			}
		}
		gwCfg.Networking = net
	}

	// Providers — map all three provider types
	providers := &config.Providers{}

	// File provider for k8s sidecar (enabled by default — opt-out via
	// spec.providers.kubernetes.enabled=false).
	if gw.Spec.KubernetesProviderEnabled() {
		providers.File = &config.FileProvider{
			Enabled:   true,
			Directory: K8sProviderDir,
			Watch:     true,
		}
	}

	if gw.Spec.Providers != nil {
		// HTTP provider
		if gw.Spec.Providers.HTTP != nil && gw.Spec.Providers.HTTP.Enabled {
			http := gw.Spec.Providers.HTTP
			providers.HTTP = &config.HTTPProvider{
				Enabled:            true,
				Endpoint:           http.Endpoint,
				Interval:           http.Interval,
				Timeout:            http.Timeout,
				Headers:            http.Headers,
				InsecureSkipVerify: http.InsecureSkipVerify,
				RetryAttempts:      http.RetryAttempts,
				RetryDelay:         http.RetryDelay,
				CacheDir:           http.CacheDir,
			}
		}

		// Git provider
		if gw.Spec.Providers.Git != nil && gw.Spec.Providers.Git.Enabled {
			git := gw.Spec.Providers.Git
			gp := &config.GitProvider{
				Enabled:  true,
				URL:      git.URL,
				Branch:   git.Branch,
				Path:     git.Path,
				Interval: git.Interval,
				CloneDir: git.CloneDir,
			}
			if git.Auth != nil {
				// Auth values come from env var substitution (${VAR}) — the
				// operator mounts git.Auth.SecretName as env vars on the pod.
				auth := &config.GitAuth{Type: git.Auth.Type}
				switch git.Auth.Type {
				case "token":
					auth.Token = "${GIT_TOKEN}"
				case "basic":
					auth.Username = "${GIT_USERNAME}"
					auth.Password = "${GIT_PASSWORD}"
				case "ssh":
					auth.SSHKeyPath = "/etc/goma/providers/git/ssh/ssh-privatekey"
				}
				gp.Auth = auth
			}
			providers.Git = gp
		}
	}

	// Only attach if at least one provider is populated
	if providers.File != nil || providers.HTTP != nil || providers.Git != nil {
		gwCfg.Providers = providers
	}

	// CertManager (ACME) configuration
	if gw.Spec.CertManager != nil {
		cm := gw.Spec.CertManager
		cfgCM := &config.CertManager{
			Provider: cm.Provider,
		}
		if cfgCM.Provider == "" {
			cfgCM.Provider = "acme"
		}
		if cm.ACME != nil {
			acme := &config.ACME{
				Email:         cm.ACME.Email,
				DirectoryURL:  cm.ACME.DirectoryURL,
				StorageFile:   "/etc/letsencrypt/acme.json",
				TermsAccepted: cm.ACME.TermsAccepted,
				ChallengeType: cm.ACME.ChallengeType,
				DNSProvider:   cm.ACME.DNSProvider,
			}
			if acme.ChallengeType == "" {
				acme.ChallengeType = "http-01"
			}
			// DNS-01 credentials are supplied via env var substitution —
			// the operator mounts cm.ACME.CredentialsSecret as env vars on the pod.
			if cm.ACME.CredentialsSecret != "" {
				acme.Credentials = &config.ACMECredentials{
					APIToken: "${GOMA_CREDENTIALS_API_TOKEN}",
				}
			}
			cfgCM.ACME = acme
		}
		cfg.CertManager = cfgCM
	}

	// Routes & Middlewares handling:
	//
	// When the k8s sidecar is enabled, routes and middlewares are delivered
	// dynamically via the FileProvider (/etc/goma/providers/k8s/goma.yaml)
	// and hot-reloaded WITHOUT restarting the gateway pod.
	//
	// In that mode we deliberately leave routes & middlewares OUT of the
	// static ConfigMap. This means:
	//   - Changes to Route/Middleware CRs only update the sidecar's output
	//     file — the ConfigMap checksum stays the same, no rolling restart.
	//   - Only changes to Gateway-level settings (tls, redis, monitoring,
	//     providers, etc.) cause a pod restart.
	//
	// When the k8s sidecar is NOT enabled, routes and middlewares are
	// bundled into the static ConfigMap and Route/Middleware changes DO
	// trigger a rolling restart.
	sidecarEnabled := gw.Spec.KubernetesProviderEnabled()

	if !sidecarEnabled {
		// Routes
		gwCfg.Routes = make([]config.Route, 0, len(routes))
		for _, r := range routes {
			gwCfg.Routes = append(gwCfg.Routes, RouteFromCR(&r))
		}

		// Middlewares
		cfg.Middlewares = make([]config.Middleware, 0, len(middlewares))
		for _, m := range middlewares {
			cfg.Middlewares = append(cfg.Middlewares, MiddlewareFromCR(&m))
		}
	} else {
		gwCfg.Routes = []config.Route{}
	}

	return cfg
}
