package controllers

import (
	"github.com/gamelife1314/go-note/common"
	"github.com/gamelife1314/go-note/models"
	"github.com/kataras/iris"
	"time"
)

func Authenticate(ctx iris.Context) {
	authToken := ctx.GetHeader("x-token")
	var user models.User
	user.ParseUserFromAuthToken(authToken)
	now := time.Now().In(common.TimeZone)
	remoteAddr := ctx.RemoteAddr()
	models.Database.Model(&user).Updates(&models.User{LastVisitedAt: &now, LastVisitedIp: &remoteAddr})
	if user.ID == 0 {
		ctx.JSON(ResponseStructure{
			Code:    Unauthenticated,
			Message: "Unauthenticated",
			Data:    map[string]interface{}{},
		})
	} else {
		ctx.Values().Set("user", user)
		ctx.Next()
	}
}
