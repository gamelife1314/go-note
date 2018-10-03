package models

import (
	"time"
)

type Article struct {
	ID                 uint       `gorm:"primary_key" json:"id"`
	Title              string     `gorm:"not null" json:"title"`
	Content            string     `gorm:"type:longtext" json:"content"`
	CreatorId          uint       `json:"creatorId"`
	Display            uint8      `gorm:"type:tinyint;default:1" json:"display"`
	IsSticky           bool       `gorm:"type:tinyint;default:0" json:"isSticky"`
	IsRecommended      bool       `gorm:"type:tinyint;default:0" json:"isRecommended"`
	IsExcellent        bool       `gorm:"type:tinyint;default:0" json:"isExcellent"`
	ViewCount          uint       `gorm:"type:int unsigned;default:0" json:"viewCount"`
	LikeCount          uint       `gorm:"type:int unsigned;default:0" json:"likeCount"`
	CreatedAt          time.Time  `json:"-"`
	CreatedAtTimestamp uint       `gorm:"-" json:"createdAt"`
	UpdatedAt          time.Time  `json:"-"`
	DeletedAt          *time.Time `sql:"index" json:"-"`

	Creator  *User     `gorm:"foreignkey:CreatorId" json:"user"`      // 文章作者
	Topics   []Topic   `gorm:"many2many:article_topic" json:"topics"` // 文章主题
	Comments []Comment `gorm:"foreignkey:article_id" json:"topics"`   // 文章主题
}

func (a *Article) New(title, content string, topicId uint, creator *User) *Article {
	var topic Topic
	Database.Where(map[string]interface{}{"id": topicId}).First(&topic)
	a.Content = content
	a.Title = title
	a.Creator = creator
	Database.Create(a)
	Database.Model(a).Association("Topics").Append([]Topic{topic})
	return a
}

func (a *Article) SetField() {
	a.CreatedAtTimestamp = uint(a.CreatedAt.Unix())
}

func (a *Article) Transform(user, topic bool) map[string]interface{} {

	if a.Comments == nil {
		Database.Preload("Comments").First(a)
	}

	result := map[string]interface{}{
		"id":            a.ID,
		"content":       a.Content,
		"creatorId":     a.CreatorId,
		"createdAt":     a.CreatedAt.Unix(),
		"isSticky":      a.IsSticky,
		"display":       a.Display,
		"isRecommended": a.IsRecommended,
		"title":         a.Title,
		"viewCount":     a.ViewCount,
		"likeCount":     a.LikeCount,
		"commentCount":  len(a.Comments),
		"isExcellent":   a.IsExcellent,
		"isSubscribe":   false,
	}

	if user == true {
		if a.Creator == nil {
			Database.Preload("Creator").First(a)
		}
		result["author"] = a.Creator.Transform()
	}

	if topic == true {
		if a.Topics == nil {
			Database.Preload("Topics").First(a)
		}
		var topics []map[string]interface{}
		for _, topic := range a.Topics {
			topics = append(topics, topic.Transform())
		}
		result["topics"] = topics
	}

	return result
}
