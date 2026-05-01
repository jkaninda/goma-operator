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
	"strings"

	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	gatewayv1alpha1 "github.com/jkaninda/goma-operator/api/v1alpha1"
	"github.com/jkaninda/goma-operator/internal/validation"
)

// RouteReconciler reconciles a Route object
type RouteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=gateway.jkaninda.dev,resources=routes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gateway.jkaninda.dev,resources=routes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gateway.jkaninda.dev,resources=routes/finalizers,verbs=update

func (r *RouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch Route
	route := &gatewayv1alpha1.Route{}
	if err := r.Get(ctx, req.NamespacedName, route); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	originalStatus := route.Status.DeepCopy()

	// Validate spec
	errs := validation.ValidateRouteSpec(&route.Spec)

	// Check that every referenced Gateway exists. The Route is considered
	// ready only if at least one referenced gateway exists AND no missing
	// gateways are reported.
	var foundGateways, missingGateways []string
	for _, gwName := range route.Spec.Gateways {
		if gwName == "" {
			continue
		}
		gw := &gatewayv1alpha1.Gateway{}
		err := r.Get(ctx, types.NamespacedName{Name: gwName, Namespace: route.Namespace}, gw)
		if apierrors.IsNotFound(err) {
			missingGateways = append(missingGateways, gwName)
		} else if err != nil {
			return ctrl.Result{}, err
		} else {
			foundGateways = append(foundGateways, gwName)
		}
	}
	if len(missingGateways) > 0 {
		errs = append(errs, fmt.Sprintf("gateway(s) not found in namespace %s: %s",
			route.Namespace, strings.Join(missingGateways, ", ")))
	}
	gatewayExists := len(foundGateways) > 0

	// Check that every referenced Middleware CR exists in the same namespace.
	// Missing middlewares block Ready and surface as a dedicated condition so
	// users can tell a typo'd reference apart from a spec error.
	var missingMiddlewares []string
	for _, mwName := range route.Spec.Middlewares {
		if mwName == "" {
			continue
		}
		mw := &gatewayv1alpha1.Middleware{}
		err := r.Get(ctx, types.NamespacedName{Name: mwName, Namespace: route.Namespace}, mw)
		if apierrors.IsNotFound(err) {
			missingMiddlewares = append(missingMiddlewares, mwName)
		} else if err != nil {
			return ctrl.Result{}, err
		}
	}
	if len(missingMiddlewares) > 0 {
		meta.SetStatusCondition(&route.Status.Conditions, metav1.Condition{
			Type:    "MiddlewaresResolved",
			Status:  metav1.ConditionFalse,
			Reason:  "MiddlewareNotFound",
			Message: fmt.Sprintf("middleware(s) not found in namespace %s: %s", route.Namespace, strings.Join(missingMiddlewares, ", ")),
		})
		errs = append(errs, fmt.Sprintf("middleware(s) not found: %s", strings.Join(missingMiddlewares, ", ")))
	} else {
		meta.SetStatusCondition(&route.Status.Conditions, metav1.Condition{
			Type:    "MiddlewaresResolved",
			Status:  metav1.ConditionTrue,
			Reason:  "AllMiddlewaresFound",
			Message: "All referenced middlewares exist",
		})
	}

	// Update status conditions (omit LastTransitionTime — meta.SetStatusCondition
	// only sets it when Status actually transitions, keeping updates idempotent).
	if len(errs) > 0 {
		route.Status.Ready = false
		meta.SetStatusCondition(&route.Status.Conditions, metav1.Condition{
			Type:    "Valid",
			Status:  metav1.ConditionFalse,
			Reason:  "ValidationFailed",
			Message: strings.Join(errs, "; "),
		})
		logger.Info("Route validation failed", "route", route.Name, "errors", errs)
	} else {
		route.Status.Ready = gatewayExists
		meta.SetStatusCondition(&route.Status.Conditions, metav1.Condition{
			Type:    "Valid",
			Status:  metav1.ConditionTrue,
			Reason:  "ValidationPassed",
			Message: "Route spec is valid",
		})
	}

	// Skip no-op status updates to avoid triggering another reconcile.
	if !equality.Semantic.DeepEqual(route.Status, *originalStatus) {
		if err := r.Status().Update(ctx, route); err != nil {
			return ctrl.Result{}, err
		}
	}

	// The GatewayReconciler watches Route events and will reconcile
	// the parent Gateway automatically via the mapRouteToGateway handler.

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Only reconcile on spec changes. Without this, status updates we write
	// here would re-trigger our own reconcile in an infinite loop.
	return ctrl.NewControllerManagedBy(mgr).
		For(&gatewayv1alpha1.Route{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Watches(
			&gatewayv1alpha1.Middleware{},
			handler.EnqueueRequestsFromMapFunc(r.mapMiddlewareToRoutes),
		).
		Named("route").
		Complete(r)
}

// mapMiddlewareToRoutes enqueues any Route in the same namespace that
// references the given Middleware by name. Triggered on Middleware
// create/update/delete so Routes with dangling references re-reconcile as
// soon as the middleware appears.
func (r *RouteReconciler) mapMiddlewareToRoutes(ctx context.Context, obj client.Object) []reconcile.Request {
	mw, ok := obj.(*gatewayv1alpha1.Middleware)
	if !ok {
		return nil
	}
	routeList := &gatewayv1alpha1.RouteList{}
	if err := r.List(ctx, routeList, client.InNamespace(mw.Namespace)); err != nil {
		return nil
	}
	var reqs []reconcile.Request
	for _, rt := range routeList.Items {
		for _, name := range rt.Spec.Middlewares {
			if name == mw.Name {
				reqs = append(reqs, reconcile.Request{
					NamespacedName: types.NamespacedName{Name: rt.Name, Namespace: rt.Namespace},
				})
				break
			}
		}
	}
	return reqs
}
