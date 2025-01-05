package tasks

import (
	"log/slog"
	"take-out/common/enum"
	"take-out/global"
	"take-out/internal/api/request"
	"take-out/internal/repository/dao"
	"time"

	"gorm.io/gorm"
)

func ProcessTimeoutOrder() {

	global.DB.Transaction(func(tx *gorm.DB) error {

		db := dao.NewTXOrderDao(tx)

		orders, err := db.QueryTimeoutOrders(
			enum.OrderStatusUnpaid,
			time.Now().Add(-time.Minute*15),
		)

		if err != nil {
			slog.Error("自动查询订单状态任务失败", "Err:", err.Error())
			return err
		}

		if len(orders) > 0 {
			var oids []uint64
			for _, order := range orders {
				oids = append(oids, order.ID)
			}
			o := request.OrderStatus{
				Status:     enum.OrderStatusCancel,
				CancelTime: time.Now(),
			}
			err = db.UpdateOrderStatus(oids, o)
			if err != nil {
				slog.Error("自动更新订单状态任务失败", "Err:", err.Error())
				return err
			}
			return err
		}

		return nil
	})

}
