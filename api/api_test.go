package api

import (
	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/hzxiao/ocss-srv/tools"
	"github.com/kataras/iris"
)

var testApp *iris.Application

func init() {
	err := db.InitDB("111.230.242.177:27017", "ocss_test")
	if err != nil {
		panic(err)
	}
	SrvAddr = "127.0.0.1:9123"
	testApp = iris.New()
	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
		On:       true,
		DocTitle: "Iris",
		DocPath:  "apidoc.html",
		BaseUrls: map[string]string{"Production": "127.0.0.1:8999", "Staging": ""},
	})
	testApp.Use(irisyaag.New()) // <- IMPORTANT, register the middleware.

	RegisterHandle(testApp)

	go func() {
		err = testApp.Run(iris.Addr(SrvAddr))
		if err != nil {
			panic(err)
		}
	}()
}

func createUserAndLogin(role int) (username, token string, err error) {
	username = tools.GenerateUniqueId()

	user := &db.User{
		Username: username,
		Role:     role,
	}
	err = db.AddUser(user)
	if err != nil {
		return
	}

	result, err := CallLogin(goutil.Map{"username": username, "password": db.DefaultPassword})
	if err != nil {
		return
	}

	token = result.GetStringP("user/token")
	return
}

func removeAll() {
	db.C(db.CollectionUser).RemoveAll(nil)
}
