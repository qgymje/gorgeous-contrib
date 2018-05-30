package logger

import (
	"git.verystar.cn/GaomingQian/gorgeous/provider"
)

type MultiLogger struct {
	loggers []provider.ILogger
}

func NewMultiLogger(loggers ...provider.ILogger) provider.ILogger {
	ml := new(MultiLogger)
	ml.loggers = loggers
	return ml
}

func (ml *MultiLogger) Debugf(format string, args ...interface{}) {
	for _, l := range ml.loggers {
		l.Debugf(format, args...)
	}
}

func (ml *MultiLogger) Infof(format string, args ...interface{}) {
	for _, l := range ml.loggers {
		l.Infof(format, args...)
	}
}

func (ml *MultiLogger) Warnf(format string, args ...interface{}) {
	for _, l := range ml.loggers {
		l.Warnf(format, args...)
	}
}

func (ml *MultiLogger) Errorf(format string, args ...interface{}) {
	for _, l := range ml.loggers {
		l.Errorf(format, args...)
	}
}
