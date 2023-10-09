package cpu

import (
	"time"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	model "github.com/hanapedia/the-bench/pkg/api/v1"
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

	cpuStressor = &CpuStressorAdapter{
		payload:     adapterConfig.Payload,
		duration:    duration,
		threadCount: threadCount,
	}

	// set destionation id
	cpuStressor.SetDestId(adapterConfig.GetId())

	return cpuStressor, err
}