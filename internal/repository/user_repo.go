package repository

import (
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/model"
)

type UserRepo interface {
	FindByOpenId(openId string) (model.User, error)
	CreateUser(user *model.User) error

	// 获取今天新增用户
	GetNewUserNumber() (int64, error)
	// 用户统计
	UserReport(dto request.ReportQuestDTO) ([]response.LocalUsertVO, error)
	// 每天新增人数
	EveryUserReport(dto request.ReportQuestDTO) ([]response.EveryUserVO, error)
}
