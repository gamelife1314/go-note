package models

import "time"

type FollowRelation struct {
	ID           uint `gorm:"primary_key"`
	SourceUser   User `gorm:"foreignkey:SourceUserId"`
	SourceUserId uint
	TargetUser   User `gorm:"foreignkey:TargetUserId"`
	TargetUserId uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index"`
}

func (fr *FollowRelation) TableName() string {
	return "follow_relations"
}
