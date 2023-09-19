package manifest

import (
	"fmt"

	"github.com/hanapedia/the-bench/pkg/operator/manifest/defaults"
	"github.com/hanapedia/the-bench/pkg/operator/object/crd"
	"github.com/hanapedia/the-bench/pkg/operator/object/factory"
)

// CreateNetworkDelay creates network delay custom resource
func CreateNetworkDelay(name string) *crd.NetworkChaos {
	networkChaosArgs := factory.NetworkChaosArgs{
		Name:            fmt.Sprintf("%s-network-delay", name),
		Namespace:       defaults.CHAOSMESH_NAMESPACE,
		TargetNamespace: defaults.NAMESPACE,
		Selector:        map[string]string{"app": name},
		Duration:        defaults.CHAOSMESH_DURATION,
		Latency:         defaults.CHAOSMESH_LATENCY,
		Jitter:          defaults.CHAOSMESH_LATENCY_JITTER,
	}
	networkChaos := factory.NewNetworkChaos(&networkChaosArgs)

	return &networkChaos
}
