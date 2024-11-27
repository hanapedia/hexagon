// hexagonconfig contains code for generating hexagon configs with common compositions such as chain, fanout, and funnel
package hexagonconfig

import (
	"time"

	v1 "github.com/hanapedia/hexagon/pkg/api/v1"
)

type CommonConfig struct {
	Version         string
	Timeout         time.Duration
	AdaptiveTimeout v1.AdaptiveTimeoutSpec
}

// GenerateChain generates ServiceUnitConfigs for chain topology of given tiers.
func GenerateChain(commonConfigs []CommonConfig) []v1.ServiceUnitConfig {
	configs := []v1.ServiceUnitConfig{}
	for tier, commonConfig := range commonConfigs[:len(commonConfigs)-1] {
		configs = append(configs, NewTrunkOrBranchService(
			commonConfig.Version,
			uint64(tier),
			0,
			tier == len(commonConfigs)-2,
			commonConfig.Timeout,
			commonConfig.AdaptiveTimeout,
		))
	}
	configs = append(configs, NewLeafService(
		commonConfigs[len(commonConfigs)-1].Version,
		uint64(len(commonConfigs)-1),
		0,
	))
	return configs
}
