package app

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	merr "github.com/micro/go-micro/v2/errors"
	"github.com/shunjiecloud/errors"
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
	var httpCode int
	var resp ErrorResponse

	in = errors.Adapt(in)
	err := in.(*errors.Error)
	innerError := err.Err
	switch innerError.(type) {
	case *merr.Error:
		//  内部含有*merr.Error
		resp.Result = innerError.(*merr.Error).Id
		resp.Msg = innerError.(*merr.Error).Detail
		httpCode = int(innerError.(*merr.Error).Code)
	case validator.ValidationErrors:
		//  内部含有validator错误
		resp.Result = "request invalid"
		resp.Msg = innerError.Error()
		httpCode = http.StatusBadRequest
	default:
		resp.Result = "Internal Server Error"
		resp.Msg = innerError.Error()
		httpCode = http.StatusInternalServerError
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
	//  保存错误到gin errors
	appCtx.GinCtx.Error(&gin.Error{
		Err:  innerError,
		Type: gin.ErrorTypePrivate,
		Meta: gin.H{
			"pos":    pos,
			"extern": externs,
		},
	})
	appCtx.GinCtx.JSON(httpCode, &resp)
}
