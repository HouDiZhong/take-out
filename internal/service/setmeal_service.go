package service

import (
	"context"
	"sort"
	"strconv"
	"take-out/common"
	"take-out/common/enum"
	"take-out/internal/api/request"
	"take-out/internal/api/response"
	"take-out/internal/model"
	"take-out/internal/repository"
)

type ISetMealService interface {
	SaveWithDish(ctx context.Context, dto request.SetMealDTO) error
	UpdateWithDish(ctx context.Context, dto request.UpSetMealDTO) error
	DeleteWithDish(ctx context.Context, ids []string) error
	PageQuery(ctx context.Context, dto request.SetMealPageQueryDTO) (*common.PageResult, error)
	OnOrClose(ctx context.Context, id uint64, status int) error
	GetByIdWithDish(ctx context.Context, dishId uint64) (response.SetMealWithDishByIdVo, error)
}

type SetMealServiceImpl struct {
	repo            repository.SetMealRepo
	setMealDishRepo repository.SetMealDishRepo
}

func (s *SetMealServiceImpl) GetByIdWithDish(ctx context.Context, setMealId uint64) (response.SetMealWithDishByIdVo, error) {
	// 获取事务
	transaction := s.repo.Transaction(ctx)
	// 开始事务
	if err := transaction.Begin(); err != nil {
		return response.SetMealWithDishByIdVo{}, err
	}
	defer func() {
		if err := recover(); err != nil {
			transaction.Rollback()
		}
	}()
	// 单独查询套餐
	setMeal, err := s.repo.GetByIdWithDish(transaction, setMealId)
	if err != nil {
		return response.SetMealWithDishByIdVo{}, err
	}
	// 查询中间表记录的套餐菜品信息
	dishList, err := s.setMealDishRepo.GetBySetMealId(transaction, setMealId)
	if err != nil {
		return response.SetMealWithDishByIdVo{}, err
	}
	// 事务提交
	if err = transaction.Commit(); err != nil {
		return response.SetMealWithDishByIdVo{}, err
	}
	// 数据组装
	setMealVo := response.SetMealWithDishByIdVo{
		Id:            setMeal.Id,
		CategoryId:    setMeal.CategoryId,
		CategoryName:  setMeal.Name,
		Description:   setMeal.Description,
		Image:         setMeal.Image,
		Name:          setMeal.Name,
		Price:         setMeal.Price,
		Status:        setMeal.Status,
		SetmealDishes: dishList,
		UpdateTime:    setMeal.UpdateTime,
	}
	return setMealVo, nil
}

func (s *SetMealServiceImpl) OnOrClose(ctx context.Context, id uint64, status int) error {
	err := s.repo.SetStatus(ctx, id, status)
	return err
}

func (s *SetMealServiceImpl) PageQuery(ctx context.Context, dto request.SetMealPageQueryDTO) (*common.PageResult, error) {
	result, err := s.repo.PageQuery(ctx, dto)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *SetMealServiceImpl) DeleteWithDish(ctx context.Context, ids []string) error {
	// 开启事务进行存储
	transaction := s.repo.Transaction(ctx)
	// 开始事务
	if err := transaction.Begin(); err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			transaction.Rollback()
		}
	}()
	// 删除套餐表里的数据
	err := s.repo.Delete(transaction, ids)
	if err != nil {
		return err
	}
	// 删除中间表里的数据
	err = s.setMealDishRepo.DeledeSetMealBatch(transaction, ids)
	if err != nil {
		return err
	}
	return transaction.Commit()
}

// SaveWithDish 保存套餐和菜品信息
func (s *SetMealServiceImpl) SaveWithDish(ctx context.Context, dto request.SetMealDTO) error {
	// 转换dto为model开启事务进行保存
	price, _ := strconv.ParseFloat(dto.Price, 64)
	setmeal := model.SetMeal{
		Id:          dto.Id,
		CategoryId:  dto.CategoryId,
		Name:        dto.Name,
		Price:       price,
		Status:      enum.ENABLE,
		Description: dto.Description,
		Image:       dto.Image,
	}
	// 开启事务进行存储
	transaction := s.repo.Transaction(ctx)
	// 开始事务
	if err := transaction.Begin(); err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			transaction.Rollback()
		}
	}()
	// 先插入套餐数据信息，并得到返回的主键id值
	err := s.repo.Insert(transaction, &setmeal)
	if err != nil {
		return err
	}
	for i := range dto.SetMealDishs {
		dto.SetMealDishs[i].SetmealId = setmeal.Id
	}
	// 向中间表插入数据
	err = s.setMealDishRepo.InsertBatch(transaction, dto.SetMealDishs)
	if err != nil {
		return err
	}
	return transaction.Commit()
}

// SaveWithDish 更新套餐和菜品信息
func (s *SetMealServiceImpl) UpdateWithDish(ctx context.Context, dto request.UpSetMealDTO) error {
	setmeal := model.SetMeal{
		Id:          dto.Id,
		CategoryId:  dto.CategoryId,
		Name:        dto.Name,
		Price:       dto.Price,
		Status:      enum.ENABLE,
		Description: dto.Description,
		Image:       dto.Image,
	}
	// 开启事务进行存储
	transaction := s.repo.Transaction(ctx)
	// 开始事务
	if err := transaction.Begin(); err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			transaction.Rollback()
		}
	}()
	// 先插入套餐数据信息，并得到返回的主键id值
	err := s.repo.Update(transaction, &setmeal)
	if err != nil {
		return err
	}
	for i := range dto.SetMealDishs {
		dto.SetMealDishs[i].SetmealId = setmeal.Id
	}
	dishList, _ := s.setMealDishRepo.GetBySetMealId(transaction, setmeal.Id)
	// 中间表数据处理
	// 两次数据是否相同
	equal, diff, add := s.MiddleTable(dishList, dto.SetMealDishs)
	if !equal { // 不相同时处理
		if add != nil {
			// 向表中插入增集
			err = s.setMealDishRepo.InsertBatch(transaction, add)
		}
		if diff != nil {
			// 将表中的差集删除
			err = s.setMealDishRepo.DeledeBatch(transaction, diff)
		}
		if err != nil {
			return err
		}
	}

	return transaction.Commit()
}

func sortSlice(slice []model.SetMealDish) {
	sort.Slice(slice, func(i, j int) bool {
		return int(slice[i].DishId) < int(slice[j].DishId)
	})
}

// 比较连个集合，并返回增集和差集
func (s *SetMealServiceImpl) MiddleTable(old, new []model.SetMealDish) (bool, []model.SetMealDish, []model.SetMealDish) {
	sortSlice(old)
	sortSlice(new)
	// 如果两个切片的长度不相等，直接返回不相等，并计算差集和增集
	if len(old) != len(new) {
		diffSet := DifferenceSet(old, new)
		additionalSet := DifferenceSet(new, old)
		return false, diffSet, additionalSet
	}

	// 如果长度相等，遍历比较每个元素
	for i := range old {
		if old[i].DishId != new[i].DishId {
			diffSet := DifferenceSet(old, new)
			additionalSet := DifferenceSet(new, old)
			return false, diffSet, additionalSet
		}
	}

	// 如果所有元素都相等，返回相等
	return true, nil, nil
}

// DifferenceSet 计算两个切片的差集
func DifferenceSet(old, new []model.SetMealDish) []model.SetMealDish {
	elementMap := make(map[uint64]bool)
	var diffSet []model.SetMealDish

	// 将 new 的所有元素添加到 elementMap 中
	for _, elem := range new {
		elementMap[elem.DishId] = true
	}

	// 遍历 old，找出不在 new 中的元素
	for _, elem := range old {
		if !elementMap[elem.DishId] {
			diffSet = append(diffSet, elem)
		}
	}

	return diffSet
}

func NewSetMealService(repo repository.SetMealRepo, setMealDishRepo repository.SetMealDishRepo) ISetMealService {
	return &SetMealServiceImpl{repo: repo, setMealDishRepo: setMealDishRepo}
}
