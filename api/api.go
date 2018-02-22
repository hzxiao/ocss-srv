package api

import (
	"github.com/kataras/iris"
)

func RegisterHandle(app *iris.Application) {
	//user
	app.Post("users", AddUser)
}

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	Err  string      `json:"err"`
}

func ResultSuccess(ctx iris.Context, data interface{}) {
	ctx.JSON(&Result{
		Code: 0,
		Data: data,
	})
}

func ResultErrByKey(ctx iris.Context, code int, key string, err error) {
	//get target value by key
	ctx.JSON(&Result{
		Code: code,
		Msg:  key,
		Err:  err.Error(),
	})
}

func ResultErrByMsg(ctx iris.Context, code int, msg string, err error) {
	ctx.JSON(&Result{
		Code: code,
		Msg:  msg,
		Err:  err.Error(),
	})
}
