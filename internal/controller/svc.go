package controller

import (
	"context"
	"strings"

	gomaprojv1beta1 "github.com/jkaninda/goma-operator/api/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// createService create K8s service
func createService(r GatewayReconciler, ctx context.Context, req ctrl.Request, gateway *gomaprojv1beta1.Gateway) error {
	l := log.FromContext(ctx)
	k8sService := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: req.Namespace,
			Labels: map[string]string{
				"app":        req.Name,
				"belongs-to": BelongsTo,
				"managed-by": req.Name,
			},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app": req.Name,
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       8080,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8080},
				}, {
					Name:       "https",
					Protocol:   corev1.ProtocolTCP,
					Port:       8443,
					TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: 8443},
				},
			},
		},
	}

	// Set Gateway instance as the owner and controller
	if err := controllerutil.SetControllerReference(gateway, k8sService, r.Scheme); err != nil {
		return err
	}

	found := &corev1.Service{}
	err := r.Get(ctx, client.ObjectKey{Namespace: req.Namespace, Name: req.Name}, found)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			l.Info("Creating a new Service", "Service.Namespace", k8sService.Namespace, "Service.Name", k8sService.Name)
			err = r.Create(ctx, k8sService)
			if err != nil {
				return err
			}
			return nil
		}
		l.Info("Failed to get Service", "Service.Namespace", k8sService.Namespace, "Service.Name", k8sService.Name)
		return err
	}

	return nil
}
