package controllers

import (
	"fmt"
	"github.com/gamelife1314/go-note/common"
	"github.com/gamelife1314/go-note/models"
	"github.com/gamelife1314/go-note/validator"
	"github.com/kataras/iris/mvc"
	"math"
	"strconv"
	"strings"
)

type UserController struct {
	BaseController
}

func (u *UserController) BeforeActivation(b mvc.BeforeActivation) {
	u.BaseController.BeforeActivation(b)
	u.PerPageLimit = 10
	b.Handle("POST", "/register", "Post")
	b.Handle("POST", "/login", "Login")
	b.Handle("GET", "/profile", "Profile", Authenticate)
	b.Handle("POST", "/profile", "UpdateProfile", Authenticate)
	b.Handle("POST", "/update/password", "UpdatePassword", Authenticate)
	b.Handle("POST", "/follow", "Follow", Authenticate)
	b.Handle("GET", "/follow/list", "Followers", Authenticate)
	b.Handle("GET", "/fans/list", "Fans", Authenticate)
	b.Handle("GET", "/dynamics/list", "Dynamics", Authenticate)
}

func (u *UserController) UpdatePassword() *ResponseStructure {
	u.ResetResponseData()
	var user = u.Ctx.Values().Get("user").(models.User)

	oldPassword := strings.Trim(u.Ctx.PostValue("oldPassword"), " ")
	password := strings.Trim(u.Ctx.PostValue("password"), " ")
	passwordConfirm := strings.Trim(u.Ctx.PostValue("passwordConfirm"), " ")

	if models.CryptPassword(oldPassword) != user.Password {
		u.ResponseStructure.Code = OldPasswordInputError
		u.ResponseStructure.Message = "旧密码输入错误"
		return u.ResponseStructure
	}

	if pass, msg := validator.Length("password", password, 6, 12); !pass {
		u.ResponseStructure.Code = UserRegisterPasswordLengthError
		u.ResponseStructure.Message = msg
		return u.ResponseStructure
	}

	if password != passwordConfirm {
		u.ResponseStructure.Code = PasswordNotEqual
		u.ResponseStructure.Message = "两次输入的密码不相等"
		return u.ResponseStructure
	}
	models.Database.Model(&user).Update("password", models.CryptPassword(password))
	u.ResponseStructure.Data["user"] = user
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
	u.ResponseStructure.Data["user"] = user
	return u.ResponseStructure
}

func (u *UserController) Profile() *ResponseStructure {
	u.ResetResponseData()
	u.ResponseStructure.Data["user"] = u.Ctx.Values().Get("user").(models.User)
	return u.ResponseStructure
}

func (u *UserController) Login() *ResponseStructure {
	u.ResetResponseData()
	credential := strings.Trim(u.Ctx.PostValue("credential"), " ")
	password := strings.Trim(u.Ctx.PostValue("password"), " ")
	var user models.User
	user.Login(credential, password)
	if user.ID != 0 {
		u.ResponseStructure.Data["user"] = user
		u.ResponseStructure.Data["authToken"] = user.AuthToken
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
	u.ResponseStructure.Data["user"] = user
	u.ResponseStructure.Data["authToken"] = user.AuthToken

	return u.ResponseStructure
}

func (u *UserController) Follow() *ResponseStructure {
	u.ResetResponseData()
	var user = u.Ctx.Values().Get("user").(models.User)

	userId, err := u.Ctx.PostValueInt("userId")

	if err != nil {
		u.ResponseStructure.Code = FollowUserNotExists
		u.ResponseStructure.Message = "userId 指定的用户不存在"
		return u.ResponseStructure
	}

	if uint(userId) == user.ID {
		u.ResponseStructure.Code = FollowUserSelf
		u.ResponseStructure.Message = "不能关注自己"
		return u.ResponseStructure
	}

	var isExists int
	models.Database.Model(&models.FollowRelation{}).Where(map[string]interface{}{
		"source_user_id": user.ID,
		"target_user_id": userId,
	}).Count(&isExists)
	if isExists != 0 {
		u.ResponseStructure.Code = FollowExists
		u.ResponseStructure.Message = "已经关注此人"
		return u.ResponseStructure
	}

	if pass, _ := validator.Exists("id", "users", strconv.Itoa(int(userId))); !pass {
		u.ResponseStructure.Code = FollowUserNotExists
		u.ResponseStructure.Message = "userId 指定的用户不存在"
		return u.ResponseStructure
	}

	user.Follow(userId)
	u.ResponseStructure.Data["followers"] = user.FansList()

	var follower models.User
	models.Database.First(&follower, userId)
	user.LikeUser(&follower)
	return u.ResponseStructure
}

func (u *UserController) Followers() *ResponseStructure {
	u.ResetResponseData()
	var user = u.Ctx.Values().Get("user").(models.User)
	u.ResponseStructure.Data["followers"] = user.FollowersList()
	return u.ResponseStructure
}

func (u *UserController) Fans() *ResponseStructure {
	u.ResetResponseData()
	var user = u.Ctx.Values().Get("user").(models.User)
	u.ResponseStructure.Data["fans"] = user.FansList()
	return u.ResponseStructure
}

func (u *UserController) Dynamics() *ResponseStructure {
	u.ResetResponseData()

	var user = u.Ctx.Values().Get("user").(models.User)
	page := u.Ctx.URLParamInt32Default("page", 1)
	offset := (page - 1) * u.PerPageLimit

	var count int
	models.Database.Model(&models.Dynamic{}).Where("user_id = ?", user.ID).Count(&count)
	pages := math.Ceil(float64(float64(count) / float64(u.PerPageLimit)))

	var dynamics []models.Dynamic
	models.Database.Where("user_id = ?", user.ID).
		Order("id desc").
		Limit(u.PerPageLimit).
		Offset(offset).
		Find(&dynamics)

	var result = make([]map[string]interface{}, 0)
	for _, article := range dynamics {
		result = append(result, article.Transform())
	}
	u.ResponseStructure.Data["dynamics"] = result
	u.ResponseStructure.Data["total_page"] = pages
	return u.ResponseStructure
}
