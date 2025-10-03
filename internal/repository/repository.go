package repository

import (
	"sync"

	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/repository/impl"
)

var (
	userRepository           repository.UserRepository
	roleRepository           repository.RoleRepository
	menuRepository           repository.MenuRepository
	apiRepository            repository.APIRepository
	deptRepository           repository.DeptRepository
	postRepository           repository.PostRepository
	dictRepository           repository.DictRepository
	loginLogRepository       repository.LoginLogRepository
	operationLogRepository   repository.OperationLogRepository
	loginAttemptRepository   repository.LoginAttemptRepository
	tokenBlacklistRepository repository.TokenBlacklistRepository
	once                     sync.Once
)

// Init 初始化所有仓储
func Init() {
	once.Do(func() {
		userRepository = impl.NewUserRepository()
		roleRepository = impl.NewRoleRepository()
		menuRepository = impl.NewMenuRepository()
		apiRepository = impl.NewAPIRepository()
		deptRepository = impl.NewDeptRepository()
		postRepository = impl.NewPostRepository()
		dictRepository = impl.NewDictRepository()

		// 初始化登录日志仓储
		loginLogRepository = impl.NewLoginLogRepository()
		// 初始化操作日志仓储
		operationLogRepository = impl.NewOperationLogRepository()
		// 初始化登录尝试仓储
		loginAttemptRepository = impl.NewLoginAttemptRepository()
		// 初始化token黑名单仓储
		tokenBlacklistRepository = impl.NewTokenBlacklistRepository()
	})
}

// GetUserRepository 获取用户仓储
func GetUserRepository() repository.UserRepository {
	return userRepository
}

// GetRoleRepository 获取角色仓储
func GetRoleRepository() repository.RoleRepository {
	return roleRepository
}

// GetMenuRepository 获取菜单仓储
func GetMenuRepository() repository.MenuRepository {
	return menuRepository
}

// GetAPIRepository 获取API仓储
func GetAPIRepository() repository.APIRepository {
	return apiRepository
}

// GetDeptRepository 获取部门仓储
func GetDeptRepository() repository.DeptRepository {
	return deptRepository
}

// GetPostRepository 获取岗位仓储
func GetPostRepository() repository.PostRepository {
	return postRepository
}

// GetDictRepository 获取字典仓储
func GetDictRepository() repository.DictRepository {
	return dictRepository
}

// GetLoginLogRepository 获取登录日志仓储
func GetLoginLogRepository() repository.LoginLogRepository {
	return loginLogRepository
}

// GetOperationLogRepository 获取操作日志仓储
func GetOperationLogRepository() repository.OperationLogRepository {
	return operationLogRepository
}

// GetLoginAttemptRepository 获取登录尝试仓储
func GetLoginAttemptRepository() repository.LoginAttemptRepository {
	return loginAttemptRepository
}

// GetTokenBlacklistRepository 获取token黑名单仓储
func GetTokenBlacklistRepository() repository.TokenBlacklistRepository {
	return tokenBlacklistRepository
}
