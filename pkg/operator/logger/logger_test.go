package logger

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAdaptoLogger(t *testing.T) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	adaptoLogger := &AdaptoLoggerType{logger: logger}

	msg := "Test info message"
	adaptoLogger.Info(msg, "key1", "value1", "key2", "value2")
}

func TestAdaptoLogger_Info(t *testing.T) {
	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetOutput(&buf)
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	adaptoLogger := &AdaptoLoggerType{logger: logger}

	msg := "Test info message"
	adaptoLogger.Info(msg, "key1", "value1", "key2", "value2")

	assert.Contains(t, buf.String(), "level=info")
	assert.Contains(t, buf.String(), "msg=\""+msg+"\"")
	assert.Contains(t, buf.String(), "key1=value1")
	assert.Contains(t, buf.String(), "key2=value2")
}

func TestConvertToLogrusFields(t *testing.T) {
	args := []interface{}{"key1", "value1", "key2", "value2", "key3", 123}
	fields := convertToLogrusFields(args...)
	fmt.Println(fields)

	expected := logrus.Fields{
		"key1": "value1",
		"key2": "value2",
		"key3": 123,
	}
	assert.Equal(t, expected, fields)
}
