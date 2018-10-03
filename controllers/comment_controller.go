package controllers

import (
	"github.com/gamelife1314/go-note/models"
	"github.com/gamelife1314/go-note/validator"
	"github.com/kataras/iris/mvc"
	"math"
	"strconv"
	"strings"
)

type CommentController struct {
	BaseController
}

func (c *CommentController) BeforeActivation(b mvc.BeforeActivation) {
	c.BaseController.BeforeActivation(b)
	c.BaseController.PerPageLimit = 10
	b.Handle("POST", "/create", "Create", Authenticate)
	b.Handle("POST", "/hate", "Hate", Authenticate)
	b.Handle("POST", "/like", "Like", Authenticate)
	b.Handle("GET", "/{articleId}/list", "List")
}

func (c *CommentController) List() *ResponseStructure {
	c.ResetResponseData()
	articleId := c.Ctx.Params().Get("articleId")
	page := c.Ctx.URLParamInt32Default("page", 1)
	offset := (page - 1) * c.PerPageLimit

	var count int
	models.Database.Model(&models.Comment{}).
		Where(map[string]interface{}{"article_id": articleId}).
		Count(&count)
	pages := math.Ceil(float64(count) / float64(c.PerPageLimit))

	var comments []models.Comment
	models.Database.
		Order("id desc").
		Limit(c.PerPageLimit).
		Offset(offset).
		Find(&comments)

	var result = make([]map[string]interface{}, 0)
	for _, article := range comments {
		result = append(result, article.Transform(false))
	}
	c.ResponseStructure.Data["comments"] = result
	c.ResponseStructure.Data["total_page"] = pages

	return c.ResponseStructure
}

func (c *CommentController) Create() *ResponseStructure {
	c.ResetResponseData()
	user := c.Ctx.Values().Get("user").(models.User)
	articleId, err := c.Ctx.PostValueInt("articleId")
	content := strings.Trim(c.Ctx.PostValue("content"), " ")
	if err != nil {
		c.ResponseStructure.Code = ArticleIdParamError
		c.ResponseStructure.Message = "articleId 应该是个整数吧"
	}
	if pass, _ := validator.Exists("id", "articles", strconv.Itoa(int(articleId))); !pass {
		c.ResponseStructure.Code = ArticleIdParamError
		c.ResponseStructure.Message = "给定的 articleId 不存在"
		return c.ResponseStructure
	}

	var comment models.Comment
	comment.New(user, content, uint(articleId))
	c.ResponseStructure.Data["comment"] = comment.Transform(false)

	var article models.Article
	models.Database.First(&article, articleId)
	user.CommentArticle(&article)

	return c.ResponseStructure
}

func (c *CommentController) Hate() *ResponseStructure {
	c.ResetResponseData()

	user := c.Ctx.Values().Get("user").(models.User)

	commentId, err := c.Ctx.PostValueInt("commentId")
	if err != nil {
		c.ResponseStructure.Code = CommentIdParamError
		c.ResponseStructure.Message = "评论id错误"
		return c.ResponseStructure
	}

	if pass, _ := validator.Exists("id", "comments", strconv.Itoa(int(commentId))); !pass {
		c.ResponseStructure.Code = CommentIdParamError
		c.ResponseStructure.Message = "给定的 commentId 不存在"
		return c.ResponseStructure
	}

	var isExists int
	models.Database.Model(&models.Dynamic{}).Where(map[string]interface{}{
		"type":      models.HateCommentDynamicType,
		"object_id": commentId,
		"user_id":   user.ID,
	}).Count(&isExists)

	var comment models.Comment
	models.Database.First(&comment, commentId)
	if isExists == 0 {
		comment.HateCount += 1
		models.Database.Save(&comment)
		user.HateComment(&comment)
	}
	c.ResponseStructure.Data["comment"] = comment.Transform(false)
	return c.ResponseStructure
}

func (c *CommentController) Like() *ResponseStructure {
	c.ResetResponseData()

	user := c.Ctx.Values().Get("user").(models.User)

	commentId, err := c.Ctx.PostValueInt("commentId")
	if err != nil {
		c.ResponseStructure.Code = CommentIdParamError
		c.ResponseStructure.Message = "评论id错误"
		return c.ResponseStructure
	}

	if pass, _ := validator.Exists("id", "comments", strconv.Itoa(int(commentId))); !pass {
		c.ResponseStructure.Code = CommentIdParamError
		c.ResponseStructure.Message = "给定的 commentId 不存在"
		return c.ResponseStructure
	}

	var isExists int
	models.Database.Model(&models.Dynamic{}).Where(map[string]interface{}{
		"type":      models.LikeCommentDynamicType,
		"object_id": commentId,
		"user_id":   user.ID,
	}).Count(&isExists)

	var comment models.Comment
	models.Database.First(&comment, commentId)
	if isExists == 0 {
		comment.LikeCount += 1
		models.Database.Save(&comment)
		user.LikeComment(&comment)
	}
	c.ResponseStructure.Data["comment"] = comment.Transform(false)
	return c.ResponseStructure
}
