package controller

import (
	"net/http"
	"take-out/common"

	"github.com/gin-gonic/gin"
)

type AddressConteroller struct{}

func NewAddressConteroller() AddressConteroller {
	return AddressConteroller{}
}

func (ac AddressConteroller) CreateAddress(c *gin.Context) {

}

func (ac AddressConteroller) AddressList(c *gin.Context) {
	c.JSON(http.StatusOK, common.Result{Code: 1, Data: "请求成功"})
}

func (ac AddressConteroller) DefaultAddress(c *gin.Context) {

}

func (ac AddressConteroller) EditAddressById(c *gin.Context) {

}

func (ac AddressConteroller) DeleteAddressById(c *gin.Context) {

}

func (ac AddressConteroller) GetAddressById(c *gin.Context) {

}

func (ac AddressConteroller) SetDefaultAddress(c *gin.Context) {

}
