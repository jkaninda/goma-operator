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

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gomaprojv1beta1 "github.com/jkaninda/goma-operator/api/v1beta1"
)

var _ = Describe("Route Controller", func() {
	Context("When reconciling a resource", func() {
		const resourceName = "test-resource"

		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      resourceName,
			Namespace: "default", // TODO(user):Modify as needed
		}
		gateway := &gomaprojv1beta1.Gateway{}

		route := &gomaprojv1beta1.Route{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind Gateway")
			err := k8sClient.Get(ctx, typeNamespacedName, gateway)
			if err != nil && errors.IsNotFound(err) {
				resource := &gomaprojv1beta1.Gateway{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: gomaprojv1beta1.GatewaySpec{
						GatewayVersion: "latest",
						Server:         gomaprojv1beta1.Server{},
						ReplicaCount:   1,
						AutoScaling: gomaprojv1beta1.AutoScaling{
							Enabled:                        false,
							MinReplicas:                    2,
							MaxReplicas:                    5,
							TargetCPUUtilizationPercentage: 80,
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
			By("creating the custom resource for the Kind Route")
			err = k8sClient.Get(ctx, typeNamespacedName, route)
			if err != nil && errors.IsNotFound(err) {
				resource := &gomaprojv1beta1.Route{
					ObjectMeta: metav1.ObjectMeta{
						Name:      resourceName,
						Namespace: "default",
					},
					Spec: gomaprojv1beta1.RouteSpec{
						Gateway: resourceName,
						Routes: []gomaprojv1beta1.RouteConfig{
							{
								Path:        "/",
								Name:        resourceName,
								Rewrite:     "/",
								Destination: "https://example.com",
								Methods:     []string{"GET", "POST"},
							},
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			// TODO(user): Cleanup logic after each test, like removing the resource instance.
			resource := &gomaprojv1beta1.Route{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance Route")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &RouteReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			// TODO(user): Add more specific assertions depending on your controller's reconciliation logic.
			// Example: If you expect a certain status condition after reconciliation, verify it here.
		})
	})
})
