package service

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// MenuService 菜单服务接口
type MenuService interface {
	// GetMenuList 获取菜单列表
	GetMenuList(req *dto.MenuQueryReq) ([]*model.Menu, error)

	// GetMenuTree 获取菜单树
	GetMenuTree(req *dto.MenuQueryReq) ([]*dto.MenuTreeResp, error)

	// GetUserMenuTree 获取用户菜单树
	GetUserMenuTree(userID uint) ([]*dto.MenuTreeResp, error)

	// GetMenuByID 根据ID获取菜单
	GetMenuByID(id uint) (*model.Menu, error)

	// CreateMenu 创建菜单
	CreateMenu(req *dto.MenuCreateReq) (uint, error)

	// UpdateMenu 更新菜单
	UpdateMenu(req *dto.MenuUpdateReq) error

	// DeleteMenu 删除菜单
	DeleteMenu(id uint) error
}
