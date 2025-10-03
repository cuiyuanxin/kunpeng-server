package casbin

import (
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/cuiyuanxin/kunpeng/pkg/database"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/cuiyuanxin/kunpeng/pkg/logger"
	"go.uber.org/zap"
)

var (
	enforcer *casbin.Enforcer
	once     sync.Once
)

// 定义RBAC模型规则
const rbacModelText = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")
`

// Init 初始化Casbin
func Init() error {
	var err error
	once.Do(func() {
		// 检查数据库连接
		db := database.GetDB()
		if db == nil {
			err = kperrors.New(kperrors.ErrSystem, nil).WithMessage("数据库连接未初始化")
			return
		}

		// 创建适配器
		adapter, adapterErr := gormadapter.NewAdapterByDB(db)
		if adapterErr != nil {
			err = kperrors.New(kperrors.ErrSystem, adapterErr)
			return
		}

		// 从字符串加载模型
		m, modelErr := model.NewModelFromString(rbacModelText)
		if modelErr != nil {
			err = kperrors.New(kperrors.ErrSystem, modelErr)
			return
		}

		// 创建执行器
		enforcer, err = casbin.NewEnforcer(m, adapter)
		if err != nil {
			err = kperrors.New(kperrors.ErrSystem, err)
			return
		}

		// 加载策略
		if err = enforcer.LoadPolicy(); err != nil {
			err = kperrors.New(kperrors.ErrSystem, err)
			return
		}

		logger.GetLogger().Info("Casbin初始化成功")
	})

	return err
}

// GetEnforcer 获取Casbin执行器
func GetEnforcer() *casbin.Enforcer {
	return enforcer
}

// Enforce 检查权限
func Enforce(sub string, obj string, act string) (bool, error) {
	if enforcer == nil {
		return false, kperrors.New(kperrors.ErrSystem, nil).WithMessage("Casbin未初始化")
	}

	// 检查权限
	result, err := enforcer.Enforce(sub, obj, act)
	if err != nil {
		logger.GetLogger().Error("Casbin权限检查失败",
			zap.Error(err),
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act))
		return false, kperrors.New(kperrors.ErrPermDenied, err)
	}

	if !result {
		logger.GetLogger().Warn("权限不足",
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act))
		return false, nil
	}

	return true, nil
}

// AddPolicy 添加策略
func AddPolicy(sub string, obj string, act string) (bool, error) {
	if enforcer == nil {
		return false, kperrors.New(kperrors.ErrSystem, nil).WithMessage("Casbin未初始化")
	}

	// 添加策略
	result, err := enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		logger.GetLogger().Error("添加策略失败",
			zap.Error(err),
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act))
		return false, kperrors.New(kperrors.ErrSystem, err)
	}

	// 保存策略
	if err := enforcer.SavePolicy(); err != nil {
		logger.GetLogger().Error("保存策略失败", zap.Error(err))
		return false, kperrors.New(kperrors.ErrSystem, err)
	}

	return result, nil
}

// RemovePolicy 删除策略
func RemovePolicy(sub string, obj string, act string) (bool, error) {
	if enforcer == nil {
		return false, kperrors.New(kperrors.ErrSystem, nil).WithMessage("Casbin未初始化")
	}

	// 删除策略
	result, err := enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		logger.GetLogger().Error("删除策略失败",
			zap.Error(err),
			zap.String("subject", sub),
			zap.String("object", obj),
			zap.String("action", act))
		return false, kperrors.New(kperrors.ErrSystem, err)
	}

	// 保存策略
	if err := enforcer.SavePolicy(); err != nil {
		logger.GetLogger().Error("保存策略失败", zap.Error(err))
		return false, kperrors.New(kperrors.ErrSystem, err)
	}

	return result, nil
}

// AddRoleForUser 为用户添加角色
func AddRoleForUser(user string, role string) (bool, error) {
	if enforcer == nil {
		return false, kperrors.New(kperrors.ErrSystem, nil).WithMessage("Casbin未初始化")
	}

	// 添加角色
	result, err := enforcer.AddRoleForUser(user, role)
	if err != nil {
		logger.GetLogger().Error("为用户添加角色失败",
			zap.Error(err),
			zap.String("user", user),
			zap.String("role", role))
		return false, kperrors.New(kperrors.ErrSystem, err)
	}

	// 保存策略
	if err := enforcer.SavePolicy(); err != nil {
		logger.GetLogger().Error("保存策略失败", zap.Error(err))
		return false, kperrors.New(kperrors.ErrSystem, err)
	}

	return result, nil
}

// DeleteRoleForUser 删除用户的角色
func DeleteRoleForUser(user string, role string) (bool, error) {
	if enforcer == nil {
		return false, kperrors.New(kperrors.ErrSystem, nil).WithMessage("Casbin未初始化")
	}

	// 删除角色
	result, err := enforcer.DeleteRoleForUser(user, role)
	if err != nil {
		logger.GetLogger().Error("删除用户角色失败",
			zap.Error(err),
			zap.String("user", user),
			zap.String("role", role))
		return false, kperrors.New(kperrors.ErrSystem, err)
	}

	// 保存策略
	if err := enforcer.SavePolicy(); err != nil {
		logger.GetLogger().Error("保存策略失败", zap.Error(err))
		return false, kperrors.New(kperrors.ErrSystem, err)
	}

	return result, nil
}

// GetRolesForUser 获取用户的所有角色
func GetRolesForUser(user string) ([]string, error) {
	if enforcer == nil {
		return nil, kperrors.New(kperrors.ErrSystem, nil).WithMessage("Casbin未初始化")
	}

	// 获取角色
	roles, err := enforcer.GetRolesForUser(user)
	if err != nil {
		logger.GetLogger().Error("获取用户角色失败",
			zap.Error(err),
			zap.String("user", user))
		return nil, kperrors.New(kperrors.ErrSystem, err)
	}

	return roles, nil
}

// GetUsersForRole 获取拥有指定角色的所有用户
func GetUsersForRole(role string) ([]string, error) {
	if enforcer == nil {
		return nil, kperrors.New(kperrors.ErrSystem, nil).WithMessage("Casbin未初始化")
	}

	// 获取用户
	users, err := enforcer.GetUsersForRole(role)
	if err != nil {
		logger.GetLogger().Error("获取角色用户失败",
			zap.Error(err),
			zap.String("role", role))
		return nil, kperrors.New(kperrors.ErrSystem, err)
	}

	return users, nil
}

// HasRoleForUser 检查用户是否拥有指定角色
func HasRoleForUser(user string, role string) (bool, error) {
	if enforcer == nil {
		return false, kperrors.New(kperrors.ErrSystem, nil).WithMessage("Casbin未初始化")
	}

	// 检查角色
	result, err := enforcer.HasRoleForUser(user, role)
	if err != nil {
		logger.GetLogger().Error("检查用户角色失败",
			zap.Error(err),
			zap.String("user", user),
			zap.String("role", role))
		return false, kperrors.New(kperrors.ErrSystem, err)
	}

	return result, nil
}

// GetPermissionsForUser 获取用户的所有权限
func GetPermissionsForUser(user string) ([][]string, error) {
	if enforcer == nil {
		return nil, kperrors.New(kperrors.ErrSystem, nil).WithMessage("Casbin未初始化")
	}

	// 获取权限
	permissions, err := enforcer.GetPermissionsForUser(user)
	if err != nil {
		logger.GetLogger().Error("获取用户权限失败",
			zap.Error(err),
			zap.String("user", user))
		return nil, kperrors.New(kperrors.ErrSystem, err)
	}

	return permissions, nil
}

// GetImplicitPermissionsForUser 获取用户的所有隐式权限（包括角色继承的权限）
func GetImplicitPermissionsForUser(user string) ([][]string, error) {
	if enforcer == nil {
		return nil, kperrors.New(kperrors.ErrSystem, nil).WithMessage("Casbin未初始化")
	}

	// 获取隐式权限
	permissions, err := enforcer.GetImplicitPermissionsForUser(user)
	if err != nil {
		logger.GetLogger().Error("获取用户隐式权限失败",
			zap.Error(err),
			zap.String("user", user))
		return nil, kperrors.New(kperrors.ErrSystem, err)
	}

	return permissions, nil
}

// LoadPolicy 重新加载策略
func LoadPolicy() error {
	if enforcer == nil {
		return kperrors.New(kperrors.ErrSystem, nil).WithMessage("Casbin未初始化")
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		logger.GetLogger().Error("加载策略失败", zap.Error(err))
		return kperrors.New(kperrors.ErrSystem, err)
	}

	return nil
}
