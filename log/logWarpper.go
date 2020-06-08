package log

import (
	"context"

	merr "github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/server"
	"github.com/shunjiecloud/errors"
	"go.uber.org/zap"
)

// logWrapper is a handler wrapper
func logWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		Info("call endpoint", zap.String("endpoint", req.Endpoint()))
		err := fn(ctx, req, rsp)
		if err != nil {
			//  先转化为sjyErr，尽量外面传进来的均为sjyErr，这样callStack可以保留。
			err = errors.Adapt(err)
			sjyErr := err.(*errors.Error)
			//  添加pos field
			pos(sjyErr)
			//  判断internel error是否为merr
			mErr, isOk := sjyErr.Err.(*merr.Error)
			if isOk == true {
				//  是merr，则根据code区分error和warn级别
				fields := make([]zap.Field, 0)
				fields = append(fields, zap.Int32("code", mErr.Code))
				fields = append(fields, zap.String("detail", mErr.Detail))
				fields = append(fields, sjyErr.Fields...)
				if mErr.Code >= 500 {
					// error
					defaultLogger.Error(mErr.Id, fields...)
				} else {
					// warn
					defaultLogger.Warn(mErr.Id, fields...)
				}
			} else {
				//  不是merr，归类为错误
				defaultLogger.Error(sjyErr.Error(), sjyErr.Fields...)
			}
		}
		return err
	}
}
