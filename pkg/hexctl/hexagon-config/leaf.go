package hexagonconfig

import (
	"fmt"

	v1 "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
)

/* - server: */
/*     action: get */
/*     variant: rest */
/*     route: get */
/*   tasks: */
/*   - adapter: */
/*       stressor: */
/*         name: cpu-stressor */
/*         variant: cpu */
/*         iters: 500 */
/*         threads: 2 */
/*     resiliency: */
/*       # callTimeout: "2.5ms" # no timeout */
/*       # taskTimeout: "1s" # no timeout */
/*       retry: */
/*         disabled: true */
/*       circuitBreaker: */
/*         disabled: true */
func NewLeafService(version string, tier, index uint64) v1.ServiceUnitConfig {
	return v1.ServiceUnitConfig{
		Version:        version,
		Name:           fmt.Sprintf("service-leaf-t%v-%v", tier, index),
		DeploymentSpec: NewDefaultDeploymentSpec(),
		AdapterConfigs: []*v1.PrimaryAdapterSpec{
			{
				ServerConfig: &v1.ServerConfig{
					Variant: constants.REST,
					Action:  constants.GET,
					Route:   "get",
				},
				TaskSpecs: []*v1.TaskSpec{
					{
						AdapterConfig: &v1.SecondaryAdapterConfig{
							StressorConfig: &v1.StressorConfig{
								Name:        "cpu-stressor",
								Variant:     constants.CPU,
								Iterations:  500,
								ThreadCount: 2,
							},
						},
						Resiliency: v1.ResiliencySpec{
							Retry:         v1.RetrySpec{Enabled: true},
							CircutBreaker: v1.CircuitBreakerSpec{Enabled: true},
						},
					},
				},
			},
		},
	}
}
