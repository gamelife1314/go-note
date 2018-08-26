package common

import (
	"github.com/betacraft/yaag/irisyaag"
	"github.com/betacraft/yaag/yaag"
	"github.com/gamelife1314/go-note/config"
	"github.com/kataras/golog"
	"github.com/kataras/iris"
	"os"
	"time"
)

var App *iris.Application
var Logger *golog.Logger
var TimeZone, _ = time.LoadLocation(config.Configuration.Other["TimeZone"].(string))

func init() {
	App = iris.New()
	Logger = App.Logger()
	Logger.SetLevel("debug")
	Logger.SetOutput(os.Stdout)

	yaag.Init(&yaag.Config{
		On:       true,
		DocTitle: "一刻",
		DocPath:  "./doc/api.html",
		BaseUrls: map[string]string{"Production": "", "Staging": ""},
	})
	App.Use(irisyaag.New())

	App.StaticWeb("/doc", "./doc")
	App.StaticWeb("/uploads", "./uploads")

	App.UseGlobal(func(context iris.Context) {
		context.Header("Access-Control-Allow-Origin", context.Host())
		context.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, PATCH, DELETE")
		context.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, x-token")
		context.Header("Access-Control-Allow-Credentials", "true")
		context.Next()
	})
}
