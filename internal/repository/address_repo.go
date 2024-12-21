package repository

import "take-out/internal/model"

type AddressRepo interface {
	CreateAddress(addressBook *model.AddressBook) error
	GetAddressListByUserId(userId uint64) ([]model.AddressBook, error)
	GetDefaultAddress(userId uint64) (model.AddressBook, error)
	EditAddressById(addressBook *model.AddressBook) error
	DeleteAddressById(uid uint64, id string) error
	GetAddressById(uid uint64, id string) (model.AddressBook, error)
	SetDefaultAddress(uid uint64, id uint64) error
	ClearDefalutAddress(uid uint64) error
}
