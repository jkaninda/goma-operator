package controller

import (
	"context"
	"encoding/json"
	"github.com/jinzhu/copier"
	"k8s.io/apimachinery/pkg/runtime"
	"reflect"
	"slices"
	"strings"

	gomaprojv1beta1 "github.com/jkaninda/goma-operator/api/v1beta1"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

func gatewayConfig(r GatewayReconciler, ctx context.Context, req ctrl.Request, gateway *gomaprojv1beta1.Gateway) GatewayConfig {
	logger := log.FromContext(ctx)
	gomaConfig := &GatewayConfig{}
	gomaConfig.Version = GatewayConfigVersion
	err := copier.Copy(&gomaConfig.Gateway, &gateway.Spec.Server)
	if err != nil {
		logger.Error(err, "failed to copy gateway spec")
	}
	// attach cert files
	if len(gateway.Spec.Server.TlsSecretName) != 0 {
		gomaConfig.Gateway.TlsCertFile = TLSCertFile
		gomaConfig.Gateway.TlsKeyFile = TLSKeyFile
	}

	labelSelector := client.MatchingLabels{}
	var middlewareNames []string
	// List ConfigMaps in the namespace with the matching label
	var routes gomaprojv1beta1.RouteList
	if err := r.List(ctx, &routes, labelSelector, client.InNamespace(req.Namespace)); err != nil {
		logger.Error(err, "Failed to list Routes")
		return *gomaConfig
	}
	var middlewares gomaprojv1beta1.MiddlewareList
	if err := r.List(ctx, &middlewares, labelSelector, client.InNamespace(req.Namespace)); err != nil {
		logger.Error(err, "Failed to list Middlewares")
		return *gomaConfig
	}
	for _, route := range routes.Items {
		logger.Info("Found Route", "Name", route.Name)
		if route.Spec.Gateway == gateway.Name {
			logger.Info("Found Route", "Name", route.Name)
			rt := Route{}
			err := copier.Copy(&rt, &route.Spec)
			if err != nil {
				logger.Error(err, "Failed to deep copy Route", "Name", route.Name)
				return *gomaConfig
			}
			rt.Name = route.Name
			gomaConfig.Gateway.Routes = append(gomaConfig.Gateway.Routes, rt)
			middlewareNames = append(middlewareNames, rt.Middlewares...)
		}
	}
	for _, mid := range middlewares.Items {
		middleware := *mapMid(mid)
		logger.Info("Adding Middleware", "Name", middleware.Name)
		if slices.Contains(middlewareNames, middleware.Name) {
			gomaConfig.Middlewares = append(gomaConfig.Middlewares, middleware)
		}

	}
	return *gomaConfig
}
func updateGatewayConfig(r RouteReconciler, ctx context.Context, req ctrl.Request, gateway gomaprojv1beta1.Gateway) (bool, error) {
	logger := log.FromContext(ctx)
	gomaConfig := &GatewayConfig{}
	gomaConfig.Version = GatewayConfigVersion
	err := copier.Copy(&gomaConfig.Gateway, &gateway.Spec.Server)
	if err != nil {
		logger.Error(err, "failed to copy gateway spec")
	}
	// attach cert files
	if len(gateway.Spec.Server.TlsSecretName) != 0 {
		gomaConfig.Gateway.TlsCertFile = TLSCertFile
		gomaConfig.Gateway.TlsKeyFile = TLSKeyFile
	}
	labelSelector := client.MatchingLabels{}
	var middlewareNames []string
	// List ConfigMaps in the namespace with the matching label
	var routes gomaprojv1beta1.RouteList
	if err := r.List(ctx, &routes, labelSelector, client.InNamespace(req.Namespace)); err != nil {
		logger.Error(err, "Failed to list Routes")
		return false, err
	}
	var middlewares gomaprojv1beta1.MiddlewareList
	if err := r.List(ctx, &middlewares, labelSelector, client.InNamespace(req.Namespace)); err != nil {
		logger.Error(err, "Failed to list Middlewares")
		return false, err
	}
	for _, route := range routes.Items {
		logger.Info("Found Route", "Name", route.Name)
		if route.Spec.Gateway == gateway.Name {
			if route.ObjectMeta.DeletionTimestamp.IsZero() {
				rt := Route{}
				err := copier.Copy(&rt, &route.Spec)
				if err != nil {
					logger.Error(err, "Failed to deep copy Route", "Name", route.Name)
					return false, err
				}
				rt.Name = route.Name
				gomaConfig.Gateway.Routes = append(gomaConfig.Gateway.Routes, rt)
				middlewareNames = append(middlewareNames, rt.Middlewares...)

			}
		}
	}
	for _, mid := range middlewares.Items {
		middleware := *mapMid(mid)
		logger.Info("Adding Middleware", "Name", middleware.Name)
		if slices.Contains(middlewareNames, middleware.Name) {
			gomaConfig.Middlewares = append(gomaConfig.Middlewares, middleware)
		}

	}

	yamlContent, err := yaml.Marshal(&gomaConfig)
	if err != nil {
		logger.Error(err, "Unable to marshal YAML")
		return false, err
	}
	// Define the desired ConfigMap
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gateway.Name,
			Namespace: req.Namespace,
			Labels: map[string]string{
				"belongs-to": BelongsTo,
				"gateway":    gateway.Name,
			},
		},
		Data: map[string]string{
			ConfigName: strings.TrimSpace(string(yamlContent)),
		},
	}
	// Check if the ConfigMap already exists
	existingConfigMap := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, existingConfigMap)
	if err != nil && client.IgnoreNotFound(err) != nil {
		logger.Error(err, "Failed to get ConfigMap")
		return false, err
	}

	if err != nil && client.IgnoreNotFound(err) == nil {
		// Create the ConfigMap if it doesn't exist
		if err = controllerutil.SetControllerReference(&gateway, configMap, r.Scheme); err != nil {
			logger.Error(err, "Failed to set controller reference")
			return false, err
		}
		if err = r.Create(ctx, configMap); err != nil {
			logger.Error(err, "Failed to create ConfigMap")
			return false, err
		}
		logger.Info("Created ConfigMap", "ConfigMap.Name", configMap.Name)
	} else {
		// Optional: Update the ConfigMap if needed
		if !reflect.DeepEqual(existingConfigMap.Data, configMap.Data) {
			logger.Info("Updating ConfigMap...", "ConfigMap.Name", configMap.Name)
			// ConfigMap data is not equal, update it
			existingConfigMap.Data = configMap.Data
			if err = r.Update(ctx, existingConfigMap); err != nil {
				return false, err
			}
			logger.Info("Updated ConfigMap", "ConfigMap.Name", configMap.Name)

		}
	}
	return true, nil

}

// mapMid converts RawExtension to struct
func mapMid(middleware gomaprojv1beta1.Middleware) *Middleware {
	mid := &Middleware{
		Name:  middleware.Name,
		Type:  middleware.Spec.Type,
		Paths: middleware.Spec.Paths,
	}

	// Mapping of middleware types to their respective struct types
	ruleMapping := map[string]interface{}{
		BasicAuth:                  &BasicRuleMiddleware{},
		OAuth:                      &OauthRulerMiddleware{},
		JWTAuth:                    &JWTRuleMiddleware{},
		RateLimit:                  &RateLimitRuleMiddleware{},
		strings.ToLower(RateLimit): &RateLimitRuleMiddleware{},
		accessPolicy:               &AccessPolicyRuleMiddleware{},
		addPrefix:                  &AddPrefixRuleMiddleware{},
		redirectRegex:              &RedirectRegexRuleMiddleware{},
		forwardAuth:                &ForwardAuthRuleMiddleware{},
	}

	rule, exists := ruleMapping[middleware.Spec.Type]
	if !exists {
		return mid
	}

	// Attempt to convert the rule to the appropriate struct
	err := ConvertRawExtensionToStruct(middleware.Spec.Rule, rule)
	if err != nil {
		return mid
	}

	mid.Rule = rule
	return mid
}

// Helper function to return a pointer to an int32
func int32Ptr(i int32) *int32 {
	return &i
}
func ConvertRawExtensionToStruct(raw runtime.RawExtension, out interface{}) error {
	// Unmarshal the raw JSON into the provided struct
	if err := json.Unmarshal(raw.Raw, out); err != nil {
		return err
	}
	return nil
}
