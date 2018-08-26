package models

import "time"

type ArticleTopic struct {
	ID        uint `gorm:"primary_key"`
	ArticleId uint
	TopicId   uint
	CreatedAt time.Time
}
