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
			1,
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

// GenerateFanout generates ServiceUnitConfigs for fanout topology of given tiers.
// []CommonConfig should contain configs for each tier. Note that it is not for each service.
// so leaves will have same config for example.
func GenerateFanout(commonConfigs []CommonConfig, fanoutDegree uint64) []v1.ServiceUnitConfig {
	configs := []v1.ServiceUnitConfig{}
	for tier, commonConfig := range commonConfigs[:len(commonConfigs)-1] {
		var fanout uint64 = 1
		isEdge := tier == len(commonConfigs)-2
		if isEdge {
			fanout = fanoutDegree
		}
		configs = append(configs, NewTrunkOrBranchService(
			commonConfig.Version,
			uint64(tier),
			0,
			isEdge,
			fanout,
			commonConfig.Timeout,
			commonConfig.AdaptiveTimeout,
		))
	}
	for index := range fanoutDegree {
		configs = append(configs, NewLeafService(
			commonConfigs[len(commonConfigs)-1].Version,
			uint64(len(commonConfigs)-1),
			index,
		))
	}
	return configs
}

// GenerateFunnel generates ServiceUnitConfigs for funnel topology of given degree.
// []CommonConfig should contain configs for each tier. Note that it is not for each service.
// so leaves will have same config for example.
func GenerateFunnel(trunkConfig, leafConfig CommonConfig, funnelDegrees uint64) []v1.ServiceUnitConfig {
	configs := []v1.ServiceUnitConfig{}
	for index := range funnelDegrees {
		configs = append(configs, NewTrunkOrBranchService(
			trunkConfig.Version,
			0,
			index,
			true,
			1,
			trunkConfig.Timeout,
			trunkConfig.AdaptiveTimeout,
		))
	}
	// add a leaf
	configs = append(configs, NewLeafService(
		leafConfig.Version,
		1,
		0,
	))
	return configs
}
