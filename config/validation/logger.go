package validation

import (
	"github.com/hanapedia/the-bench/config/model"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()

	// Set log formatting to use the text formatter with timestamps
	logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	// Set log level to show all log levels
	logger.SetLevel(logrus.TraceLevel)
}

func PrintErrors(cve model.ConfigValidationError) {
	for _, err := range cve.ServiceUnitFieldErrors {
		logger.Errorf(err.Error())
	}
	for _, err := range cve.AdapterFieldErrors {
		logger.Errorf(err.Error())
	}
	for _, err := range cve.MappingErrors {
		logger.Errorf(err.Error())
	}

}
