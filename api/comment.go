package api

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/kataras/iris/context"
	log "github.com/sirupsen/logrus"
	"strconv"
)

func AddComment(ctx context.Context) {
	var c db.Comment
	err := ctx.ReadJSON(&c)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	c.UID = ctx.Values().GetString("uid")
	role, err := ctx.Values().GetFloat64("role")
	if err != nil {
		log.Errorf("[AddComment] get role uid(%v) error(%v)", c.UID, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	c.Role = int(role)
	err = db.AddComment(&c)
	if err != nil {
		log.Errorf("[AddComment] add comment(%v) error(%v)", goutil.Struct2Json(c), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"comment": c,
	})
}

func DelComment(ctx context.Context) {
	id := ctx.Params().Get("id")
	err := db.DelComment([]string{id})
	if err != nil {
		log.Errorf("[DelComment] del comment(%v) error(%v)", id, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"id": id,
	})
}

func ListComment(ctx context.Context) {
	argMap, err := CheckURLArg(ctx.FormValues(), []*Arg{
		{Key: "tcid", Type: "string"},
		{Key: "status", Type: "int", DefaultValue: strconv.Itoa(db.StatusNormal)},
		{Key: "page", Type: "int", DefaultValue: "1"},
		{Key: "pageSize", Type: "int", DefaultValue: "20"},
		{Key: "sort", Type: "string"},
	})
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}
	cond := TakeByKeys(argMap, "tcid", "status")
	limit := int(argMap.GetInt64("pageSize"))
	skip := int(argMap.GetInt64("page")) * limit
	var sort []string
	if argMap.Exist("sort") {
		sort = append(sort, argMap.GetString("sort"))
	}

	commentList, total, err := db.ListComment(cond, sort, skip, limit)
	if err != nil {
		log.Errorf("[ListComment] get by(%v) error(%v)", argMap, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	WriteResultSuccess(ctx, goutil.Map{
		"commentList": commentList,
		"total":       total,
	})
}

func AddChildComment(ctx context.Context) {
	id := ctx.Params().Get("id")
	var child goutil.Map
	err := ctx.ReadJSON(&child)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	c, err := db.UpdateChildComment(id, "add", child)
	if err != nil {
		log.Errorf("[AddChildComment] add child comment(%v) child(%v) error(%v)", id, goutil.Struct2Json(child), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"comment": c,
	})
}

func DelChildComment(ctx context.Context) {
	id := ctx.Params().Get("id")
	childId := ctx.Params().Get("childId")

	c, err := db.UpdateChildComment(id, "del", goutil.Map{"id": childId})
	if err != nil {
		log.Errorf("[DelChildComment] del child comment(%v) child(%v) error(%v)", id, childId, err)
		WriteResultWithSrvErr(ctx, err)
		return
	}

	WriteResultSuccess(ctx, goutil.Map{
		"comment": c,
	})
}
