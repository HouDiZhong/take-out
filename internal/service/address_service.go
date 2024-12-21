package service

import (
	"take-out/common/e"
	"take-out/common/enum"
	"take-out/internal/api/request"
	"take-out/internal/model"
	"take-out/internal/repository"

	"github.com/gin-gonic/gin"
)

type AddressService interface {
	CreateAddressBook(c *gin.Context, dto request.AddressBookDTO) error
	GetAddressListByUserId(c *gin.Context) ([]model.AddressBook, error)
	GetDefaultAddress(c *gin.Context) (model.AddressBook, error)
	EditAddressById(c *gin.Context, dto request.AddressBookDTO) error
	DeleteAddressById(c *gin.Context, id string) error
	GetAddressById(c *gin.Context, id string) (model.AddressBook, error)
	SetDefaultAddress(c *gin.Context, id uint64) error
}

type AddressServiceImpl struct {
	repo repository.AddressRepo
}

func NewAddressService(repo repository.AddressRepo) AddressService {
	return AddressServiceImpl{repo: repo}
}

func (a AddressServiceImpl) CreateAddressBook(c *gin.Context, dto request.AddressBookDTO) error {
	address := model.AddressBook{
		Consignee:    dto.Consignee,
		Sex:          dto.Sex,
		Phone:        dto.Phone,
		ProvinceCode: dto.ProvinceCode,
		ProvinceName: dto.ProvinceName,
		CityCode:     dto.CityCode,
		CityName:     dto.CityName,
		DistrictCode: dto.DistrictCode,
		DistrictName: dto.DistrictName,
		Detail:       dto.Detail,
		Label:        dto.Label,
		IsDefault:    dto.IsDefault,
	}
	if userId, exists := c.Get(enum.CurrentId); exists {
		address.UserID = userId.(uint64)
	}
	if address.IsDefault {
		if err := a.repo.ClearDefalutAddress(address.UserID); err != nil {
			return err
		}
	}
	return a.repo.CreateAddress(&address)
}

func (a AddressServiceImpl) GetAddressListByUserId(c *gin.Context) ([]model.AddressBook, error) {
	if userId, exists := c.Get(enum.CurrentId); exists {
		return a.repo.GetAddressListByUserId(userId.(uint64))
	}

	return nil, e.Error_ACCOUNT_NOT_FOUND
}

func (a AddressServiceImpl) GetDefaultAddress(c *gin.Context) (model.AddressBook, error) {
	if userId, exists := c.Get(enum.CurrentId); exists {
		return a.repo.GetDefaultAddress(userId.(uint64))
	}

	return model.AddressBook{}, e.Error_ACCOUNT_NOT_FOUND
}

func (a AddressServiceImpl) EditAddressById(c *gin.Context, dto request.AddressBookDTO) error {
	address := model.AddressBook{
		ID:           dto.ID,
		Consignee:    dto.Consignee,
		Sex:          dto.Sex,
		Phone:        dto.Phone,
		ProvinceCode: dto.ProvinceCode,
		ProvinceName: dto.ProvinceName,
		CityCode:     dto.CityCode,
		CityName:     dto.CityName,
		DistrictCode: dto.DistrictCode,
		DistrictName: dto.DistrictName,
		Detail:       dto.Detail,
		Label:        dto.Label,
		IsDefault:    dto.IsDefault,
	}
	if userId, exists := c.Get(enum.CurrentId); exists {
		address.UserID = userId.(uint64)
	}
	if address.IsDefault {
		if err := a.repo.ClearDefalutAddress(address.UserID); err != nil {
			return err
		}
	}
	return a.repo.EditAddressById(&address)
}

func (a AddressServiceImpl) DeleteAddressById(c *gin.Context, id string) error {
	if userId, exists := c.Get(enum.CurrentId); exists {
		return a.repo.DeleteAddressById(userId.(uint64), id)
	}

	return e.Error_ACCOUNT_NOT_FOUND
}

func (a AddressServiceImpl) GetAddressById(c *gin.Context, id string) (model.AddressBook, error) {
	if userId, exists := c.Get(enum.CurrentId); exists {
		return a.repo.GetAddressById(userId.(uint64), id)
	}

	return model.AddressBook{}, e.Error_ACCOUNT_NOT_FOUND
}

func (a AddressServiceImpl) SetDefaultAddress(c *gin.Context, id uint64) error {
	if userId, exists := c.Get(enum.CurrentId); exists {
		return a.repo.SetDefaultAddress(userId.(uint64), id)
	}

	return e.Error_ACCOUNT_NOT_FOUND
}
