package repository

import "take-out/internal/model"

type UserRepo interface {
	FindByOpenId(openId string) (model.User, error)
	CreateUser(user *model.User) error
}
