package controllers

import (
	"fmt"
	"github.com/gamelife1314/go-note/common"
	"github.com/gamelife1314/go-note/config"
	"github.com/gamelife1314/go-note/models"
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"io"
	"os"
	"path"
	"strings"
)

type UploadController struct {
	BaseController
}

func CustomContentLengthLimiter(limit int64) func(ctx iris.Context) {
	return func(ctx iris.Context) {
		if ctx.GetContentLength() > limit<<20 {
			ctx.JSON(ResponseStructure{
				Code:    UploadFileExceedLimit,
				Message: fmt.Sprintf("上传文件超过限制：%dM", limit),
				Data:    EmptyData{},
			})
			return
		}
		ctx.Next()
	}
}

func CreateUserDirMiddleware(ctx iris.Context) {
	user := ctx.Values().Get("user").(models.User)
	userDir := fmt.Sprintf("%s/%d", config.Configuration.Other["UploadDir"].(string), user.ID)
	_, err := os.Stat(userDir)
	if err != nil && os.IsNotExist(err) {
		os.Mkdir(userDir, os.ModePerm)
	}
	ctx.Values().Set("userDir", userDir)
	ctx.Next()
}

func (u *UploadController) BeforeActivation(b mvc.BeforeActivation) {
	u.BaseController.BeforeActivation(b)

	b.Handle("POST", "/avatar", "Avatar", CustomContentLengthLimiter(5), Authenticate, CreateUserDirMiddleware)
}

func (u *UploadController) Avatar() *ResponseStructure {
	user := u.Ctx.Values().Get("user").(models.User)
	file, info, err := u.Ctx.FormFile("file")
	if err != nil {
		u.ResponseStructure.Code = ServerInternalError
		u.ResponseStructure.Message = "服务器内部错误"
		return u.ResponseStructure
	}
	defer file.Close()

	userDir := u.Ctx.Values().Get("userDir").(string)
	fileExt := path.Ext(info.Filename)
	fileName := fmt.Sprintf("%s/%s%s", userDir, common.Md5(common.GenerateRandomString(128)), fileExt)

	out, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		u.ResponseStructure.Code = ServerInternalError
		u.ResponseStructure.Message = "服务器内部错误"
		return u.ResponseStructure
	}
	defer out.Close()
	io.Copy(out, file)

	if user.Avatar != nil {
		if _, err := os.Stat(*user.Avatar); err == nil {
			os.Remove(*user.Avatar)
		}
	}

	models.Database.Model(&user).Updates(map[string]interface{}{"avatar": strings.TrimLeft(fileName, ".")})
	u.ResponseStructure.Data = user

	return u.ResponseStructure
}
