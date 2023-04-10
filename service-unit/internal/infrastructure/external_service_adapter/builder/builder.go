package builder

import (
	"errors"
	"strings"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/external_service_adapter/rest"
	"github.com/hanapedia/the-bench/service-unit/pkg/shared"
)

type AdapterDetails struct {
	serviceName   string
	protocol      string
	action        string
	handlerName string
}

func (adapterDetails AdapterDetails) RestAdapterBuilder() (core.ExternalServiceAdapter, error) {
	var err error = nil
	var externalServiceAdapter core.ExternalServiceAdapter
	switch adapterDetails.action {
	case "read":
		externalServiceAdapter = rest.RestReadAdapter{URL: adapterDetails.serviceName}
	case "write":
		externalServiceAdapter = rest.RestWriteAdapter{URL: adapterDetails.serviceName}
	default:
		err = errors.New("No matching protocol found")
	}
	return externalServiceAdapter, err
}

func newAdapterDetails(id string) (AdapterDetails, error) {
	idSubstring := strings.Split(id, ".")
	var err error = nil
	if len(idSubstring) != 4 {
		err = errors.New("Invalid adapter id")
	}
	return AdapterDetails{
		serviceName:   idSubstring[shared.ServiceNameIndex],
		protocol:      idSubstring[shared.ProtocolIndex],
		action:        idSubstring[shared.ActionIndex],
		handlerName: idSubstring[shared.AdapterNameIndex],
	}, err
}

func NewExternalServiceAdapterFromID(id string) (core.ExternalServiceAdapter, error) {
	adapterDetails, err := newAdapterDetails(id)
	var externalServiceAdapter core.ExternalServiceAdapter
	switch adapterDetails.protocol {
	case "rest":
		externalServiceAdapter, err = adapterDetails.RestAdapterBuilder()
	default:
		err = errors.New("No matching protocol found")
	}

	return externalServiceAdapter, err
}
