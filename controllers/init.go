package controllers

import (
	"github.com/gamelife1314/go-note/common"
	"github.com/kataras/iris/mvc"
)

type EmptyData struct{}

type ResponseStructure struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func init() {

	mvc.New(common.App.Party("/user")).Handle(new(UserController))
	mvc.New(common.App.Party("/upload")).Handle(new(UploadController))
}
