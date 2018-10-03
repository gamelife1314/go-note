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
	AuthToken *string   `gorm:"-" json:"-"`
	Followers []User    `gorm:"many2many:follow_relations;association_jointable_foreignkey:target_user_id;jointable_foreignkey:source_user_id"`
	Fans      []User    `gorm:"many2many:follow_relations;association_jointable_foreignkey:source_user_id;jointable_foreignkey:target_user_id"`
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

func (u *User) Transform() map[string]interface{} {
	var result = map[string]interface{}{
		"address":       u.Address,
		"avatar":        u.Avatar,
		"company":       u.Company,
		"createdAt":     u.CreatedAt.Unix(),
		"email":         u.Email,
		"emailIsActive": u.EmailIsActive,
		"facebook":      u.Facebook,
		"gender":        u.Gender,
		"github":        u.Github,
		"homepage":      u.Homepage,
		"instagram":     u.Instagram,
		"isLocked":      u.IsLocked,
		"lastVisitedIp": u.LastVisitedIp,
		"motto":         u.Motto,
		"nickname":      u.Nickname,
		"steam":         u.Steam,
		"telegram":      u.Telegram,
		"twitter":       u.Twitter,
		"uid":           EncryptModelId(u.ID),
		"id":            u.ID,
	}

	if u.LastVisitedAt != nil {
		result["lastVisitedAt"] = u.LastVisitedAt.Unix()
	}

	return result
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

func (u *User) FollowersList() (followers []map[string]interface{}) {
	var f []User
	Database.Model(u).Related(&f, "Followers")
	for _, user := range f {
		followers = append(followers, user.Transform())
	}
	return followers
}

func (u *User) FansList() (fans []map[string]interface{}) {
	var f []User
	Database.Model(u).Related(&f, "Fans")
	for _, user := range f {
		fans = append(fans, user.Transform())
	}
	return fans
}

func (u *User) Follow(followerId int) *User {
	followRelation := FollowRelation{
		SourceUserId: u.ID,
		TargetUserId: uint(followerId),
	}
	Database.NewRecord(followRelation)
	Database.Create(&followRelation)
	return u
}

func (u *User) LikeArticle(article *Article) *Dynamic {

	var dynamic Dynamic

	Database.Where(map[string]interface{}{
		"type":      LikeArticleDynamicType,
		"object_id": article.ID,
		"user_id":   u.ID,
	}).First(&dynamic)

	if dynamic.ID != 0 {
		return &dynamic
	}

	data, _ := json.Marshal(map[string]interface{}{
		"article": article.Transform(false, false),
	})

	dynamic = Dynamic{
		Type:     LikeArticleDynamicType,
		ObjectId: article.ID,
		UserId:   u.ID,
		Data:     string(data),
	}

	Database.NewRecord(dynamic)
	Database.Create(&dynamic)

	return &dynamic
}

func (u *User) LikeUser(user *User) *Dynamic {

	var dynamic Dynamic

	Database.Where(map[string]interface{}{
		"type":      FocusUserDynamicType,
		"object_id": user.ID,
		"user_id":   u.ID,
	}).First(&dynamic)

	if dynamic.ID != 0 {
		return &dynamic
	}

	data, _ := json.Marshal(map[string]interface{}{
		"user": user.Transform(),
	})

	dynamic = Dynamic{
		Type:     FocusUserDynamicType,
		ObjectId: user.ID,
		UserId:   u.ID,
		Data:     string(data),
	}

	Database.NewRecord(dynamic)
	Database.Create(&dynamic)

	return &dynamic
}

func (u *User) LikeComment(comment *Comment) *Dynamic {

	var dynamic Dynamic

	Database.Where(map[string]interface{}{
		"type":      LikeCommentDynamicType,
		"object_id": comment.ID,
		"user_id":   u.ID,
	}).First(&dynamic)

	if dynamic.ID != 0 {
		return &dynamic
	}

	data, _ := json.Marshal(map[string]interface{}{
		"comment": comment.Transform(false),
	})

	dynamic = Dynamic{
		Type:     LikeCommentDynamicType,
		ObjectId: comment.ID,
		UserId:   u.ID,
		Data:     string(data),
	}

	Database.NewRecord(dynamic)
	Database.Create(&dynamic)

	return &dynamic
}

func (u *User) HateComment(comment *Comment) *Dynamic {

	var dynamic Dynamic

	Database.Where(map[string]interface{}{
		"type":      HateCommentDynamicType,
		"object_id": comment.ID,
		"user_id":   u.ID,
	}).First(&dynamic)

	if dynamic.ID != 0 {
		return &dynamic
	}

	data, _ := json.Marshal(map[string]interface{}{
		"comment": comment.Transform(false),
	})

	dynamic = Dynamic{
		Type:     HateCommentDynamicType,
		ObjectId: comment.ID,
		UserId:   u.ID,
		Data:     string(data),
	}

	Database.NewRecord(dynamic)
	Database.Create(&dynamic)

	return &dynamic
}

func (u *User) CommentArticle(article *Article) *Dynamic {
	var dynamic Dynamic

	data, _ := json.Marshal(map[string]interface{}{
		"article": article.Transform(true, true),
	})

	dynamic = Dynamic{
		Type:     CommentArticleDynamicType,
		ObjectId: article.ID,
		UserId:   u.ID,
		Data:     string(data),
	}

	Database.NewRecord(dynamic)
	Database.Create(&dynamic)

	return &dynamic
}
