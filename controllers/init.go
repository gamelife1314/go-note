package controllers

import (
	"github.com/gamelife1314/go-note/common"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func init() {

	mvc.New(common.App.Party("/user").AllowMethods(iris.MethodOptions)).Handle(new(UserController))
	mvc.New(common.App.Party("/upload").AllowMethods(iris.MethodOptions)).Handle(new(UploadController))
	mvc.New(common.App.Party("/site").AllowMethods(iris.MethodOptions)).Handle(new(SiteController))
	mvc.New(common.App.Party("/topic").AllowMethods(iris.MethodOptions)).Handle(new(TopicController))
}
