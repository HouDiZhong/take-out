package repository

import (
	"context"
	"take-out/common"
	"take-out/global/tx"
	"take-out/internal/api/request"
	"take-out/internal/model"
)

type DishRepo interface {
	Transaction(ctx context.Context) tx.Transaction
	Insert(db tx.Transaction, dish *model.Dish) error
	PageQuery(ctx context.Context, dto *request.DishPageQueryDTO) (*common.PageResult, error)
	GetById(ctx context.Context, id uint64) (*model.Dish, error)
	List(ctx context.Context, categoryId uint64) ([]model.Dish, error)
	OnOrClose(ctx context.Context, id uint64, status int) error
	Update(db tx.Transaction, dish model.Dish) error
	Delete(db tx.Transaction, id uint64) error
}
