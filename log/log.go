package log

import (
	"fmt"

	"github.com/shunjiecloud/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var defaultLogger *zap.Logger

func init() {
	var err error
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	defaultLogger, err = config.Build()
	if err != nil {
		panic(err)
	}
}

func pos(err *errors.Error) {
	//  add pos into fields
	pos := make([]string, 0)
	for _, f := range err.StackFrames() {
		pos = append(pos, fmt.Sprintf("%v %v", f.File, f.LineNumber))
	}
	err.Fields = append(err.Fields, zap.Strings("pos", pos))
}

func Warn(err *errors.Error) {
	//  pos
	pos(err)
	defaultLogger.Warn(err.Error(), err.Fields...)
}

func Error(err *errors.Error) {
	//  pos
	pos(err)
	defaultLogger.Error(err.Error(), err.Fields...)
}

func Info(err *errors.Error) {
	//  pos
	pos(err)
	defaultLogger.Info(err.Error(), err.Fields...)
}
