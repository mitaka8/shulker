package matrix

import (
	log "github.com/sirupsen/logrus"
	"maunium.net/go/mautrix/crypto"
)

type logger struct{}

var _ crypto.Logger = &logger{}

func (f logger) Error(message string, args ...interface{}) {
	log.Errorf(message, args...)
}

func (f logger) Warn(message string, args ...interface{}) {
	log.Warnf(message, args...)
}

func (f logger) Debug(message string, args ...interface{}) {
	log.Debugf(message, args...)
}

func (f logger) Trace(message string, args ...interface{}) {
	log.Tracef(message, args...)
}
