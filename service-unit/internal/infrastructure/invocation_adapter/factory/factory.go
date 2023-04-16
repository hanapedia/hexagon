package factory

import (
	"errors"
	"strings"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/pkg/constants"
)

type InvocationAdapterDetails struct {
	serviceName   string
	protocol      string
	action        string
	handlerName string
}

func newInvocationAdapterDetails(id string) (InvocationAdapterDetails, error) {
	idSubstring := strings.Split(id, ".")
	var err error
	if len(idSubstring) != 4 {
		err = errors.New("Invalid adapter id")
	}
	return InvocationAdapterDetails{
		serviceName:   idSubstring[constants.ServiceNameIndex],
		protocol:      idSubstring[constants.ProtocolIndex],
		action:        idSubstring[constants.ActionIndex],
		handlerName: idSubstring[constants.AdapterNameIndex],
	}, err
}

func NewInvocationAdapterFromID(id string) (core.InvocationAdapter, error) {
	invocationAdapterDetails, err := newInvocationAdapterDetails(id)
	var invocationAdapter core.InvocationAdapter
	switch invocationAdapterDetails.protocol {
	case "rest":
		invocationAdapter, err = invocationAdapterDetails.RestInvocationAdapterFactory()
	default:
		err = errors.New("No matching protocol found")
	}

	return invocationAdapter, err
}
