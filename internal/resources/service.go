package resources

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	gatewayv1alpha1 "github.com/jkaninda/goma-operator/api/v1alpha1"
)

// BuildService creates a Service for the gateway.
func BuildService(gw *gatewayv1alpha1.Gateway) *corev1.Service {
	spec := gw.Spec.Service

	ports := []corev1.ServicePort{
		{
			Name:       "http",
			Port:       gw.Spec.HTTPPort(),
			TargetPort: intstr.FromString("http"),
			Protocol:   corev1.ProtocolTCP,
		},
		{
			Name:       "https",
			Port:       gw.Spec.HTTPSPort(),
			TargetPort: intstr.FromString("https"),
			Protocol:   corev1.ProtocolTCP,
		},
	}

	// NodePort (only applied when type=NodePort)
	if gw.Spec.ServiceType() == corev1.ServiceTypeNodePort && spec != nil {
		if spec.HTTPNodePort > 0 {
			ports[0].NodePort = spec.HTTPNodePort
		}
		if spec.HTTPSNodePort > 0 {
			ports[1].NodePort = spec.HTTPSNodePort
		}
	}

	// Merge operator labels with any user-supplied labels.
	labels := CommonLabels(gw.Name)
	if spec != nil {
		for k, v := range spec.Labels {
			labels[k] = v
		}
	}

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gw.Name,
			Namespace: gw.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Type:     gw.Spec.ServiceType(),
			Selector: SelectorLabels(gw.Name),
			Ports:    ports,
		},
	}

	if spec != nil {
		if len(spec.Annotations) > 0 {
			svc.Annotations = spec.Annotations
		}
		svc.Spec.LoadBalancerIP = spec.LoadBalancerIP
		svc.Spec.LoadBalancerSourceRanges = spec.LoadBalancerSourceRanges
		svc.Spec.LoadBalancerClass = spec.LoadBalancerClass
		svc.Spec.ExternalTrafficPolicy = spec.ExternalTrafficPolicy
		svc.Spec.SessionAffinity = spec.SessionAffinity
		svc.Spec.IPFamilyPolicy = spec.IPFamilyPolicy
		svc.Spec.IPFamilies = spec.IPFamilies
	}

	return svc
}
