package logger

import (
	"fmt"

	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
)

type SentryLogger struct {
	client *raven.Client
	tag    map[string]string
}

func NewSentryLogger(appName string, dsn string) *SentryLogger {
	l := new(SentryLogger)
	l.client = raven.DefaultClient
	l.client.SetDSN(dsn)
	l.tag = make(map[string]string)

	l.tag["app"] = appName
	return l
}

func (l *SentryLogger) Errorf(format string, args ...interface{}) {
	l.tag["err_level"] = "error"
	err := fmt.Errorf(format, args...)
	err = errors.Wrap(err, "")
	l.client.CaptureMessage(fmt.Sprintf("%+v", err), l.tag)
}

func (l *SentryLogger) Warnf(format string, args ...interface{}) {
	l.tag["err_level"] = "fatal"
	l.client.CaptureMessage(fmt.Sprintf(format, args...), l.tag)
}

func (l *SentryLogger) Infof(format string, args ...interface{}) {
}

func (l *SentryLogger) Debugf(format string, args ...interface{}) {
}
