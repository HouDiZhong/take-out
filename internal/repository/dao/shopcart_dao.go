package dao

import (
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/model"
	"take-out/internal/repository"

	"gorm.io/gorm"
)

var modelShopCart []model.ShopCart

type ShopCartDao struct {
	db *gorm.DB
}

func NewShopCartRepo(db *gorm.DB) repository.ShopCartRepo {
	return ShopCartDao{db: db}
}

func (s ShopCartDao) AddShopCart(shopCart model.ShopCart) error {
	return s.db.Model(&modelShopCart).Create(&shopCart).Error
}

func (s ShopCartDao) GetShopCartAll(uid uint64) ([]response.ShopCartVO, error) {
	var shopCart []response.ShopCartVO
	err := s.db.Model(&modelShopCart).
		Where("user_id = ?", uid).
		Find(&shopCart).Error
	return shopCart, err
}

func (s ShopCartDao) GetShopCart(uid uint64, shopCart request.ShopCartDTO) ([]response.ShopCartVO, error) {
	var shopCarts []response.ShopCartVO
	db := s.DynamicSplicingWhere(uid, shopCart)
	err := db.Model(&modelShopCart).Find(&shopCarts).Error
	return shopCarts, err
}

func (s ShopCartDao) UpdateShopCart(uid uint64, shopCart request.ShopCartDTO, number int) error {
	db := s.DynamicSplicingWhere(uid, shopCart)
	return db.Model(&modelShopCart).Update("number", number).Error
}

func (s ShopCartDao) DynamicSplicingWhere(uid uint64, shopCart request.ShopCartDTO) *gorm.DB {
	db := s.db.Where("user_id = ?", uid)
	if shopCart.DishID != "" {
		db = db.Where("dish_id = ?", shopCart.DishID)
	}
	if shopCart.SetmealID != "" {
		db = db.Where("setmeal_id = ?", shopCart.SetmealID)
	}
	if shopCart.DishFlavor != "" {
		db = db.Where("dish_flavor = ?", shopCart.DishFlavor)
	}
	return db
}

func (s ShopCartDao) DelShopCart(uid uint64, shopCart request.ShopCartDTO) error {
	db := s.DynamicSplicingWhere(uid, shopCart)
	return db.Model(&modelShopCart).Delete(&modelShopCart).Error
}

func (s ShopCartDao) CleanShopCart(uid uint64) error {
	return s.db.Model(&modelShopCart).Where("user_id = ?", uid).Delete(&modelShopCart).Error
}
