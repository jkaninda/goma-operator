package controller

import (
	"context"
	"fmt"
	"time"

	gomaprojv1beta1 "github.com/jkaninda/goma-operator/api/v1beta1"
	v1 "k8s.io/api/apps/v1"
	av1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// createUpdateDeployment creates Kubernetes deployment
func createUpdateDeployment(r GatewayReconciler, ctx context.Context, req ctrl.Request, gateway gomaprojv1beta1.Gateway, imageName string) error {
	logger := log.FromContext(ctx)
	var volumes []corev1.Volume
	var volumeMounts []corev1.VolumeMount

	volumes = append(volumes, corev1.Volume{
		Name: "config",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: req.Name,
				},
			},
		},
	})
	volumeMounts = append(volumeMounts, corev1.VolumeMount{
		Name:      "config",
		MountPath: ConfigPath,
		ReadOnly:  true,
	})
	if len(gateway.Spec.Server.TlsSecretName) != 0 {
		volumes = append(volumes, corev1.Volume{
			Name: req.Name,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: gateway.Spec.Server.TlsSecretName,
				},
			},
		})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      req.Name,
			ReadOnly:  true,
			MountPath: CertsPath,
		})

	}
	// check if ReplicaCount is defined
	if gateway.Spec.ReplicaCount != 0 {
		ReplicaCount = gateway.Spec.ReplicaCount
	}
	// Define the desired Deployment
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels:    gateway.Labels,
		},
		Spec: v1.DeploymentSpec{
			Replicas: int32Ptr(ReplicaCount), // Set desired replicas
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app":        req.Name,
					"belongs-to": BelongsTo,
					"managed-by": gateway.Name,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app":        req.Name,
						"belongs-to": BelongsTo,
						"managed-by": gateway.Name,
					},
					Annotations: map[string]string{
						"updated-at": time.Now().Format(time.RFC3339),
					},
				},
				Spec: corev1.PodSpec{
					Affinity: gateway.Spec.Affinity,
					Containers: []corev1.Container{
						{
							Name:            "gateway",
							Image:           imageName,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Command:         []string{"/usr/local/bin/goma", "server"},
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
								},
							},
							ReadinessProbe: &corev1.Probe{
								InitialDelaySeconds: 15,
								PeriodSeconds:       10,
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/readyz",
										Port: intstr.FromInt32(8080),
									},
								},
							},
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: 30,
								PeriodSeconds:       30,
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Path: "/healthz",
										Port: intstr.FromInt32(8080),
									},
								},
							},
							Resources:    gateway.Spec.Resources,
							VolumeMounts: volumeMounts,
						},
					},
					Volumes: volumes,
				},
			},
		},
	}

	// Check if the Deployment already exists
	var existingDeployment v1.Deployment
	err := r.Get(ctx, types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, &existingDeployment)
	if err != nil && client.IgnoreNotFound(err) != nil {
		logger.Error(err, "Failed to get Deployment")
		return err
	}
	if err != nil && client.IgnoreNotFound(err) == nil {
		logger.Info("Creating a new Deployment")
		// Create the Deployment if it doesn't exist
		if err = controllerutil.SetControllerReference(&gateway, deployment, r.Scheme); err != nil {
			logger.Error(err, "Failed to set controller reference")
			return err
		}
		if err = r.Create(ctx, deployment); err != nil {
			logger.Error(err, "Failed to create Deployment")
			return err
		}
		logger.Info("Created Deployment", "Deployment.Name", deployment.Name)
	} else {
		logger.Info("Deployment already exists", "Deployment.Name", deployment.Name)
		// Update the Deployment if the spec has changed
		if !equalDeploymentSpec(existingDeployment.Spec, deployment.Spec, gateway.Spec.AutoScaling.Enabled) {
			logger.Info("Updating Deployment", "Deployment.Name", deployment.Name)
			existingDeployment.Spec = deployment.Spec
			if err = r.Update(ctx, &existingDeployment); err != nil {
				logger.Error(err, "Failed to update Deployment")
				return err
			}
			logger.Info("Updated Deployment", "Deployment.Name", deployment.Name)
		}
	}

	// check if hpa is enabled
	if gateway.Spec.AutoScaling.Enabled {
		err = createHpa(r, ctx, req, &gateway)
		if err != nil {
			logger.Error(err, "Failed to create HorizontalPodAutoscaler")
			return err
		}
	} else {
		// Check if the hpa already exists
		var existHpa av1.HorizontalPodAutoscaler
		err = r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, &existHpa)
		if err != nil && client.IgnoreNotFound(err) != nil {
			logger.Error(err, "Failed to get HorizontalPodAutoscaler")
			return err
		}
		if err == nil {
			// Delete the HorizontalPodAutoscaler
			if err = r.Delete(ctx, &existHpa); err != nil {
				logger.Error(err, "Failed to delete HorizontalPodAutoscaler")
				return err
			}
			logger.Info("Deleted HorizontalPodAutoscaler successfully", "HorizontalPodAutoscaler.Name", req.Name)
		}
	}

	return nil
}

// Helper function to compare Deployment specs
func equalDeploymentSpec(existing, desired v1.DeploymentSpec, autoScalingEnabled bool) bool {
	if existing.Template.Spec.Containers[0].Image != desired.Template.Spec.Containers[0].Image {
		return false
	}

	if !autoScalingEnabled {
		if existing.Replicas == nil || *existing.Replicas != *desired.Replicas {
			return false
		}
	}
	return true
}
func restartDeployment(r client.Client, ctx context.Context, req ctrl.Request, gateway *gomaprojv1beta1.Gateway) error {
	logger := log.FromContext(ctx)
	// Fetch the Deployment
	var deployment v1.Deployment
	if err := r.Get(ctx, types.NamespacedName{Name: gateway.Name, Namespace: req.Namespace}, &deployment); err != nil {
		logger.Error(err, "Failed to get Deployment", "name", gateway.Name, "namespace", req.Name)
		return client.IgnoreNotFound(err)
	}

	// Add or update an annotation to trigger a rolling update
	if deployment.Spec.Template.ObjectMeta.Annotations == nil {
		deployment.Spec.Template.ObjectMeta.Annotations = make(map[string]string)
	}

	deployment.Spec.Template.ObjectMeta.Annotations["restarted-at"] = time.Now().Format(time.RFC3339)
	deployment.Spec.Template.ObjectMeta.Annotations["updated-at"] = time.Now().Format(time.RFC3339)
	// Update the Deployment
	if err := r.Update(ctx, &deployment); err != nil {
		logger.Error(err, "Failed to update Deployment for restart", "name", gateway.Name, "namespace", req.Name)
		return err
	}

	logger.Info("Successfully restarted Deployment", "name", gateway.Name, "namespace", req.Name)
	return nil
}

// currentReplicas returns current replicas
func currentReplicas(ctx context.Context, c client.Client, hpaName, namespace string) (int32, error) {
	hpa := &av1.HorizontalPodAutoscaler{}
	// Retrieve the HPA resource
	err := c.Get(ctx, types.NamespacedName{Name: hpaName, Namespace: namespace}, hpa)
	if err != nil {
		return 0, fmt.Errorf("failed to get HPA: %w", err)
	}
	// Access the current replicas in the status field
	replicas := hpa.Status.CurrentReplicas
	return replicas, nil
}
