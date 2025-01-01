package repository

import (
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/model"
)

type ShopCartRepo interface {
	AddShopCart(shopCart model.ShopCart) error
	GetShopCartAll(uid uint64) ([]response.ShopCartVO, error)
	DelShopCart(uid uint64, shopCart request.ShopCartDTO) error
	UpdateShopCart(uid uint64, shopCart request.ShopCartDTO, number int) error
	GetShopCart(uid uint64, shopCart request.ShopCartDTO) ([]response.ShopCartVO, error)
	CleanShopCart(uid uint64) error
}
