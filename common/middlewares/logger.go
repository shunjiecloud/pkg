package middlewares

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	f := func(param gin.LogFormatterParams) string {
		// your custom format
		log := gin.H{}
		log["client_ip"] = param.ClientIP
		log["time"] = param.TimeStamp.Format(time.RFC1123)
		log["method"] = param.Method
		log["path"] = param.Path
		log["status"] = param.StatusCode
		log["latency"] = fmt.Sprintf("%s", param.Latency)
		log["user_agent"] = param.Request.UserAgent()
		log["msg"] = param.ErrorMessage
		l, err := json.Marshal(log)
		if err != nil {
			return ""
		}
		return string(l) + "\n"
	}
	return gin.LoggerWithConfig(gin.LoggerConfig{
		Formatter: f,
	})
}
