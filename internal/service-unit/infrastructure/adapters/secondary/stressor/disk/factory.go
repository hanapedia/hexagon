package disk

import (
	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func DiskStressorAdapterFactory(adapterConfig *model.StressorConfig, client secondary.SecondaryAdapterClient) (secondary.SecodaryPort, error) {
	var diskStressor secondary.SecodaryPort
	var err error

	iters := adapterConfig.Iterations
	if adapterConfig.Iterations <= 0 {
		iters = 1
	}

	threadCount := adapterConfig.ThreadCount
	if adapterConfig.ThreadCount <= 0 {
		threadCount = 1
	}

	diskStressor = &diskStressorAdapter{
		payloadSize: model.GetPayloadSize(adapterConfig.Payload),
		iterations:  iters,
		threadCount: threadCount,

	}

	// set destionation id
	diskStressor.SetDestId(adapterConfig.GetId())

	logger.Logger.Debugf("Initialized cpu stressor adapter: %s", adapterConfig.GetId())
	return diskStressor, err
}
