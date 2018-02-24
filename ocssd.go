package main

import (
	"github.com/hzxiao/goutil/util"
	"github.com/hzxiao/ocss-srv/api"
	"github.com/hzxiao/ocss-srv/config"
	"github.com/hzxiao/ocss-srv/db"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var err error
	//config
	err = config.InitConfig("config", "./config")
	if err != nil {
		panic(err)
	}
	config.PrintAll()
	//db
	err = db.InitDB(config.GetString("db.conn"), config.GetString("db.name"))
	if err != nil {
		panic(err)
	}
	//api
	app := iris.New()
	app.WrapRouter(cors.WrapNext(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}))

	api.RegisterHandle(app)
	app.Get("/demo", func(ctx iris.Context) {
		res := util.Map{
			"key": "value",
		}
		ctx.JSON(res)
	})
	app.Run(iris.Addr(config.GetString("server.port")), iris.WithCharset("UTF-8"))
}
