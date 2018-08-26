package controllers

import (
	"fmt"
	"github.com/gamelife1314/go-note/common"
	"github.com/gamelife1314/go-note/models"
	"github.com/gamelife1314/go-note/validator"
	"github.com/kataras/iris/mvc"
	"strings"
)

type UserController struct {
	BaseController
}

func (u *UserController) BeforeActivation(b mvc.BeforeActivation) {
	u.BaseController.BeforeActivation(b)
	b.Handle("POST", "/register", "Post")
	b.Handle("POST", "/login", "Login")
	b.Handle("GET", "/profile", "Profile", Authenticate)
	b.Handle("POST", "/profile", "UpdateProfile", Authenticate)
	b.Handle("POST", "/update/password", "UpdatePassword", Authenticate)
}

func (u *UserController) UpdatePassword() *ResponseStructure {
	u.ResetResponseData()
	var user = u.Ctx.Values().Get("user").(models.User)

	password := strings.Trim(u.Ctx.PostValue("password"), " ")
	passwordConfirm := strings.Trim(u.Ctx.PostValue("passwordConfirm"), " ")

	if password != passwordConfirm {
		u.ResponseStructure.Code = PasswordNotEqual
		u.ResponseStructure.Message = "两次输入的密码不相等"
		return u.ResponseStructure
	}

	if pass, msg := validator.Length("password", password, 6, 12); !pass {
		u.ResponseStructure.Code = UserRegisterPasswordLengthError
		u.ResponseStructure.Message = msg
		return u.ResponseStructure
	}
	models.Database.Model(&user).Update("password", models.CryptPassword(password))
	u.ResponseStructure.Data = user
	return u.ResponseStructure
}

func (u *UserController) UpdateProfile() *ResponseStructure {
	u.ResetResponseData()
	var user = u.Ctx.Values().Get("user").(models.User)

	updateFields := []string{"motto", "homepage", "company", "address", "github",
		"twitter", "facebook", "instagram", "telegram", "telegram", "gender",
		"steam"}
	updates := map[string]interface{}{}

	for field, value := range u.Ctx.FormValues() {
		if common.InStringArray(field, updateFields) {
			updates[field] = value[0]
		}
	}
	models.Database.Model(&user).Updates(updates)
	u.ResponseStructure.Data = user
	return u.ResponseStructure
}

func (u *UserController) Profile() *ResponseStructure {
	u.ResetResponseData()
	u.ResponseStructure.Data = u.Ctx.Values().Get("user").(models.User)
	return u.ResponseStructure
}

func (u *UserController) Login() *ResponseStructure {
	u.ResetResponseData()
	credential := strings.Trim(u.Ctx.PostValue("credential"), " ")
	password := strings.Trim(u.Ctx.PostValue("password"), " ")
	var user models.User
	user.Login(credential, password)
	if user.ID != 0 {
		u.ResponseStructure.Data = user
	} else {
		u.ResponseStructure.Code = UserLoginError
		u.ResponseStructure.Message = fmt.Sprintf("用户名（或邮箱）与密码不匹配")
	}
	return u.ResponseStructure
}

func (u *UserController) Post() *ResponseStructure {
	u.ResetResponseData()

	nickname := strings.Trim(u.Ctx.PostValue("nickname"), " ")
	password := strings.Trim(u.Ctx.PostValue("password"), " ")
	email := strings.Trim(u.Ctx.PostValue("email"), " ")
	fmt.Println("hell", nickname)

	if pass, msg := validator.Length("nickname", nickname, 3, 12); !pass {
		u.ResponseStructure.Code = UserRegisterNicknameLengthError
		u.ResponseStructure.Message = msg
		return u.ResponseStructure
	}

	userModel := &models.User{}
	if pass, msg := validator.Unique("nickname", userModel.TableName(), nickname); !pass {
		u.ResponseStructure.Code = UserRegisterNicknameUniqueError
		u.ResponseStructure.Message = msg
		return u.ResponseStructure
	}

	if pass, msg := validator.Length("password", password, 6, 12); !pass {
		u.ResponseStructure.Code = UserRegisterPasswordLengthError
		u.ResponseStructure.Message = msg
		return u.ResponseStructure
	}

	emailRegexp := `^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z0-9]{2,6}$`
	if pass, _ := validator.Regexp("email", email, emailRegexp); !pass {
		u.ResponseStructure.Code = UserRegisterEmailFormatError
		u.ResponseStructure.Message = "不是有效的邮箱"
		return u.ResponseStructure
	}

	if pass, msg := validator.Unique("email", userModel.TableName(), email); !pass {
		u.ResponseStructure.Code = UserRegisterEmailUniqueError
		u.ResponseStructure.Message = msg
		return u.ResponseStructure
	}

	user := models.NewUser(nickname, password, email)
	u.ResponseStructure.Data = user

	return u.ResponseStructure
}
