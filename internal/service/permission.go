package service

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/pkg/auth"
	"github.com/cuiyuanxin/kunpeng/internal/database"
	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// PermissionService 权限服务
type PermissionService struct {
	db     *gorm.DB
	casbin *auth.CasbinEnforcer
}

// NewPermissionService 创建权限服务
func NewPermissionService(db *gorm.DB) *PermissionService {
	// 使用全局的Casbin权限管理器
	casbinEnforcer := database.GetCasbin()
	if casbinEnforcer == nil {
		// 如果获取失败，记录错误但不阻止服务启动
		fmt.Printf("Failed to get Casbin enforcer from database manager\n")
	}

	return &PermissionService{
		db:     db,
		casbin: casbinEnforcer,
	}
}

// CheckUserPermission 检查用户权限
func (s *PermissionService) CheckUserPermission(userID uint, permissionCode string) (bool, error) {
	return s.casbin.CheckUserPermission(userID, permissionCode)
}

// CheckPermission 检查用户权限
func (s *PermissionService) CheckPermission(userID uint, resource, action string) (bool, error) {
	if s.casbin == nil {
		return false, fmt.Errorf("casbin enforcer not initialized")
	}

	userSubject := fmt.Sprintf("user:%d", userID)
	return s.casbin.Enforce(userSubject, resource, action)
}

// CheckUserPermissionByPath 根据路径检查用户权限
func (s *PermissionService) CheckUserPermissionByPath(userID uint, path, method string) (bool, error) {
	// 首先使用RESTful API权限检查
	hasRESTfulPermission, err := s.casbin.CheckUserPermissionByPath(userID, path, method)
	if err != nil {
		return false, fmt.Errorf("failed to check RESTful permission: %w", err)
	}
	
	// 如果RESTful权限检查通过，直接返回
	if hasRESTfulPermission {
		return true, nil
	}
	
	// 如果RESTful权限检查不通过，回退到传统权限检查
	var permission model.AdminPermission
	err = s.db.Where("url = ? AND method = ? AND status = 1", path, method).First(&permission).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 如果没有找到对应的权限配置，默认拒绝访问
			return false, nil
		}
		return false, err
	}

	return s.CheckUserPermission(userID, permission.Code)
}

// SyncUserRoles 同步用户角色到Casbin
func (s *PermissionService) SyncUserRoles(userID uint, roleIDs []uint) error {
	return s.casbin.SyncUserRoles(userID, roleIDs)
}

// SyncRolePermissions 同步角色权限到Casbin
func (s *PermissionService) SyncRolePermissions(roleID uint, permissionIDs []uint) error {
	// 获取权限代码列表
	var permissions []model.AdminPermission
	err := s.db.Where("id IN ? AND status = 1", permissionIDs).Find(&permissions).Error
	if err != nil {
		return fmt.Errorf("failed to get permissions: %w", err)
	}

	permissionCodes := make([]string, len(permissions))
	for i, perm := range permissions {
		permissionCodes[i] = perm.Code
	}

	return s.casbin.SyncRolePermissions(roleID, permissionCodes)
}

// InitializePermissions 初始化权限数据到Casbin
func (s *PermissionService) InitializePermissions() error {
	// 获取所有用户及其角色
	var users []model.AdminUser
	err := s.db.Preload("Roles").Where("status = 1").Find(&users).Error
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	// 同步用户角色
	for _, user := range users {
		roleIDs := make([]uint, len(user.Roles))
		for i, role := range user.Roles {
			roleIDs[i] = role.ID
		}
		err := s.SyncUserRoles(user.ID, roleIDs)
		if err != nil {
			return fmt.Errorf("failed to sync user roles for user %d: %w", user.ID, err)
		}
	}

	// 获取所有角色及其权限
	var roles []model.AdminRole
	err = s.db.Preload("Permissions").Where("status = 1").Find(&roles).Error
	if err != nil {
		return fmt.Errorf("failed to get roles: %w", err)
	}

	// 同步角色权限
	for _, role := range roles {
		permissionIDs := make([]uint, len(role.Permissions))
		for i, perm := range role.Permissions {
			permissionIDs[i] = perm.ID
		}
		err := s.SyncRolePermissions(role.ID, permissionIDs)
		if err != nil {
			return fmt.Errorf("failed to sync role permissions for role %d: %w", role.ID, err)
		}
	}

	return nil
}

// GetUserRoles 获取用户角色
func (s *PermissionService) GetUserRoles(userID uint) ([]string, error) {
	return s.casbin.GetRolesForUser(fmt.Sprintf("user:%d", userID))
}

// GetUserPermissions 获取用户权限
func (s *PermissionService) GetUserPermissions(userID uint) ([][]string, error) {
	return s.casbin.GetPermissionsForUser(fmt.Sprintf("user:%d", userID))
}

// GetRolePermissions 获取角色权限
func (s *PermissionService) GetRolePermissions(roleID uint) ([][]string, error) {
	return s.casbin.GetPermissionsForRole(fmt.Sprintf("role:%d", roleID))
}

// AddUserRole 为用户添加角色
func (s *PermissionService) AddUserRole(userID, roleID uint) error {
	userStr := fmt.Sprintf("user:%d", userID)
	roleStr := fmt.Sprintf("role:%d", roleID)
	
	_, err := s.casbin.AddRoleForUser(userStr, roleStr)
	if err != nil {
		return fmt.Errorf("failed to add role for user: %w", err)
	}

	// 同时更新数据库
	err = s.db.Exec("INSERT IGNORE INTO admin_user_roles (admin_user_id, admin_role_id) VALUES (?, ?)", userID, roleID).Error
	if err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}

	return nil
}

// RemoveUserRole 移除用户角色
func (s *PermissionService) RemoveUserRole(userID, roleID uint) error {
	userStr := fmt.Sprintf("user:%d", userID)
	roleStr := fmt.Sprintf("role:%d", roleID)
	
	_, err := s.casbin.DeleteRoleForUser(userStr, roleStr)
	if err != nil {
		return fmt.Errorf("failed to remove role for user: %w", err)
	}

	// 同时更新数据库
	err = s.db.Exec("DELETE FROM admin_user_roles WHERE admin_user_id = ? AND admin_role_id = ?", userID, roleID).Error
	if err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}

	return nil
}

// AddRolePermission 为角色添加权限
func (s *PermissionService) AddRolePermission(roleID, permissionID uint) error {
	// 获取权限信息
	var permission model.AdminPermission
	err := s.db.First(&permission, permissionID).Error
	if err != nil {
		return fmt.Errorf("failed to get permission: %w", err)
	}

	roleStr := fmt.Sprintf("role:%d", roleID)
	_, err = s.casbin.AddPolicy(roleStr, permission.Code, "*")
	if err != nil {
		return fmt.Errorf("failed to add policy: %w", err)
	}

	// 同时更新数据库
	err = s.db.Exec("INSERT IGNORE INTO admin_role_permissions (admin_role_id, admin_permission_id) VALUES (?, ?)", roleID, permissionID).Error
	if err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}

	return nil
}

// RemoveRolePermission 移除角色权限
func (s *PermissionService) RemoveRolePermission(roleID, permissionID uint) error {
	// 获取权限信息
	var permission model.AdminPermission
	err := s.db.First(&permission, permissionID).Error
	if err != nil {
		return fmt.Errorf("failed to get permission: %w", err)
	}

	roleStr := fmt.Sprintf("role:%d", roleID)
	_, err = s.casbin.RemovePolicy(roleStr, permission.Code, "*")
	if err != nil {
		return fmt.Errorf("failed to remove policy: %w", err)
	}

	// 同时更新数据库
	err = s.db.Exec("DELETE FROM admin_role_permissions WHERE admin_role_id = ? AND admin_permission_id = ?", roleID, permissionID).Error
	if err != nil {
		return fmt.Errorf("failed to update database: %w", err)
	}

	return nil
}

// GetUserMenuPermissions 获取用户菜单权限
func (s *PermissionService) GetUserMenuPermissions(userID uint) ([]model.AdminPermission, error) {
	// 获取用户所有权限
	var permissions []model.AdminPermission
	err := s.db.Table("admin_permissions").
		Joins("JOIN admin_role_permissions ON admin_permissions.id = admin_role_permissions.admin_permission_id").
		Joins("JOIN admin_user_roles ON admin_role_permissions.admin_role_id = admin_user_roles.admin_role_id").
		Where("admin_user_roles.admin_user_id = ? AND admin_permissions.status = 1 AND admin_permissions.type IN (1)", userID).
		Group("admin_permissions.id").
		Order("admin_permissions.sort_order ASC").
		Find(&permissions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user menu permissions: %w", err)
	}

	return permissions, nil
}

// GetUserButtonPermissions 获取用户按钮权限
func (s *PermissionService) GetUserButtonPermissions(userID uint) ([]string, error) {
	var permissions []model.AdminPermission
	err := s.db.Table("admin_permissions").
		Joins("JOIN admin_role_permissions ON admin_permissions.id = admin_role_permissions.admin_permission_id").
		Joins("JOIN admin_user_roles ON admin_role_permissions.admin_role_id = admin_user_roles.admin_role_id").
		Where("admin_user_roles.admin_user_id = ? AND admin_permissions.status = 1 AND admin_permissions.type = 2", userID).
		Group("admin_permissions.id").
		Find(&permissions).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get user button permissions: %w", err)
	}

	permissionCodes := make([]string, len(permissions))
	for i, perm := range permissions {
		permissionCodes[i] = perm.Code
	}

	return permissionCodes, nil
}

// ReloadPolicy 重新加载权限策略
func (s *PermissionService) ReloadPolicy() error {
	return s.casbin.LoadPolicy()
}

// AddRESTfulPermission 为角色添加RESTful API权限
func (s *PermissionService) AddRESTfulPermission(roleID uint, path, method string) error {
	return s.casbin.AddRoleRESTfulPolicy(roleID, path, method)
}

// RemoveRESTfulPermission 移除角色的RESTful API权限
func (s *PermissionService) RemoveRESTfulPermission(roleID uint, path, method string) error {
	return s.casbin.RemoveRoleRESTfulPolicy(roleID, path, method)
}

// SyncRoleRESTfulPermissions 同步角色的RESTful API权限
func (s *PermissionService) SyncRoleRESTfulPermissions(roleID uint, permissions []map[string]string) error {
	roleStr := fmt.Sprintf("role:%d", roleID)
	
	// 删除角色的所有RESTful权限
	policies := s.casbin.GetFilteredPolicy(0, roleStr)
	
	for _, policy := range policies {
		if len(policy) >= 4 && policy[3] == "allow" {
			_, err := s.casbin.RemoveRESTfulPolicy(policy[0], policy[1], policy[2])
			if err != nil {
				return fmt.Errorf("failed to remove RESTful policy: %w", err)
			}
		}
	}
	
	// 添加新的RESTful权限
	for _, perm := range permissions {
		path, pathOk := perm["path"]
		method, methodOk := perm["method"]
		if pathOk && methodOk {
			err := s.AddRESTfulPermission(roleID, path, method)
			if err != nil {
				return fmt.Errorf("failed to add RESTful permission: %w", err)
			}
		}
	}
	
	return nil
}

// GetRoleRESTfulPermissions 获取角色的RESTful权限
func (s *PermissionService) GetRoleRESTfulPermissions(roleID uint) ([]map[string]string, error) {
	var permissions []map[string]string

	// 从Casbin获取角色权限
	policies := s.casbin.GetFilteredPolicy(0, fmt.Sprintf("role:%d", roleID))

	for _, policy := range policies {
		if len(policy) >= 3 {
			permissions = append(permissions, map[string]string{
				"path":   policy[1],
				"method": policy[2],
			})
		}
	}

	return permissions, nil
}

// GetAll 获取所有权限
func (s *PermissionService) GetAll() ([]model.AdminPermission, error) {
	var permissions []model.AdminPermission
	if err := s.db.Find(&permissions).Error; err != nil {
		return nil, fmt.Errorf("failed to get all permissions: %w", err)
	}
	return permissions, nil
}