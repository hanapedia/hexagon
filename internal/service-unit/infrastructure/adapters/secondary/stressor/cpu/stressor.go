package cpu

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
)

type cpuStressorAdapter struct {
	payloadSize int64
	iterations  int
	threadCount int
	secondary.SecondaryPortBase
}

func (csa *cpuStressorAdapter) Call(ctx context.Context) secondary.SecondaryPortCallResult {
	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start the desired number of goroutines
	for i := 0; i < csa.threadCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stressCPU(ctx, csa.iterations)
		}()
	}

	// Wait for all goroutines to finish
	wg.Wait()

	if ctx.Err() != nil {
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error: fmt.Errorf("CPU stressor Call timeout exceeded"),
		}
	}

	// prepare payload
	payload := utils.GenerateRandomString(csa.payloadSize)

	return secondary.SecondaryPortCallResult{
		Payload: &payload,
		Error:   nil,
	}
}

// stressCPU generates artifical stress to cpu
func stressCPU(ctx context.Context, iter int) {
	for i := 0; i < iter; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			_ = math.Sqrt(rand.Float64())
		}
	}
}
