package cpu

import (
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
)

func CpuStressorAdapterFactory(adapterConfig *model.StressorConfig) (ports.SecodaryPort, error) {
	var cpuStressor ports.SecodaryPort
	var err error

	duration, err := time.ParseDuration(adapterConfig.Duration)
	if err != nil {
		return cpuStressor, err
	}

	threadCount := adapterConfig.ThreadCount
	if adapterConfig.ThreadCount <= 0 {
		threadCount = 1
	}

	cpuStressor = &cpuStressorAdapter{
		payloadSize: model.GetPayloadSize(adapterConfig.Payload),
		duration:    duration,
		threadCount: threadCount,
	}

	// set destionation id
	cpuStressor.SetDestId(adapterConfig.GetId())

	return cpuStressor, err
}
