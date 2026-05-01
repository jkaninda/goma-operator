// Package validation provides shared validation logic for CRD specs,
// used by both reconcilers and webhooks.
package validation

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	gatewayv1alpha1 "github.com/jkaninda/goma-operator/api/v1alpha1"
)

// ValidateBackendEndpoint returns an error string if the endpoint is not a
// well-formed absolute URL with an http/https scheme. Returns "" on success.
func ValidateBackendEndpoint(endpoint string) string {
	if endpoint == "" {
		return "backend endpoint is required"
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Sprintf("backend endpoint %q is not a valid URL: %v", endpoint, err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Sprintf("backend endpoint %q must use http or https scheme", endpoint)
	}
	if u.Host == "" {
		return fmt.Sprintf("backend endpoint %q is missing a host", endpoint)
	}
	return ""
}

var validMethods = map[string]bool{
	http.MethodGet:     true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodPatch:   true,
	http.MethodDelete:  true,
	http.MethodHead:    true,
	http.MethodOptions: true,
	http.MethodConnect: true,
	http.MethodTrace:   true,
}

// ValidateRouteSpec validates a RouteSpec and returns a list of errors.
func ValidateRouteSpec(spec *gatewayv1alpha1.RouteSpec) []string {
	var errs []string

	if len(spec.Gateways) == 0 {
		errs = append(errs, "spec.gateways must contain at least one gateway name")
	} else {
		for i, gw := range spec.Gateways {
			if gw == "" {
				errs = append(errs, fmt.Sprintf("spec.gateways[%d] must not be empty", i))
			}
		}
	}
	if spec.Path == "" {
		errs = append(errs, "spec.path is required")
	}

	hasTarget := spec.Target != ""
	hasBackends := len(spec.Backends) > 0

	if !hasTarget && !hasBackends {
		errs = append(errs, "either spec.target or spec.backends must be set")
	}
	if hasTarget && hasBackends {
		errs = append(errs, "spec.target and spec.backends are mutually exclusive")
	}

	for _, m := range spec.Methods {
		if !validMethods[strings.ToUpper(m)] {
			errs = append(errs, fmt.Sprintf("invalid HTTP method: %s", m))
		}
	}

	if spec.Target != "" {
		if msg := ValidateBackendEndpoint(spec.Target); msg != "" {
			errs = append(errs, strings.Replace(msg, "backend endpoint", "spec.target", 1))
		}
	}
	for i, b := range spec.Backends {
		if msg := ValidateBackendEndpoint(b.Endpoint); msg != "" {
			errs = append(errs, fmt.Sprintf("spec.backends[%d]: %s", i, msg))
		}
	}

	if spec.HealthCheck != nil && spec.HealthCheck.Path == "" {
		errs = append(errs, "healthCheck.path is required when healthCheck is set")
	}

	return errs
}
