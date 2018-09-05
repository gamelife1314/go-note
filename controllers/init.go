package controllers

import (
	"github.com/gamelife1314/go-note/common"
	"github.com/kataras/iris/mvc"
)

func init() {

	mvc.New(common.App.Party("/user")).Handle(new(UserController))
	mvc.New(common.App.Party("/upload")).Handle(new(UploadController))
	mvc.New(common.App.Party("/site")).Handle(new(SiteController))
	mvc.New(common.App.Party("/topic")).Handle(new(TopicController))
}
