package converter

import (
	"fmt"

	gatewayv1alpha1 "github.com/jkaninda/goma-operator/api/v1alpha1"
	"github.com/jkaninda/goma-operator/internal/config"
)

// RouteFromCR converts a Route CRD to an internal config Route.
func RouteFromCR(cr *gatewayv1alpha1.Route) config.Route {
	r := config.Route{
		Name:           cr.Name,
		Path:           cr.Spec.Path,
		Rewrite:        cr.Spec.Rewrite,
		Target:         cr.Spec.Target,
		Methods:        cr.Spec.Methods,
		Hosts:          cr.Spec.Hosts,
		Priority:       cr.Spec.Priority,
		Enabled:        cr.Spec.Enabled,
		Middlewares:    cr.Spec.Middlewares,
		DisableMetrics: cr.Spec.DisableMetrics,
	}

	// Backends
	if len(cr.Spec.Backends) > 0 {
		r.Backends = make([]config.Backend, 0, len(cr.Spec.Backends))
		for _, b := range cr.Spec.Backends {
			be := config.Backend{
				Endpoint:  b.Endpoint,
				Weight:    b.Weight,
				Exclusive: b.Exclusive,
			}
			if len(b.Match) > 0 {
				be.Match = make([]config.BackendMatch, 0, len(b.Match))
				for _, m := range b.Match {
					be.Match = append(be.Match, config.BackendMatch{
						Source:   m.Source,
						Name:     m.Name,
						Operator: m.Operator,
						Value:    m.Value,
					})
				}
			}
			r.Backends = append(r.Backends, be)
		}
	}

	// Health check
	if cr.Spec.HealthCheck != nil {
		r.HealthCheck = &config.HealthCheck{
			Path:            cr.Spec.HealthCheck.Path,
			Interval:        cr.Spec.HealthCheck.Interval,
			Timeout:         cr.Spec.HealthCheck.Timeout,
			HealthyStatuses: cr.Spec.HealthCheck.HealthyStatuses,
		}
	}

	// Security
	if cr.Spec.Security != nil {
		sec := &config.Security{
			ForwardHostHeaders:      cr.Spec.Security.ForwardHostHeaders,
			EnableExploitProtection: cr.Spec.Security.EnableExploitProtection,
		}
		if cr.Spec.Security.TLS != nil {
			secTLS := &config.SecurityTLS{
				InsecureSkipVerify: cr.Spec.Security.TLS.InsecureSkipVerify,
			}
			if cr.Spec.Security.TLS.RootCAsSecret != "" {
				secTLS.RootCAs = fmt.Sprintf("%s/%s/ca.crt", certsBasePath, cr.Spec.Security.TLS.RootCAsSecret)
			}
			if cr.Spec.Security.TLS.ClientCertSecret != "" {
				secTLS.ClientCert = fmt.Sprintf("%s/%s/tls.crt", certsBasePath, cr.Spec.Security.TLS.ClientCertSecret)
				secTLS.ClientKey = fmt.Sprintf("%s/%s/tls.key", certsBasePath, cr.Spec.Security.TLS.ClientCertSecret)
			}
			sec.TLS = secTLS
		}
		r.Security = sec
	}

	// Per-route serving TLS certificate
	if cr.Spec.TLS != nil && cr.Spec.TLS.SecretName != "" {
		r.TLS = &config.RouteTLSCert{
			Certificate: config.RouteTLSCertPair{
				Cert: fmt.Sprintf("%s/%s/tls.crt", certsBasePath, cr.Spec.TLS.SecretName),
				Key:  fmt.Sprintf("%s/%s/tls.key", certsBasePath, cr.Spec.TLS.SecretName),
			},
		}
	}

	// Maintenance
	if cr.Spec.Maintenance != nil {
		r.Maintenance = &config.Maintenance{
			Enabled: cr.Spec.Maintenance.Enabled,
			Body:    cr.Spec.Maintenance.Body,
			Status:  cr.Spec.Maintenance.Status,
		}
	}

	return r
}
