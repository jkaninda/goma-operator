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
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	gatewayv1alpha1 "github.com/jkaninda/goma-operator/api/v1alpha1"
	"github.com/jkaninda/goma-operator/internal/validation"
)

// MiddlewareReconciler reconciles a Middleware object
type MiddlewareReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=gateway.jkaninda.dev,resources=middlewares,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gateway.jkaninda.dev,resources=middlewares/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gateway.jkaninda.dev,resources=middlewares/finalizers,verbs=update

func (r *MiddlewareReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Fetch Middleware
	mw := &gatewayv1alpha1.Middleware{}
	if err := r.Get(ctx, req.NamespacedName, mw); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	originalStatus := mw.Status.DeepCopy()

	// Validate spec
	errs := validation.ValidateMiddlewareSpec(&mw.Spec)

	// Find all Routes referencing this middleware
	routeList := &gatewayv1alpha1.RouteList{}
	if err := r.List(ctx, routeList, client.InNamespace(mw.Namespace)); err != nil {
		return ctrl.Result{}, err
	}

	var referencedBy []string
	for _, route := range routeList.Items {
		for _, mwName := range route.Spec.Middlewares {
			if mwName == mw.Name {
				referencedBy = append(referencedBy, route.Name)
				break
			}
		}
	}
	sort.Strings(referencedBy)

	// Update status
	mw.Status.ReferencedBy = referencedBy

	if len(errs) > 0 {
		mw.Status.Ready = false
		meta.SetStatusCondition(&mw.Status.Conditions, metav1.Condition{
			Type:    "Valid",
			Status:  metav1.ConditionFalse,
			Reason:  "ValidationFailed",
			Message: strings.Join(errs, "; "),
		})
		logger.Info("Middleware validation failed", "middleware", mw.Name, "errors", errs)
	} else {
		mw.Status.Ready = true
		meta.SetStatusCondition(&mw.Status.Conditions, metav1.Condition{
			Type:    "Valid",
			Status:  metav1.ConditionTrue,
			Reason:  "ValidationPassed",
			Message: "Middleware spec is valid",
		})
	}

	// Skip no-op status updates to avoid triggering another reconcile.
	if !equality.Semantic.DeepEqual(mw.Status, *originalStatus) {
		if err := r.Status().Update(ctx, mw); err != nil {
			return ctrl.Result{}, err
		}
	}

	// The GatewayReconciler watches Middleware events and will reconcile
	// affected Gateways automatically via the mapMiddlewareToGateway handler.

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MiddlewareReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// Only reconcile on spec changes to avoid status-update feedback loops.
	return ctrl.NewControllerManagedBy(mgr).
		For(&gatewayv1alpha1.Middleware{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Named("middleware").
		Complete(r)
}
