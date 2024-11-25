package disk

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
)

type DiskStressorClient struct {
	file *os.File
	mu   sync.Mutex // Mutex to synchronize file access
}

func (dsc *DiskStressorClient) Close() {
	dsc.file.Close()
}

func NewDiskStressorClient(id string) *DiskStressorClient {
	filePath := filepath.Join(constants.DISK_STRESSOR_TMP_FILEPATH, id)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.Logger.WithField("filePath", filePath).Fatalf("Failed to open file")
	}
	return &DiskStressorClient{
		file: file,
	}
}
