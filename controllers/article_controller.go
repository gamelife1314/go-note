package controllers

import (
	"github.com/gamelife1314/go-note/models"
	"github.com/gamelife1314/go-note/validator"
	"github.com/kataras/iris/mvc"
	"math"
	"strconv"
	"strings"
)

type ArticleController struct {
	BaseController
}

func (a *ArticleController) BeforeActivation(b mvc.BeforeActivation) {
	a.BaseController.BeforeActivation(b)
	a.BaseController.PerPageLimit = 10

	b.Handle("POST", "/create", "Create", Authenticate)
	b.Handle("POST", "/like", "Like", Authenticate)
	b.Handle("GET", "/{articleId:int min(1)}", "Id")
	b.Handle("GET", "/active/list", "ActiveList")
	b.Handle("GET", "/excellent/list", "ExcellentList")
	b.Handle("GET", "/no-comments/list", "NoCommentsList")
	b.Handle("GET", "/latest/list", "LatestList")
}

func (a *ArticleController) Like() *ResponseStructure {
	a.ResetResponseData()
	user := a.Ctx.Values().Get("user").(models.User)

	articleId, err := a.Ctx.PostValueInt("articleId")
	if err != nil {
		a.ResponseStructure.Code = ArticleIdParamError
		a.ResponseStructure.Message = "帖子id错误"
		return a.ResponseStructure
	}

	if pass, _ := validator.Exists("id", "articles", strconv.Itoa(int(articleId))); !pass {
		a.ResponseStructure.Code = ArticleIdParamError
		a.ResponseStructure.Message = "给定的 articleId 不存在"
		return a.ResponseStructure
	}

	var isExists int
	models.Database.Model(&models.Dynamic{}).Where(map[string]interface{}{
		"type":      models.LikeArticleDynamicType,
		"object_id": articleId,
		"user_id":   user.ID,
	}).Count(&isExists)

	var article models.Article
	models.Database.First(&article, articleId)
	if isExists == 0 {
		article.LikeCount += 1
		models.Database.Save(&article)
		user.LikeArticle(&article)
	}
	a.ResponseStructure.Data["article"] = article.Transform(true, true)

	return a.ResponseStructure
}

func (a *ArticleController) LatestList() *ResponseStructure {
	a.ResetResponseData()
	page := a.Ctx.URLParamInt32Default("page", 1)
	offset := (page - 1) * a.PerPageLimit

	var count int
	models.Database.Model(&models.Article{}).Count(&count)
	pages := math.Ceil(float64(count) / float64(a.PerPageLimit))

	var articles []models.Article
	models.Database.
		Order("id desc").
		Limit(a.PerPageLimit).
		Offset(offset).
		Find(&articles)

	var result = make([]map[string]interface{}, 0)
	for _, article := range articles {
		result = append(result, article.Transform(true, true))
	}
	a.ResponseStructure.Data["articles"] = result
	a.ResponseStructure.Data["total_page"] = pages
	return a.ResponseStructure
}

func (a *ArticleController) NoCommentsList() *ResponseStructure {
	a.ResetResponseData()

	page := a.Ctx.URLParamInt32Default("page", 1)
	offset := (page - 1) * a.PerPageLimit

	var count int
	models.Database.Model(&models.Article{}).Count(&count)
	pages := math.Ceil(float64(count) / float64(a.PerPageLimit))

	var articles []models.Article
	models.Database.
		Select("articles.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN `comments` ON `comments`.`article_id` = `articles`.`id`").
		Group("articles.id").
		Order("comment_count desc").Having("comment_count > 0").
		Limit(a.PerPageLimit).
		Offset(offset).
		Preload("Creator").
		Preload("Topics").
		Preload("Comments").
		Find(&articles)

	var result = make([]map[string]interface{}, 0)
	for _, article := range articles {
		result = append(result, article.Transform(true, true))
	}
	a.ResponseStructure.Data["articles"] = result
	a.ResponseStructure.Data["total_page"] = pages
	return a.ResponseStructure
}

func (a *ArticleController) ExcellentList() *ResponseStructure {
	a.ResetResponseData()

	page := a.Ctx.URLParamInt32Default("page", 1)
	offset := (page - 1) * a.PerPageLimit

	var count int
	models.Database.Model(&models.Article{}).Count(&count)
	pages := math.Ceil(float64(count) / float64(a.PerPageLimit))

	var articles []models.Article
	models.Database.
		Where(map[string]interface{}{"is_excellent": 1}).
		Limit(a.PerPageLimit).
		Offset(offset).
		Find(&articles)

	var result = make([]map[string]interface{}, 0)
	for _, article := range articles {
		result = append(result, article.Transform(true, true))
	}
	a.ResponseStructure.Data["articles"] = result
	a.ResponseStructure.Data["total_page"] = pages

	return a.ResponseStructure
}

func (a *ArticleController) ActiveList() *ResponseStructure {
	a.ResetResponseData()

	page := a.Ctx.URLParamInt32Default("page", 1)
	offset := (page - 1) * a.PerPageLimit

	var count int
	models.Database.Model(&models.Article{}).Count(&count)
	pages := math.Ceil(float64(count) / float64(a.PerPageLimit))

	var articles []models.Article
	models.Database.
		Select("articles.*, COUNT(comments.id) as comment_count").
		Joins("LEFT JOIN `comments` ON `comments`.`article_id` = `articles`.`id`").
		Group("articles.id").
		Order("comment_count desc").
		Limit(a.PerPageLimit).
		Offset(offset).
		Preload("Creator").
		Preload("Topics").
		Preload("Comments").
		Find(&articles)

	var result = make([]map[string]interface{}, 0)
	for _, article := range articles {
		result = append(result, article.Transform(true, true))
	}
	a.ResponseStructure.Data["articles"] = result
	a.ResponseStructure.Data["total_page"] = pages
	return a.ResponseStructure
}

func (a *ArticleController) Id() *ResponseStructure {
	a.ResetResponseData()
	articleId, err := a.Ctx.Params().GetInt("articleId")
	if err != nil {
		a.ResponseStructure.Code = ArticleIdParamError
		a.ResponseStructure.Message = "帖子id错误"
		return a.ResponseStructure
	}

	var article models.Article
	models.Database.First(&article, articleId)
	a.ResponseStructure.Data["article"] = article.Transform(true, true)
	return a.ResponseStructure
}

func (a *ArticleController) Create() *ResponseStructure {
	a.ResetResponseData()

	user := a.Ctx.Values().Get("user").(models.User)

	title := strings.Trim(a.Ctx.PostValue("title"), " ")
	content := strings.Trim(a.Ctx.PostValue("content"), " ")
	topicId, err := a.Ctx.PostValueInt("topicId")
	if err != nil {
		a.ResponseStructure.Code = ArticleTopicNotExists
		a.ResponseStructure.Message = "topic id 指定的topic不存在"
		return a.ResponseStructure
	}

	if pass, msg := validator.Length("title", title, 8, 128); !pass {
		a.ResponseStructure.Code = ArticleTitleLengthError
		a.ResponseStructure.Message = msg
		return a.ResponseStructure
	}

	if pass, msg := validator.Length("content", content, 24, 65535); !pass {
		a.ResponseStructure.Code = ArticleTitleContentError
		a.ResponseStructure.Message = msg
		return a.ResponseStructure
	}

	if pass, msg := validator.Exists("id", "topics", strconv.Itoa(int(topicId))); !pass {
		a.ResponseStructure.Code = ArticleTopicNotExists
		a.ResponseStructure.Message = msg
		return a.ResponseStructure
	}

	var article models.Article
	article.New(title, content, uint(topicId), &user)
	article.SetField()
	article.Creator.FillRelatedFields()
	a.ResponseStructure.Data["article"] = article.Transform(true, true)

	return a.ResponseStructure
}
