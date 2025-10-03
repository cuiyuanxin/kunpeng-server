package service

import (
	"sync"

	"github.com/cuiyuanxin/kunpeng/internal/interfaces/service"
	"github.com/cuiyuanxin/kunpeng/internal/repository"
	"github.com/cuiyuanxin/kunpeng/internal/service/impl"
)

var (
	userService           service.UserService
	roleService           service.RoleService
	menuService           service.MenuService
	apiService            service.APIService
	deptService           service.DeptService
	postService           service.PostService
	dictService           service.DictService
	loginLogService       service.LoginLogService
	operationLogService   service.OperationLogService
	loginAttemptService   service.LoginAttemptService
	tokenBlacklistService service.TokenBlacklistService
	once                  sync.Once
)

// GetUserService 获取用户服务
func GetUserService() service.UserService {
	once.Do(initService)
	return userService
}

// GetRoleService 获取角色服务
func GetRoleService() service.RoleService {
	once.Do(initService)
	return roleService
}

// GetMenuService 获取菜单服务
func GetMenuService() service.MenuService {
	once.Do(initService)
	return menuService
}

// GetAPIService 获取API服务
func GetAPIService() service.APIService {
	once.Do(initService)
	return apiService
}

// GetDeptService 获取部门服务
func GetDeptService() service.DeptService {
	once.Do(initService)
	return deptService
}

// GetPostService 获取岗位服务
func GetPostService() service.PostService {
	once.Do(initService)
	return postService
}

// GetDictService 获取字典服务
func GetDictService() service.DictService {
	once.Do(initService)
	return dictService
}

// GetLoginLogService 获取登录日志服务
func GetLoginLogService() service.LoginLogService {
	once.Do(initService)
	return loginLogService
}

// GetOperationLogService 获取操作日志服务
func GetOperationLogService() service.OperationLogService {
	once.Do(initService)
	return operationLogService
}

// GetLoginAttemptService 获取登录尝试服务
func GetLoginAttemptService() service.LoginAttemptService {
	once.Do(initService)
	return loginAttemptService
}

// GetTokenBlacklistService 获取token黑名单服务
func GetTokenBlacklistService() service.TokenBlacklistService {
	once.Do(initService)
	return tokenBlacklistService
}

// initService 初始化服务
func initService() {
	// 初始化登录尝试服务
	loginAttemptService = impl.NewLoginAttemptService(repository.GetLoginAttemptRepository())
	// 初始化token黑名单服务
	tokenBlacklistService = impl.NewTokenBlacklistService(repository.GetTokenBlacklistRepository())

	// 初始化用户服务（需要依赖登录尝试服务）
	userService = impl.NewUserService(loginAttemptService)

	// 初始化其他服务
	roleService = &impl.RoleServiceImpl{}
	menuService = &impl.MenuServiceImpl{}
	apiService = &impl.APIServiceImpl{}
	deptService = &impl.DeptServiceImpl{}
	postService = &impl.PostServiceImpl{}
	dictService = &impl.DictServiceImpl{}
	loginLogService = &impl.LoginLogServiceImpl{}
	operationLogService = &impl.OperationLogServiceImpl{}
}
