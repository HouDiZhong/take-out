package repository

import (
	"take-out/internal/api/response"
	"take-out/internal/model"
)

type ShopCartRepo interface {
	AddShopCart(shopCart model.ShopCart) error
	GetShopCart(uid uint64) ([]response.ShopCartVO, error)
	DelShopCartByDishId(uid uint64, did string) error
	DelShopCartBySetMealId(uid uint64, sid string) error
	CleanShopCart(uid uint64) error
}
