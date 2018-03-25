package api

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/kataras/iris/context"
	log "github.com/sirupsen/logrus"
	"fmt"
)

func ListNotice(ctx context.Context) {
	argMap, err := CheckURLArg(ctx.FormValues(), []*Arg{
		{Key: "status", Type: "int"},
		{Key: "page", Type: "int", DefaultValue: "0"},
		{Key: "pageSize", Type: "int", DefaultValue: "200"},
		{Key: "sort", Type: "string"},
	})
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}
	cond := TakeByKeys(argMap, "status")
	cond["uid"] = ctx.Values().Get("uid")

	limit := int(argMap.GetInt64("pageSize"))
	skip := int(argMap.GetInt64("page")) * limit
	var sort []string
	if argMap.Exist("sort") {
		sort = append(sort, argMap.GetString("sort"))
	}

	noticeList, total, err := db.ListNotice(cond, sort, skip, limit)
	if err != nil {
		log.Errorf("[ListNotice] get by(%v) error(%v)", argMap, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	WriteResultSuccess(ctx, goutil.Map{
		"noticeList": noticeList,
		"total":      total,
	})
}

func UpdateNotice(ctx context.Context) {
	id := ctx.Params().Get("id")
	var notice db.Notice
	err := ctx.ReadJSON(&notice)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	notice.ID = id
	err = db.UpdateNotice(&notice)
	if err != nil {
		log.Errorf("[UpdateNotice] update by(%v) error(%v)", notice, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	WriteResultSuccess(ctx, "OK")
}

func CountNotice(ctx context.Context) {
	uid := fmt.Sprintf("%v", ctx.Values().Get("uid"))
	res, err := db.CountNoticeDiffStatus(uid)
	if err != nil {
		log.Errorf("[CountNotice] get by(%v) error(%v)", uid, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	WriteResultSuccess(ctx, goutil.Map{
		"result": res,
	})
}
