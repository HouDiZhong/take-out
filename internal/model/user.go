package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         int64     `json:"id" gorm:"primaryKey;AUTO_INCREMENT"`
	OpenID     string    `json:"openid" gorm:"column:openid"`
	Name       string    `json:"name"`
	Phone      string    `json:"phone"`
	Sex        string    `json:"sex"`
	IDNumber   string    `json:"id_number"`
	Avatar     string    `json:"avatar"`
	CreateTime time.Time `json:"create_time"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreateTime = time.Now()

	return nil
}
