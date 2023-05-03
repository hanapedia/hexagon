package validation

import (
	"github.com/hanapedia/the-bench/config/model"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func init() {
	Logger = logrus.New()

	// Set log formatting to use the text formatter with timestamps
	Logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	// Set log level to show all log levels
	Logger.SetLevel(logrus.TraceLevel)
}

func PrintErrors(cve model.ConfigValidationError) {
	for _, err := range cve.ServiceUnitFieldErrors {
		Logger.Errorf(err.Error())
	}
	for _, err := range cve.AdapterFieldErrors {
		Logger.Errorf(err.Error())
	}
	for _, err := range cve.MappingErrors {
		Logger.Errorf(err.Error())
	}

}
