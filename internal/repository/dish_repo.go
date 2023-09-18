package repository

import (
	"context"
	"gorm.io/gorm"
	"take-out/common"
	"take-out/internal/api/request"
	"take-out/internal/model"
)

type DishRepo interface {
	Transaction(ctx context.Context) *gorm.DB
	Insert(db *gorm.DB, dish *model.Dish) error
	PageQuery(ctx context.Context, dto *request.DishPageQueryDTO) (*common.PageResult, error)
}