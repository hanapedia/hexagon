package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

// SecondaryAdapterMock should mock SecodaryPort from
// "github.com/hanapedia/hexagon/internal/service-unit/application/ports"
type SecondaryAdapterMock struct {
	mu        sync.Mutex
	client    ports.SecondaryAdapterClient
	name      string
	duration  time.Duration
	failCount int
	ports.SecondaryPortBase
}

type SecondaryAdapterClientMock struct {
	client interface{}
}

// Call mock implementation
func (sacm *SecondaryAdapterMock) Call(ctx context.Context) ports.SecondaryPortCallResult {
	sacm.mu.Lock()
	defer sacm.mu.Unlock()
	logger.Logger.Infof("Mock calling %s", sacm.name)
	timer := time.NewTimer(sacm.duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ports.SecondaryPortCallResult{Payload: nil, Error: fmt.Errorf("Mocking Call Timed out")}
	case <-timer.C:
		if sacm.failCount > 0 {
			sacm.failCount--
			return ports.SecondaryPortCallResult{Payload: nil, Error: fmt.Errorf("Mocking Call Fail")}
		}
	}
	return ports.SecondaryPortCallResult{Payload: nil, Error: nil}
}

// Close mock implementation
func (sacm SecondaryAdapterClientMock) Close() {
	return
}

// NewSecondaryAdapter creates mocked implementation of ports.SecondaryPort
func NewSecondaryAdapter(name string, duration time.Duration, failCount int) ports.SecodaryPort {
	return &SecondaryAdapterMock{
		client:    NewSecondaryAdapterClient(),
		name:      name,
		duration:  duration,
		failCount: failCount,
	}
}

func NewSecondaryAdapterClient() ports.SecondaryAdapterClient {
	return SecondaryAdapterClientMock{client: nil}
}
