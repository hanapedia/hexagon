package factory

import (
	"errors"
	"strings"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/pkg/shared"
)

type ServiceAdapterDetails struct {
	serviceName   string
	protocol      string
	action        string
	handlerName string
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
		serviceAdapter, err = serviceAdapterDetails.RestServiceAdapterFactory()
	default:
		err = errors.New("No matching protocol found")
	}

	return serviceAdapter, err
}
