package factory

import (
	"errors"
	"strings"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/internal/infrastructure/service_adapter/rest"
	"github.com/hanapedia/the-bench/service-unit/pkg/shared"
)

type ServiceAdapterDetails struct {
	serviceName   string
	protocol      string
	action        string
	handlerName string
}

func (serviceAdapterDetails ServiceAdapterDetails) RestAdapterBuilder() (core.ServiceAdapter, error) {
	var err error = nil
	var serviceAdapter core.ServiceAdapter
	switch serviceAdapterDetails.action {
	case "read":
		serviceAdapter = rest.RestReadAdapter{URL: serviceAdapterDetails.serviceName}
	case "write":
		serviceAdapter = rest.RestWriteAdapter{URL: serviceAdapterDetails.serviceName}
	default:
		err = errors.New("No matching protocol found")
	}
	return serviceAdapter, err
}

func newServiceAdapterDetails(id string) (ServiceAdapterDetails, error) {
	idSubstring := strings.Split(id, ".")
	var err error = nil
	if len(idSubstring) != 4 {
		err = errors.New("Invalid adapter id")
	}
	return ServiceAdapterDetails{
		serviceName:   idSubstring[shared.ServiceNameIndex],
		protocol:      idSubstring[shared.ProtocolIndex],
		action:        idSubstring[shared.ActionIndex],
		handlerName: idSubstring[shared.AdapterNameIndex],
	}, err
}

func NewServiceAdapterFromID(id string) (core.ServiceAdapter, error) {
	serviceAdapterDetails, err := newServiceAdapterDetails(id)
	var serviceAdapter core.ServiceAdapter
	switch serviceAdapterDetails.protocol {
	case "rest":
		serviceAdapter, err = serviceAdapterDetails.RestAdapterBuilder()
	default:
		err = errors.New("No matching protocol found")
	}

	return serviceAdapter, err
}
