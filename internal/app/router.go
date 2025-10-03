package app

import (
	"github.com/cuiyuanxin/kunpeng/internal/controller"
	"github.com/cuiyuanxin/kunpeng/internal/middleware"
	"github.com/cuiyuanxin/kunpeng/pkg/config"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// RegisterRoutes 注册路由
func (a *App) RegisterRoutes() {
	// 注册Swagger
	a.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 健康检查
	a.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"app":    config.GetAppConfig().Name,
			"env":    config.GetAppConfig().Mode,
		})
	})

	// API v1
	v1 := a.engine.Group("/api/v1")
	{
		// 无需认证的接口
		v1.POST("/login", controller.GetUserController().Login)
		v1.POST("/refresh-token", controller.GetUserController().RefreshToken)
		v1.POST("/logout", controller.GetUserController().Logout)

		// 需要认证的接口
		auth := v1.Group("")
		auth.Use(middleware.JWT())
		{
			// 用户相关接口
			auth.GET("/user/info", controller.GetUserController().GetUserInfo)
			auth.GET("/users", controller.GetUserController().GetUserList)
			auth.GET("/users/:id", controller.GetUserController().GetUserByID)
			auth.POST("/users", controller.GetUserController().CreateUser)
			auth.PUT("/users", controller.GetUserController().UpdateUser)
			auth.DELETE("/users/:id", controller.GetUserController().DeleteUser)
			auth.DELETE("/users/batch", controller.GetUserController().BatchDeleteUser)
			auth.PUT("/users/status", controller.GetUserController().ChangeUserStatus)
			auth.PUT("/users/:id/password/reset", controller.GetUserController().ResetUserPassword)
			auth.PUT("/users/password", controller.GetUserController().ChangePassword)

			// 角色相关接口
			auth.GET("/roles", controller.GetRoleController().GetRoleList)
			auth.GET("/roles/all", controller.GetRoleController().GetAllRoles)
			auth.GET("/roles/:id", controller.GetRoleController().GetRoleByID)
			auth.POST("/roles", controller.GetRoleController().CreateRole)
			auth.PUT("/roles", controller.GetRoleController().UpdateRole)
			auth.DELETE("/roles/:id", controller.GetRoleController().DeleteRole)
			auth.DELETE("/roles/batch", controller.GetRoleController().BatchDeleteRole)
			auth.PUT("/roles/status", controller.GetRoleController().ChangeRoleStatus)
			auth.GET("/roles/:id/menus", controller.GetRoleController().GetRoleMenus)
			auth.PUT("/roles/menus", controller.GetRoleController().AssignMenus)
			auth.GET("/roles/:id/apis", controller.GetRoleController().GetRoleAPIs)
			auth.PUT("/roles/apis", controller.GetRoleController().UpdateRoleAPIs)

			// 菜单相关接口
			auth.GET("/menus", controller.GetMenuController().GetMenuList)
			auth.GET("/menus/tree", controller.GetMenuController().GetMenuTree)
			auth.GET("/menus/user", controller.GetMenuController().GetUserMenuTree)
			auth.GET("/menus/:id", controller.GetMenuController().GetMenuByID)
			auth.POST("/menus", controller.GetMenuController().CreateMenu)
			auth.PUT("/menus", controller.GetMenuController().UpdateMenu)
			auth.DELETE("/menus/:id", controller.GetMenuController().DeleteMenu)

			// API相关接口
			auth.GET("/apis", controller.GetAPIController().GetAPIList)
			auth.GET("/apis/all", controller.GetAPIController().GetAllAPIs)
			auth.GET("/apis/:id", controller.GetAPIController().GetAPIByID)
			auth.POST("/apis", controller.GetAPIController().CreateAPI)
			auth.PUT("/apis", controller.GetAPIController().UpdateAPI)
			auth.DELETE("/apis/:id", controller.GetAPIController().DeleteAPI)
			auth.DELETE("/apis/batch", controller.GetAPIController().BatchDeleteAPI)

			// 部门相关接口
			auth.GET("/depts", controller.GetDeptController().GetDeptList)
			auth.GET("/depts/tree", controller.GetDeptController().GetDeptTree)
			auth.GET("/depts/:id", controller.GetDeptController().GetDeptByID)
			auth.POST("/depts", controller.GetDeptController().CreateDept)
			auth.PUT("/depts", controller.GetDeptController().UpdateDept)
			auth.DELETE("/depts/:id", controller.GetDeptController().DeleteDept)

			// 岗位相关接口
			auth.GET("/posts", controller.GetPostController().GetPostList)
			auth.GET("/posts/all", controller.GetPostController().GetAllPosts)
			auth.GET("/posts/:id", controller.GetPostController().GetPostByID)
			auth.POST("/posts", controller.GetPostController().CreatePost)
			auth.PUT("/posts", controller.GetPostController().UpdatePost)
			auth.DELETE("/posts/:id", controller.GetPostController().DeletePost)
			auth.DELETE("/posts/batch", controller.GetPostController().BatchDeletePost)

			// 字典相关接口
			auth.GET("/dict/types", controller.GetDictController().GetDictTypeList)
			auth.GET("/dict/types/all", controller.GetDictController().GetAllDictTypes)
			auth.GET("/dict/types/:id", controller.GetDictController().GetDictTypeByID)
			auth.POST("/dict/types", controller.GetDictController().CreateDictType)
			auth.PUT("/dict/types", controller.GetDictController().UpdateDictType)
			auth.DELETE("/dict/types/:id", controller.GetDictController().DeleteDictType)
			auth.DELETE("/dict/types/batch", controller.GetDictController().BatchDeleteDictType)
			auth.GET("/dict/data", controller.GetDictController().GetDictDataList)
			auth.GET("/dict/data/type/:dict_type", controller.GetDictController().GetDictDataByType)
			auth.GET("/dict/data/:id", controller.GetDictController().GetDictDataByID)
			auth.POST("/dict/data", controller.GetDictController().CreateDictData)
			auth.PUT("/dict/data", controller.GetDictController().UpdateDictData)
			auth.DELETE("/dict/data/:id", controller.GetDictController().DeleteDictData)
			auth.DELETE("/dict/data/batch", controller.GetDictController().BatchDeleteDictData)

			// 登录日志相关接口
			auth.GET("/login-logs", controller.GetLoginLogController().GetLoginLogList)
			auth.GET("/login-logs/:id", controller.GetLoginLogController().GetLoginLogByID)
			auth.DELETE("/login-logs/:id", controller.GetLoginLogController().DeleteLoginLog)
			auth.DELETE("/login-logs/batch", controller.GetLoginLogController().BatchDeleteLoginLog)
			auth.POST("/login-logs/clean", controller.GetLoginLogController().CleanOldLogs)

			// 操作日志相关接口
			auth.GET("/operation-logs", controller.GetOperationLogController().GetOperationLogList)
			auth.GET("/operation-logs/:id", controller.GetOperationLogController().GetOperationLogByID)
			auth.DELETE("/operation-logs/:id", controller.GetOperationLogController().DeleteOperationLog)
			auth.DELETE("/operation-logs/batch", controller.GetOperationLogController().BatchDeleteOperationLog)
			auth.POST("/operation-logs/clean", controller.GetOperationLogController().CleanOldLogs)
		}
	}
}
