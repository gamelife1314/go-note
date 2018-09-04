package controllers

import (
	"github.com/gamelife1314/go-note/models"
	"github.com/kataras/iris/mvc"
)

type SiteController struct {
	BaseController
}

func (s *SiteController) BeforeActivation(b mvc.BeforeActivation) {
	s.BaseController.BeforeActivation(b)

	b.Handle("GET", "/hot/topics", "HotTopics")
	b.Handle("GET", "/active/users", "ActiveUsers")
	b.Handle("GET", "/latest/users", "LatestUsers")
}

func (s *SiteController) HotTopics() *ResponseStructure {
	s.ResetResponseData()
	return s.ResponseStructure
}

func (s *SiteController) ActiveUsers() *ResponseStructure {
	s.ResetResponseData()
	var users []models.User
	models.Database.Order("last_visited_at desc").Limit(10).Find(&users)
	for index := range users {
		users[index].FillRelatedFields()
	}
	s.ResponseStructure.Data["users"] = users
	return s.ResponseStructure
}

func (s *SiteController) LatestUsers() *ResponseStructure {
	s.ResetResponseData()
	var users []models.User
	models.Database.Order("created_at desc").Limit(10).Find(&users)
	for index := range users {
		users[index].FillRelatedFields()
	}
	s.ResponseStructure.Data["users"] = users
	return s.ResponseStructure
}
