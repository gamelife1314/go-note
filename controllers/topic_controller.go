package controllers

import (
	"github.com/gamelife1314/go-note/models"
	"github.com/kataras/iris/mvc"
)

type TopicController struct {
	BaseController
}

func (t *TopicController) BeforeActivation(b mvc.BeforeActivation) {
	t.BaseController.BeforeActivation(b)
	b.Handle("GET", "/list", "List")
}

func (t *TopicController) List() *ResponseStructure {
	t.ResetResponseData()
	t.ResponseStructure.Data["topics"] = models.TopicsByLevel()
	return t.ResponseStructure
}
