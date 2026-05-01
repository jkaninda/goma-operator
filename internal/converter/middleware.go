package converter

import (
	"encoding/json"

	gatewayv1alpha1 "github.com/jkaninda/goma-operator/api/v1alpha1"
	"github.com/jkaninda/goma-operator/internal/config"
)

// MiddlewareFromCR converts a Middleware CRD to an internal config Middleware.
func MiddlewareFromCR(cr *gatewayv1alpha1.Middleware) config.Middleware {
	m := config.Middleware{
		Name:  cr.Name,
		Type:  cr.Spec.Type,
		Paths: cr.Spec.Paths,
	}

	// Unmarshal the raw JSON rule into a generic interface{} for YAML output.
	if cr.Spec.Rule != nil && cr.Spec.Rule.Raw != nil {
		var rule interface{}
		if err := json.Unmarshal(cr.Spec.Rule.Raw, &rule); err == nil {
			m.Rule = rule
		}
	}

	return m
}
