package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Database *gorm.DB

func DatabaseInit() {
	Database.LogMode(true)
	Database.AutoMigrate(
		&User{},
		&Dynamic{},
		&FollowRelation{},
		&Article{},
		&Comment{},
		&Topic{},
		&ArticleTopic{},
	)

	Database.Model(&Topic{}).RemoveIndex("name")
	Database.Model(&Topic{}).AddUniqueIndex("idx_name_parent_id", "name", "parent_id")
	InitTopics()
}
