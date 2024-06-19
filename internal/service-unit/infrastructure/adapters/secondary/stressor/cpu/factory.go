package cpu

import (
	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func CpuStressorAdapterFactory(adapterConfig *model.StressorConfig) (ports.SecodaryPort, error) {
	var cpuStressor ports.SecodaryPort
	var err error

	iters := adapterConfig.Iterations
	if adapterConfig.Iterations <= 0 {
		iters = 1
	}

	threadCount := adapterConfig.ThreadCount
	if adapterConfig.ThreadCount <= 0 {
		threadCount = 1
	}

	cpuStressor = &cpuStressorAdapter{
		payloadSize: model.GetPayloadSize(adapterConfig.Payload),
		iterations:  iters,
		threadCount: threadCount,
	}

	// set destionation id
	cpuStressor.SetDestId(adapterConfig.GetId())

	logger.Logger.Debugf("Initialized cpu stressor adapter: %s", adapterConfig.GetId())
	return cpuStressor, err
}
