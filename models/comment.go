package models

import "time"

type Comment struct {
	ID            uint   `gorm:"primary_key"`
	Content       string `gorm:"type:text;not null"`
	ArticleId     uint
	CommentatorId uint
	LikeCount     uint `gorm:"type:int unsigned;default:0"`
	HateCount     uint `gorm:"type:int unsigned;default:0"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`

	Article     *Article
	Commentator *User `gorm:"foreignkey:CommentatorId"`
}

func (c *Comment) New(user User, content string, articleId uint) *Comment {
	if c.ID == 0 {
		c.Content = content
		c.ArticleId = articleId
		c.CommentatorId = user.ID
		Database.Create(c)
	}
	return c
}

func (c *Comment) Transform(article bool) map[string]interface{} {

	data := map[string]interface{}{
		"id":        c.ID,
		"content":   c.Content,
		"articleId": c.ArticleId,
		"likeCount": c.LikeCount,
		"hateCount": c.HateCount,
		"createdAt": c.CreatedAt.Unix(),
	}

	if c.Commentator == nil {
		Database.Preload("Commentator").First(c)
	}

	data["commentor"] = c.Commentator.Transform()

	if article {
		if c.Article == nil {
			Database.Preload("Article").First(c)
		}
		data["article"] = c.Article.Transform(true, true)
	}

	return data
}
