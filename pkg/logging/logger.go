package logging

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogConfig allows customization of logging behavior
type LogConfig struct {
	ServiceName string
	Debug       bool
	WriteToFile bool
}

// Logger wraps zerolog.Logger for additional functionality
type Logger struct {
	zerolog.Logger
}

// NewLogger creates a new configured logger
func NewLogger(config LogConfig) Logger {
	// Configure base logger
	zerolog.TimeFieldFormat = time.RFC3339Nano
	
	// Create multi-writer for console and optional file output
	var writers []io.Writer
	
	// Always write to console
	writers = append(writers, zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
	})

	// Optional file logging
	if config.WriteToFile {
		fileWriter := &lumberjack.Logger{
			Filename:   fmt.Sprintf("logs/%s.log", config.ServiceName),
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
			Compress:   true,
		}
		writers = append(writers, fileWriter)
	}

	// Create multi-writer
	multiWriter := io.MultiWriter(writers...)

	// Create logger
	logger := zerolog.New(multiWriter).
		With().
		Timestamp().
		Str("service", config.ServiceName).
		Logger()

	// Set global log level
	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	return Logger{logger}
}

// WithRequestID adds a request ID to the logger context
func (l Logger) WithRequestID(requestID string) *zerolog.Logger {
	logger := l.With().Str("request_id", requestID).Logger()
	return &logger
}

// ErrorWithContext logs an error with additional context
func (l Logger) ErrorWithContext(err error, msg string, fields map[string]interface{}) {
	event := l.Error().Err(err)
	for k, v := range fields {
		event = event.Interface(k, v)
	}
	event.Msg(msg)
}
