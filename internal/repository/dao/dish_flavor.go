package dao

import (
	"take-out/global/tx"
	"take-out/internal/model"
	"take-out/internal/repository"
)

type DishFlavorDao struct {
}

func (d *DishFlavorDao) Update(transactions tx.Transaction, flavor model.DishFlavor) error {
	db, err := tx.GetGormDB(transactions)
	if err != nil {
		return err
	}
	err = db.Model(&model.DishFlavor{Id: flavor.Id}).Updates(flavor).Error
	return err
}

func (d *DishFlavorDao) InsertBatch(transactions tx.Transaction, flavor []model.DishFlavor) error {
	db, err := tx.GetGormDB(transactions)
	if err != nil {
		return err
	}
	// 批量插入口味数据
	err = db.Create(&flavor).Error
	return err
}

func (d *DishFlavorDao) DeleteByDishId(transactions tx.Transaction, dishId uint64) error {
	db, err := tx.GetGormDB(transactions)
	if err != nil {
		return err
	}
	// 根据dishId删除对应的口味数据
	err = db.Where("dish_id = ?", dishId).Delete(&model.DishFlavor{}).Error
	return err
}

func (d *DishFlavorDao) GetByDishId(db tx.Transaction, dishId uint64) ([]model.DishFlavor, error) {
	//TODO implement me
	panic("implement me")
}

// NewDishFlavorDao db 是上个事务创建出来的
func NewDishFlavorDao() repository.DishFlavorRepo {
	return &DishFlavorDao{}
}
