package resources

import (
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gatewayv1alpha1 "github.com/jkaninda/goma-operator/api/v1alpha1"
)

// ServiceAccountName returns the ServiceAccount name for a gateway.
func ServiceAccountName(gatewayName string) string {
	return gatewayName + "-gateway"
}

// RoleName returns the Role name for a gateway.
func RoleName(gatewayName string) string {
	return gatewayName + "-gateway"
}

// BuildServiceAccount creates a ServiceAccount for the gateway pod.
// The sidecar (goma-k8s-provider) uses this SA to watch Route/Middleware CRDs.
func BuildServiceAccount(gw *gatewayv1alpha1.Gateway) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ServiceAccountName(gw.Name),
			Namespace: gw.Namespace,
			Labels:    CommonLabels(gw.Name),
		},
	}
}

// BuildRole creates a Role granting the sidecar read access to Route/Middleware
// CRDs in its own namespace.
func BuildRole(gw *gatewayv1alpha1.Gateway) *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RoleName(gw.Name),
			Namespace: gw.Namespace,
			Labels:    CommonLabels(gw.Name),
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"gateway.jkaninda.dev"},
				Resources: []string{"routes", "middlewares"},
				Verbs:     []string{"get", "list", "watch"},
			},
			// For ACME certificate persistence (goma-k8s-provider ACME sync).
			{
				APIGroups: []string{""},
				Resources: []string{"secrets"},
				Verbs:     []string{"get", "list", "watch", "create", "update", "patch"},
			},
		},
	}
}

// BuildRoleBinding creates a RoleBinding that binds the gateway's
// ServiceAccount to the Role.
func BuildRoleBinding(gw *gatewayv1alpha1.Gateway) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RoleName(gw.Name),
			Namespace: gw.Namespace,
			Labels:    CommonLabels(gw.Name),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      ServiceAccountName(gw.Name),
				Namespace: gw.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: rbacv1.GroupName,
			Kind:     "Role",
			Name:     RoleName(gw.Name),
		},
	}
}
