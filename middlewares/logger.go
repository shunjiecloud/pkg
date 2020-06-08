package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	merr "github.com/micro/go-micro/v2/errors"
	"github.com/shunjiecloud/pkg/log"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		fields := make([]zap.Field, 0)

		//  Stop timer
		end := time.Now()
		fields = append(fields, zap.String("time", end.Format(time.RFC1123)))
		fields = append(fields, zap.String("latency", fmt.Sprintf("%s", end.Sub(start))))

		fields = append(fields, zap.String("client_ip", c.ClientIP()))
		fields = append(fields, zap.String("method", c.Request.Method))
		fields = append(fields, zap.Int("status", c.Writer.Status()))
		fields = append(fields, zap.String("user_agent", c.Request.UserAgent()))

		if raw != "" {
			path = path + "?" + raw
		}
		fields = append(fields, zap.String("path", path))

		//  errors
		errs := c.Errors.ByType(gin.ErrorTypePrivate)
		if len(errs) == 0 {
			//  log info
			log.Info("call request", fields...)
			return
		}

		for _, e := range errs {
			result := ""
			msg := ""
			fieldsEachErr := fields
			var httpCode int
			innErr := e.Err
			switch innErr.(type) {
			case *merr.Error:
				//  innErr为*merr.Error
				result = innErr.(*merr.Error).Id
				fieldsEachErr = append(fieldsEachErr, zap.String("detail", innErr.(*merr.Error).Detail))
				httpCode = int(innErr.(*merr.Error).Code)
			case validator.ValidationErrors:
				//  innErr为validator错误
				result = "request invalid"
				msg = e.Error()
				httpCode = http.StatusBadRequest
			default:
				result = "Internal Server Error"
				msg = e.Error()
				httpCode = http.StatusInternalServerError
			}

			//  meta
			metas := e.Meta.(gin.H)
			for k, v := range metas {
				fieldsEachErr = append(fieldsEachErr, zap.String(k, v))
			}

			//  log error
			if httpCode >= 500 {
				// error
				log.ErrorString(result, fieldsEachErr...)
			} else {
				// warn
				log.ErrorString(result, fieldsEachErr...)
			}
		}

	}
}
