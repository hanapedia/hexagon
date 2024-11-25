package disk

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/hanapedia/hexagon/internal/service-unit/application/ports/secondary"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
)

type diskStressorAdapter struct {
	payloadSize int64
	iterations  int
	threadCount int

	client *DiskStressorClient

	secondary.SecondaryPortBase
}

func (dsa *diskStressorAdapter) Call(ctx context.Context) secondary.SecondaryPortCallResult {
	dsa.client.mu.Lock()
	defer dsa.client.mu.Unlock()

	// prepare payload
	payload := utils.GenerateRandomString(dsa.payloadSize)

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start the desired number of goroutines
	for i := 0; i < dsa.threadCount; i++ {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			stressDisk(ctx, dsa.client.file, dsa.iterations, []byte(p))
		}(payload)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	if ctx.Err() != nil {
		return secondary.SecondaryPortCallResult{
			Payload: nil,
			Error:   fmt.Errorf("Disk stressor Call timeout exceeded"),
		}
	}

	return secondary.SecondaryPortCallResult{
		Payload: &payload,
		Error:   nil,
	}
}

// stressDisk simulates disk I/O stress by reading and writing small chunks iteratively
func stressDisk(ctx context.Context, file *os.File, iter int, payload []byte) error {
	for i := 0; i < iter; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Read from the file
			_, err := file.Read(payload)
			if err != nil && err != os.ErrClosed && err != io.EOF {
				return fmt.Errorf("error reading file: %v", err)
			}

			// Write back to the file
			_, err = file.Write(payload)
			if err != nil {
				return fmt.Errorf("error writing file: %v", err)
			}
		}
	}

	return nil
}
