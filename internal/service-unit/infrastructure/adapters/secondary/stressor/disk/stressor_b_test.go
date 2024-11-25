package disk

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
)

/*
## Results
Payload sizes
	SMALL:  1024,
	MEDIUM: 4096,
	LARGE:  16384,

goos: darwin
goarch: arm64
pkg: github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/stressor/disk
cpu: Apple M1
BenchmarkStressDiskPayloadSizeB/small-8                    58884             20176 ns/op
BenchmarkStressDiskPayloadSizeB/medium-8                   52226             25377 ns/op
BenchmarkStressDiskPayloadSizeB/large-8                    35560            245688 ns/op
*/
func BenchmarkStressDiskPayloadSizeB(b *testing.B) {
	iter := 10
	payloads := map[string]string{
		string(constants.SMALL):  utils.GenerateRandomString(constants.PayloadSizeMap[constants.SMALL]),
		string(constants.MEDIUM): utils.GenerateRandomString(constants.PayloadSizeMap[constants.MEDIUM]),
		string(constants.LARGE):  utils.GenerateRandomString(constants.PayloadSizeMap[constants.LARGE]),
	}

	// Open the file once for the duration of the benchmark
	file, err := os.OpenFile(constants.DISK_STRESSOR_TMP_FILEPATH, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		b.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Iterate over the payloads and create sub-benchmarks
	for k, p := range payloads {
		b.Run(k, func(b *testing.B) {
			// Truncate the file to zero size
			if err := file.Truncate(0); err != nil {
				b.Fatalf("Failed to truncate file: %v", err)
			}

			// Reset the file pointer to the beginning
			if _, err := file.Seek(0, 0); err != nil {
				b.Fatalf("Failed to reset file pointer: %v", err)
			}

			for i := 0; i < b.N; i++ {
				// Perform disk stress for the current payload
				err := stressDisk(context.Background(), file, iter, []byte(p))
				if err != nil {
					b.Fatalf("Failed during stressDisk: %v", err)
				}
			}
		})
	}
}

/*
## Results
Iters 10, 20, 40, 160

goos: darwin
goarch: arm64
pkg: github.com/hanapedia/hexagon/internal/service-unit/infrastructure/adapters/secondary/stressor/disk
cpu: Apple M1
BenchmarkStressDiskItersB/iters:_0-8                       58228             20157 ns/op
BenchmarkStressDiskItersB/iters:_1-8                       29820             40238 ns/op
BenchmarkStressDiskItersB/iters:_2-8                       15003            104261 ns/op
BenchmarkStressDiskItersB/iters:_3-8                        3349            374009 ns/op
*/
func BenchmarkStressDiskItersB(b *testing.B) {
	iters := []int{
		10, 20, 40, 160,
	}
	p := utils.GenerateRandomString(constants.PayloadSizeMap[constants.SMALL])

	// Open the file once for the duration of the benchmark
	file, err := os.OpenFile(constants.DISK_STRESSOR_TMP_FILEPATH, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		b.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Iterate over the payloads and create sub-benchmarks
	for i, iter := range iters {
		b.Run(fmt.Sprintf("iters: %v", i), func(b *testing.B) {
			// Truncate the file to zero size
			if err := file.Truncate(0); err != nil {
				b.Fatalf("Failed to truncate file: %v", err)
			}

			// Reset the file pointer to the beginning
			if _, err := file.Seek(0, 0); err != nil {
				b.Fatalf("Failed to reset file pointer: %v", err)
			}

			for i := 0; i < b.N; i++ {
				// Perform disk stress for the current payload
				err := stressDisk(context.Background(), file, iter, []byte(p))
				if err != nil {
					b.Fatalf("Failed during stressDisk: %v", err)
				}
			}
		})
	}
}
