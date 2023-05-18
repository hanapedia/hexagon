package logger

import (
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
