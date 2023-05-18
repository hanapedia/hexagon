package v1

type IngressAdapterStatusTypes string

const (
	// IngressAdapterStatusReady indicates that the all of the egress adapters are ready.
	IngressAdapterStatusReady IngressAdapterStatusTypes = "Ready"

	// IngressAdapterStatusPending indicates that the at least one egress adatper is not ready.
	IngressAdapterStatusPending IngressAdapterStatusTypes = "Pending"

	// IngressAdapterStatusCircularDependencyFound indicates that there is circular dependency in the ingress adapter call graph.
	IngressAdapterStatusCircularDependencyFound IngressAdapterStatusTypes = "CircularDependencyFound"
)
