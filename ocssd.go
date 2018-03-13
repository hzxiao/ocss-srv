package main

import (
	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
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

	yaag.Init(&yaag.Config{ // <- IMPORTANT, init the middleware.
		On:       true,
		DocTitle: "Iris",
		DocPath:  "apidoc.html",
		BaseUrls: map[string]string{"Production": "127.0.0.1:8999", "Staging": ""},
	})
	app.Use(irisyaag.New()) // <- IMPORTANT, register the middleware.

	Cors(app)

	api.RegisterHandle(app)

	app.StaticWeb("/", config.GetString("file.webPath"))
	app.Run(iris.Addr(config.GetString("server.port")), iris.WithCharset("UTF-8"))
}

func Cors(app *iris.Application) {
	app.WrapRouter(cors.WrapNext(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "DELETE", "PUT"},
		AllowedHeaders:   []string{"x-requested-with", "Content-Type", "origin", "authorization", "accept", "client-security-token"},
	}))
}
