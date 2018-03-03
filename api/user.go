package api

import (
	"bytes"
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/juju/errors"
	"github.com/kataras/iris"
	log "github.com/sirupsen/logrus"
)

func Login(ctx iris.Context) {
	var info goutil.Map
	err := ctx.ReadJSON(&info)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	user, err := db.VerifyUser(info.GetString("username"), info.GetString("password"))
	if err != nil {
		log.Error("[Login] ", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	if user == nil {
		WriteResultErrByMsg(ctx, CodeUserNotFound, "用户名不存在或密码错误", nil)
		return
	}
	switch user.Status {
	case db.UserStatsForbid:
		WriteResultErrByMsg(ctx, CodeForbid, "禁止登录", nil)
	case db.UserStatsDelete:
		WriteResultErrByMsg(ctx, CodeDeleted, "无效的用户", nil)
	default:
		userMap := goutil.Struct2Map(user)
		userMap.Set("token", NewToken(user.Username, user.Role))
		WriteResultSuccess(ctx, goutil.Map{"user": userMap})
	}
}

func CallLogin(info goutil.Map) (goutil.Map, error) {
	result, err := tools.HttpPost(fmt.Sprintf("http://%v/login", SrvAddr), "",
		"application/json", bytes.NewBufferString(goutil.Struct2Json(info)))
	if err != nil {
		return result, err
	}

	if result.GetInt64("code") == 0 {
		return result.GetMap("data"), nil
	}

	return result, errors.New(result.GetString("msg") + " " + result.GetString("err"))
}

func AddUser(ctx iris.Context) {
	var user db.User
	err := ctx.ReadJSON(&user)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	err = db.AddUser(&user)
	if err != nil {
		log.Printf("[AddUser] add user(%v) error(%v)", goutil.Struct2Json(user), err)
		WriteResultErrByKey(ctx, 3, "srv-err", err)
		return
	}

	user.Password = ""
	WriteResultSuccess(ctx, goutil.Map{
		"user": user,
	})
}

func CallAddUser(token string, user *db.User) (goutil.Map, error) {
	result, err := tools.HttpPost(fmt.Sprintf("http://%v/users", SrvAddr), token,
		"application/json", bytes.NewBufferString(goutil.Struct2Json(user)))
	if err != nil {
		return result, err
	}

	return handleACallResult(result)
}

func UpdateUser(ctx iris.Context) {
	username := ctx.Params().Get("username")
	var user db.User
	err := ctx.ReadJSON(&user)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	user.Username = username
	err = db.UpdateUser(&user)
	if err != nil {
		log.Printf("[UpdateUser] add user(%v) error(%v)", goutil.Struct2Json(user), err)
		WriteResultErrByKey(ctx, 3, "srv-err", err)
		return
	}

	user.Password = ""
	WriteResultSuccess(ctx, goutil.Map{
		"user": user,
	})
}

func CallUpdateUser(token string, username string, user *db.User) (goutil.Map, error) {
	result, err := tools.HttpPut(fmt.Sprintf("http://%v/users/%v", SrvAddr, username), token,
		"application/json", bytes.NewBufferString(goutil.Struct2Json(user)))
	if err != nil {
		return result, err
	}
	return handleACallResult(result)
}

func GetUser(ctx iris.Context) {
	username := ctx.Params().Get("username")
	log.Println(username)
	user, err := db.FindUserByUsername(username)
	if err == nil {
		user.Password = ""
		WriteResultSuccess(ctx, goutil.Map{"user": user})
		return
	}
	if err == db.ErrNotFound {
		WriteResultErrByMsg(ctx, CodeUserNotFound, "用户名不存在", nil)
	} else {
		WriteResultErrByKey(ctx, 3, "srv-err", err)
	}
}

func CallGetUser(token string, username string) (goutil.Map, error) {
	result, err := tools.HttpGet(fmt.Sprintf("http://%v/users/%v", SrvAddr, username), token)
	if err != nil {
		return result, err
	}

	return handleACallResult(result)
}