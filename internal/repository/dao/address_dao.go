package dao

import (
	"take-out/internal/model"
	"take-out/internal/repository"

	"gorm.io/gorm"
)

type AddressDao struct {
	db *gorm.DB
}

func NewAddressDao(db *gorm.DB) repository.AddressRepo {
	return &AddressDao{db: db}
}

func (a *AddressDao) CreateAddress(addressBook *model.AddressBook) error {
	return a.db.Create(&addressBook).Error
}

func (a *AddressDao) GetAddressListByUserId(userId uint64) ([]model.AddressBook, error) {
	var addList = []model.AddressBook{}
	err := a.db.Model(&model.AddressBook{}).Where("user_id = ?", userId).Find(&addList).Error
	return addList, err
}

func (a *AddressDao) GetDefaultAddress(userId uint64) (model.AddressBook, error) {
	var defaAddress = model.AddressBook{}
	err := a.db.Model(&model.AddressBook{}).Where("user_id = ? and is_default = ?", userId, 1).First(&defaAddress).Error
	return defaAddress, err
}

func (a *AddressDao) EditAddressById(addressBook *model.AddressBook) error {
	return a.db.Model(&model.AddressBook{}).Where("id = ?", addressBook.ID).Updates(&addressBook).Error
}

func (a *AddressDao) DeleteAddressById(uid uint64, id string) error {
	return a.db.Model(&model.AddressBook{}).Where("user_id = ? and id = ?", uid, id).Delete(&model.AddressBook{}).Error
}

func (a *AddressDao) GetAddressById(uid uint64, id string) (model.AddressBook, error) {
	var address = model.AddressBook{}
	err := a.db.Model(&model.AddressBook{}).Where("user_id = ? and id = ?", uid, id).First(&address).Error

	return address, err
}

func (a *AddressDao) ClearDefalutAddress(uid uint64) error {
	return a.db.Model(&model.AddressBook{}).Where("user_id = ? and is_default = ?", uid, 1).Update("is_default", 0).Error
}

func (a *AddressDao) SetDefaultAddress(uid uint64, id uint64) error {
	err := a.ClearDefalutAddress(uid)
	if err != nil {
		return err
	}
	return a.db.Model(&model.AddressBook{}).Where("user_id = ? and id = ?", uid, id).Update("is_default", 1).Error
}
