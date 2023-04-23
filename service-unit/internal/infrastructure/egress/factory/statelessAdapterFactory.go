package factory

import (
	"errors"

	"github.com/hanapedia/the-bench/config/constants"
	"github.com/hanapedia/the-bench/config/model"
	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
)

func statelesEgressAdapterFactory(adapterConfig model.StatelessEgressConfig) (core.EgressAdapter, error) {
	switch adapterConfig.Variant {
	case constants.REST:
		return restEgressAdapterFactory(adapterConfig)
	default:
		err := errors.New("No matching protocol found")
		return nil, err
	}

}
