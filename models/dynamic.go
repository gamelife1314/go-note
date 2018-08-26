package models

import "time"

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
