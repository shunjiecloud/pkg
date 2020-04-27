package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shunjiecloud/errors"
	ec "github.com/shunjiecloud/pkg/errcode"
	"gopkg.in/go-playground/validator.v9"
)

//  OkResponse
func (appCtx *AppContext) OkResponse() {
	ok := struct {
		Result string `json:"result"`
	}{
		Result: "ok",
	}
	appCtx.GinCtx.JSON(200, ok)
}

// Response Json
func (appCtx *AppContext) WriteJson(httpCode int, data interface{}) {
	appCtx.GinCtx.JSON(httpCode, data)
}

//  ErrorResponse
type ErrorResponse struct {
	Result string `json:"result"`
	Msg    string `json:"msg,omitempty"`
}

func (appCtx *AppContext) WriteError(in error) {
	in = errors.Adapt(in)
	var httpCode int
	httpCode = http.StatusInternalServerError
	var resp ErrorResponse
	//  是否为*errors.Error
	err, isOk := in.(*errors.Error)
	if isOk == true {
		innerError := err.Err
		switch innerError.(type) {
		case ec.Error:
			//  内部含有code.Error
			resp.Msg = innerError.(ec.Error).Message()
			httpCode = innerError.(ec.Error).HttpCode
		case validator.ValidationErrors:
			httpCode = http.StatusBadRequest
		default:
		}
		//  externs
		externs := make([]string, 0)
		for _, f := range err.Fields {
			externs = append(externs, fmt.Sprintf("%v:%v", f.Key, f.String))
		}
		//  pos
		pos := make([]string, 0)
		for _, f := range err.StackFrames() {
			pos = append(pos, fmt.Sprintf("%v %v", f.File, f.LineNumber))
		}
		appCtx.GinCtx.Error(&gin.Error{
			Err:  innerError,
			Type: gin.ErrorTypePrivate,
			Meta: gin.H{
				"pos":    pos,
				"extern": externs,
			},
		})
	}
	resp.Result = in.Error()
	appCtx.GinCtx.JSON(httpCode, &resp)
}
