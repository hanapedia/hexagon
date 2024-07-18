package stressor

import (
	"errors"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/stressor/cpu"
	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
)

func NewSecondaryAdapter(adapterConfig *model.StressorConfig) (secondary.SecodaryPort, error) {
	switch adapterConfig.Variant {
	case constants.CPU:
		return cpu.CpuStressorAdapterFactory(adapterConfig)
	default:
		err := errors.New("No matching protocol found when creating stressor adapter.")
		return nil, err
	}

}
