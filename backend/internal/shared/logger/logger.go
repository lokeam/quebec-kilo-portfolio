package logger

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type AlertLevel string
type Environment string

const (
	EnvDev  Environment =  "dev"
	EnvProd Environment = "prod"
	EnvTest Environment = "test"
)

const (
	// Alert levels
	LevelDebug AlertLevel = "debug"
	LevelInfo  AlertLevel = "info"
	LevelWarn  AlertLevel = "warn"
	LevelError AlertLevel = "error"

	// Color codes
	yellow = "\x1b[33m"
	red = "\x1b[31m"
	green = "\x1b[32m"
	purple  = "\x1b[35m"
	blue = "\x1b[34m"
	magenta = "\x1b[95m"
	cyan = "\x1b[36m"
	gray = "\x1b[90m"
	reset = "\x1b[0m"
)

type customConsoleEncoder struct {
	zapcore.Encoder
}

type LoggerInterface interface {
	Info(msg string, fields map[string]any)
	Debug(msg string, fields map[string]any)
	Warn(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)
}

// Logger is a wrapper around zap
type Logger struct {
	output   io.Writer
	env      Environment
	level    AlertLevel
	zap      *zap.Logger
	debug    bool
}

// Option allows for a functional approach to configuring Logger
type Option func(*Logger)

/* --- Extendable 'With' methods to be used by function caller ---*/

// WithOutput sets the output destination
func WithOutput(w io.Writer) Option {
	return func(l *Logger) {
		l.output = w
	}
}

// WithEnv sets the dev or prod environment
func WithEnv(env Environment) Option {
	return func(l *Logger) {
		l.env = env
	}
}

// WithAlertLevel sets the logging alert "level"
func WithAlertLevel(level AlertLevel) Option {
	return func(l *Logger) {
		l.level = level
	}
}

func WithDebug(debug bool) Option {
	return func(l *Logger) {
		l.debug = debug
	}
}

func NewLogger(options ...Option) (*Logger, error) {
	// Begin with safe defaults
	logger := &Logger{
			output: os.Stdout, // Defaults to console
			env:    EnvProd,   // Defaults to prod
			level:  LevelInfo, // Defaults to info level
			debug:  false,
	}

	// Apply any custom options
	for _, opt := range options {
			opt(logger)
	}

	// Configure Zap based on environment
	var config zap.Config
	var core zapcore.Core

	if logger.env == EnvDev {
			// Development config - pretty printing
			config = zap.NewDevelopmentConfig()

			// Add enhanced color encoding
			config.EncoderConfig = zapcore.EncoderConfig{
				TimeKey:          "timestamp",
				LevelKey:         "level",
				NameKey:          "logger",
				CallerKey:        "caller",
				MessageKey:       "msg",
				StacktraceKey:    "stacktrace",
				LineEnding:       zapcore.DefaultLineEnding,

				// Color encoding for specific components
				EncodeLevel:      zapcore.CapitalColorLevelEncoder,
				EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) { // Gray color for timestamp
					enc.AppendString(gray + t.Format(time.RFC3339) + reset)
				},
				EncodeDuration:   zapcore.SecondsDurationEncoder,
				EncodeCaller:   func(
					caller zapcore.EntryCaller,
					enc zapcore.PrimitiveArrayEncoder,
					) {
					enc.AppendString("\x1b[36m" + caller.TrimmedPath() + "\x1b[0m") // Cyan color for caller
				},
			}

			// Create core with console encoder for readable output
			encoder := customConsoleEncoder{zapcore.NewConsoleEncoder(config.EncoderConfig)}

			core = zapcore.NewCore(
				encoder,
				zapcore.AddSync(logger.output),
				zap.NewAtomicLevelAt(zapcore.DebugLevel),
			)
	} else {
			// Production config - JSON format
			config = zap.NewProductionConfig()
			config.EncoderConfig = zapcore.EncoderConfig{
					TimeKey:        "timestamp",
					LevelKey:       "level",
					NameKey:        "logger",
					CallerKey:      "caller",
					MessageKey:     "msg",
					StacktraceKey:  "stacktrace",
					LineEnding:     zapcore.DefaultLineEnding,
					EncodeLevel:    zapcore.LowercaseLevelEncoder,
					EncodeTime:     zapcore.ISO8601TimeEncoder,
					EncodeDuration: zapcore.SecondsDurationEncoder,
					EncodeCaller:   zapcore.ShortCallerEncoder,
			}

			// Create core with JSON encoder for structured output
			core = zapcore.NewCore(
				zapcore.NewJSONEncoder(config.EncoderConfig),
				zapcore.AddSync(logger.output),
				zap.NewAtomicLevelAt(zapcore.InfoLevel),
			)
	}

	// Create the logger with the configured core
	logger.zap = zap.New(
			core,
			// zap.AddCaller(), TODO: Debug Zapcore bug where call defaults to logger instead of error file source
	)

	return logger, nil
}

// Helper fn - Convert fields to zap fields
func (l *Logger) convertToZapFields(fields map[string]any) []zap.Field {
	if fields == nil {
			return nil
	}

	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {

		// If we're in dev mode, use human readable values
		if l.env == EnvDev {
			switch v := value.(type) {
			case nil:
					zapFields = append(zapFields, zap.String(key, "<nil>"))
			case []interface{}:
					zapFields = append(zapFields, zap.String(key, fmt.Sprint(v)))
			case map[string]interface{}:
					zapFields = append(zapFields, zap.String(key, fmt.Sprintf("map[%v]", v)))
			case struct{ Name string }:
					zapFields = append(zapFields, zap.String(key, fmt.Sprintf("{%s}", v.Name)))
			default:
				zapFields = append(zapFields, zap.String(key, fmt.Sprintf("%v", v)))
			}
		} else {
			zapFields = append(zapFields, zap.Any(key, value))
		}

	}
	return zapFields
}

// Logger methods by alert level (Info, Debug, Warn, Error)
func (logger *Logger) Info(msg string, fields map[string]any) {
	logger.zap.Info(msg, logger.convertToZapFields(fields)...)
}

func (logger *Logger) Debug(msg string, fields map[string]any) {
	logger.zap.Debug(msg, logger.convertToZapFields(fields)...)
}

func (logger *Logger) Warn(msg string, fields map[string]any) {
	logger.zap.Warn(msg, logger.convertToZapFields(fields)...)
}

func (logger *Logger) Error(msg string, fields map[string]any) {
	logger.zap.Error(msg, logger.convertToZapFields(fields)...)
}

// LogMiddleware creates a middleware that logs HTTP requests
func (l *Logger) LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Stat timer
		start := time.Now()

		// Create a response wrapper to grab the status code
		responseWrapper := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Setup recovery in case of panic
		defer func() {
			if err := recover(); err != nil {
				// Log panic
				l.Error("panic", map[string]any{
					"error":   err,
					"method":  r.Method,
					"path":   r.URL.Path,
					"status": 500,
				})

				// Re-panic after logging
				panic(err)
			}
		}()

		// Process request
		next.ServeHTTP(responseWrapper, r)

		// Log request details
		requestFields := map[string]any{
			// Request info
			"method":      r.Method,
			"path":        r.URL.Path,
			"query":       r.URL.Query().Encode(),

			// Response info
			"status":     responseWrapper.Status(),
			"duration":   time.Since(start).String(),

			// Client info
			"remote_addr": r.RemoteAddr,
			"user_agent": r.UserAgent(),
		}

		// Log at appropriate level based on status code
		switch {
		case responseWrapper.Status() >= 500:
			l.Error("http request", requestFields)
		case responseWrapper.Status() >= 400:
			l.Warn("http request", requestFields)
		default:
			l.Info("http request", requestFields)
		}
	})
}

// Cleanup releases resources
func (l *Logger) Cleanup() error {
	return l.zap.Sync()
}

// Checks log level + colorizes the message
func (c customConsoleEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	if entry.Level == zapcore.ErrorLevel {
		entry.Message = yellow + entry.Message + reset
	}
	if entry.Level == zapcore.InfoLevel {
		entry.Message = green + entry.Message + reset
	}
	if entry.Level == zapcore.DebugLevel {
		entry.Message = cyan + entry.Message + reset
	}

	return c.Encoder.EncodeEntry(entry, fields)
}
