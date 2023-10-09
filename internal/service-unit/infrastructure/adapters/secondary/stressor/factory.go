package stressor

import (
	"errors"

	"github.com/hanapedia/the-bench/internal/service-unit/application/ports"
	"github.com/hanapedia/the-bench/internal/service-unit/infrastructure/adapters/secondary/stressor/cpu"
	model "github.com/hanapedia/the-bench/pkg/api/v1"
	"github.com/hanapedia/the-bench/pkg/operator/constants"
)

func NewSecondaryAdapter(adapterConfig *model.StressorConfig) (ports.SecodaryPort, error) {
	switch adapterConfig.Variant {
	case constants.CPU:
		return cpu.CpuStressorAdapterFactory(adapterConfig)
	default:
		err := errors.New("No matching protocol found when creating producer adapter.")
		return nil, err
	}

}
