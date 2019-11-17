package logging

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// NewLogger returns a logging object
func NewLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:          false,
		DisableLevelTruncation: true,
		FullTimestamp:          false,
	})
	logger.SetOutput(os.Stdout)

	// Get verbosity/log level from OS var
	if logLevel, ok := os.LookupEnv("LOG_LEVEL"); ok {
		switch strings.ToLower(logLevel) {
		case "debug":
			logger.Level = logrus.DebugLevel
		case "info":
			logger.Level = logrus.InfoLevel
		case "warning":
			logger.Level = logrus.WarnLevel
		case "error":
			logger.Level = logrus.ErrorLevel
		case "fatal":
			logger.Level = logrus.FatalLevel
		default:
			logger.Level = logrus.InfoLevel
		}
	} else {
		// Default to InfoLevel if env var not set
		logger.SetLevel(logrus.InfoLevel)
	}

	return logger
}

// TestLogger creates a dummy logger for test cases
var TestLogger = func() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:          false,
		DisableLevelTruncation: true,
		FullTimestamp:          false,
	})
	logger.Out = ioutil.Discard
	return logger
}
