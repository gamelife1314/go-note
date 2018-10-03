package models

import (
	"time"
)

const FocusUserDynamicType = 1      // 关注用户，用户id
const LikeCommentDynamicType = 2    // 点赞评论，评论id
const LikeArticleDynamicType = 3    // 点赞文章，文章id
const CommentArticleDynamicType = 4 // 评论文章，文章id
const HateCommentDynamicType = 5    // 讨厌评论，评论id

type Dynamic struct {
	ID        uint  `gorm:"primary_key"`
	Type      uint8 `gorm:"type:tinyint;index"`
	ObjectId  uint
	Data      string `gorm:"type:text"`
	UserId    uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	User User
}

func (dynamic *Dynamic) TableName() string {
	return "dynamics"
}

func (dynamic *Dynamic) Transform() map[string]interface{} {

	data := map[string]interface{}{
		"id":        dynamic.ID,
		"type":      dynamic.Type,
		"createdAt": dynamic.CreatedAt.Unix(),
	}

	if dynamic.Type == FocusUserDynamicType {
		var user User
		Database.First(&user, dynamic.ObjectId)
		data["data"] = map[string]interface{}{
			"user": user.Transform(),
		}
	}

	if dynamic.Type == LikeCommentDynamicType || dynamic.Type == HateCommentDynamicType {
		var comment Comment
		Database.First(&comment, dynamic.ObjectId)
		data["data"] = map[string]interface{}{
			"comment": comment.Transform(false),
		}
	}

	if dynamic.Type == LikeArticleDynamicType || dynamic.Type == CommentArticleDynamicType {
		var article Article
		Database.First(&article, dynamic.ObjectId)
		data["data"] = map[string]interface{}{
			"article": article.Transform(true, true),
		}
	}

	return data
}
