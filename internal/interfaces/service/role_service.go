package service

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// RoleService 角色服务接口
type RoleService interface {
	// GetRoleList 获取角色列表
	GetRoleList(req *dto.RolePageReq) (*dto.PageResp, error)

	// GetAllRoles 获取所有角色
	GetAllRoles() ([]*model.Role, error)

	// GetRoleByID 根据ID获取角色
	GetRoleByID(id uint) (*model.Role, error)

	// CreateRole 创建角色
	CreateRole(req *dto.RoleCreateReq) (uint, error)

	// UpdateRole 更新角色
	UpdateRole(req *dto.RoleUpdateReq) error

	// DeleteRole 删除角色
	DeleteRole(id uint) error

	// BatchDeleteRole 批量删除角色
	BatchDeleteRole(ids []uint) error

	// ChangeRoleStatus 修改角色状态
	ChangeRoleStatus(req *dto.StatusReq) error

	// GetRoleMenus 获取角色菜单
	GetRoleMenus(roleID uint) (*dto.MenuRoleResp, error)

	// UpdateRoleMenus 更新角色菜单
	UpdateRoleMenus(req *dto.RoleMenuReq) error

	// GetRoleAPIs 获取角色API
	GetRoleAPIs(roleID uint) (*dto.APIRoleResp, error)

	// UpdateRoleAPIs 更新角色API
	UpdateRoleAPIs(req *dto.RoleAPIReq) error
}
