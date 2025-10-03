package repository

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	// 根据ID获取角色
	FindByID(id uint) (*model.Role, error)

	// 获取所有角色
	FindAll() ([]*model.Role, error)

	// 获取角色列表
	FindList(req *dto.RolePageReq) ([]*model.Role, int64, error)

	// 创建角色
	Create(role *model.Role) error

	// 更新角色
	Update(role *model.Role) error

	// 删除角色
	Delete(id uint) error

	// 批量删除角色
	BatchDelete(ids []uint) error

	// 更新角色状态
	UpdateStatus(id uint, status int) error

	// 获取角色菜单
	GetRoleMenus(roleID uint) ([]uint, error)

	// 更新角色菜单
	UpdateRoleMenus(roleID uint, menuIDs []uint) error

	// 获取角色API
	GetRoleAPIs(roleID uint) ([]uint, error)

	// 更新角色API
	UpdateRoleAPIs(roleID uint, apiIDs []uint) error
}
