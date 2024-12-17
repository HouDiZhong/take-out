package repository

import (
	"take-out/global/tx"
	"take-out/internal/model"
)

type SetMealDishRepo interface {
	InsertBatch(db tx.Transaction, setmealDishs []model.SetMealDish) error
	DeledeBatch(db tx.Transaction, setmealDishs []model.SetMealDish) error
	DeledeSetMealBatch(db tx.Transaction, ids []string) error
	GetBySetMealId(db tx.Transaction, SetMealId uint64) ([]model.SetMealDish, error)
	GetBySetDishId(db tx.Transaction, DishId uint64) ([]model.SetMealDish, error)
}
