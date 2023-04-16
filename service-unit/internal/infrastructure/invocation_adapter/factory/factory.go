package factory

import (
	"errors"
	"strings"

	"github.com/hanapedia/the-bench/service-unit/internal/domain/core"
	"github.com/hanapedia/the-bench/service-unit/pkg/shared"
)

type InvocationAdapterDetails struct {
	serviceName   string
	protocol      string
	action        string
	handlerName string
}

func newInvocationAdapterDetails(id string) (InvocationAdapterDetails, error) {
	idSubstring := strings.Split(id, ".")
	var err error = nil
	if len(idSubstring) != 4 {
		err = errors.New("Invalid adapter id")
	}
	return InvocationAdapterDetails{
		serviceName:   idSubstring[shared.ServiceNameIndex],
		protocol:      idSubstring[shared.ProtocolIndex],
		action:        idSubstring[shared.ActionIndex],
		handlerName: idSubstring[shared.AdapterNameIndex],
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
