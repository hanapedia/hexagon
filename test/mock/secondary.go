package mock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

// SecondaryAdapterMock should mock SecodaryPort from
// "github.com/hanapedia/hexagon/internal/service-unit/application/ports"
type SecondaryAdapterMock struct {
	mu        sync.Mutex
	client    secondary.SecondaryAdapterClient
	name      string
	duration  time.Duration
	failCount int
	secondary.SecondaryPortBase
}

type SecondaryAdapterClientMock struct {
	client interface{}
}

// Call mock implementation
func (sacm *SecondaryAdapterMock) Call(ctx context.Context) secondary.SecondaryPortCallResult {
	sacm.mu.Lock()
	defer sacm.mu.Unlock()
	logger.Logger.Infof("Mock calling %s", sacm.name)
	timer := time.NewTimer(sacm.duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return secondary.SecondaryPortCallResult{Payload: nil, Error: fmt.Errorf("Mocking Call Timed out")}
	case <-timer.C:
		if sacm.failCount > 0 {
			sacm.failCount--
			return secondary.SecondaryPortCallResult{Payload: nil, Error: fmt.Errorf("Mocking Call Fail")}
		}
	}
	return secondary.SecondaryPortCallResult{Payload: nil, Error: nil}
}

// Close mock implementation
func (sacm SecondaryAdapterClientMock) Close() {
	return
}

// NewSecondaryAdapter creates mocked implementation of ports.SecondaryPort
func NewSecondaryAdapter(name string, duration time.Duration, failCount int) secondary.SecodaryPort {
	return &SecondaryAdapterMock{
		client:    NewSecondaryAdapterClient(),
		name:      name,
		duration:  duration,
		failCount: failCount,
	}
}

func NewSecondaryAdapterClient() secondary.SecondaryAdapterClient {
	return SecondaryAdapterClientMock{client: nil}
}
