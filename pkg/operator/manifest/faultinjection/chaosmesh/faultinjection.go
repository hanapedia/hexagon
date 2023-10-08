package manifest

import (
	"fmt"

	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/object/crd"
	"github.com/hanapedia/the-bench/pkg/operator/object/factory"
)

// CreateNetworkDelay creates network delay custom resource
func CreateNetworkDelay(name string) *crd.NetworkChaos {
	networkChaosArgs := factory.NetworkChaosArgs{
		Name:            fmt.Sprintf("%s-network-delay", name),
		Namespace:       model.CHAOSMESH_NAMESPACE,
		TargetNamespace: model.NAMESPACE,
		Selector:        map[string]string{"app": name},
		Duration:        model.CHAOSMESH_DURATION,
		Latency:         model.CHAOSMESH_LATENCY,
		Jitter:          model.CHAOSMESH_LATENCY_JITTER,
	}
	networkChaos := factory.NewNetworkChaos(&networkChaosArgs)

	return &networkChaos
}
