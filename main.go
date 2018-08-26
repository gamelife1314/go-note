package main

import (
	"flag"
	"fmt"
	"github.com/gamelife1314/go-note/common"
	"github.com/gamelife1314/go-note/config"
	_ "github.com/gamelife1314/go-note/controllers"
	"github.com/gamelife1314/go-note/models"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
)

var port = flag.Int("port", 8989, "Http Server Listen Port!")
var addr = flag.String("addr", "0.0.0.0", "Http Server Listen Addr!")

func init() {
	var err error
	models.Database, err = gorm.Open("mysql", config.Configuration.Other["DSN"].(string))
	if err != nil {
		panic("Error Happened when Connect to Database!")
	}

	models.Database.LogMode(true)
	models.Database.AutoMigrate(
		&models.User{},
		&models.Dynamic{},
		&models.FollowRelation{},
		&models.Article{},
		&models.Comment{},
		&models.Topic{},
		&models.ArticleTopic{},
	)
}

func main() {
	flag.Parse()
	common.App.Run(iris.Addr(fmt.Sprintf("%s:%d", *addr, *port)))
}
