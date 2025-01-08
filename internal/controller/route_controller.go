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
	gomaprojv1beta1 "github.com/jkaninda/goma-operator/api/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// RouteReconciler reconciles a Route object
type RouteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=gomaproj.github.io,resources=routes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gomaproj.github.io,resources=routes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gomaproj.github.io,resources=routes/finalizers,verbs=update
// +kubebuilder:rbac:groups=gomaproj.github.io,resources=middlewares,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gomaproj.github.io,resources=middlewares/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gomaproj.github.io,resources=middlewares/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Route object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *RouteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	// Fetch the custom resource
	route := &gomaprojv1beta1.Route{}
	if err := r.Get(ctx, req.NamespacedName, route); err != nil {
		logger.Error(err, "Unable to fetch CustomResource")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	var gateway gomaprojv1beta1.Gateway
	if err := r.Get(ctx, types.NamespacedName{Name: route.Spec.Gateway, Namespace: route.Namespace}, &gateway); err != nil {
		logger.Error(err, "Failed to fetch Gateway")
		return ctrl.Result{}, err
	}

	// Handle finalizer logic
	if route.ObjectMeta.DeletionTimestamp.IsZero() {
		// Object is not being deleted: ensure the finalizer is added
		if !controllerutil.ContainsFinalizer(route, FinalizerName) {
			controllerutil.AddFinalizer(route, FinalizerName)
			if err := r.Update(ctx, route); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to add finalizer: %w", err)
			}
		}
	} else {
		// Object is being deleted: handle finalization logic
		if controllerutil.ContainsFinalizer(route, FinalizerName) {
			// Execute finalization steps
			completed, err := updateGatewayConfig(*r, ctx, req, gateway)
			if err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to update gateway config: %w", err)
			}
			if completed {
				if err := restartDeployment(r.Client, ctx, req, &gateway); err != nil {
					return ctrl.Result{}, fmt.Errorf("failed to restart deployment: %w", err)
				}
			}
			// Remove the finalizer after successful cleanup
			controllerutil.RemoveFinalizer(route, FinalizerName)
			if err := r.Update(ctx, route); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to remove finalizer: %w", err)
			}
		}
		// No further reconciliation needed for a deleted object
		return ctrl.Result{}, nil
	}

	ok, err := updateGatewayConfig(*r, ctx, req, gateway)
	if err != nil {
		return ctrl.Result{}, err
	}
	if ok {
		if err = restartDeployment(r.Client, ctx, req, &gateway); err != nil {
			logger.Error(err, "Failed to restart Deployment")
			return ctrl.Result{}, err

		}
	}

	logger.Info("Reconciliation complete")
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RouteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gomaprojv1beta1.Route{}).
		Named("route").
		Complete(r)
}
