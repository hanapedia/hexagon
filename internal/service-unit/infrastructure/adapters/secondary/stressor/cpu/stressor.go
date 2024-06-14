package cpu

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
)

type cpuStressorAdapter struct {
	payloadSize int64
	duration    time.Duration
	threadCount int
	ports.SecondaryPortBase
}

func (csa *cpuStressorAdapter) Call(ctx context.Context) ports.SecondaryPortCallResult {
	// prepare payload
	payload := utils.GenerateRandomString(csa.payloadSize)

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start the desired number of goroutines
	for i := 0; i < csa.threadCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stressCPU(ctx, csa.duration)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	return ports.SecondaryPortCallResult{
		Payload: &payload,
		Error:   nil,
	}
}

// stressCPU generates artifical stress to cpu
func stressCPU(ctx context.Context, duration time.Duration) {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			return
		default:
			_ = math.Sqrt(rand.Float64())
		}
	}
}
