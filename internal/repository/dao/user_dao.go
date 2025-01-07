package dao

import (
	"take-out/common/utils"
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/model"
	"take-out/internal/repository"

	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repository.UserRepo {
	return &UserDao{db: db}
}

func (u *UserDao) FindByOpenId(openId string) (model.User, error) {
	user := model.User{}
	err := u.db.Model(&user).Where("openid = ?", openId).Find(&user).Error
	return user, err
}

func (u *UserDao) CreateUser(user *model.User) error {
	return u.db.Create(&user).Error
}

func (u *UserDao) GetNewUserNumber() (int64, error) {
	var number int64
	err := u.db.Model(&model.User{}).Where("create_time = ?", utils.ToDay()).Count(&number).Error
	return number, err
}

func (o *UserDao) UserReport(dto request.ReportQuestDTO) ([]response.LocalUsertVO, error) {
	var dbData []response.LocalUsertVO
	err := o.db.Raw(`
		SELECT date, NewUserCount, MAX(TotalUserCount) AS TotalUserCount
		FROM (
			SELECT
				DATE(create_time) AS date,
				COUNT(*) OVER (PARTITION BY DATE(create_time)) AS NewUserCount,
				COUNT(*) OVER (ORDER BY DATE(create_time) ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS TotalUserCount
			FROM user
			WHERE create_time >= ? AND create_time < ?
		) AS daily_stats
		GROUP BY date, NewUserCount
		ORDER BY date;
	`, dto.Begin, dto.End).Scan(&dbData).Error

	return dbData, err
}

func (o *UserDao) EveryUserReport(dto request.ReportQuestDTO) ([]response.EveryUserVO, error) {
	var dbData []response.EveryUserVO
	err := o.db.Raw(`
		SELECT
			DATE(create_time) AS Times,
			COUNT(*) as NewUsers
		FROM user
		WHERE create_time >= ? AND create_time < ?
		GROUP BY Times;
	`, dto.Begin, dto.End).Scan(&dbData).Error

	return dbData, err
}
