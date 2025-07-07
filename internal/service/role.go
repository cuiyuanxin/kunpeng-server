package service

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// RoleService 角色服务
type RoleService struct {
	db *gorm.DB
}

// NewRoleService 创建角色服务
func NewRoleService(db *gorm.DB) *RoleService {
	return &RoleService{
		db: db,
	}
}

// GetAll 获取所有角色
func (s *RoleService) GetAll() ([]model.Role, error) {
	var roles []model.Role
	if err := s.db.Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("failed to get roles: %w", err)
	}
	return roles, nil
}

// GetByID 根据ID获取角色
func (s *RoleService) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := s.db.First(&role, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return &role, nil
}

// GetByName 根据名称获取角色
func (s *RoleService) GetByName(name string) (*model.Role, error) {
	var role model.Role
	if err := s.db.Where("name = ?", name).First(&role).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return &role, nil
}

// Create 创建角色
func (s *RoleService) Create(req *model.RoleCreateRequest) (*model.Role, error) {
	// 检查角色名是否已存在
	var existRole model.Role
	if err := s.db.Where("name = ?", req.Name).First(&existRole).Error; err == nil {
		return nil, errors.New("role name already exists")
	}

	// 创建角色
	role := &model.Role{
		Name:        req.Name,
		Description: req.Description,
		Status:      1, // 默认启用
	}

	if err := s.db.Create(role).Error; err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return role, nil
}

// Update 更新角色
func (s *RoleService) Update(id uint, req *model.RoleUpdateRequest) (*model.Role, error) {
	role, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段

	if req.Description != nil && *req.Description != "" {
		role.Description = req.Description
	}
	if req.Status != nil {
		role.Status = *req.Status
	}

	if err := s.db.Save(role).Error; err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return role, nil
}

// Delete 删除角色
func (s *RoleService) Delete(id uint) error {
	if err := s.db.Delete(&model.Role{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}
	return nil
}

// AssignPermissions 为角色分配权限
func (s *RoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	role, err := s.GetByID(roleID)
	if err != nil {
		return err
	}

	// 清除现有权限
	if err := s.db.Model(role).Association("Permissions").Clear(); err != nil {
		return fmt.Errorf("failed to clear role permissions: %w", err)
	}

	// 分配新权限
	var permissions []model.Permission
	if err := s.db.Find(&permissions, permissionIDs).Error; err != nil {
		return fmt.Errorf("failed to find permissions: %w", err)
	}

	if err := s.db.Model(role).Association("Permissions").Append(permissions); err != nil {
		return fmt.Errorf("failed to assign permissions: %w", err)
	}

	return nil
}

// GetPermissions 获取角色的权限
func (s *RoleService) GetPermissions(roleID uint) ([]model.Permission, error) {
	role, err := s.GetByID(roleID)
	if err != nil {
		return nil, err
	}

	var permissions []model.Permission
	if err := s.db.Model(role).Association("Permissions").Find(&permissions); err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	return permissions, nil
}