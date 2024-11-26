package stressor

import (
	"errors"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/stressor/cpu"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/stressor/disk"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

func NewSecondaryAdapter(adapterConfig *model.StressorConfig, client secondary.SecondaryAdapterClient) (secondary.SecodaryPort, error) {
	switch adapterConfig.Variant {
	case constants.CPU:
		return cpu.CpuStressorAdapterFactory(adapterConfig)
	case constants.DISK:
		return disk.DiskStressorAdapterFactory(adapterConfig, client)
	default:
		err := errors.New("No matching protocol found when creating stressor adapter.")
		return nil, err
	}

}

func NewClient(adapterConfig *model.StressorConfig) secondary.SecondaryAdapterClient {
	switch adapterConfig.Variant {
	case constants.CPU:
		cpuStressorClient := cpu.NewCpuStressorClient()
		return cpuStressorClient
	case constants.DISK:
		diskStressorClient := disk.NewDiskStressorClient(adapterConfig)
		return diskStressorClient
	default:
		logger.Logger.Fatalf("invalid protocol")
		return nil
	}
}
