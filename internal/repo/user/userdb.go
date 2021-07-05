package user

import (
	"goo/internal/types"
	"goo/internal/utils"
	"context"
	"errors"
	"gorm.io/gorm"
	"strings"
)

type UserRepo interface {
	GetCurrentUserInfo(ctx context.Context, adminMember *AdminMember) error
	VerifyAdminMember(ctx context.Context, adminMember *AdminMember) (userType string, err error)
}

type AdminMember struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(255);not null"`
	Password string `gorm:"column:password;type:varchar(255);not null"`
	Email    string `gorm:"unique;column:email;type:varchar(255);not null"`
}

func (u *AdminMember) TableName() string {
	return "admin_member"
}

type UserDB struct {
	db *gorm.DB
}

func NewUserDB(db *gorm.DB) *UserDB {
	return &UserDB{db: db}
}

func (u *UserDB) GetCurrentUserInfo(ctx context.Context, adminMember *AdminMember) error {
	return u.db.First(adminMember, adminMember.ID).Error
}

func (u *UserDB) VerifyAdminMember(ctx context.Context, adminMember *AdminMember) (userType string, err error) {
	email := adminMember.Email
	password := adminMember.Password
	index := strings.Index(email, ":")
	if index == -1 {
		return userType, errors.New("格式错误")
	}

	userType = email[:index]
	if ok, _ := utils.InArray(userType, types.UserType); !ok {
		return userType, errors.New("格式错误")
	}
	email = email[index+1:]
	return userType, u.db.Where("email = ?", email).Where("password = ?", password).First(adminMember).Error
}
