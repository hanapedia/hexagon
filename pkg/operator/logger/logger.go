package logger

import (
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

var AdaptoLogger *AdaptoLoggerType

type AdaptoLoggerType struct {
	logger *logrus.Logger
}

// Info logs an informational message
func (l *AdaptoLoggerType) Info(msg string, args ...interface{}) {
	fields := convertToLogrusFields(args)
	l.logger.WithFields(fields).Info(msg)
}

// Error logs an error message
func (l *AdaptoLoggerType) Error(msg string, args ...interface{}) {
	fields := convertToLogrusFields(args)
	l.logger.WithFields(fields).Error(msg)
}

// Debug logs a debug message
func (l *AdaptoLoggerType) Debug(msg string, args ...interface{}) {
	fields := convertToLogrusFields(args)
	l.logger.WithFields(fields).Debug(msg)
}

// convertToLogrusFields converts a key-value sequence to logrus.Fields
func convertToLogrusFields(args ...interface{}) logrus.Fields {
	fields := make(logrus.Fields)

	for i := 0; i < len(args)-1; i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue // skip if key is not a string
		}
		fields[key] = args[i+1]
	}
	return fields
}

func init() {
	Logger = logrus.New()

	// Set log formatting to use the text formatter with timestamps
	Logger.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp: true,
	})

	// Set log level to show all log levels
	Logger.SetLevel(logrus.TraceLevel)

	AdaptoLogger = &AdaptoLoggerType{logger: Logger}
}
