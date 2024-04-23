package mock

import (
	"context"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
)

// SecondaryAdapterMock should mock SecodaryPort from
// "github.com/hanapedia/hexagon/internal/service-unit/application/ports"
type SecondaryAdapterMock struct {
	client ports.SecondaryAdapterClient
	ports.SecondaryPortBase
}

type SecondaryAdapterClientMock struct {
	client interface{}
}

// Call mock implementation
func (sacm SecondaryAdapterMock) Call(context.Context) ports.SecondaryPortCallResult {
	return ports.SecondaryPortCallResult{Payload: nil, Error: nil}
}

// Close mock implementation
func (sacm SecondaryAdapterClientMock) Close() {
	return
}

// NewSecondaryAdapter creates mocked implementation of ports.SecondaryPort
func NewSecondaryAdapter() ports.SecodaryPort {
	return &SecondaryAdapterMock{client: NewSecondaryAdapterClient()}
}

func NewSecondaryAdapterClient() ports.SecondaryAdapterClient {
	return SecondaryAdapterClientMock{client: nil}
}
