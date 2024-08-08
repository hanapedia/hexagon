package v1

type ConfigKind string

var (
	ServiceUnitKind   ConfigKind = "ServiceUnit"
	ClusterConfigKind ConfigKind = "ClusterConfig"
)

type ConfigTemplate struct {
	ApiVersion string     `json:"apiVersion,omitempty"`
	Kind       ConfigKind `json:"kind,omitempty"`
}
