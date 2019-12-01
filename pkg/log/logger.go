package log

import (
	"github.com/sirupsen/logrus"
)

// Logger is abstraction of logger library.
type Logger struct {
	impl          *logrus.Logger
	defaultFields Fields
}

// Fields of log message
type Fields map[string]interface{}

func NewWithFields(fields Fields) *Logger {
	impl := logrus.New()
	impl.SetFormatter(&logrus.JSONFormatter{})
	return &Logger{
		impl:          impl,
		defaultFields: fields,
	}
}

// New creates logger.
func New() *Logger {
	return NewWithFields(Fields{})
}

func (l *Logger) Panicf(msg string, fields Fields) {
	l.impl.WithFields(logrus.Fields(l.defaultFields)).WithFields(logrus.Fields(fields)).Panic(msg)
}

func (l *Logger) Panic(v ...interface{}) {
	l.impl.Panic(v...)
}

func (l *Logger) Fatalf(msg string, fields Fields) {
	l.impl.WithFields(logrus.Fields(l.defaultFields)).WithFields(logrus.Fields(fields)).Fatal(msg)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.impl.Fatal(v...)
}

func (l *Logger) Errorf(msg string, fields Fields) {
	l.impl.WithFields(logrus.Fields(l.defaultFields)).WithFields(logrus.Fields(fields)).Error(msg)
}

func (l *Logger) Error(v ...interface{}) {
	l.impl.Error(v...)
}

func (l *Logger) Warnf(msg string, fields Fields) {
	l.impl.WithFields(logrus.Fields(l.defaultFields)).WithFields(logrus.Fields(fields)).Warn(msg)
}

func (l *Logger) Warn(v ...interface{}) {
	l.impl.Warn(v...)
}

func (l *Logger) Infof(msg string, fields Fields) {
	l.impl.WithFields(logrus.Fields(l.defaultFields)).WithFields(logrus.Fields(fields)).Info(msg)
}

func (l *Logger) Info(v ...interface{}) {
	l.impl.Info(v...)
}

func (l *Logger) Debugf(msg string, fields Fields) {
	l.impl.WithFields(logrus.Fields(l.defaultFields)).WithFields(logrus.Fields(fields)).Debug(msg)
}

func (l *Logger) Debug(v ...interface{}) {
	l.impl.Debug(v...)
}
