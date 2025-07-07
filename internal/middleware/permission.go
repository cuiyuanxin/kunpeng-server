package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/cuiyuanxin/kunpeng/internal/service"
)

// PermissionMiddleware 权限检查中间件
func PermissionMiddleware(permissionService *service.PermissionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(uint)
		if !ok {
			// 尝试从字符串转换
			if userIDStr, ok := userIDInterface.(string); ok {
				if id, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
					userID = uint(id)
				} else {
					response.Unauthorized(c, "用户ID格式错误")
					c.Abort()
					return
				}
			} else {
				response.Unauthorized(c, "用户ID格式错误")
				c.Abort()
				return
			}
		}

		// 检查是否为超级管理员
		isSuperAdmin, exists := c.Get("is_super_admin")
		if exists && isSuperAdmin.(bool) {
			// 超级管理员跳过权限检查
			c.Next()
			return
		}

		// 获取请求路径和方法
		path := c.Request.URL.Path
		method := c.Request.Method

		// 检查权限
		hasPermission, err := permissionService.CheckUserPermissionByPath(userID, path, method)
		if err != nil {
			response.ServerError(c, "权限检查失败")
			c.Abort()
			return
		}

		if !hasPermission {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequirePermission 要求特定权限的中间件
func RequirePermission(permissionService *service.PermissionService, permissionCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(uint)
		if !ok {
			// 尝试从字符串转换
			if userIDStr, ok := userIDInterface.(string); ok {
				if id, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
					userID = uint(id)
				} else {
					response.Unauthorized(c, "用户ID格式错误")
					c.Abort()
					return
				}
			} else {
				response.Unauthorized(c, "用户ID格式错误")
				c.Abort()
				return
			}
		}

		// 检查是否为超级管理员
		isSuperAdmin, exists := c.Get("is_super_admin")
		if exists && isSuperAdmin.(bool) {
			// 超级管理员跳过权限检查
			c.Next()
			return
		}

		// 检查特定权限
		hasPermission, err := permissionService.CheckUserPermission(userID, permissionCode)
		if err != nil {
			response.ServerError(c, "权限检查失败")
			c.Abort()
			return
		}

		if !hasPermission {
			response.Forbidden(c, "权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole 要求特定角色的中间件
func RequireRole(permissionService *service.PermissionService, roleCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取用户ID
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			response.Unauthorized(c, "未登录")
			c.Abort()
			return
		}

		userID, ok := userIDInterface.(uint)
		if !ok {
			// 尝试从字符串转换
			if userIDStr, ok := userIDInterface.(string); ok {
				if id, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
					userID = uint(id)
				} else {
					response.Unauthorized(c, "用户ID格式错误")
					c.Abort()
					return
				}
			} else {
				response.Unauthorized(c, "用户ID格式错误")
				c.Abort()
				return
			}
		}

		// 检查是否为超级管理员
		isSuperAdmin, exists := c.Get("is_super_admin")
		if exists && isSuperAdmin.(bool) {
			// 超级管理员跳过权限检查
			c.Next()
			return
		}

		// 获取用户角色
		userRoles, err := permissionService.GetUserRoles(userID)
		if err != nil {
			response.ServerError(c, "角色检查失败")
			c.Abort()
			return
		}

		// 检查是否有指定角色
		hasRole := false
		for _, role := range userRoles {
			if role == roleCode {
				hasRole = true
				break
			}
		}

		if !hasRole {
			response.Forbidden(c, "角色权限不足")
			c.Abort()
			return
		}

		c.Next()
	}
}