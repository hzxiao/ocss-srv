package api

import (
	"bytes"
	"fmt"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/kataras/iris/context"
	log "github.com/sirupsen/logrus"
)

func Login(ctx context.Context) {
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

	return handleACallResult(result)
}

func AddUser(ctx context.Context) {
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

func UpdateUser(ctx context.Context) {
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

func GetUser(ctx context.Context) {
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

func UpdateUserPassword(ctx context.Context)  {
	var data goutil.Map
	err := ctx.ReadJSON(&data)
	if err != nil {
		WriteResultWithArgErr(ctx, err)
		return
	}

	user, err := db.VerifyUser(data.GetString("username"), data.GetString("oldPass"))
	if err != nil {
		log.Error("[UpdateUserPassword] ", err)
		WriteResultWithSrvErr(ctx, err)
		return
	}
	if user == nil {
		WriteResultErrByMsg(ctx, CodeUserNotFound, "原密码错误", nil)
		return
	}

	updateUser := &db.User{
		Username: data.GetString("username"),
		Password: data.GetString("newPass"),
	}

	err = db.UpdateUser(updateUser)
	if err != nil {
		log.Printf("[UpdateUserPassword] update user(%v) error(%v)", goutil.Struct2Json(updateUser), err)
		WriteResultErrByKey(ctx, 3, "srv-err", err)
		return
	}

	WriteResultSuccess(ctx, "OK")
}