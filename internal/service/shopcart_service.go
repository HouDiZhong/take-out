package service

import (
	"errors"
	"strconv"
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/model"
	"take-out/internal/repository"

	"github.com/gin-gonic/gin"
)

type ShopCartService interface {
	GetShopCartAll(uid uint64) ([]response.ShopCartVO, error)
	AddShopCart(c *gin.Context, uid uint64, shopCart request.ShopCartDTO) error
	DelShopCart(uid uint64, shopCart request.ShopCartDTO) error
	ClearnShopCart(uid uint64) error
}

type ShopCartServiceImpl struct {
	dishRepo repository.DishRepo
	smRepo   repository.SetMealRepo
	repo     repository.ShopCartRepo
}

func NewShopCartService(dishRepo repository.DishRepo, smRepo repository.SetMealRepo, repo repository.ShopCartRepo) ShopCartService {
	return ShopCartServiceImpl{dishRepo: dishRepo, smRepo: smRepo, repo: repo}
}

func (s ShopCartServiceImpl) GetShopCartAll(uid uint64) ([]response.ShopCartVO, error) {
	return s.repo.GetShopCartAll(uid)
}

func (s ShopCartServiceImpl) AddShopCart(c *gin.Context, uid uint64, shopCart request.ShopCartDTO) error {
	shopCartItem, err := s.repo.GetShopCart(uid, shopCart)
	if err != nil {
		return err
	} else {
		if len(shopCartItem) > 0 {
			return s.repo.UpdateShopCart(uid, shopCart, shopCartItem[0].Number+1)
		}
	}
	if shopCart.DishID != "" {
		id, _ := strconv.ParseUint(shopCart.DishID, 10, 64)
		dishVo, _ := s.dishRepo.GetById(c, id)
		shopCart := model.ShopCart{
			Name:       dishVo.Name,
			DishId:     dishVo.Id,
			Image:      dishVo.Image,
			DishFlavor: shopCart.DishFlavor,
			Number:     1,
			Amount:     dishVo.Price,
			UserId:     uid,
		}
		return s.repo.AddShopCart(shopCart)
	} else if shopCart.SetmealID != "" {
		setMealVo, _ := s.smRepo.QuerySetMealDesById(shopCart.SetmealID)
		shopCart := model.ShopCart{
			Name:       setMealVo.Name,
			SetmealId:  setMealVo.Id,
			Image:      setMealVo.Image,
			DishFlavor: shopCart.DishFlavor,
			Number:     1,
			Amount:     setMealVo.Price,
			UserId:     uid,
		}
		return s.repo.AddShopCart(shopCart)
	}
	return errors.New("没有id")
}

func (s ShopCartServiceImpl) DelShopCart(uid uint64, shopCart request.ShopCartDTO) error {
	if shopCart.DishID != "" && shopCart.SetmealID != "" {
		return s.repo.DelShopCart(uid, shopCart)
	}
	return errors.New("没有id")
}

func (s ShopCartServiceImpl) ClearnShopCart(uid uint64) error {
	return s.repo.CleanShopCart(uid)
}
