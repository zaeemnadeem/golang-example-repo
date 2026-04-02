package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

// Init initializes the global logger based on environment
func Init(env string) *zap.Logger {
	once.Do(func() {
		var config zap.Config
		if env == "production" {
			config = zap.NewProductionConfig()
		} else {
			config = zap.NewDevelopmentConfig()
			config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		// Ensure we always write to stdout in our containers
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stderr"}

		logger, err := config.Build()
		if err != nil {
			// Fallback if something went terribly wrong
			logger = zap.NewExample()
		}
		globalLogger = logger
		zap.ReplaceGlobals(logger) // Set global zap logger
	})
	return globalLogger
}

// Get returns the initialized global logger, or a no-op logger if not initialized.
func Get() *zap.Logger {
	if globalLogger == nil {
		return zap.NewNop()
	}
	return globalLogger
}
