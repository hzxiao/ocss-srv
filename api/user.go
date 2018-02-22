package api

import (
	"github.com/hzxiao/goutil/util"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/kataras/iris"
	"log"
)

func AddUser(ctx iris.Context) {
	user := &db.User{Username: "aaa"}
	err := db.AddUser(user)
	if err != nil {
		log.Printf("[AddUser] add user(%v) error(%v)", util.Struct2Json(user), err)
		ResultErrByKey(ctx, 3, "srv-err", err)
		return
	}

	ResultSuccess(ctx, util.Map{
		"user": user,
	})
}
