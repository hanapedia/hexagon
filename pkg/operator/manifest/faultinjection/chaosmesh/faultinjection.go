package manifest

import (
	"fmt"

	"github.com/hanapedia/hexagon/pkg/api/defaults"
	"github.com/hanapedia/hexagon/pkg/operator/object/crd"
	"github.com/hanapedia/hexagon/pkg/operator/object/factory"
)

// CreateNetworkDelay creates network delay custom resource
func CreateNetworkDelay(name string) *crd.NetworkChaos {
	networkChaosArgs := factory.NetworkChaosArgs{
		Name:            fmt.Sprintf("%s-network-delay", name),
		Namespace:       defaults.CHAOSMESH_NAMESPACE,
		TargetNamespace: defaults.NAMESPACE,
		Selector:        map[string]string{factory.AppLabel: name},
		Duration:        defaults.CHAOSMESH_DURATION,
		Latency:         defaults.CHAOSMESH_LATENCY,
		Jitter:          defaults.CHAOSMESH_LATENCY_JITTER,
	}
	networkChaos := factory.NewNetworkChaos(&networkChaosArgs)

	return &networkChaos
}
