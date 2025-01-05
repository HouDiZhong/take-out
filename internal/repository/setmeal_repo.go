package repository

import (
	"context"
	"take-out/common"
	"take-out/global/tx"
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/model"
)

type SetMealRepo interface {
	Transaction(ctx context.Context) tx.Transaction
	Insert(db tx.Transaction, meal *model.SetMeal) error
	// UpData 动态修改
	Update(ctx tx.Transaction, meal *model.SetMeal) error
	Delete(ctx tx.Transaction, ids []string) error
	PageQuery(ctx context.Context, dto request.SetMealPageQueryDTO) (*common.PageResult, error)
	SetStatus(ctx context.Context, id uint64, status int) error
	GetByIdWithDish(db tx.Transaction, dishId uint64) (model.SetMeal, error)

	QueryListById(cId string) ([]model.SetMeal, error)
	SetMealDishById(cId string) ([]response.SetMealDish, error)

	QuerySetMealDesById(sId string) (model.SetMeal, error)

	QuerySetMealDesStatusNumber() (response.SetmealAndDishVO, error)
}
