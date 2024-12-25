package response

import "time"

type ShopCartVO struct {
	Id         uint64    `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	Name       string    `json:"name"`
	Image      string    `json:"image"`
	UserId     uint64    `json:"user_id"`
	DishId     uint64    `json:"dish_id"`
	SetmealId  uint64    `json:"setmeal_id"`
	DishFlavor string    `json:"dish_flavor"`
	Number     int       `json:"number"`
	Amount     float64   `json:"amount"`
	CreateTime time.Time `json:"create_tiem"`
}
