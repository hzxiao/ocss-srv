package api

import (
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/kataras/iris/context"
)

func GetAllDept(ctx context.Context) {
	deptList, err := db.FindAllDept()
	if err != nil {
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"deptList": deptList,
	})
}

func CallGetAllDept() (goutil.Map, error) {
	result, err := tools.HttpGet(fmt.Sprintf("http://%v/depts/", SrvAddr), "")
	if err != nil {
		return result, err
	}

	return handleACallResult(result)
}

func GetAllMajor(ctx context.Context) {
	argMap, err := CheckURLArg(ctx.FormValues(), []*Arg{
		{Key: "deptId", Type: "string"},
	})
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	majorList, err := db.FindAllMajor(argMap.GetString("deptId"))
	if err != nil {
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"majorList": majorList,
	})
}

func CallGetAllMajor() (goutil.Map, error) {
	result, err := tools.HttpGet(fmt.Sprintf("http://%v/majors/", SrvAddr), "")
	if err != nil {
		return result, err
	}

	return handleACallResult(result)
}
