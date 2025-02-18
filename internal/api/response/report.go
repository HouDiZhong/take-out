package response

type SalesTop10ReportVO struct {
	NameList   string `json:"nameList"`
	NumberList string `json:"numberList"`
}

type OrderReportVO struct {
	DateList            string `json:"dateList"`            // 日期列表
	OrderCountList      string `json:"orderCountList"`      // 订单数列表
	ValidOrderCountList string `json:"validOrderCountList"` // 有效订单数列表
	OrderNuberReportVO
}

type OrderNuberReportVO struct {
	OrderCompletionRate float64 `json:"orderCompletionRate"` // 订单完成率
	ValidOrderCount     int     `json:"validOrderCount"`     // 有效订单数
	TotalOrderCount     int     `json:"totalOrderCount"`     // 订单总数
}

type LocalOrderVO struct {
	Date            string
	ValidOrderCount string
	TotalOrderCount string
}

type TurnoverReportVO struct {
	DateList     string `json:"dateList"`     // 日期列表
	TurnoverList string `json:"turnoverList"` // 营业额
}

type LocalTurnoverVO struct {
	Date          string
	TurnoverCount string
}

type UserReportVO struct {
	DateList      string `json:"dateList"`      // 日期列表
	NewUserList   string `json:"newUserList"`   // 新增用户
	TotalUserList string `json:"totalUserList"` // 总用户量
}

type LocalUsertVO struct {
	Date           string
	NewUserCount   string
	TotalUserCount string
}

type ExcelVO struct {
	Times              string  // 日期
	Turnovers          float64 // 营业额
	ValidOrders        int     // 有效订单
	OrderStatusNumbers float64 // 订单完成率
	UnitPrices         float64 // 平均客单价
	TotalOrders        int     // 所有订单
}

type EveryUserVO struct {
	Times    string
	NewUsers int
}

type ExcelsVO struct {
	Date string
	ExcelVO
	EveryUserVO
}

func (l ExcelVO) GetDate() string {
	return l.Times[:10]
}
func (l ExcelsVO) GetDate() string {
	return l.Date[:10]
}
func (l EveryUserVO) GetDate() string {
	return l.Times[:10]
}
func (l LocalTurnoverVO) GetDate() string {
	return l.Date[:10]
}
func (l LocalOrderVO) GetDate() string {
	return l.Date[:10]
}

func (l LocalUsertVO) GetDate() string {
	return l.Date[:10]
}

type ILocalOrder interface {
	GetDate() string
}
