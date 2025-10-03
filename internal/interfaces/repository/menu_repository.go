package repository

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// MenuRepository 菜单仓储接口
type MenuRepository interface {
	// 根据ID获取菜单
	FindByID(id uint) (*model.Menu, error)

	// 获取菜单列表
	FindList(req *dto.MenuListReq) ([]*model.Menu, error)

	// 获取菜单树
	FindTree() ([]*dto.MenuTreeResp, error)

	// 获取用户菜单树
	FindUserMenuTree(userID uint) ([]*dto.MenuTreeResp, error)

	// 创建菜单
	Create(menu *model.Menu) error

	// 更新菜单
	Update(menu *model.Menu) error

	// 删除菜单
	Delete(id uint) error

	// 检查菜单是否有子菜单
	HasChildren(id uint) (bool, error)

	// 根据角色ID获取菜单ID列表
	FindMenuIDsByRoleID(roleID uint) ([]uint, error)
}
