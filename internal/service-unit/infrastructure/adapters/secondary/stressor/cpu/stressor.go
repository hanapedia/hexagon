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
	payload, err := utils.GenerateRandomString(csa.payloadSize)
	if err != nil {
		return ports.SecondaryPortCallResult{
			Payload: nil,
			Error:   err,
		}
	}

	// Convert the duration to CPU time
	targetTime := currentCPUTime() + csa.duration.Nanoseconds()

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start the desired number of goroutines
	for i := 0; i < csa.threadCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stressCPU(targetTime)
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
func stressCPU(targetTime int64) {
	for {
		math.Sqrt(rand.Float64())
		// If the current CPU time exceeds the target, exit the goroutine
		if currentCPUTime() >= targetTime {
			break
		}
	}
}

// currentCPUTime gets current CPU time in nanoseconds
func currentCPUTime() int64 {
	return time.Now().UnixNano()
}
