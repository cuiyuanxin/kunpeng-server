package impl

import (
	"errors"

	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/pkg/database"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"gorm.io/gorm"
)

// MenuServiceImpl 菜单服务实现
type MenuServiceImpl struct{}

// GetMenuList 获取菜单列表
func (s *MenuServiceImpl) GetMenuList(req *dto.MenuQueryReq) ([]*model.Menu, error) {
	var menus []*model.Menu
	db := database.GetDB().Model(&model.Menu{})

	// 构建查询条件
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	// 查询菜单
	err := db.Order("sort ASC").Find(&menus).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	return menus, nil
}

// GetMenuTree 获取菜单树
func (s *MenuServiceImpl) GetMenuTree(req *dto.MenuQueryReq) ([]*dto.MenuTreeResp, error) {
	// 获取菜单列表
	menus, err := s.GetMenuList(req)
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	return s.buildMenuTree(menus, 0), nil
}

// GetUserMenuTree 获取用户菜单树
func (s *MenuServiceImpl) GetUserMenuTree(userID uint) ([]*dto.MenuTreeResp, error) {
	// 获取用户信息
	var user model.User
	err := database.GetDB().First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kperrors.New(kperrors.ErrUserNotFound, err)
		}
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 获取用户角色
	var role model.Role
	err = database.GetDB().First(&role, user.RoleID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kperrors.New(kperrors.ErrRoleNotFound, err)
		}
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 获取角色菜单
	var roleMenus []model.RoleMenu
	err = database.GetDB().Where("role_id = ?", role.ID).Find(&roleMenus).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 提取菜单ID
	menuIDs := make([]uint, 0, len(roleMenus))
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	// 获取菜单
	var menus []*model.Menu
	err = database.GetDB().Where("id IN ? AND status = ? AND visible = ?", menuIDs, 1, 1).Order("sort ASC").Find(&menus).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 构建菜单树
	return s.buildMenuTree(menus, 0), nil
}

// GetMenuByID 根据ID获取菜单
func (s *MenuServiceImpl) GetMenuByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := database.GetDB().First(&menu, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kperrors.New(kperrors.ErrMenuNotFound, err)
		}
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return &menu, nil
}

// CreateMenu 创建菜单
func (s *MenuServiceImpl) CreateMenu(req *dto.MenuCreateReq) (uint, error) {
	// 检查父菜单是否存在
	if req.ParentID != 0 {
		var parentMenu model.Menu
		err := database.GetDB().First(&parentMenu, req.ParentID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, kperrors.New(kperrors.ErrMenuParentNotFound, err)
			}
			return 0, kperrors.New(kperrors.ErrDatabase, err)
		}

		// 检查父菜单状态
		if parentMenu.Status != 1 {
			return 0, kperrors.New(kperrors.ErrMenuParentDisabled, nil)
		}
	}

	// 创建菜单
	menu := model.Menu{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Type:       req.Type,
		Path:       req.Path,
		Component:  req.Component,
		Permission: req.Permission,
		Icon:       req.Icon,
		Sort:       req.Sort,
		Visible:    req.Visible,
		Status:     req.Status,
		Remark:     req.Remark,
	}

	err := database.GetDB().Create(&menu).Error
	if err != nil {
		return 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	return menu.ID, nil
}

// UpdateMenu 更新菜单
func (s *MenuServiceImpl) UpdateMenu(req *dto.MenuUpdateReq) error {
	// 检查菜单是否存在
	var menu model.Menu
	err := database.GetDB().First(&menu, req.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrMenuNotFound, err)
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 检查父菜单是否存在
	if req.ParentID != 0 {
		var parentMenu model.Menu
		err := database.GetDB().First(&parentMenu, req.ParentID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return kperrors.New(kperrors.ErrMenuParentNotFound, err)
			}
			return kperrors.New(kperrors.ErrDatabase, err)
		}

		// 检查父菜单状态
		if parentMenu.Status != 1 {
			return kperrors.New(kperrors.ErrMenuParentDisabled, nil)
		}

		// 检查是否将菜单设置为自己的子菜单
		if req.ParentID == req.ID {
			return kperrors.New(kperrors.ErrMenuParentInvalid, nil)
		}

		// 检查是否将菜单设置为自己子菜单的子菜单
		var childMenus []*model.Menu
		err = database.GetDB().Where("parent_id = ?", req.ID).Find(&childMenus).Error
		if err != nil {
			return kperrors.New(kperrors.ErrDatabase, err)
		}
		childIDs := make([]uint, 0, len(childMenus))
		for _, childMenu := range childMenus {
			childIDs = append(childIDs, childMenu.ID)
		}
		if len(childIDs) > 0 && contains(childIDs, req.ParentID) {
			return kperrors.New(kperrors.ErrMenuParentInvalid, nil)
		}
	}

	// 更新菜单
	updates := map[string]interface{}{
		"parent_id":  req.ParentID,
		"name":       req.Name,
		"type":       req.Type,
		"path":       req.Path,
		"component":  req.Component,
		"permission": req.Permission,
		"icon":       req.Icon,
		"sort":       req.Sort,
		"visible":    req.Visible,
		"status":     req.Status,
		"remark":     req.Remark,
	}

	err = database.GetDB().Model(&menu).Updates(updates).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// DeleteMenu 删除菜单
func (s *MenuServiceImpl) DeleteMenu(id uint) error {
	// 检查菜单是否存在
	var menu model.Menu
	err := database.GetDB().First(&menu, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrMenuNotFound, err)
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 检查是否有子菜单
	var childCount int64
	err = database.GetDB().Model(&model.Menu{}).Where("parent_id = ?", id).Count(&childCount).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	if childCount > 0 {
		return kperrors.New(kperrors.ErrMenuHasChildren, nil)
	}

	// 删除菜单
	err = database.GetDB().Delete(&menu).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 删除角色菜单关联
	err = database.GetDB().Where("menu_id = ?", id).Delete(&model.RoleMenu{}).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// buildMenuTree 构建菜单树
func (s *MenuServiceImpl) buildMenuTree(menus []*model.Menu, parentID uint) []*dto.MenuTreeResp {
	var tree []*dto.MenuTreeResp
	for _, menu := range menus {
		if menu.ParentID == parentID {
			node := &dto.MenuTreeResp{
				ID:         menu.ID,
				ParentID:   menu.ParentID,
				Name:       menu.Name,
				Type:       menu.Type,
				Path:       menu.Path,
				Component:  menu.Component,
				Permission: menu.Permission,
				Icon:       menu.Icon,
				Sort:       menu.Sort,
				Visible:    menu.Visible,
				Status:     menu.Status,
				Children:   s.buildMenuTree(menus, menu.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// contains 检查切片是否包含某个元素
func contains(slice []uint, item uint) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
