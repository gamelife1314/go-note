package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type BaseController struct {
	ResponseStructure *ResponseStructure
	Ctx               iris.Context
}

func (base *BaseController) ResetResponseData() {
	base.ResponseStructure.Code = 200
	base.ResponseStructure.Message = "ok"
	base.ResponseStructure.Data = &EmptyData{}
}

func (base *BaseController) BeforeActivation(b mvc.BeforeActivation) {
	base.ResponseStructure = &ResponseStructure{
		Code:    200,
		Message: "ok",
		Data:    &EmptyData{},
	}
}
