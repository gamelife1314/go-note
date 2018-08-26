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

	Article     Article
	Commentator User `gorm:"foreignkey:CommentatorId"`
}
