package model

type AddressBook struct {
	ID           uint64 `json:"id" grom:"primaryKey;autoIncrement"`
	UserID       uint64 `json:"user_id"`
	Consignee    string `json:"consignee"`
	Sex          string `json:"sex"`
	Phone        string `json:"phone"`
	ProvinceCode string `json:"province_code"`
	ProvinceName string `json:"province_name"`
	CityCode     string `json:"city_code"`
	CityName     string `json:"city_name"`
	DistrictCode string `json:"district_code"`
	DistrictName string `json:"district_name"`
	Detail       string `json:"detail"`
	Label        string `json:"label"`
	IsDefault    bool   `json:"is_default"`
}
