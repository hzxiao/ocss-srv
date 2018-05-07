package api

import (
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/kataras/iris/context"
	log "github.com/sirupsen/logrus"
	"strconv"
	"github.com/kataras/iris/core/errors"
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

	if commentList == nil || len(commentList) == 0 {
		WriteResultSuccess(ctx, goutil.Map{
			"commentList": nil,
			"total":       0,
		})
	}

	var cmtList []goutil.Map

	for i := range commentList {
		cmt := goutil.Struct2Map(commentList[i])
		uid := commentList[i].UID
		role := commentList[i].Role
		userInfo,err := GetUserInfoForComment(goutil.Map{"role":role,"uid":uid})
		if err != nil {
			log.Errorf("[ListComment] get by(%v) error(%v)", argMap, err)
			WriteResultWithSrvErr(ctx, err)
			return
		}
		cmt.Set("role", userInfo.GetString("role"))
		cmt.Set("name", userInfo.GetString("name"))
		cmt.Set("childTotal", 0)
		if commentList[i].Children != nil && len(commentList[i].Children) > 0{
			children := commentList[i].Children
			for j := range children {
				uid := children[j].GetString("uid")
				role := children[j].GetInt64("role")
				userInfo,err := GetUserInfoForComment(goutil.Map{"role":role,"uid":uid})
				if err != nil {
					log.Errorf("[ListComment] get by(%v) error(%v)", argMap, err)
					WriteResultWithSrvErr(ctx, err)
					return
				}
				children[j].Set("role", userInfo.GetString("role"))
				children[j].Set("name", userInfo.GetString("name"))
			}
			cmt.Set("children", children)
			cmt.Set("childTotal", len(children))
		}
		cmtList = append(cmtList,cmt)

	}
	log.Printf("[ListComment] by argMap(%)  get cmtList(%v) ", argMap, cmtList)

	WriteResultSuccess(ctx, goutil.Map{
		"commentList": cmtList,
		"total":       total,
	})
}
// 根据角色获取用户信息
func GetUserInfoForComment(userParam goutil.Map)(goutil.Map, error)  {
	userInfo:=goutil.Map{}
	uid := userParam.GetString("uid")
	switch userParam.GetInt64("role") {
	case 1:
		userInfo.Set("role", AdminCN)
		userInfo.Set("name", AdminCN)
		break
	case 2:
		teacher,err := db.LoadTeacher(uid)
		if err != nil {
			log.Errorf("[GetUserInfoForComment] get by(%v) error(%v)", userParam, err)
			return nil,err
		}
		if teacher != nil {
			userInfo.Set("name", teacher.Name)
		}
		userInfo.Set("role", TeacherCN)
		break
	case 3:
		student,err := db.LoadStudent(uid)
		if err != nil {
			log.Errorf("[GetUserInfoForComment] get by(%v) error(%v)", userParam, err)
			return nil,err
		}
		if student != nil {
			userInfo.Set("name", student.Name)
		}
		userInfo.Set("role", StudentCN)
		break
	}
	return userInfo,nil
}
func AddChildComment(ctx context.Context) {
	id := ctx.Params().Get("id")
	var child goutil.Map
	err := ctx.ReadJSON(&child)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}
	if len(child) == 0 {
		log.Errorf("[AddChildComment] add child comment(%v) child(%v) error(%v)", id, goutil.Struct2Json(child), err)
		WriteResultWithSrvErr(ctx, errors.New("child is empty"))
		return
	}
	child.Set("uid", ctx.Values().GetString("uid"))
	role, err := ctx.Values().GetFloat64("role")
	if err != nil {
		log.Errorf("[AddComment] get role uid(%v) error(%v)", ctx.Values().GetString("uid"), err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	child.Set("role", int(role))
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
