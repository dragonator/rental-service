// Package logger contains a sligthly extended version of the Zap logger.
package logger

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	_samplingConfigInitial    = 100
	_samplingConfigThereafter = 100
)

var _errInvalidLogLevel = fmt.Errorf("invalid log level")

// Logger is a logger implementation.
type Logger struct {
	*zap.SugaredLogger
	atomic zap.AtomicLevel
}

// NewLogger is a constructor func for Logger.
func NewLogger(logLevel string) *Logger {
	atom := zap.NewAtomicLevel()

	newLevel, err := parseLogLevel(logLevel)
	if err == nil {
		atom.SetLevel(newLevel)
	}

	encodingConfig := zapcore.EncoderConfig{
		TimeKey:     "time",
		LevelKey:    "level",
		NameKey:     "logger",
		MessageKey:  "message",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.UTC().Format(time.RFC3339))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	config := zap.Config{
		Level: atom,
		Sampling: &zap.SamplingConfig{
			Initial:    _samplingConfigInitial,
			Thereafter: _samplingConfigThereafter,
		},
		Encoding:         "json",
		EncoderConfig:    encodingConfig,
		OutputPaths:      []string{"stdout", "server.log"},
		ErrorOutputPaths: []string{"stderr", "server.log"},
	}

	core, err := config.Build()
	if err != nil {
		panic(fmt.Errorf("building logger config: %w", err))
	}

	logger := zap.New(core.Core(), zap.AddCaller())

	sugared := logger.Sugar()

	return &Logger{
		SugaredLogger: sugared,
		atomic:        atom,
	}
}

// SetLogLevel sets the logger level.
func (l *Logger) SetLogLevel(level string) error {
	newLevel, err := parseLogLevel(level)
	if err != nil {
		return fmt.Errorf("parsing log level: %w", err)
	}

	l.atomic.SetLevel(newLevel)

	return nil
}

func parseLogLevel(level string) (zapcore.Level, error) {
	level = strings.ToUpper(level)

	var newLevel zapcore.Level

	switch level {
	case "ERROR":
		newLevel = zap.ErrorLevel
	case "WARN":
		newLevel = zap.WarnLevel
	case "DEBUG":
		newLevel = zap.DebugLevel
	case "INFO":
		newLevel = zap.InfoLevel
	default:
		return newLevel, fmt.Errorf("%w: %s", _errInvalidLogLevel, level)
	}

	return newLevel, nil
}
