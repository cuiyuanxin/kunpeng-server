package controller

import (
	"sync"
)

var (
	userController         UserController
	roleController         RoleController
	menuController         MenuController
	apiController          APIController
	deptController         DeptController
	postController         PostController
	dictController         DictController
	loginLogController     LoginLogController
	operationLogController OperationLogController
	once                   sync.Once
)

// GetUserController 获取用户控制器
func GetUserController() *UserController {
	once.Do(initController)
	return &userController
}

// GetRoleController 获取角色控制器
func GetRoleController() *RoleController {
	once.Do(initController)
	return &roleController
}

// GetMenuController 获取菜单控制器
func GetMenuController() *MenuController {
	once.Do(initController)
	return &menuController
}

// GetAPIController 获取API控制器
func GetAPIController() *APIController {
	once.Do(initController)
	return &apiController
}

// GetDeptController 获取部门控制器
func GetDeptController() *DeptController {
	once.Do(initController)
	return &deptController
}

// GetPostController 获取岗位控制器
func GetPostController() *PostController {
	once.Do(initController)
	return &postController
}

// GetDictController 获取字典控制器
func GetDictController() *DictController {
	once.Do(initController)
	return &dictController
}

// GetLoginLogController 获取登录日志控制器
func GetLoginLogController() *LoginLogController {
	once.Do(initController)
	return &loginLogController
}

// GetOperationLogController 获取操作日志控制器
func GetOperationLogController() *OperationLogController {
	once.Do(initController)
	return &operationLogController
}

// initController 初始化控制器
func initController() {
	userController = UserController{}
	roleController = RoleController{}
	menuController = MenuController{}
	apiController = APIController{}
	deptController = DeptController{}
	postController = PostController{}
	dictController = DictController{}
	loginLogController = LoginLogController{}
	operationLogController = OperationLogController{}
}
