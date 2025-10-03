package impl

import (
	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// RoleRepositoryImpl 角色仓储实现
type RoleRepositoryImpl struct {
	BaseRepository
}

// NewRoleRepository 创建角色仓储实例
func NewRoleRepository() repository.RoleRepository {
	return &RoleRepositoryImpl{
		BaseRepository: NewBaseRepository(),
	}
}

// FindByID 根据ID查找角色
func (r *RoleRepositoryImpl) FindByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return &role, nil
}

// FindAll 查找所有启用的角色
func (r *RoleRepositoryImpl) FindAll() ([]*model.Role, error) {
	var roles []*model.Role
	err := r.db.Where("status = ?", 1).Order("sort ASC").Find(&roles).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return roles, nil
}

// FindList 分页查询角色列表
func (r *RoleRepositoryImpl) FindList(req *dto.RolePageReq) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	db := r.db.Model(&model.Role{})

	// 构建查询条件
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		db = db.Where("code LIKE ?", "%"+req.Code+"%")
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	// 时间范围查询
	if req.BeginTime != "" && req.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.BeginTime, req.EndTime)
	}

	// 查询总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 分页查询
	err = db.Order("sort ASC").
		Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&roles).Error
	if err != nil {
		return nil, 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	return roles, total, nil
}

// Create 创建角色
func (r *RoleRepositoryImpl) Create(role *model.Role) error {
	err := r.db.Create(role).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// Update 更新角色
func (r *RoleRepositoryImpl) Update(role *model.Role) error {
	err := r.db.Save(role).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// Delete 删除角色
func (r *RoleRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&model.Role{}, id).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// BatchDelete 批量删除角色
func (r *RoleRepositoryImpl) BatchDelete(ids []uint) error {
	err := r.db.Delete(&model.Role{}, ids).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// UpdateStatus 更新角色状态
func (r *RoleRepositoryImpl) UpdateStatus(id uint, status int) error {
	err := r.db.Model(&model.Role{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// GetRoleMenus 获取角色菜单权限
func (r *RoleRepositoryImpl) GetRoleMenus(roleID uint) ([]uint, error) {
	var menuIDs []uint
	err := r.db.Table("role_menus").Where("role_id = ?", roleID).Pluck("menu_id", &menuIDs).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 如果没有找到记录，返回空切片而不是错误
	if len(menuIDs) == 0 {
		return []uint{}, nil
	}

	return menuIDs, nil
}

// UpdateRoleMenus 更新角色菜单权限
func (r *RoleRepositoryImpl) UpdateRoleMenus(roleID uint, menuIDs []uint) error {
	// 开启事务
	tx := r.db.Begin()
	if tx.Error != nil {
		return kperrors.New(kperrors.ErrDatabase, tx.Error)
	}

	// 删除原有的角色菜单关联
	if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
		tx.Rollback()
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 插入新的角色菜单关联
	if len(menuIDs) > 0 {
		var roleMenus []model.RoleMenu
		for _, menuID := range menuIDs {
			roleMenus = append(roleMenus, model.RoleMenu{
				RoleID: roleID,
				MenuID: menuID,
			})
		}

		if err := tx.Create(&roleMenus).Error; err != nil {
			tx.Rollback()
			return kperrors.New(kperrors.ErrDatabase, err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// GetRoleAPIs 获取角色API权限
func (r *RoleRepositoryImpl) GetRoleAPIs(roleID uint) ([]uint, error) {
	var apiIDs []uint
	err := r.db.Table("role_apis").Where("role_id = ?", roleID).Pluck("api_id", &apiIDs).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 如果没有找到记录，返回空切片而不是错误
	if len(apiIDs) == 0 {
		return []uint{}, nil
	}

	return apiIDs, nil
}

// UpdateRoleAPIs 更新角色API权限
func (r *RoleRepositoryImpl) UpdateRoleAPIs(roleID uint, apiIDs []uint) error {
	// 开启事务
	tx := r.db.Begin()
	if tx.Error != nil {
		return kperrors.New(kperrors.ErrDatabase, tx.Error)
	}

	// 删除原有的角色API关联
	if err := tx.Where("role_id = ?", roleID).Delete(&model.RoleAPI{}).Error; err != nil {
		tx.Rollback()
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 插入新的角色API关联
	if len(apiIDs) > 0 {
		var roleAPIs []model.RoleAPI
		for _, apiID := range apiIDs {
			roleAPIs = append(roleAPIs, model.RoleAPI{
				RoleID: roleID,
				APIID:  apiID,
			})
		}

		if err := tx.Create(&roleAPIs).Error; err != nil {
			tx.Rollback()
			return kperrors.New(kperrors.ErrDatabase, err)
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}
