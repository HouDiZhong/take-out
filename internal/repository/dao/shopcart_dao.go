package dao

import (
	"take-out/internal/api/response"
	"take-out/internal/model"
	"take-out/internal/repository"

	"gorm.io/gorm"
)

type ShopCartDao struct {
	db *gorm.DB
}

func NewShopCartRepo(db *gorm.DB) repository.ShopCartRepo {
	return ShopCartDao{db: db}
}

func (s ShopCartDao) AddShopCart(shopCart model.ShopCart) error {
	return s.db.Model(&shopCart).Create(&shopCart).Error
}

func (s ShopCartDao) GetShopCart(uid uint64) ([]response.ShopCartVO, error) {
	var shopCart []response.ShopCartVO
	err := s.db.Model(&model.ShopCart{}).
		Where("user_id = ?", uid).
		Find(&shopCart).Error
	return shopCart, err
}

func (s ShopCartDao) DelShopCartByDishId(uid uint64, did string) error {
	return s.db.Model(&model.ShopCart{}).Where("user_id = ? and dish_id = ?", uid, did).Delete(&model.ShopCart{}).Error
}

func (s ShopCartDao) DelShopCartBySetMealId(uid uint64, sid string) error {
	return s.db.Model(&model.ShopCart{}).Where("user_id = ? and setmeal_id = ?", uid, sid).Delete(&model.ShopCart{}).Error
}

func (s ShopCartDao) CleanShopCart(uid uint64) error {
	return s.db.Model(&model.ShopCart{}).Where("user_id = ?", uid).Delete(&model.ShopCart{}).Error
}
