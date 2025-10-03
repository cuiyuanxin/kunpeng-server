package impl

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/internal/repository"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// RoleServiceImpl 角色服务实现
type RoleServiceImpl struct{}

// GetRoleList 获取角色列表
func (s *RoleServiceImpl) GetRoleList(req *dto.RolePageReq) (*dto.PageResp, error) {
	// 查询角色列表
	roles, total, err := repository.GetRoleRepository().FindList(req)
	if err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		totalPages++
	}

	return &dto.PageResp{
		List:       roles,
		Total:      total,
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAllRoles 获取所有角色
func (s *RoleServiceImpl) GetAllRoles() ([]*model.Role, error) {
	return repository.GetRoleRepository().FindAll()
}

// GetRoleByID 根据ID获取角色
func (s *RoleServiceImpl) GetRoleByID(id uint) (*model.Role, error) {
	return repository.GetRoleRepository().FindByID(id)
}

// CreateRole 创建角色
func (s *RoleServiceImpl) CreateRole(req *dto.RoleCreateReq) (uint, error) {
	// 创建角色
	role := model.Role{
		Name:   req.Name,
		Code:   req.Code,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	}

	err := repository.GetRoleRepository().Create(&role)
	if err != nil {
		return 0, err
	}

	return role.ID, nil
}

// UpdateRole 更新角色
func (s *RoleServiceImpl) UpdateRole(req *dto.RoleUpdateReq) error {
	// 检查角色是否存在
	role, err := repository.GetRoleRepository().FindByID(req.ID)
	if err != nil {
		return err
	}

	// 更新角色信息
	role.Name = req.Name
	role.Code = req.Code
	role.Sort = req.Sort
	role.Status = req.Status
	role.Remark = req.Remark

	return repository.GetRoleRepository().Update(role)
}

// DeleteRole 删除角色
func (s *RoleServiceImpl) DeleteRole(id uint) error {
	// 检查角色是否存在
	_, err := repository.GetRoleRepository().FindByID(id)
	if err != nil {
		return err
	}

	// 检查角色是否有关联用户
	err = repository.GetUserRepository().FindByRoleID(id)
	if err == nil {
		return kperrors.New(kperrors.ErrRoleHasUsers, nil)
	}

	// 删除角色
	return repository.GetRoleRepository().Delete(id)
}

// BatchDeleteRole 批量删除角色
func (s *RoleServiceImpl) BatchDeleteRole(ids []uint) error {
	// 检查角色是否有关联用户
	for _, id := range ids {
		err := repository.GetUserRepository().FindByRoleID(id)
		if err == nil {
			return kperrors.New(kperrors.ErrRoleHasUsers, nil)
		}
	}

	// 删除角色
	return repository.GetRoleRepository().BatchDelete(ids)
}

// ChangeRoleStatus 修改角色状态
func (s *RoleServiceImpl) ChangeRoleStatus(req *dto.StatusReq) error {
	return repository.GetRoleRepository().UpdateStatus(req.ID, int(req.Status))
}

// GetRoleMenus 获取角色菜单
func (s *RoleServiceImpl) GetRoleMenus(roleID uint) (*dto.MenuRoleResp, error) {
	// 检查角色是否存在
	_, err := repository.GetRoleRepository().FindByID(roleID)
	if err != nil {
		return nil, err
	}

	// 获取角色菜单
	menuIDs, err := repository.GetRoleRepository().GetRoleMenus(roleID)
	if err != nil {
		return nil, err
	}

	return &dto.MenuRoleResp{
		MenuIDs: menuIDs,
	}, nil
}

// UpdateRoleMenus 更新角色菜单
func (s *RoleServiceImpl) UpdateRoleMenus(req *dto.RoleMenuReq) error {
	// 检查角色是否存在
	_, err := repository.GetRoleRepository().FindByID(req.RoleID)
	if err != nil {
		return err
	}

	// 更新角色菜单
	return repository.GetRoleRepository().UpdateRoleMenus(req.RoleID, req.MenuIDs)
}

// GetRoleAPIs 获取角色API
func (s *RoleServiceImpl) GetRoleAPIs(roleID uint) (*dto.APIRoleResp, error) {
	// 检查角色是否存在
	_, err := repository.GetRoleRepository().FindByID(roleID)
	if err != nil {
		return nil, err
	}

	// 获取角色API
	apiIDs, err := repository.GetRoleRepository().GetRoleAPIs(roleID)
	if err != nil {
		return nil, err
	}

	return &dto.APIRoleResp{
		APIIDs: apiIDs,
	}, nil
}

// UpdateRoleAPIs 更新角色API
func (s *RoleServiceImpl) UpdateRoleAPIs(req *dto.RoleAPIReq) error {
	// 检查角色是否存在
	_, err := repository.GetRoleRepository().FindByID(req.RoleID)
	if err != nil {
		return err
	}

	// 更新角色API
	return repository.GetRoleRepository().UpdateRoleAPIs(req.RoleID, req.APIIDs)
}
