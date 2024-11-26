package disk

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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
BenchmarkStressDiskPayloadSizeB/small-8                     7405            159164 ns/op
BenchmarkStressDiskPayloadSizeB/medium-8                    6418            167628 ns/op
BenchmarkStressDiskPayloadSizeB/large-8                     6600            178190 ns/op
*/
func BenchmarkStressDiskPayloadSizeB(b *testing.B) {
	iter := 10
	payloads := map[string]int64{
		string(constants.SMALL):  constants.PayloadSizeMap[constants.SMALL],
		string(constants.MEDIUM): constants.PayloadSizeMap[constants.MEDIUM],
		string(constants.LARGE):  constants.PayloadSizeMap[constants.LARGE],
	}

	readFilePath := filepath.Join(constants.DISK_STRESSOR_TMP_FILEPATH, fmt.Sprintf("%s.%s", "read", "id"))
	readFile, err := os.OpenFile(readFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		b.Fatalf("Failed to open read file: %v", err)
	}
	writeFilePath := filepath.Join(constants.DISK_STRESSOR_TMP_FILEPATH, fmt.Sprintf("%s.%s", "write", "id"))
	writeFile, err := os.OpenFile(writeFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		b.Fatalf("Failed to open write file: %v", err)
	}

	defer readFile.Close()
	defer writeFile.Close()

	// Iterate over the payloads and create sub-benchmarks
	for k, p := range payloads {
		b.Run(k, func(b *testing.B) {
			// write initial data
			if err := readFile.Truncate(0); err != nil {
				b.Fatalf("failed to truncate readFile: %v", err)
			}
			if _, err := readFile.Seek(0, 0); err != nil {
				b.Fatalf("failed to seek readFile: %v", err)
			}
			if err := writeFile.Truncate(0); err != nil {
				b.Fatalf("failed to truncate writeFile: %v", err)
			}
			if _, err := writeFile.Seek(0, 0); err != nil {
				b.Fatalf("failed to seek writeFile: %v", err)
			}
			payload := utils.GenerateRandomString(p)
			_, err = readFile.Write([]byte(payload))
			if err != nil {
				b.Fatalf("error writing initial data to file: %v", err)
			}
			_, err = writeFile.Write([]byte(payload))
			if err != nil {
				b.Fatalf("error writing initial data to file: %v", err)
			}

			for i := 0; i < b.N; i++ {
				// Perform disk stress for the current payload
				err := stressDisk(context.Background(), readFile, writeFile, iter, p)
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
BenchmarkStressDiskItersB/iters:_0-8                        7326            158623 ns/op
BenchmarkStressDiskItersB/iters:_1-8                        3706            317450 ns/op
BenchmarkStressDiskItersB/iters:_2-8                        1862            638363 ns/op
BenchmarkStressDiskItersB/iters:_3-8                         468           2527196 ns/op
*/
func BenchmarkStressDiskItersB(b *testing.B) {
	iters := []int{
		10, 20, 40, 160,
	}
	p := constants.PayloadSizeMap[constants.SMALL]

	readFilePath := filepath.Join(constants.DISK_STRESSOR_TMP_FILEPATH, fmt.Sprintf("%s.%s", "read", "id"))
	readFile, err := os.OpenFile(readFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		b.Fatalf("Failed to open read file: %v", err)
	}
	writeFilePath := filepath.Join(constants.DISK_STRESSOR_TMP_FILEPATH, fmt.Sprintf("%s.%s", "write", "id"))
	writeFile, err := os.OpenFile(writeFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		b.Fatalf("Failed to open write file: %v", err)
	}

	defer readFile.Close()
	defer writeFile.Close()

	// Iterate over the payloads and create sub-benchmarks
	for i, iter := range iters {
		b.Run(fmt.Sprintf("iters: %v", i), func(b *testing.B) {
			if err := readFile.Truncate(0); err != nil {
				b.Fatalf("failed to truncate readFile: %v", err)
			}
			if _, err := readFile.Seek(0, 0); err != nil {
				b.Fatalf("failed to seek readFile: %v", err)
			}
			if err := writeFile.Truncate(0); err != nil {
				b.Fatalf("failed to truncate writeFile: %v", err)
			}
			if _, err := writeFile.Seek(0, 0); err != nil {
				b.Fatalf("failed to seek writeFile: %v", err)
			}
			// write initial data
			payload := utils.GenerateRandomString(p)
			_, err = readFile.Write([]byte(payload))
			if err != nil {
				b.Fatalf("error writing initial data to file: %v", err)
			}
			_, err = writeFile.Write([]byte(payload))
			if err != nil {
				b.Fatalf("error writing initial data to file: %v", err)
			}

			for i := 0; i < b.N; i++ {
				// Perform disk stress for the current payload
				err := stressDisk(context.Background(), readFile, writeFile, iter, p)
				if err != nil {
					b.Fatalf("Failed during stressDisk: %v", err)
				}
			}
		})
	}
}
