package repository

import (
	"take-out/global/tx"
	"take-out/internal/model"
)

type DishFlavorRepo interface {
	// InsertBatch 批量插入菜品口味
	InsertBatch(db tx.Transaction, flavor []model.DishFlavor) error
	// DeleteByDishId 根据菜品id删除口味数据
	DeleteByDishId(db tx.Transaction, dishId uint64) error
	// GetByDishId 根据菜品id查询对应的口味数据
	GetByDishId(db tx.Transaction, dishId uint64) ([]model.DishFlavor, error)
	// Update 修改口味数据
	Update(db tx.Transaction, flavor model.DishFlavor) error
}
