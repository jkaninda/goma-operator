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
	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/apps/v1"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"strings"
)

// GatewayReconciler reconciles a Gateway object
type GatewayReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=gomaproj.github.io,resources=gateways,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gomaproj.github.io,resources=gateways/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gomaproj.github.io,resources=gateways/finalizers,verbs=update
// +kubebuilder:rbac:groups=gomaproj.github.io,resources=middlewares,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=gomaproj.github.io,resources=middlewares/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=gomaproj.github.io,resources=middlewares/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=events,verbs=create;update;patch;
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=autoscaling,resources=horizontalpodautoscalers,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the Gateway object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	imageName := AppImageName
	// Fetch the custom resource
	gateway := &gomaprojv1beta1.Gateway{}
	if err := r.Get(ctx, req.NamespacedName, gateway); err != nil {
		logger.Error(err, "Unable to fetch Gateway")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if the object is being deleted and if so, handle it
	if gateway.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(gateway, FinalizerName) {
			controllerutil.AddFinalizer(gateway, FinalizerName)
			err := r.Update(ctx, gateway)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		if controllerutil.ContainsFinalizer(gateway, FinalizerName) {
			// Once finalization is done, remove the finalizer
			if err := r.finalize(ctx, gateway); err != nil {
				return ctrl.Result{}, err
			}

			controllerutil.RemoveFinalizer(gateway, FinalizerName)
			err := r.Update(ctx, gateway)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}

	if gateway.Spec.GatewayVersion != "" {
		imageName = fmt.Sprintf("%s:%s", AppImageName, gateway.Spec.GatewayVersion)
	}
	if gateway.Spec.ReplicaCount != 0 {
		ReplicaCount = gateway.Spec.ReplicaCount
	}
	gomaConfig := gatewayConfig(*r, ctx, req, gateway)
	yamlContent, err := yaml.Marshal(&gomaConfig)
	if err != nil {
		logger.Error(err, "Unable to marshal YAML")
		return ctrl.Result{}, err
	}
	// Define the desired ConfigMap
	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
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
	var existingConfigMap corev1.ConfigMap
	err = r.Get(ctx, types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace}, &existingConfigMap)
	if err != nil && client.IgnoreNotFound(err) != nil {
		logger.Error(err, "Failed to get ConfigMap")
		return ctrl.Result{}, err
	}

	if err != nil && client.IgnoreNotFound(err) == nil {
		// Create the ConfigMap if it doesn't exist
		if err := controllerutil.SetControllerReference(gateway, configMap, r.Scheme); err != nil {
			logger.Error(err, "Failed to set controller reference")
			return ctrl.Result{}, err
		}
		if err := r.Create(ctx, configMap); err != nil {
			addCondition(&gateway.Status, "ConfigMapNotReady", metav1.ConditionFalse, "ConfigMapNotReady", "Failed to add configMap for Gateway")
			logger.Error(err, "Failed to create ConfigMap")
			return ctrl.Result{}, err
		}
		logger.Info("Created ConfigMap", "ConfigMap.Name", configMap.Name)
	} else {
		// Optional: Update the ConfigMap if needed
		if !equalConfigMapData(existingConfigMap.Data, configMap.Data) {
			existingConfigMap.Data = configMap.Data
			if err := r.Update(ctx, &existingConfigMap); err != nil {
				logger.Error(err, "Failed to update ConfigMap")
				addCondition(&gateway.Status, "ConfigMapReady", metav1.ConditionFalse, "ConfigMapReady", "Failed to update ConfigMap for Gateway")
				return ctrl.Result{}, err
			}
			if err = restartDeployment(r.Client, ctx, req, gateway); err != nil {
				logger.Error(err, "Failed to restart Deployment")
				return ctrl.Result{}, err

			}
			logger.Info("Updated ConfigMap", "ConfigMap.Name", configMap.Name)
		}

	}
	err = createUpdateDeployment(*r, ctx, req, *gateway, imageName)
	if err != nil {
		addCondition(&gateway.Status, "DeploymentNotReady", metav1.ConditionFalse, "DeploymentNotReady", "Failed to created deployment for Gateway")
		logger.Error(err, "Failed to create Deployment")
		return ctrl.Result{}, err
	}

	err = createService(*r, ctx, req, gateway)
	if err != nil {
		addCondition(&gateway.Status, "ServiceNotReady", metav1.ConditionFalse, "ServiceNotReady", "Failed to create Service for Gateway")
		logger.Error(err, "Failed to create Service")
		return ctrl.Result{}, err
	}

	addCondition(&gateway.Status, "GatewayReady", metav1.ConditionTrue, "AllSubresourcesReady", "All subresources are ready")
	logger.Info("All Subresources ready")

	// Update the Status
	gateway.Status.Replicas = gateway.Spec.ReplicaCount
	if gateway.Spec.AutoScaling.Enabled {
		replicas, err := currentReplicas(ctx, r.Client, gateway.Name, gateway.Namespace)
		if err != nil {
			logger.Error(err, "Failed to get current replicas")
		}
		gateway.Status.Replicas = replicas

	}
	gateway.Status.Routes = int32(len(gomaConfig.Gateway.Routes))
	if err = r.updateStatus(ctx, gateway); err != nil {
		logger.Error(err, "Failed to update resource status")
		return ctrl.Result{}, err
	}
	logger.Info("Successfully updated resource status")
	return ctrl.Result{}, nil
}
func (r *GatewayReconciler) updateStatus(ctx context.Context, gateway *gomaprojv1beta1.Gateway) error {
	return r.Client.Status().Update(ctx, gateway)
}

func addCondition(status *gomaprojv1beta1.GatewayStatus, condType string, statusType metav1.ConditionStatus, reason, message string) {
	for i, existingCondition := range status.Conditions {
		if existingCondition.Type == condType {
			// Condition already exists, update it
			status.Conditions[i].Status = statusType
			status.Conditions[i].Reason = reason
			status.Conditions[i].Message = message
			status.Conditions[i].LastTransitionTime = metav1.Now()
			return
		}
	}

	// The Condition does not exist, add it
	condition := metav1.Condition{
		Type:               condType,
		Status:             statusType,
		Reason:             reason,
		Message:            message,
		LastTransitionTime: metav1.Now(),
	}
	status.Conditions = append(status.Conditions, condition)
}

func (r *GatewayReconciler) finalize(ctx context.Context, gateway *gomaprojv1beta1.Gateway) error {
	logger := log.FromContext(ctx)
	logger.Info("Finalizing Gateway", "Name", gateway.Name, "Namespace", gateway.Namespace)
	// Delete the ConfigMap
	configMap := &corev1.ConfigMap{}
	err := r.Get(ctx, client.ObjectKey{Namespace: gateway.Namespace, Name: gateway.Name}, configMap)
	if err != nil {
		logger.Error(err, "Failed to get Deployment")
		return err
	}
	logger.Info("Deleting ConfigMap...", "Name", configMap.Name)
	err = r.Delete(ctx, configMap)
	if err != nil {
		logger.Error(err, "Failed to delete Deployment")
		return err
	}

	// Delete the Deployment
	deployment := &v1.Deployment{}
	err = r.Get(ctx, client.ObjectKey{Namespace: gateway.Namespace, Name: gateway.Name}, deployment)
	if err != nil {
		logger.Error(err, "Failed to get Deployment")
		return err
	}
	logger.Info("Deleting Deployment...", "Name", deployment.Name)
	err = r.Delete(ctx, deployment)
	if err != nil {
		logger.Error(err, "Failed to delete Deployment")
		return err
	}

	if gateway.Spec.AutoScaling.Enabled {
		// Delete the HorizontalPodAutoscaler
		hpa := &autoscalingv1.HorizontalPodAutoscaler{}
		err = r.Get(ctx, client.ObjectKey{Namespace: gateway.Namespace, Name: gateway.Name}, hpa)
		if err != nil {
			logger.Error(err, "Failed to get HorizontalPodAutoscaler")
			return err
		}
		logger.Info("Deleting HorizontalPodAutoscaler...", "Name", hpa.Name)
		err = r.Delete(ctx, hpa)
		if err != nil {
			logger.Error(err, "Failed to delete HorizontalPodAutoscaler")
			return err
		}

	}
	logger.Info("Deleted Deployment", "Name", deployment.Name, "Namespace", deployment.Namespace)

	// Delete the Service
	service := &corev1.Service{}
	err = r.Get(ctx, client.ObjectKey{Namespace: gateway.Namespace, Name: gateway.Name}, service)
	if err != nil {
		logger.Error(err, "Failed to get Service")
		return err
	}
	logger.Info("Deleting Service...", "Name", service.Name)
	err = r.Delete(ctx, service)
	if err != nil {
		logger.Error(err, "Failed to delete Service")
		return err
	}

	logger.Info("Deleted Service", "Name", service.Name, "Namespace", service.Namespace)
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	pred := predicate.GenerationChangedPredicate{}
	return ctrl.NewControllerManagedBy(mgr).
		For(&gomaprojv1beta1.Gateway{}).
		WithEventFilter(pred).
		Owns(&corev1.ConfigMap{}).                      // Watch ConfigMaps created by the controller
		Owns(&v1.Deployment{}).                         // Watch Deployments created by the controller
		Owns(&corev1.Service{}).                        // Watch Services created by the controller
		Owns(&autoscalingv1.HorizontalPodAutoscaler{}). // Watch HorizontalPodAutoscaler created by the controller
		Named("gateway").
		Complete(r)
}
