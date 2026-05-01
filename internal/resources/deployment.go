package resources

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	gatewayv1alpha1 "github.com/jkaninda/goma-operator/api/v1alpha1"
)

const (
	configMountPath     = "/etc/goma"
	certsMountPath      = "/etc/goma/certs"
	routeCertsMountPath = "/etc/goma/route-certs"
	k8sProviderPath     = "/etc/goma/providers/k8s"
	acmeMountPath       = "/etc/letsencrypt"
)

// BuildDeployment creates a Deployment for the gateway.
// routes is the set of Route CRs attached to this gateway — used to mount
// per-route TLS Secrets as volumes so the gateway can serve their certificates.
func BuildDeployment(gw *gatewayv1alpha1.Gateway, routes []gatewayv1alpha1.Route) *appsv1.Deployment {
	replicas := int32(1)
	if gw.Spec.Replicas != nil {
		replicas = *gw.Spec.Replicas
	}

	labels := CommonLabels(gw.Name)
	podLabels := CommonLabels(gw.Name)

	// Volumes
	volumes := []corev1.Volume{
		{
			Name: "config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: ConfigMapName(gw.Name),
					},
				},
			},
		},
	}

	// Volume mounts for gateway container
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "config",
			MountPath: configMountPath + "/" + ConfigFileName,
			SubPath:   ConfigFileName,
		},
	}

	// TLS secret volumes
	for _, tls := range gw.Spec.Server.TLS {
		volName := "tls-" + tls.SecretName
		volumes = append(volumes, corev1.Volume{
			Name: volName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: tls.SecretName,
				},
			},
		})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      volName,
			MountPath: fmt.Sprintf("%s/%s", certsMountPath, tls.SecretName),
			ReadOnly:  true,
		})
	}

	// Per-route serving TLS Secrets — each route with spec.tls.secretName
	// gets its own volume + mount so the gateway can present per-host certs.
	// Dedup by Secret name in case multiple routes share the same cert.
	mountedRouteTLS := make(map[string]bool)
	for _, route := range routes {
		if route.Spec.TLS == nil || route.Spec.TLS.SecretName == "" {
			continue
		}
		sn := route.Spec.TLS.SecretName
		if mountedRouteTLS[sn] {
			continue
		}
		mountedRouteTLS[sn] = true
		volName := "route-tls-" + sn
		volumes = append(volumes, corev1.Volume{
			Name: volName,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: sn,
				},
			},
		})
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      volName,
			MountPath: fmt.Sprintf("%s/%s", certsMountPath, sn),
			ReadOnly:  true,
		})
	}

	containers := []corev1.Container{
		{
			Name:    "goma-gateway",
			Image:   gw.Spec.Image,
			Command: []string{"/usr/local/bin/goma"},
			Args:    []string{"server"},
			Ports: []corev1.ContainerPort{
				{Name: "http", ContainerPort: 8080, Protocol: corev1.ProtocolTCP},
				{Name: "https", ContainerPort: 8443, Protocol: corev1.ProtocolTCP},
			},
			StartupProbe: &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/readyz",
						Port: intstr.FromInt32(8080),
					},
				},
				PeriodSeconds:    2,
				FailureThreshold: 30, // allow up to 60s for cold start
			},
			ReadinessProbe: &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/readyz",
						Port: intstr.FromInt32(8080),
					},
				},
				PeriodSeconds:    5,
				TimeoutSeconds:   2,
				SuccessThreshold: 1,
				FailureThreshold: 3,
			},
			LivenessProbe: &corev1.Probe{
				ProbeHandler: corev1.ProbeHandler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: "/healthz",
						Port: intstr.FromInt32(8080),
					},
				},
				PeriodSeconds:    10,
				TimeoutSeconds:   3,
				FailureThreshold: 3,
			},
			Resources:    gw.Spec.Resources,
			VolumeMounts: volumeMounts,
		},
	}

	// Add k8s-provider sidecar (enabled by default — opt out via
	// spec.providers.kubernetes.enabled=false).
	if gw.Spec.KubernetesProviderEnabled() {
		k8sVolName := "k8s-provider"
		acmeVolName := "acme-data"
		routeCertsVolName := "route-certs"

		volumes = append(volumes,
			corev1.Volume{
				Name:         k8sVolName,
				VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
			},
			corev1.Volume{
				Name:         acmeVolName,
				VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
			},
			corev1.Volume{
				Name:         routeCertsVolName,
				VolumeSource: corev1.VolumeSource{EmptyDir: &corev1.EmptyDirVolumeSource{}},
			},
		)

		// Add shared volume mounts to gateway container
		containers[0].VolumeMounts = append(containers[0].VolumeMounts,
			corev1.VolumeMount{Name: k8sVolName, MountPath: k8sProviderPath},
			corev1.VolumeMount{Name: acmeVolName, MountPath: acmeMountPath},
			corev1.VolumeMount{Name: routeCertsVolName, MountPath: routeCertsMountPath},
		)

		sidecarImage := gw.Spec.KubernetesProviderImage()

		containers = append(containers, corev1.Container{
			Name:  "k8s-provider",
			Image: sidecarImage,
			Env: []corev1.EnvVar{
				{Name: "GOMA_K8S_GATEWAY", Value: gw.Name},
				{
					Name: "GOMA_K8S_NAMESPACE",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{FieldPath: "metadata.namespace"},
					},
				},
				{Name: "GOMA_K8S_OUTPUT_DIR", Value: k8sProviderPath},
			},
			VolumeMounts: []corev1.VolumeMount{
				{Name: k8sVolName, MountPath: k8sProviderPath},
				{Name: acmeVolName, MountPath: acmeMountPath},
				{Name: routeCertsVolName, MountPath: routeCertsMountPath},
			},
		})
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gw.Name,
			Namespace: gw.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{MatchLabels: SelectorLabels(gw.Name)},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: podLabels},
				Spec: corev1.PodSpec{
					ServiceAccountName: ServiceAccountName(gw.Name),
					Containers:         containers,
					Volumes:            volumes,
					ImagePullSecrets:   gw.Spec.ImagePullSecrets,
					Affinity:           gw.Spec.Affinity,
				},
			},
		},
	}

	return dep
}
