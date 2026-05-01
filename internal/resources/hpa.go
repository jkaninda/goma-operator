package resources

import (
	autoscalingv2 "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gatewayv1alpha1 "github.com/jkaninda/goma-operator/api/v1alpha1"
)

// BuildHPA creates a HorizontalPodAutoscaler for the gateway.
func BuildHPA(gw *gatewayv1alpha1.Gateway) *autoscalingv2.HorizontalPodAutoscaler {
	as := gw.Spec.AutoScaling

	hpa := &autoscalingv2.HorizontalPodAutoscaler{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gw.Name,
			Namespace: gw.Namespace,
			Labels:    CommonLabels(gw.Name),
		},
		Spec: autoscalingv2.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscalingv2.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       gw.Name,
			},
			MaxReplicas: as.MaxReplicas,
		},
	}

	if as.MinReplicas != nil {
		hpa.Spec.MinReplicas = as.MinReplicas
	}

	var metrics []autoscalingv2.MetricSpec

	if as.TargetCPUUtilizationPercentage != nil {
		metrics = append(metrics, autoscalingv2.MetricSpec{
			Type: autoscalingv2.ResourceMetricSourceType,
			Resource: &autoscalingv2.ResourceMetricSource{
				Name: "cpu",
				Target: autoscalingv2.MetricTarget{
					Type:               autoscalingv2.UtilizationMetricType,
					AverageUtilization: as.TargetCPUUtilizationPercentage,
				},
			},
		})
	}

	if as.TargetMemoryUtilizationPercentage != nil {
		metrics = append(metrics, autoscalingv2.MetricSpec{
			Type: autoscalingv2.ResourceMetricSourceType,
			Resource: &autoscalingv2.ResourceMetricSource{
				Name: "memory",
				Target: autoscalingv2.MetricTarget{
					Type:               autoscalingv2.UtilizationMetricType,
					AverageUtilization: as.TargetMemoryUtilizationPercentage,
				},
			},
		})
	}

	if len(metrics) > 0 {
		hpa.Spec.Metrics = metrics
	}

	return hpa
}
