package disk

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	model "github.com/hanapedia/hexagon/pkg/api/v1"
	"github.com/hanapedia/hexagon/pkg/operator/constants"
	"github.com/hanapedia/hexagon/pkg/operator/logger"
	"github.com/hanapedia/hexagon/pkg/service-unit/utils"
)

type DiskStressorClient struct {
	readFile  *os.File
	writeFile *os.File
	mu        sync.Mutex // Mutex to synchronize file access
}

func (dsc *DiskStressorClient) Close() {
	dsc.readFile.Close()
	dsc.writeFile.Close()
}

func NewDiskStressorClient(adapterConfig *model.StressorConfig) *DiskStressorClient {
	readFilePath := filepath.Join(constants.DISK_STRESSOR_TMP_FILEPATH, fmt.Sprintf("%s.%s", "read", adapterConfig.GetGroupByKey()))
	readFile, err := os.OpenFile(readFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.Logger.WithField("readFilePath", readFilePath).Fatalf("Failed to open read file: %v", err)
	}
	writeFilePath := filepath.Join(constants.DISK_STRESSOR_TMP_FILEPATH, fmt.Sprintf("%s.%s", "write", adapterConfig.GetGroupByKey()))
	writeFile, err := os.OpenFile(writeFilePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		logger.Logger.WithField("writeFilePath", writeFilePath).Fatalf("Failed to open write file: %v", err)
	}

	// write initial data
	payload := utils.GenerateRandomString(model.GetPayloadSize(adapterConfig.Payload))
	_, err = readFile.Write([]byte(payload))
	if err != nil {
		logger.Logger.WithField("readFilePath", readFilePath).Fatalf("error writing initial data to file: %v", err)
	}
	_, err = writeFile.Write([]byte(payload))
	if err != nil {
		logger.Logger.WithField("writeFilePath", writeFilePath).Fatalf("error writing initial data to file: %v", err)
	}
	return &DiskStressorClient{
		readFile:  readFile,
		writeFile: writeFile,
	}
}
