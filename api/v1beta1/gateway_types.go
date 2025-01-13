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

package v1beta1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GatewaySpec defines the desired configuration and state of the Gateway deployment.
type GatewaySpec struct {
	// ImageName specifies the image of the Goma Gateway to use.
	// This value is used to pull the desired image for the Goma Gateway service.
	//Default: jkaninda/goma-gateway:latest
	ImageName string `json:"imageName,omitempty"` // The version tag of the Goma Gateway image
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	ImagePullSecrets []v1.LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,15,rep,name=imagePullSecrets"`
	// Deprecated
	// Please use ImageName instead
	GatewayVersion string `json:"gatewayVersion,omitempty"` // The version tag of the Goma Gateway image
	// Server contains the configuration for the Gateway server.
	// It includes settings related to the server such as port, protocol, and other gateway-related configurations.
	Server Server `json:"server,omitempty"` // Gateway server configuration

	// ReplicaCount defines the number of replicas for the Gateway deployment.
	// This field determines how many instances of the Gateway service will run in the cluster.
	// A higher count provides better availability and fault tolerance.
	ReplicaCount int32 `json:"replicaCount,omitempty"` // Number of replicas for the Gateway deployment

	// AutoScaling defines the settings for enabling or disabling pod autoscaling.
	// When enabled, it allows Kubernetes to automatically adjust the number of Gateway pods based on usage metrics (e.g., CPU or memory).
	AutoScaling AutoScaling `json:"autoScaling,omitempty"` // Auto-scaling configuration for Gateway pods

	// Resources defines the resource requests and limits for the Gateway deployment.
	// It specifies the amount of CPU and memory the Gateway pods should reserve, as well as any limits.
	// This helps ensure proper resource allocation and avoids resource contention.
	Resources v1.ResourceRequirements `json:"resources,omitempty"` // Resource requirements for Gateway pods

	// Affinity defines node affinity rules for pod scheduling.
	// It allows you to specify rules for how pods should be placed on specific nodes, based on labels, zones, or other factors.
	Affinity *v1.Affinity `json:"affinity,omitempty"` // Affinity rules for pod scheduling
}

// GatewayStatus defines the observed state of Gateway.
type GatewayStatus struct {
	Replicas   int32              `json:"replicas,omitempty"`
	Routes     int32              `json:"routes,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Gateway is the Schema for the gateways API.
type Gateway struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GatewaySpec   `json:"spec,omitempty"`
	Status GatewayStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// GatewayList contains a list of Gateway.
type GatewayList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gateway `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Gateway{}, &GatewayList{})
}
