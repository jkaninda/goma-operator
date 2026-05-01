package resources

const (
	// LabelApp is the standard app label.
	LabelApp = "app"
	// LabelManagedBy identifies the operator.
	LabelManagedBy = "app.kubernetes.io/managed-by"
	// LabelPartOf identifies the system this belongs to.
	LabelPartOf = "app.kubernetes.io/part-of"
)

// CommonLabels returns the standard set of labels for all resources.
func CommonLabels(gatewayName string) map[string]string {
	return map[string]string{
		LabelApp:       gatewayName,
		LabelManagedBy: "goma-operator",
		LabelPartOf:    "goma-gateway",
	}
}

// SelectorLabels returns the minimal labels used for pod selection.
func SelectorLabels(gatewayName string) map[string]string {
	return map[string]string{
		LabelApp: gatewayName,
	}
}
