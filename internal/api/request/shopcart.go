package request

type ShopCartDTO struct {
	DishFlavor string `json:"dishFlavor"`
	DishID     string `json:"dishId"`
	SetmealID  string `json:"setmealId"`
}
