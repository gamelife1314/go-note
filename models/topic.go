package models

import "time"

type Topic struct {
	ID        uint   `gorm:"primary_key"`
	Name      string `gorm:"type:char(24);not null;unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	Articles []Article `gorm:"many2many:article_topic"`
}
