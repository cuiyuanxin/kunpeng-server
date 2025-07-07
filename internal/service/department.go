package service

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// DepartmentService 部门服务
type DepartmentService struct {
	db *gorm.DB
}

// NewDepartmentService 创建部门服务
func NewDepartmentService(db *gorm.DB) *DepartmentService {
	return &DepartmentService{
		db: db,
	}
}

// GetAll 获取所有部门
func (s *DepartmentService) GetAll() ([]model.Department, error) {
	var departments []model.Department
	if err := s.db.Find(&departments).Error; err != nil {
		return nil, fmt.Errorf("failed to get departments: %w", err)
	}
	return departments, nil
}

// GetByID 根据ID获取部门
func (s *DepartmentService) GetByID(id uint) (*model.Department, error) {
	var department model.Department
	if err := s.db.First(&department, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, fmt.Errorf("failed to get department: %w", err)
	}
	return &department, nil
}

// GetByName 根据名称获取部门
func (s *DepartmentService) GetByName(name string) (*model.Department, error) {
	var department model.Department
	if err := s.db.Where("name = ?", name).First(&department).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("department not found")
		}
		return nil, fmt.Errorf("failed to get department: %w", err)
	}
	return &department, nil
}

// Create 创建部门
func (s *DepartmentService) Create(req *model.DepartmentCreateRequest) (*model.Department, error) {
	// 检查部门名是否已存在
	var existDept model.Department
	if err := s.db.Where("name = ?", req.Name).First(&existDept).Error; err == nil {
		return nil, errors.New("department name already exists")
	}

	// 创建部门
	department := &model.Department{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		ManagerID:   req.ManagerID,
		Status:      1, // 默认启用
	}

	if err := s.db.Create(department).Error; err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}

	return department, nil
}

// Update 更新部门
func (s *DepartmentService) Update(id uint, req *model.DepartmentUpdateRequest) (*model.Department, error) {
	department, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Name != nil && *req.Name != "" {
		department.Name = *req.Name
	}
	if req.Description != nil && *req.Description != "" {
		department.Description = req.Description
	}
	if req.ParentID != nil {
		department.ParentID = *req.ParentID
	}
	if req.Status != nil {
		department.Status = *req.Status
	}

	if err := s.db.Save(department).Error; err != nil {
		return nil, fmt.Errorf("failed to update department: %w", err)
	}

	return department, nil
}

// Delete 删除部门
func (s *DepartmentService) Delete(id uint) error {
	// 检查是否有子部门
	var count int64
	if err := s.db.Model(&model.Department{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check child departments: %w", err)
	}
	if count > 0 {
		return errors.New("cannot delete department with child departments")
	}

	// 检查是否有用户
	if err := s.db.Model(&model.User{}).Where("department_id = ?", id).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check department users: %w", err)
	}
	if count > 0 {
		return errors.New("cannot delete department with users")
	}

	if err := s.db.Delete(&model.Department{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete department: %w", err)
	}
	return nil
}

// GetTree 获取部门树结构
func (s *DepartmentService) GetTree() ([]model.Department, error) {
	var departments []model.Department
	if err := s.db.Order("parent_id, id").Find(&departments).Error; err != nil {
		return nil, fmt.Errorf("failed to get departments: %w", err)
	}

	// 构建树结构
	return s.buildTree(departments, nil), nil
}

// buildTree 构建树结构
func (s *DepartmentService) buildTree(departments []model.Department, parentID *uint) []model.Department {
	var tree []model.Department
	for _, dept := range departments {
		if (parentID == nil && dept.ParentID == 0) || (parentID != nil && dept.ParentID == *parentID) {
			dept.Children = s.buildTree(departments, &dept.ID)
			tree = append(tree, dept)
		}
	}
	return tree
}

// GetUsers 获取部门用户
func (s *DepartmentService) GetUsers(id uint) ([]model.User, error) {
	var users []model.User
	if err := s.db.Where("department_id = ?", id).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to get department users: %w", err)
	}
	return users, nil
}