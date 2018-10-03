package models

import "time"

type ArticleTopic struct {
	ID        uint `gorm:"primary_key"`
	ArticleId uint
	TopicId   uint
	CreatedAt time.Time
}

func (a *ArticleTopic) TableName() string {
	return "article_topic"
}
