package models

import "time"

type Article struct {
	ID            uint   `gorm:"primary_key"`
	Title         string `gorm:"not null"`
	Content       string `gorm:"type:longtext"`
	CreatorId     uint
	Display       uint8 `gorm:"type:tinyint;default:1"`
	IsSticky      bool  `gorm:"type:tinyint;default:0"`
	IsRecommended bool  `gorm:"type:tinyint;default:0"`
	ViewCount     uint  `gorm:"type:int unsigned;default:0"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`

	Creator User    `gorm:"foreignkey:CreatorId"`    // 文章作者
	Topics  []Topic `gorm:"many2many:article_topic"` // 文章主题
}
