package validation

import (
	"fmt"

	gatewayv1alpha1 "github.com/jkaninda/goma-operator/api/v1alpha1"
)

var supportedMiddlewareTypes = map[string]bool{
	"basic":            true,
	"basicAuth":        true,
	"jwt":              true,
	"jwtAuth":          true,
	"oauth":            true,
	"oauth2":           true,
	"rateLimit":        true,
	"access":           true,
	"accessLog":        true,
	"accessPolicy":     true,
	"addPrefix":        true,
	"redirect":         true,
	"redirectRegex":    true,
	"rewriteRegex":     true,
	"forwardAuth":      true,
	"httpCache":        true,
	"redirectScheme":   true,
	"bodyLimit":        true,
	"responseHeaders":  true,
	"requestHeaders":   true,
	"errorInterceptor": true,
	"ldap":             true,
	"ldapAuth":         true,
	"userAgentBlock":   true,
}

// ValidateMiddlewareSpec validates a MiddlewareSpec and returns a list of errors.
func ValidateMiddlewareSpec(spec *gatewayv1alpha1.MiddlewareSpec) []string {
	var errs []string

	if spec.Type == "" {
		errs = append(errs, "spec.type is required")
	} else if !supportedMiddlewareTypes[spec.Type] {
		errs = append(errs, fmt.Sprintf("unsupported middleware type: %s", spec.Type))
	}

	return errs
}
