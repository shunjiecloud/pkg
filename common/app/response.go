package app

import (
	"net/http"

	ec "github.com/shunjieyun/account/errcode"
	"github.com/shunjieyun/errors"
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

//  Response Error
func (appCtx *AppContext) WriteError(in error) {
	var httpCode int
	httpCode = http.StatusInternalServerError
	var resp ErrorResponse
	//  是否为*errors.Error
	err, isOk := in.(*errors.Error)
	if isOk == true {
		innerErr, isOk := err.Err.(ec.Error)
		if isOk == true {
			//  内部含有ec.Error
			resp.Result = innerErr.Result
			resp.Msg = innerErr.Message()
			httpCode = innerErr.HttpCode
		} else {
			//  内部无ec.Error
			resp.Result = err.Err.Error()
		}
	} else {
		resp.Result = in.Error()
	}
	appCtx.GinCtx.JSON(httpCode, &resp)
}
