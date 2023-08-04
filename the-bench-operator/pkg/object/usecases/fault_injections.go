package usecases

import (
	"fmt"

	"github.com/hanapedia/the-bench/the-bench-operator/pkg/object/crd"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/object/factory"
)

// CreateNetworkDelay creates network delay custom resource
func CreateNetworkDelay(name string) *crd.NetworkChaos {
	networkChaosArgs := factory.NetworkChaosArgs{
		Name:            fmt.Sprintf("%s-network-delay", name),
		Namespace:       factory.CHAOSMESH_NAMESPACE,
		TargetNamespace: factory.NAMESPACE,
		Selector:        map[string]string{"app": name},
		Duration:        factory.CHAOSMESH_DURATION,
		Latency:         factory.CHAOSMESH_LATENCY,
		Jitter:          factory.CHAOSMESH_LATENCY_JITTER,
	}
	networkChaos := factory.NetworkChaosFactory(&networkChaosArgs)

	return &networkChaos
}
