package log

import (
	"github.com/shunjiecloud/errors"
	"go.uber.org/zap"
)

var defaultLogger *zap.Logger

func init() {
	var err error
	defaultLogger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func Warn(err *errors.Error) {
	defaultLogger.Warn(err.Error(), err.Fields...)
}

func Error(err *errors.Error) {
	defaultLogger.Error(err.Error(), err.Fields...)
}

func Info(err *errors.Error) {
	defaultLogger.Info(err.Error(), err.Fields...)
}
