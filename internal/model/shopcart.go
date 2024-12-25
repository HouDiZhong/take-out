package model

import (
	"time"

	"gorm.io/gorm"
)

type ShopCart struct {
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

func (e *ShopCart) BeforeCreate(tx *gorm.DB) error {
	// 自动填充 创建时间、创建人、更新时间、更新用户
	e.CreateTime = time.Now()
	// 从上下文获取用户信息
	/* value := tx.Statement.Context.Value(enum.CurrentId)
	fmt.Printf("---------------------------------------------uid: %v\n", value)
	if uid, ok := value.(uint64); ok {
		e.UserId = uid
	} */
	return nil
}

func (e *ShopCart) TableName() string {
	return "shopping_cart"
}
