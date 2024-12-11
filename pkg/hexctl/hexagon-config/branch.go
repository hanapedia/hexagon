package hexagonconfig

import (
	"fmt"
	"time"

	v1 "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
)

/* adapters: */
/* - server: */
/*     action: read */
/*     variant: rest */
/*     route: get */
/*   tasks: */
/*   - adapter: */
/*       stressor: */
/*         name: cpu-stressor */
/*         variant: cpu */
/*         iters: 2 */
/*         threads: 1 */
/*   - adapter: */
/*       invocation: */
/*         variant: rest */
/*         service: NEXT BRANCH OR LEAF */
/*         action: get */
/*         route: get */
/*     resiliency: */
/*       isCritical: true */
/*       callTimeout: {{ .Timeout }} */
/*       adaptiveCallTimeout: */
/*         rto: {{ .RTO }} */
/*         min: 1ms */
/*         max: {{ .MaxTimeout }} */
/*         slo: {{ .SLO }} */
/*         capacity: {{ .Capacity }} */
/*         interval: {{ .Interval }} */
/*         kMargin: {{ .KMargin }} */
/*       retry: */
/*         disabled: true */
/*       circuitBreaker: */
/*         disabled: true */
func NewTrunkOrBranchService(version string, tier, index uint64, isEdge bool, fanout uint64, timeout time.Duration, adaptiveTimeout v1.AdaptiveTimeoutSpec) v1.ServiceUnitConfig {
	this := "trunk"
	primaryRoute := DEFAULT_GW_ROUTE
	if tier > 0 {
		this = "branch"
		primaryRoute = DEFAULT_GET_ROUTE
	}
	next := "branch"
	if isEdge {
		next = "leaf"
	}

	// build fanout adapters
	taskSpecs := []*v1.TaskSpec{}
	taskSpecs = append(taskSpecs, &v1.TaskSpec{
		AdapterConfig: &v1.SecondaryAdapterConfig{
			StressorConfig: &v1.StressorConfig{
				Name:        "cpu-stressor",
				Variant:     constants.CPU,
				Iterations:  2,
				ThreadCount: 1,
			},
		},
		Resiliency: v1.ResiliencySpec{
			Retry:         v1.RetrySpec{Disabled: true},
			CircutBreaker: v1.CircuitBreakerSpec{Disabled: true},
		},
	})
	for nextIndex := range fanout {
		taskSpecs = append(taskSpecs,
			&v1.TaskSpec{
				AdapterConfig: &v1.SecondaryAdapterConfig{
					InvocationConfig: &v1.InvocationConfig{
						Service: fmt.Sprintf("service-%s-t%v-%v", next, tier+1, nextIndex),
						Variant: constants.REST,
						Action:  constants.GET,
						Route:   DEFAULT_GET_ROUTE,
					},
				},
				Resiliency: v1.ResiliencySpec{
					IsCritical:          true,
					CallTimeout:         timeout.String(),
					AdaptiveCallTimeout: adaptiveTimeout,
					Retry:               v1.RetrySpec{Disabled: true},
					CircutBreaker:       v1.CircuitBreakerSpec{Disabled: true},
				},
			},
		)
	}

	return v1.ServiceUnitConfig{
		Version:        version,
		Name:           fmt.Sprintf("service-%s-t%v-%v", this, tier, index),
		DeploymentSpec: NewDefaultDeploymentSpec(),
		AdapterConfigs: []*v1.PrimaryAdapterSpec{
			{
				ServerConfig: &v1.ServerConfig{
					Variant: constants.REST,
					Action:  constants.GET,
					Route:   primaryRoute,
				},
				TaskSpecs: taskSpecs,
			},
		},
	}
}
