package service

import (
	"context"
	"take-out/common"
	"take-out/common/e"
	"take-out/common/enum"
	"take-out/common/utils"
	"take-out/global"
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/model"
	"take-out/internal/repository"

	"github.com/gin-gonic/gin"
)

type IEmployeeService interface {
	Login(context.Context, request.EmployeeLogin) (*response.EmployeeLogin, error)
	Logout(ctx *gin.Context) error
	EditPassword(context.Context, request.EmployeeEditPassword) error
	CreateEmployee(ctx context.Context, employee request.EmployeeDTO) error
	PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error)
	SetStatus(ctx context.Context, id uint64, status int) error
	UpdateEmployee(ctx context.Context, dto request.EmployeeDTO) error
	GetById(ctx context.Context, id uint64) (model.Employee, error)
}

type EmployeeImpl struct {
	repo repository.EmployeeRepo
}

func (ei *EmployeeImpl) GetById(ctx context.Context, id uint64) (model.Employee, error) {
	employee, err := ei.repo.GetById(ctx, id)
	employee.Password = "***"
	return *employee, err
}

func (ei *EmployeeImpl) UpdateEmployee(ctx context.Context, dto request.EmployeeDTO) error {
	// 构建model实体进行更新
	err := ei.repo.Update(ctx, model.Employee{
		Id:       dto.Id,
		Username: dto.UserName,
		Name:     dto.Name,
		Phone:    dto.Phone,
		Sex:      dto.Sex,
		IdNumber: dto.IdNumber,
	})
	return err
}

func (ei *EmployeeImpl) SetStatus(ctx context.Context, id uint64, status int) error {
	// 设置用户状态,构造实体
	entity := model.Employee{Id: id, Status: status}
	err := ei.repo.UpdateStatus(ctx, entity)
	return err
}

func (ei *EmployeeImpl) PageQuery(ctx context.Context, dto request.EmployeePageQueryDTO) (*common.PageResult, error) {
	// 分页查询
	pageResult, err := ei.repo.PageQuery(ctx, dto)
	// 屏蔽敏感信息
	if employees, ok := pageResult.Records.([]model.Employee); ok {
		// 替换敏感信息
		for key := range employees {
			employees[key].Password = "****"
			employees[key].IdNumber = "****"
			employees[key].Phone = "****"
		}
		// 重新赋值
		pageResult.Records = employees
	}

	return pageResult, err
}

func (ei *EmployeeImpl) CreateEmployee(ctx context.Context, employeeDTO request.EmployeeDTO) error {
	var err error
	// 1.新增员工,构建员工基础信息
	entity := model.Employee{
		Id:       employeeDTO.Id,
		IdNumber: employeeDTO.IdNumber,
		Name:     employeeDTO.Name,
		Phone:    employeeDTO.Phone,
		Sex:      employeeDTO.Sex,
		Username: employeeDTO.UserName,
	}
	// 新增用户为启用状态
	entity.Status = enum.ENABLE
	// 新增用户初始密码为123456
	entity.Password = utils.MD5V("123456", "", 0)
	// 新增用户
	err = ei.repo.Insert(ctx, entity)
	return err
}

func (ei *EmployeeImpl) EditPassword(ctx context.Context, employeeEdit request.EmployeeEditPassword) error {
	// 1.获取员工信息
	employee, err := ei.repo.GetById(ctx, employeeEdit.EmpId)
	if err != nil {
		return err

	}
	// 校验用户老密码
	if employee == nil {
		return e.Error_ACCOUNT_NOT_FOUND
	}
	oldHashPassword := utils.MD5V(employeeEdit.OldPassword, "", 0)
	if employee.Password != oldHashPassword {
		return e.Error_PASSWORD_ERROR
	}
	// 修改员工密码
	newHashPassword := utils.MD5V(employeeEdit.NewPassword, "", 0) // 使用新密码生成哈希值
	err = ei.repo.Update(ctx, model.Employee{
		Id:       employeeEdit.EmpId,
		Password: newHashPassword,
	})
	return err
}

func (ei *EmployeeImpl) Logout(ctx *gin.Context) error {
	id, exists := ctx.Get(enum.CurrentId)

	if exists {
		_, err := utils.DeleteRedisToken(id.(uint64), global.Config.Jwt.Admin.Secret)
		return err
	}

	return e.Error_ACCOUNT_NOT_FOUND
}

func (ei *EmployeeImpl) Login(ctx context.Context, employeeLogin request.EmployeeLogin) (*response.EmployeeLogin, error) {
	// 1.查询用户是否存在
	employee, err := ei.repo.GetByUserName(ctx, employeeLogin.UserName)
	if err != nil || employee == nil {
		return nil, e.Error_ACCOUNT_NOT_FOUND
	}
	// 2.校验密码
	password := utils.MD5V(employeeLogin.Password, "", 0)
	if password != employee.Password {
		return nil, e.Error_PASSWORD_ERROR
	}
	// 3.校验状态
	if employee.Status == enum.DISABLE {
		return nil, e.Error_ACCOUNT_LOCKED
	}
	// 生成Token
	jwtConfig := global.Config.Jwt.Admin
	token, err := utils.GenerateToken(employee.Id, jwtConfig)
	if err != nil {
		return nil, err
	}
	// 4.构造返回数据
	resp := response.EmployeeLogin{
		Id:       employee.Id,
		Name:     employee.Name,
		Token:    token,
		UserName: employee.Username,
	}
	return &resp, nil
}

func NewEmployeeService(repo repository.EmployeeRepo) IEmployeeService {
	return &EmployeeImpl{repo: repo}
}
