package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap.Logger instance
func NewLogger(env, level, path string) *zap.Logger {
	// Using zap's preset constructors is the simplest way to get a feel for the
	// package, but they don't allow much customization.
	var logger *zap.Logger

	if env != "production" {
		logger = zap.Must(developmentLogger(level))
	} else {
		logger = zap.Must(productionLogger(level))
	}

	return logger
}

// set production logger defaults
func productionLogger(l string) (*zap.Logger, error) {
	level, _ := zap.ParseAtomicLevel(strings.ToLower(l))

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(level.Level()),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build()
}

// set development logger defaults
func developmentLogger(l string) (*zap.Logger, error) {
	level, _ := zap.ParseAtomicLevel(strings.ToLower(l))

	encoderCfg := zap.NewDevelopmentEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(level.Level()),
		Development: true,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    encoderCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return cfg.Build()
}
