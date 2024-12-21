package dao

import (
	"take-out/internal/model"
	"take-out/internal/repository"

	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repository.UserRepo {
	return &UserDao{db: db}
}

func (u *UserDao) FindByOpenId(openId string) (model.User, error) {
	user := model.User{}
	err := u.db.Model(&user).Where("openid = ?", openId).Find(&user).Error
	return user, err
}

func (u *UserDao) CreateUser(user *model.User) error {
	return u.db.Create(&user).Error
}
