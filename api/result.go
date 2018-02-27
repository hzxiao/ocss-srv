package api

import "github.com/kataras/iris"

type Result struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
	Err  string      `json:"err"`
}

func WriteResultSuccess(ctx iris.Context, data interface{}) {
	ctx.JSON(&Result{
		Code: CodeSuccess,
		Data: data,
	})
	ctx.StopExecution()
}

func WriteResultErrByKey(ctx iris.Context, code int, key string, err error) {
	//get target value by key
	result := &Result{
		Code: code,
		Msg:  key,
	}
	if err != nil {
		result.Err = err.Error()
	}
	writeResult(ctx, result)
}

func WriteResultErrByMsg(ctx iris.Context, code int, msg string, err error) {
	result := &Result{
		Code: code,
		Msg:  msg,
	}
	if err != nil {
		result.Err = err.Error()
	}
	writeResult(ctx, result)
}

func WriteResultWithArgErr(ctx iris.Context, err error) {
	result := &Result{
		Code: CodeArgErr,
		Msg:  "参数错误",
	}
	if err != nil {
		result.Err = err.Error()
	}
	writeResult(ctx, result)
}

func WriteResultWithSrvErr(ctx iris.Context, err error) {
	result := &Result{
		Code: CodeSrvErr,
		Msg:  "服务器错误",
	}
	if err != nil {
		result.Err = err.Error()
	}
	writeResult(ctx, result)
}

func writeResult(ctx iris.Context, result *Result) {
	ctx.JSON(result)
	ctx.StopExecution()
}
