package logger

import (
	"github.com/sirupsen/logrus"
)

// InitLogger init logger
func InitLogger() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000000",
	})
}

// Infof log info
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Info log info
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Errorf log error
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Error log error
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Debugf log debug
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Debug log debug
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Warnf log warn
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

// Warn log warn
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}
