package enum

type CommonStatus = int

const (
	ENABLE  CommonStatus = 1 // 启用
	DISABLE CommonStatus = 0 // 禁用
)

// 1待付款 2待接单 3已接单 4派送中 5已完成 6已取消
const (
	OrderStatusUnpaid int8 = iota + 1 // 未支付
	OrderStatusWait                   // 待接单
	OrderStatusAccept                 // 已接单
	OrderStatusSend                   // 派送中
	OrderStatusFinish                 // 已完成
	OrderStatusCancel                 // 已取消
)

const TimeLayout = "2006-01-02 15:04:05"
