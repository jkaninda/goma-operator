/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"sort"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	gatewayv1alpha1 "github.com/jkaninda/goma-operator/api/v1alpha1"
	"github.com/jkaninda/goma-operator/internal/converter"
	"github.com/jkaninda/goma-operator/internal/resources"
)

const (
	gatewayFinalizer = "gateway.jkaninda.dev/finalizer"
)

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=gateway.jkaninda.dev,resources=gateways,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gateway.jkaninda.dev,resources=gateways/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gateway.jkaninda.dev,resources=gateways/finalizers,verbs=update
// +kubebuilder:rbac:groups=gateway.jkaninda.dev,resources=routes,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.jkaninda.dev,resources=routes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gateway.jkaninda.dev,resources=middlewares,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=configmaps;services;secrets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles;rolebindings,verbs=get;list;watch;create;update;patch;delete

func (r *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch Gateway
	gw := &gatewayv1alpha1.Gateway{}
	if err := r.Get(ctx, req.NamespacedName, gw); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Snapshot the original status so we can detect real changes at the end
	// and skip no-op status updates (which would re-trigger reconcile).
	originalStatus := gw.Status.DeepCopy()

	// Handle deletion
	if !gw.DeletionTimestamp.IsZero() {
		if controllerutil.ContainsFinalizer(gw, gatewayFinalizer) {
			controllerutil.RemoveFinalizer(gw, gatewayFinalizer)
			if err := r.Update(ctx, gw); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	// Add finalizer
	if !controllerutil.ContainsFinalizer(gw, gatewayFinalizer) {
		controllerutil.AddFinalizer(gw, gatewayFinalizer)
		if err := r.Update(ctx, gw); err != nil {
			return ctrl.Result{}, err
		}
	}

	// List Routes for this gateway
	routeList := &gatewayv1alpha1.RouteList{}
	if err := r.List(ctx, routeList, client.InNamespace(gw.Namespace)); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to list routes: %w", err)
	}

	// Filter routes that reference this gateway in their spec.gateways list
	var gatewayRoutes []gatewayv1alpha1.Route
	referencedMiddlewares := make(map[string]bool)
	for _, route := range routeList.Items {
		for _, gwName := range route.Spec.Gateways {
			if gwName == gw.Name {
				gatewayRoutes = append(gatewayRoutes, route)
				for _, mwName := range route.Spec.Middlewares {
					referencedMiddlewares[mwName] = true
				}
				break
			}
		}
	}

	// Also include middlewares referenced by the Gateway's own monitoring
	// section (metrics endpoint protection) so they end up in the generated
	// config even when no Route uses them.
	if gw.Spec.Server.Monitoring != nil && gw.Spec.Server.Monitoring.Middleware != nil {
		for _, mwName := range gw.Spec.Server.Monitoring.Middleware.Metrics {
			referencedMiddlewares[mwName] = true
		}
	}

	// Sort routes by priority (desc) then name (asc)
	sort.Slice(gatewayRoutes, func(i, j int) bool {
		if gatewayRoutes[i].Spec.Priority != gatewayRoutes[j].Spec.Priority {
			return gatewayRoutes[i].Spec.Priority > gatewayRoutes[j].Spec.Priority
		}
		return gatewayRoutes[i].Name < gatewayRoutes[j].Name
	})

	// List and filter middlewares referenced by routes
	mwList := &gatewayv1alpha1.MiddlewareList{}
	if err := r.List(ctx, mwList, client.InNamespace(gw.Namespace)); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to list middlewares: %w", err)
	}

	var gatewayMiddlewares []gatewayv1alpha1.Middleware
	for _, mw := range mwList.Items {
		if referencedMiddlewares[mw.Name] {
			gatewayMiddlewares = append(gatewayMiddlewares, mw)
		}
	}

	// Sort middlewares by name
	sort.Slice(gatewayMiddlewares, func(i, j int) bool {
		return gatewayMiddlewares[i].Name < gatewayMiddlewares[j].Name
	})

	// Build config
	cfg := converter.GatewayConfigFromCRs(gw, gatewayRoutes, gatewayMiddlewares)

	// ServiceAccount + Role + RoleBinding
	// The gateway pod needs a ServiceAccount with permissions to watch
	// Route/Middleware CRDs (used by the goma-k8s-provider sidecar) and
	// read/write the ACME secret.
	if err := r.reconcileRBAC(ctx, gw); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to reconcile RBAC: %w", err)
	}

	// ConfigMap
	desiredCM, checksum, err := resources.BuildConfigMap(gw, cfg)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to build configmap: %w", err)
	}
	if err := controllerutil.SetControllerReference(gw, desiredCM, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	existingCM := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: desiredCM.Name, Namespace: desiredCM.Namespace}, existingCM)
	configChanged := false
	if apierrors.IsNotFound(err) {
		logger.Info("Creating ConfigMap", "name", desiredCM.Name)
		if err := r.Create(ctx, desiredCM); err != nil {
			return ctrl.Result{}, err
		}
		configChanged = true
	} else if err != nil {
		return ctrl.Result{}, err
	} else if !equality.Semantic.DeepEqual(existingCM.Data, desiredCM.Data) {
		existingCM.Data = desiredCM.Data
		logger.Info("Updating ConfigMap", "name", desiredCM.Name)
		if err := r.Update(ctx, existingCM); err != nil {
			return ctrl.Result{}, err
		}
		configChanged = true
	}

	// Deployment
	desiredDep := resources.BuildDeployment(gw, gatewayRoutes)
	if err := controllerutil.SetControllerReference(gw, desiredDep, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	existingDep := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: desiredDep.Name, Namespace: desiredDep.Namespace}, existingDep)
	if apierrors.IsNotFound(err) {
		if configChanged {
			// Set restart annotation on initial creation
			if desiredDep.Spec.Template.Annotations == nil {
				desiredDep.Spec.Template.Annotations = make(map[string]string)
			}
			desiredDep.Spec.Template.Annotations["goma.jkaninda.dev/restarted-at"] = time.Now().Format(time.RFC3339)
		}
		logger.Info("Creating Deployment", "name", desiredDep.Name)
		if err := r.Create(ctx, desiredDep); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		// Update deployment spec — replace the entire PodSpec and Labels so
		// newly-added fields (ServiceAccountName, etc.) propagate to existing
		// Deployments on operator upgrade. Preserve any pre-existing template
		// annotations that we don't manage.
		existingDep.Spec.Replicas = desiredDep.Spec.Replicas
		existingDep.Spec.Template.Spec = desiredDep.Spec.Template.Spec
		existingDep.Spec.Template.Labels = desiredDep.Spec.Template.Labels

		// Rolling restart on config change
		if configChanged {
			if existingDep.Spec.Template.Annotations == nil {
				existingDep.Spec.Template.Annotations = make(map[string]string)
			}
			existingDep.Spec.Template.Annotations["goma.jkaninda.dev/restarted-at"] = time.Now().Format(time.RFC3339)
		}

		logger.Info("Updating Deployment", "name", existingDep.Name)
		if err := r.Update(ctx, existingDep); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Service
	desiredSvc := resources.BuildService(gw)
	if err := controllerutil.SetControllerReference(gw, desiredSvc, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	existingSvc := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: desiredSvc.Name, Namespace: desiredSvc.Namespace}, existingSvc)
	if apierrors.IsNotFound(err) {
		logger.Info("Creating Service", "name", desiredSvc.Name)
		if err := r.Create(ctx, desiredSvc); err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		// Update the Service while preserving fields that are immutable
		// after creation (clusterIP, clusterIPs, and the auto-assigned
		// nodePorts when the user didn't pin them).
		preservedClusterIP := existingSvc.Spec.ClusterIP
		preservedClusterIPs := existingSvc.Spec.ClusterIPs

		// Carry over auto-assigned NodePorts if the user didn't set them
		// explicitly, to avoid churn on every reconcile.
		for i := range desiredSvc.Spec.Ports {
			if desiredSvc.Spec.Ports[i].NodePort == 0 {
				for _, ep := range existingSvc.Spec.Ports {
					if ep.Name == desiredSvc.Spec.Ports[i].Name && ep.NodePort != 0 {
						desiredSvc.Spec.Ports[i].NodePort = ep.NodePort
					}
				}
			}
		}

		existingSvc.Labels = desiredSvc.Labels
		existingSvc.Annotations = desiredSvc.Annotations
		existingSvc.Spec = desiredSvc.Spec
		existingSvc.Spec.ClusterIP = preservedClusterIP
		existingSvc.Spec.ClusterIPs = preservedClusterIPs

		if err := r.Update(ctx, existingSvc); err != nil {
			return ctrl.Result{}, err
		}
	}

	// HPA
	if gw.Spec.AutoScaling != nil && gw.Spec.AutoScaling.Enabled {
		desiredHPA := resources.BuildHPA(gw)
		if err := controllerutil.SetControllerReference(gw, desiredHPA, r.Scheme); err != nil {
			return ctrl.Result{}, err
		}

		existingHPA := &autoscalingv2.HorizontalPodAutoscaler{}
		err = r.Get(ctx, types.NamespacedName{Name: desiredHPA.Name, Namespace: desiredHPA.Namespace}, existingHPA)
		if apierrors.IsNotFound(err) {
			logger.Info("Creating HPA", "name", desiredHPA.Name)
			if err := r.Create(ctx, desiredHPA); err != nil {
				return ctrl.Result{}, err
			}
		} else if err != nil {
			return ctrl.Result{}, err
		} else {
			existingHPA.Spec = desiredHPA.Spec
			if err := r.Update(ctx, existingHPA); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// Delete HPA if autoscaling is disabled
		existingHPA := &autoscalingv2.HorizontalPodAutoscaler{}
		err = r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, existingHPA)
		if err == nil {
			logger.Info("Deleting HPA (autoscaling disabled)", "name", gw.Name)
			if err := r.Delete(ctx, existingHPA); err != nil && !apierrors.IsNotFound(err) {
				return ctrl.Result{}, err
			}
		}
	}

	// NOTE: Route statuses are owned by RouteReconciler (which triggers on
	// Route spec changes). Updating Route status here would cause an infinite
	// reconcile loop because Route status updates trigger the Route watcher,
	// which re-enqueues the Gateway via mapRouteToGateway.

	// Update Gateway status
	// Fetch current deployment status
	currentDep := &appsv1.Deployment{}
	depReplicas := int32(0)
	depReady := int32(0)
	if err := r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, currentDep); err == nil {
		depReplicas = currentDep.Status.Replicas
		depReady = currentDep.Status.ReadyReplicas
	}

	gw.Status.Replicas = depReplicas
	gw.Status.ReadyReplicas = depReady
	gw.Status.Routes = int32(len(gatewayRoutes))
	gw.Status.Middlewares = int32(len(gatewayMiddlewares))
	gw.Status.ConfigChecksum = checksum

	// Compute external addresses from the Service
	gw.Status.Addresses = r.computeGatewayAddresses(ctx, gw)

	available := metav1.ConditionFalse
	if depReady > 0 {
		available = metav1.ConditionTrue
	}
	// NOTE: omit LastTransitionTime — meta.SetStatusCondition only sets it
	// when Status actually transitions, which is required for idempotent
	// reconciliation (avoids infinite status-update → reconcile loops).
	meta.SetStatusCondition(&gw.Status.Conditions, metav1.Condition{
		Type:    "Available",
		Status:  available,
		Reason:  "DeploymentStatus",
		Message: fmt.Sprintf("%d/%d replicas ready", depReady, depReplicas),
	})
	meta.SetStatusCondition(&gw.Status.Conditions, metav1.Condition{
		Type:    "ConfigSynced",
		Status:  metav1.ConditionTrue,
		Reason:  "ConfigGenerated",
		Message: fmt.Sprintf("Config generated from %d routes, %d middlewares", len(gatewayRoutes), len(gatewayMiddlewares)),
	})

	// Only update status if something actually changed, otherwise the update
	// bumps resourceVersion and triggers another reconcile.
	if !equality.Semantic.DeepEqual(gw.Status, *originalStatus) {
		if err := r.Status().Update(ctx, gw); err != nil {
			return ctrl.Result{}, err
		}
	}

	logger.Info("Reconciled gateway",
		"routes", len(gatewayRoutes),
		"middlewares", len(gatewayMiddlewares),
		"configChanged", configChanged,
	)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// GenerationChangedPredicate ignores status-only updates so that
	// RouteReconciler / MiddlewareReconciler status writes don't trigger
	// the Gateway to reconcile (which would cause an infinite loop).
	specChangedOnly := builder.WithPredicates(predicate.GenerationChangedPredicate{})

	return ctrl.NewControllerManagedBy(mgr).
		For(&gatewayv1alpha1.Gateway{}, specChangedOnly).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.ConfigMap{}).
		Owns(&corev1.Service{}).
		Owns(&autoscalingv2.HorizontalPodAutoscaler{}).
		Watches(&gatewayv1alpha1.Route{},
			handler.EnqueueRequestsFromMapFunc(r.mapRouteToGateway),
			specChangedOnly,
		).
		Watches(&gatewayv1alpha1.Middleware{},
			handler.EnqueueRequestsFromMapFunc(r.mapMiddlewareToGateway),
			specChangedOnly,
		).
		Named("gateway").
		Complete(r)
}

// reconcileRBAC ensures the ServiceAccount, Role, and RoleBinding exist for
// the gateway pod. The gateway pod's sidecar (goma-k8s-provider) requires
// these to watch Route/Middleware CRDs from within the cluster.
func (r *GatewayReconciler) reconcileRBAC(ctx context.Context, gw *gatewayv1alpha1.Gateway) error {
	logger := log.FromContext(ctx)

	// ServiceAccount
	desiredSA := resources.BuildServiceAccount(gw)
	if err := controllerutil.SetControllerReference(gw, desiredSA, r.Scheme); err != nil {
		return err
	}
	existingSA := &corev1.ServiceAccount{}
	err := r.Get(ctx, types.NamespacedName{Name: desiredSA.Name, Namespace: desiredSA.Namespace}, existingSA)
	if apierrors.IsNotFound(err) {
		logger.Info("Creating ServiceAccount", "name", desiredSA.Name)
		if err := r.Create(ctx, desiredSA); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Role
	desiredRole := resources.BuildRole(gw)
	if err := controllerutil.SetControllerReference(gw, desiredRole, r.Scheme); err != nil {
		return err
	}
	existingRole := &rbacv1.Role{}
	err = r.Get(ctx, types.NamespacedName{Name: desiredRole.Name, Namespace: desiredRole.Namespace}, existingRole)
	if apierrors.IsNotFound(err) {
		logger.Info("Creating Role", "name", desiredRole.Name)
		if err := r.Create(ctx, desiredRole); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if !equality.Semantic.DeepEqual(existingRole.Rules, desiredRole.Rules) {
		existingRole.Rules = desiredRole.Rules
		logger.Info("Updating Role", "name", desiredRole.Name)
		if err := r.Update(ctx, existingRole); err != nil {
			return err
		}
	}

	// RoleBinding
	desiredRB := resources.BuildRoleBinding(gw)
	if err := controllerutil.SetControllerReference(gw, desiredRB, r.Scheme); err != nil {
		return err
	}
	existingRB := &rbacv1.RoleBinding{}
	err = r.Get(ctx, types.NamespacedName{Name: desiredRB.Name, Namespace: desiredRB.Namespace}, existingRB)
	if apierrors.IsNotFound(err) {
		logger.Info("Creating RoleBinding", "name", desiredRB.Name)
		if err := r.Create(ctx, desiredRB); err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else if !equality.Semantic.DeepEqual(existingRB.Subjects, desiredRB.Subjects) ||
		!equality.Semantic.DeepEqual(existingRB.RoleRef, desiredRB.RoleRef) {
		existingRB.Subjects = desiredRB.Subjects
		existingRB.RoleRef = desiredRB.RoleRef
		logger.Info("Updating RoleBinding", "name", desiredRB.Name)
		if err := r.Update(ctx, existingRB); err != nil {
			return err
		}
	}

	return nil
}

// computeGatewayAddresses resolves the set of addresses at which the gateway
// is reachable, based on the Service type and status:
//
//   - LoadBalancer → entries from .status.loadBalancer.ingress[] (IP or hostname)
//   - NodePort     → the Service's external IPs if set, else the cluster IP
//   - ClusterIP    → the Service's .spec.clusterIP (in-cluster only)
//
// Returns nil if the Service does not exist yet or has no reachable addresses.
func (r *GatewayReconciler) computeGatewayAddresses(ctx context.Context, gw *gatewayv1alpha1.Gateway) []gatewayv1alpha1.GatewayAddress {
	svc := &corev1.Service{}
	if err := r.Get(ctx, types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}, svc); err != nil {
		return nil
	}

	var addrs []gatewayv1alpha1.GatewayAddress

	switch svc.Spec.Type {
	case corev1.ServiceTypeLoadBalancer:
		for _, ing := range svc.Status.LoadBalancer.Ingress {
			if ing.IP != "" {
				addrs = append(addrs, gatewayv1alpha1.GatewayAddress{
					Type:  "IPAddress",
					Value: ing.IP,
				})
			}
			if ing.Hostname != "" {
				addrs = append(addrs, gatewayv1alpha1.GatewayAddress{
					Type:  "Hostname",
					Value: ing.Hostname,
				})
			}
		}
	case corev1.ServiceTypeNodePort:
		for _, ip := range svc.Spec.ExternalIPs {
			addrs = append(addrs, gatewayv1alpha1.GatewayAddress{
				Type:  "IPAddress",
				Value: ip,
			})
		}
		if len(addrs) == 0 && svc.Spec.ClusterIP != "" && svc.Spec.ClusterIP != corev1.ClusterIPNone {
			addrs = append(addrs, gatewayv1alpha1.GatewayAddress{
				Type:  "IPAddress",
				Value: svc.Spec.ClusterIP,
			})
		}
	default:
		if svc.Spec.ClusterIP != "" && svc.Spec.ClusterIP != corev1.ClusterIPNone {
			addrs = append(addrs, gatewayv1alpha1.GatewayAddress{
				Type:  "IPAddress",
				Value: svc.Spec.ClusterIP,
			})
		}
	}

	return addrs
}

// mapRouteToGateway maps a Route event to each Gateway listed in its spec.
func (r *GatewayReconciler) mapRouteToGateway(ctx context.Context, obj client.Object) []reconcile.Request {
	route, ok := obj.(*gatewayv1alpha1.Route)
	if !ok {
		return nil
	}
	requests := make([]reconcile.Request, 0, len(route.Spec.Gateways))
	for _, gwName := range route.Spec.Gateways {
		if gwName == "" {
			continue
		}
		requests = append(requests, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: gwName, Namespace: route.Namespace},
		})
	}
	return requests
}

// mapMiddlewareToGateway maps a Middleware event to all Gateways that reference it via routes.
func (r *GatewayReconciler) mapMiddlewareToGateway(ctx context.Context, obj client.Object) []reconcile.Request {
	mw, ok := obj.(*gatewayv1alpha1.Middleware)
	if !ok {
		return nil
	}

	// Find all routes in this namespace referencing this middleware
	routeList := &gatewayv1alpha1.RouteList{}
	if err := r.List(ctx, routeList, client.InNamespace(mw.Namespace)); err != nil {
		return nil
	}

	gateways := make(map[string]bool)
	for _, route := range routeList.Items {
		for _, mwName := range route.Spec.Middlewares {
			if mwName == mw.Name {
				for _, gwName := range route.Spec.Gateways {
					if gwName != "" {
						gateways[gwName] = true
					}
				}
				break
			}
		}
	}

	requests := make([]reconcile.Request, 0, len(gateways))
	for gwName := range gateways {
		requests = append(requests, reconcile.Request{
			NamespacedName: types.NamespacedName{Name: gwName, Namespace: mw.Namespace},
		})
	}
	return requests
}
