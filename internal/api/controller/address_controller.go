package controller

import (
	"net/http"
	"take-out/common"
	"take-out/common/e"
	"take-out/global"
	"take-out/internal/api/request"
	"take-out/internal/service"

	"github.com/gin-gonic/gin"
)

type AddressConteroller struct {
	service service.AddressService
}

func NewAddressConteroller(service service.AddressService) AddressConteroller {
	return AddressConteroller{service: service}
}

func (ac AddressConteroller) CreateAddress(c *gin.Context) {
	var addressBook request.AddressBookDTO
	if err := c.ShouldBindJSON(&addressBook); err != nil {
		c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Data: "参数错误"})
		return
	}
	err := ac.service.CreateAddressBook(c, addressBook)
	if err != nil {
		global.Log.Error("添加地址失败", err.Error())
		c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Data: "添加地址失败"})
		return
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: "请求成功"})
}

func (ac AddressConteroller) AddressList(c *gin.Context) {
	addList, err := ac.service.GetAddressListByUserId(c)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: addList})
}

func (ac AddressConteroller) DefaultAddress(c *gin.Context) {
	addDefa, err := ac.service.GetDefaultAddress(c)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: addDefa})
}

func (ac AddressConteroller) EditAddressById(c *gin.Context) {
	var addressBook request.AddressBookDTO
	if err := c.ShouldBindJSON(&addressBook); err != nil {
		c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Data: "参数错误"})
		return
	}
	err := ac.service.EditAddressById(c, addressBook)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: "请求成功"})
}

func (ac AddressConteroller) DeleteAddressById(c *gin.Context) {
	id := c.Query("id")
	err := ac.service.DeleteAddressById(c, id)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: "请求成功"})
}

func (ac AddressConteroller) GetAddressById(c *gin.Context) {
	id := c.Param("id")
	address, err := ac.service.GetAddressById(c, id)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Data: address, Msg: "请求成功"})
}

func (ac AddressConteroller) SetDefaultAddress(c *gin.Context) {
	var AddressIDDTO request.AddressIDDTO
	if err := c.ShouldBindJSON(&AddressIDDTO); err != nil {
		c.JSON(http.StatusBadRequest, common.Result{Code: e.ERROR, Data: "参数错误"})
		return
	}
	err := ac.service.SetDefaultAddress(c, AddressIDDTO.ID)
	if err != nil {
		c.JSON(http.StatusOK, common.Result{Code: e.ERROR, Msg: err.Error()})
	}
	c.JSON(http.StatusOK, common.Result{Code: e.SUCCESS, Msg: "请求成功"})
}
