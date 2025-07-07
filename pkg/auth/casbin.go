package auth

import (
	"fmt"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// CasbinEnforcer Casbin权限执行器
type CasbinEnforcer struct {
	enforcer *casbin.Enforcer
	db       *gorm.DB
}

// getRBACModelText 获取RBAC模型配置文本
func getRBACModelText() string {
	return `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && regexMatch(r.act, p.act)
`
}

// NewCasbinEnforcer 创建Casbin权限执行器
func NewCasbinEnforcer(db *gorm.DB) (*CasbinEnforcer, error) {
	// 创建Gorm适配器
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin adapter: %w", err)
	}

	// 从字符串创建模型
	m, err := model.NewModelFromString(getRBACModelText())
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin model: %w", err)
	}

	// 创建Casbin执行器
	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	// 启用自动保存
	enforcer.EnableAutoSave(true)

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, fmt.Errorf("failed to load casbin policy: %w", err)
	}

	return &CasbinEnforcer{
		enforcer: enforcer,
		db:       db,
	}, nil
}

// Enforce 检查权限
func (c *CasbinEnforcer) Enforce(sub, obj, act string) (bool, error) {
	return c.enforcer.Enforce(sub, obj, act)
}

// AddPolicy 添加权限策略
func (c *CasbinEnforcer) AddPolicy(sub, obj, act string) (bool, error) {
	return c.enforcer.AddPolicy(sub, obj, act, "allow")
}

// RemovePolicy 移除权限策略
func (c *CasbinEnforcer) RemovePolicy(sub, obj, act string) (bool, error) {
	return c.enforcer.RemovePolicy(sub, obj, act, "allow")
}

// AddRoleForUser 为用户添加角色
func (c *CasbinEnforcer) AddRoleForUser(user, role string) (bool, error) {
	return c.enforcer.AddRoleForUser(user, role)
}

// DeleteRoleForUser 删除用户角色
func (c *CasbinEnforcer) DeleteRoleForUser(user, role string) (bool, error) {
	return c.enforcer.DeleteRoleForUser(user, role)
}

// DeleteRolesForUser 删除用户的所有角色
func (c *CasbinEnforcer) DeleteRolesForUser(user string) (bool, error) {
	return c.enforcer.DeleteRolesForUser(user)
}

// GetRolesForUser 获取用户的所有角色
func (c *CasbinEnforcer) GetRolesForUser(user string) ([]string, error) {
	return c.enforcer.GetRolesForUser(user)
}

// GetUsersForRole 获取角色的所有用户
func (c *CasbinEnforcer) GetUsersForRole(role string) ([]string, error) {
	return c.enforcer.GetUsersForRole(role)
}

// GetPermissionsForUser 获取用户的所有权限
func (c *CasbinEnforcer) GetPermissionsForUser(user string) ([][]string, error) {
	return c.enforcer.GetPermissionsForUser(user)
}

// GetPermissionsForRole 获取角色的所有权限
func (c *CasbinEnforcer) GetPermissionsForRole(role string) ([][]string, error) {
	// Casbin v2中没有GetPermissionsForRole方法，使用GetFilteredPolicy替代
	return c.enforcer.GetFilteredPolicy(0, role)
}

// HasRoleForUser 检查用户是否有指定角色
func (c *CasbinEnforcer) HasRoleForUser(user, role string) (bool, error) {
	return c.enforcer.HasRoleForUser(user, role)
}



// AddPermissionForUser 为用户添加权限
func (c *CasbinEnforcer) AddPermissionForUser(user, obj, act string) (bool, error) {
	return c.enforcer.AddPermissionForUser(user, obj, act)
}

// DeletePermissionForUser 删除用户权限
func (c *CasbinEnforcer) DeletePermissionForUser(user, obj, act string) (bool, error) {
	return c.enforcer.DeletePermissionForUser(user, obj, act)
}

// DeletePermissionsForUser 删除用户的所有权限
func (c *CasbinEnforcer) DeletePermissionsForUser(user string) (bool, error) {
	return c.enforcer.DeletePermissionsForUser(user)
}

// SyncUserRoles 同步用户角色到Casbin
func (c *CasbinEnforcer) SyncUserRoles(userID uint, roleIDs []uint) error {
	userStr := fmt.Sprintf("user:%d", userID)
	
	// 删除用户的所有角色
	_, err := c.DeleteRolesForUser(userStr)
	if err != nil {
		return fmt.Errorf("failed to delete roles for user: %w", err)
	}

	// 添加新的角色
	for _, roleID := range roleIDs {
		roleStr := fmt.Sprintf("role:%d", roleID)
		_, err := c.AddRoleForUser(userStr, roleStr)
		if err != nil {
			return fmt.Errorf("failed to add role for user: %w", err)
		}
	}

	return nil
}

// SyncRolePermissions 同步角色权限到Casbin
func (c *CasbinEnforcer) SyncRolePermissions(roleID uint, permissionCodes []string) error {
	roleStr := fmt.Sprintf("role:%d", roleID)
	
	// 删除角色的所有权限
	policies, err := c.enforcer.GetFilteredPolicy(0, roleStr)
	if err != nil {
		return fmt.Errorf("failed to get filtered policy: %w", err)
	}
	for _, policy := range policies {
		if len(policy) >= 4 {
			_, err := c.RemovePolicy(policy[0], policy[1], policy[2])
			if err != nil {
				return fmt.Errorf("failed to remove policy: %w", err)
			}
		} else if len(policy) >= 3 {
			// 兼容旧格式的策略
			_, err := c.enforcer.RemovePolicy(policy[0], policy[1], policy[2])
			if err != nil {
				return fmt.Errorf("failed to remove policy: %w", err)
			}
		}
	}

	// 添加新的权限
	for _, permCode := range permissionCodes {
		_, err := c.AddPolicy(roleStr, permCode, "*")
		if err != nil {
			return fmt.Errorf("failed to add policy: %w", err)
		}
	}

	return nil
}

// CheckUserPermission 检查用户是否有指定权限
func (c *CasbinEnforcer) CheckUserPermission(userID uint, permissionCode string) (bool, error) {
	userStr := fmt.Sprintf("user:%d", userID)
	return c.Enforce(userStr, permissionCode, "*")
}

// CheckUserPermissionByPath 检查用户对指定路径和HTTP方法的权限
func (c *CasbinEnforcer) CheckUserPermissionByPath(userID uint, path, method string) (bool, error) {
	userStr := fmt.Sprintf("user:%d", userID)
	return c.Enforce(userStr, path, method)
}

// AddRESTfulPolicy 添加RESTful API权限策略
func (c *CasbinEnforcer) AddRESTfulPolicy(sub, path, method string) (bool, error) {
	return c.enforcer.AddPolicy(sub, path, method, "allow")
}

// RemoveRESTfulPolicy 移除RESTful API权限策略
func (c *CasbinEnforcer) RemoveRESTfulPolicy(sub, path, method string) (bool, error) {
	return c.enforcer.RemovePolicy(sub, path, method, "allow")
}

// AddRoleRESTfulPolicy 为角色添加RESTful API权限
func (c *CasbinEnforcer) AddRoleRESTfulPolicy(roleID uint, path, method string) error {
	roleStr := fmt.Sprintf("role:%d", roleID)
	_, err := c.AddRESTfulPolicy(roleStr, path, method)
	if err != nil {
		return fmt.Errorf("failed to add RESTful policy for role: %w", err)
	}
	return nil
}

// RemoveRoleRESTfulPolicy 移除角色的RESTful API权限
func (c *CasbinEnforcer) RemoveRoleRESTfulPolicy(roleID uint, path, method string) error {
	roleStr := fmt.Sprintf("role:%d", roleID)
	_, err := c.RemoveRESTfulPolicy(roleStr, path, method)
	if err != nil {
		return fmt.Errorf("failed to remove RESTful policy for role: %w", err)
	}
	return nil
}

// GetFilteredPolicy 获取过滤后的策略
func (c *CasbinEnforcer) GetFilteredPolicy(fieldIndex int, fieldValues ...string) [][]string {
	policies, _ := c.enforcer.GetFilteredPolicy(fieldIndex, fieldValues...)
	return policies
}

// LoadPolicy 重新加载策略
func (c *CasbinEnforcer) LoadPolicy() error {
	return c.enforcer.LoadPolicy()
}

// SavePolicy 保存策略到数据库
func (c *CasbinEnforcer) SavePolicy() error {
	return c.enforcer.SavePolicy()
}