package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Init init
func Init(l *zap.Logger) {
	logger = l
}

// Info ...
func Info(msg string, fileds ...zapcore.Field) {
	logger.Info(msg, fileds...)
}

// Fatal ...
func Fatal(msg string, fileds ...zapcore.Field) {
	logger.Fatal(msg, fileds...)
}

// Warn ...
func Warn(msg string, fileds ...zapcore.Field) {
	logger.Warn(msg, fileds...)
}

// Debug ...
func Debug(msg string, fileds ...zapcore.Field) {
	logger.Debug(msg, fileds...)
}

// Error ...
func Error(msg string, fileds ...zapcore.Field) {
	logger.Error(msg, fileds...)
}
