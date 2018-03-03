package api

import (
	"github.com/kataras/iris/context"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/hzxiao/goutil"
	"fmt"
	"github.com/hzxiao/ocss-srv/tools"
)

func GetAllDept(ctx context.Context) {
	deptList, err := db.FindAllDept()
	if err != nil {
		WriteResultWithSrvErr(ctx,err)
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
	majorList, err := db.FindAllMajor()
	if err != nil {
		WriteResultWithSrvErr(ctx,err)
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