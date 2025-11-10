package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger interface for application logging
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	With(fields ...Field) Logger
	WithContext(ctx context.Context) Logger

	Sync() error
}

// Field represents a log field
type Field = zapcore.Field

// Common field constructors
var (
	String   = zap.String
	Int      = zap.Int
	Int64    = zap.Int64
	Float64  = zap.Float64
	Bool     = zap.Bool
	Duration = zap.Duration
	Time     = zap.Time
	ZError   = zap.Error
	Any      = zap.Any
	Strings  = zap.Strings
	Uint64   = zap.Uint64
)

// zapLogger wraps zap.Logger to implement our Logger interface
type zapLogger struct {
	logger *zap.Logger
	name   string
}

// Config for logger initialization
type Config struct {
	Level      string
	LogDir     string
	MaxSize    int // MB
	MaxBackups int
	MaxAge     int // days
	Compress   bool
	Console    bool
}

// DefaultConfig returns default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Level:      "info",
		LogDir:     "./logs",
		MaxSize:    100,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
		Console:    true,
	}
}

// New creates a new logger instance
func New(name string, cfg *Config) (Logger, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// Parse log level
	level, err := zapcore.ParseLevel(cfg.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}

	// Ensure log directory exists
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Create encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "",
		LevelKey:       "level",
		NameKey:        "",
		CallerKey:      "",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create cores for different outputs
	var cores []zapcore.Core

	// File output with rotation
	logFile := filepath.Join(cfg.LogDir, fmt.Sprintf("%s.log", time.Now().Format("2006-01-02")))
	fileWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}
	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	cores = append(cores, zapcore.NewCore(fileEncoder, zapcore.AddSync(fileWriter), level))

	// Console output (if enabled)
	if cfg.Console {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
	}

	// Create logger
	core := zapcore.NewTee(cores...)
	zapLog := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	if name != "" {
		zapLog = zapLog.Named(name)
	}

	return &zapLogger{
		logger: zapLog,
		name:   name,
	}, nil
}

// Debug logs a debug message
func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info message
func (l *zapLogger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

// Warn logs a warning message
func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

// Error logs an error message
func (l *zapLogger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

// Fatal logs a fatal message and exits
func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

// With creates a child logger with additional fields
func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		logger: l.logger.With(fields...),
		name:   l.name,
	}
}

// WithContext creates a logger with context information
func (l *zapLogger) WithContext(ctx context.Context) Logger {
	// Extract common context values
	fields := []Field{}

	if requestID := ctx.Value("request_id"); requestID != nil {
		fields = append(fields, String("request_id", fmt.Sprint(requestID)))
	}

	if userID := ctx.Value("user_id"); userID != nil {
		fields = append(fields, String("user_id", fmt.Sprint(userID)))
	}

	return l.With(fields...)
}

// Sync flushes any buffered log entries
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

// Global logger instance
var globalLogger Logger

// InitGlobal initializes the global logger
func InitGlobal(name string, cfg *Config) error {
	logger, err := New(name, cfg)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// Global returns the global logger instance
func Global() Logger {
	if globalLogger == nil {
		// Create default logger if not initialized
		logger, _ := New("default", DefaultConfig())
		globalLogger = logger
	}
	return globalLogger
}

// Helper functions for global logger
func Debug(msg string, fields ...Field) {
	Global().Debug(msg, fields...)
}

func Info(msg string, fields ...Field) {
	Global().Info(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	Global().Warn(msg, fields...)
}

func Error(msg string, fields ...Field) {
	Global().Error(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	Global().Fatal(msg, fields...)
}

func With(fields ...Field) Logger {
	return Global().With(fields...)
}

func Sync() error {
	return Global().Sync()
}
