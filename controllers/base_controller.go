package controllers

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type EmptyData struct{}

type ResponseStructure struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type BaseController struct {
	ResponseStructure *ResponseStructure
	Ctx               iris.Context
}

func (base *BaseController) ResetResponseData() {
	base.ResponseStructure.Code = 200
	base.ResponseStructure.Message = "ok"
	base.ResponseStructure.Data = map[string]interface{}{}
}

func (base *BaseController) BeforeActivation(b mvc.BeforeActivation) {
	base.ResponseStructure = &ResponseStructure{
		Code:    200,
		Message: "ok",
		Data:    map[string]interface{}{},
	}
}
