package impl

import (
	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// MenuRepositoryImpl 菜单仓储实现
type MenuRepositoryImpl struct {
	BaseRepository
}

// NewMenuRepository 创建菜单仓储
func NewMenuRepository() repository.MenuRepository {
	return &MenuRepositoryImpl{
		BaseRepository: NewBaseRepository(),
	}
}

// FindByID 根据ID获取菜单
func (r *MenuRepositoryImpl) FindByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.First(&menu, id).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}
	return &menu, nil
}

// FindList 获取菜单列表
func (r *MenuRepositoryImpl) FindList(req *dto.MenuListReq) ([]*model.Menu, error) {
	var menus []*model.Menu

	db := r.db.Model(&model.Menu{})

	// 构建查询条件
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	// 查询
	err := db.Order("sort ASC").Find(&menus).Error
	if err != nil {
		return nil, errors.New(errors.ErrDatabase, err)
	}

	return menus, nil
}

// FindTree 获取菜单树
func (r *MenuRepositoryImpl) FindTree() ([]*dto.MenuTreeResp, error) {
	var menus []*model.Menu

	// 查询所有菜单
	err := r.db.Order("sort ASC").Find(&menus).Error
	if err != nil {
		return nil, errors.New(errors.ErrDatabase, err)
	}

	// 构建菜单树
	return buildMenuTree(menus, 0), nil
}

// FindUserMenuTree 获取用户菜单树
func (r *MenuRepositoryImpl) FindUserMenuTree(userID uint) ([]*dto.MenuTreeResp, error) {
	var menus []*model.Menu

	// 查询用户角色
	var user model.User
	err := r.db.First(&user, userID).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}

	// 如果是超级管理员，返回所有菜单
	if user.RoleID == 1 {
		err = r.db.Where("status = ? AND type IN (0, 1)", 1).Order("sort ASC").Find(&menus).Error
		if err != nil {
			return nil, errors.New(errors.ErrDatabase, err)
		}
	} else {
		// 查询角色菜单
		var menuIDs []uint
		err = r.db.Model(&model.RoleMenu{}).
			Select("menu_id").
			Where("role_id = ?", user.RoleID).
			Pluck("menu_id", &menuIDs).Error
		if err != nil {
			return nil, errors.New(errors.ErrDatabase, err)
		}

		// 查询菜单
		err = r.db.Where("id IN ? AND status = ? AND type IN (0, 1)", menuIDs, 1).
			Order("sort ASC").
			Find(&menus).Error
		if err != nil {
			return nil, errors.New(errors.ErrDatabase, err)
		}
	}

	// 构建菜单树
	return buildMenuTree(menus, 0), nil
}

// Create 创建菜单
func (r *MenuRepositoryImpl) Create(menu *model.Menu) error {
	err := r.db.Create(menu).Error
	if err != nil {
		return errors.New(errors.ErrDatabase, err)
	}
	return nil
}

// Update 更新菜单
func (r *MenuRepositoryImpl) Update(menu *model.Menu) error {
	err := r.db.Save(menu).Error
	if err != nil {
		return errors.New(errors.ErrDatabase, err)
	}
	return nil
}

// Delete 删除菜单
func (r *MenuRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&model.Menu{}, id).Error
	if err != nil {
		return errors.New(errors.ErrDatabase, err)
	}
	return nil
}

// HasChildren 检查菜单是否有子菜单
func (r *MenuRepositoryImpl) HasChildren(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ErrDatabase, err)
	}
	return count > 0, nil
}

// FindMenuIDsByRoleID 根据角色ID获取菜单ID列表
func (r *MenuRepositoryImpl) FindMenuIDsByRoleID(roleID uint) ([]uint, error) {
	var menuIDs []uint
	err := r.db.Model(&model.RoleMenu{}).
		Select("menu_id").
		Where("role_id = ?", roleID).
		Pluck("menu_id", &menuIDs).Error
	if err != nil {
		return nil, errors.New(errors.ErrDatabase, err)
	}
	return menuIDs, nil
}

// buildMenuTree 构建菜单树
func buildMenuTree(menus []*model.Menu, parentID uint) []*dto.MenuTreeResp {
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
				IsCache:    menu.IsCache,
				IsFrame:    menu.IsFrame,
				Children:   buildMenuTree(menus, menu.ID),
			}
			tree = append(tree, node)
		}
	}

	return tree
}
