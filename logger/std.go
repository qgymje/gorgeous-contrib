package logger

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
)

type StdLogger struct {
}

func NewStdLogger() *StdLogger {
	log.SetFlags(log.Ltime | log.Lshortfile)
	return &StdLogger{}
}

func (l *StdLogger) Debugf(format string, args ...interface{}) {
	log.Printf("[Debug]"+format, args...)
}

func (l *StdLogger) Infof(format string, args ...interface{}) {
	log.Printf("[Info]"+format, args...)
}

func (l *StdLogger) Warnf(format string, args ...interface{}) {
	log.Printf("[Warn]"+format, args...)
}

func (l *StdLogger) Errorf(format string, args ...interface{}) {
	err := fmt.Errorf("[Error]"+format, args...)
	log.Printf("%+v", errors.Wrap(err, ""))
}
