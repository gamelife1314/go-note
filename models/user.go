package models

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gamelife1314/go-note/common"
	"github.com/gamelife1314/go-note/config"
	"github.com/kataras/iris/core/errors"
	"time"
)

const MaleGender = 1
const FemaleGender = 2

const UserUnLocked = 0
const UserLocked = 1

type User struct {
	ID                     uint       `gorm:"primary_key" json:"-"`
	UID                    string     `gorm:"-" json:"uid"`
	Nickname               string     `gorm:"type:char(24);not null;unique_index" json:"nickname"`
	Gender                 uint8      `gorm:"type:tinyint;default:1" json:"gender"`
	Email                  string     `gorm:"unique_index;not null" json:"email"`
	Motto                  *string    `gorm:"default: null" json:"motto"`    // 个性签名
	Homepage               *string    `gorm:"default: null" json:"homepage"` // 个人主页
	Company                *string    `gorm:"default: null" json:"company"`  // 公司
	Address                *string    `gorm:"default: null" json:"address"`  // 住址
	Avatar                 *string    `gorm:"default: null" json:"avatar"`   // 头像
	Password               string     `gorm:"not null;type:varchar(255)" json:"-"`
	Github                 *string    `gorm:"default: null" json:"github"`
	Twitter                *string    `gorm:"default: null" json:"twitter"`
	Facebook               *string    `gorm:"default: null" json:"facebook"`
	Instagram              *string    `gorm:"default: null" json:"instagram"`
	Telegram               *string    `gorm:"default: null" json:"telegram"`
	Steam                  *string    `gorm:"default: null" json:"steam"`
	EmailIsActive          uint8      `gorm:"default:0;not null" json:"emailIsActive"`
	IsLocked               uint8      `gorm:"default:0;not null" json:"isLocked"`
	LastVisitedAt          *time.Time `gorm:"default: null" json:"-"`
	LastVisitedAtTimestamp *uint      `gorm:"-" json:"lastVisitedAt"`
	LastVisitedIp          *string    `gorm:"default: null" json:"lastVisitedIp"`
	CreatedAt              time.Time  `json:"-"`
	CreatedAtTimestamp     uint       `gorm:"-" json:"createdAt"`
	UpdatedAt              time.Time  `json:"-"`
	DeletedAt              *time.Time `sql:"index" json:"-"`

	Dynamics  []Dynamic `json:"-"`                             // 关于某人的动态
	Articles  []Article `gorm:"foreignkey:CreatorId" json:"-"` // 某人发布的文章
	AuthToken *string   `gorm:"-" json:"authToken"`
}

type UserAuthToken struct {
	ExpiredAt int64 `json:"expiredAt"`
	UserId    uint  `json:"userId"`
}

func (u *User) TableName() string {
	return "users"
}

func NewUser(nickname, password, email string) *User {
	user := User{
		Nickname: nickname,
		Password: CryptPassword(password),
		Email:    email,
	}
	Database.NewRecord(user)
	Database.Create(&user)
	user.FillRelatedFields()

	return &user
}

func (u *User) Login(credential, password string) {
	password = CryptPassword(password)
	Database.Where("nickname = ?", credential).Or("email = ?", credential).First(u)
	u.FillRelatedFields()
}

func (u *User) FillRelatedFields() {
	u.UID = EncryptModelId(u.ID)
	u.CreatedAtTimestamp = uint(u.CreatedAt.Unix())
	u.AuthToken, _ = u.GenerateAuthToken()
	if u.LastVisitedAt != nil {
		ts := uint(u.LastVisitedAt.Unix())
		u.LastVisitedAtTimestamp = &ts
	}
}

func (u *User) GenerateAuthToken() (*string, error) {
	var authToken = UserAuthToken{
		UserId:    u.ID,
		ExpiredAt: time.Now().In(common.TimeZone).Unix() + config.Configuration.Other["TokenExpires"].(int64),
	}
	if jsonBytes, err := json.Marshal(authToken); err == nil {
		var tokenString string
		tokenString = base64.StdEncoding.EncodeToString(jsonBytes)
		return &tokenString, nil
	} else {
		return nil, err
	}
}

func (u *User) ParseUserFromAuthToken(token string) error {
	if data, err := base64.StdEncoding.DecodeString(token); err != nil {
		return err
	} else {
		var userAuthToken UserAuthToken
		err = json.Unmarshal([]byte(data), &userAuthToken)
		if userAuthToken.ExpiredAt < time.Now().In(common.TimeZone).Unix() {
			return errors.New("该 Token 已经过期")
		} else {
			Database.Where("is_locked = ?", UserUnLocked).First(u, userAuthToken.UserId)
			if u.ID == userAuthToken.UserId {
				u.FillRelatedFields()
				return nil
			} else {
				return errors.New(fmt.Sprintf("用户不存在"))
			}
		}
	}
}
