package factory

import (
	"errors"

	model "github.com/hanapedia/the-bench/the-bench-operator/api/v1"
	"github.com/hanapedia/the-bench/the-bench-operator/pkg/constants"
	"github.com/hanapedia/the-bench/service-unit/stateless/internal/domain/core"
)

func statelesEgressAdapterFactory(adapterConfig model.StatelessEgressAdapterConfig) (core.EgressAdapter, error) {
	switch adapterConfig.Variant {
	case constants.REST:
		return restEgressAdapterFactory(adapterConfig)
	default:
		err := errors.New("No matching protocol found")
		return nil, err
	}

}
