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
			stressDisk(ctx, dsa.client.readFile, dsa.client.writeFile, dsa.iterations, dsa.payloadSize)
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

func stressDisk(ctx context.Context, srcFile, dstFile *os.File, iter int, chunkSize int64) error {
	buffer := make([]byte, chunkSize)

	for i := 0; i < iter; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Copy data from srcFile to dstFile
			if _, err := srcFile.Seek(0, 0); err != nil {
				return fmt.Errorf("failed to seek srcFile: %v", err)
			}
			if err := dstFile.Truncate(0); err != nil {
				return fmt.Errorf("failed to truncate dstFile: %v", err)
			}
			if _, err := dstFile.Seek(0, 0); err != nil {
				return fmt.Errorf("failed to seek dstFile: %v", err)
			}

			// Perform the copy
			_, err := io.CopyBuffer(dstFile, srcFile, buffer)
			if err != nil {
				return fmt.Errorf("error during io.Copy: %v", err)
			}

			// Swap src and dst for the next iteration
			srcFile, dstFile = dstFile, srcFile
		}
	}
	return nil
}
