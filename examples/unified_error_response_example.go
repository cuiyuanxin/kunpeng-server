package examples

// 统一错误处理和响应示例
// 展示如何使用融合后的response和errors包

import (
	"github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserController 用户控制器示例
type UserController struct {
	userService *UserService
}

// UserService 用户服务示例
type UserService struct {
	db *gorm.DB
}

// DemoUser 演示用户模型（用于错误处理示例）
type DemoUser struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Username string `json:"username" gorm:"uniqueIndex"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Status   string `json:"status" gorm:"default:'active'"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=20"`
	Email    string `json:"email" binding:"omitempty,email"`
	Status   string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ===== 控制器层示例 =====

// GetUser 获取用户信息
func (ctrl *UserController) GetUser(c *gin.Context) {
	userID := c.Param("id")
	
	// 使用服务层方法，自动处理错误
	user, err := ctrl.userService.GetUserByID(userID)
	response.SuccessOrError(c, user, err)
}

// CreateUser 创建用户
func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	
	// 参数绑定验证
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	
	// 调用服务层创建用户
	user, err := ctrl.userService.CreateUser(req.Username, req.Email)
	response.SuccessOrError(c, user, err)
}

// UpdateUser 更新用户
func (ctrl *UserController) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var req UpdateUserRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}
	
	user, err := ctrl.userService.UpdateUser(userID, &req)
	response.SuccessOrError(c, user, err)
}

// DeleteUser 删除用户
func (ctrl *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	
	err := ctrl.userService.DeleteUser(userID)
	if err != nil {
		response.HandleBusinessError(c, err)
		return
	}
	
	response.Success(c, gin.H{"message": "用户删除成功"})
}

// GetUsers 获取用户列表（分页）
func (ctrl *UserController) GetUsers(c *gin.Context) {
	page := 1
	size := 10
	
	// 这里可以从查询参数获取page和size
	// page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	// size, _ = strconv.Atoi(c.DefaultQuery("size", "10"))
	
	users, total, err := ctrl.userService.GetUsers(page, size)
	response.SuccessPageOrError(c, users, total, page, size, err)
}

// ===== 服务层示例 =====

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(userID string) (*DemoUser, error) {
	if userID == "" {
		return nil, errors.ErrInvalidParams.WithDetails("用户ID不能为空")
	}
	
	var user DemoUser
	err := s.db.First(&user, "id = ?", userID).Error
	if err != nil {
		return nil, errors.ConvertGormError(err)
	}
	
	// 检查用户状态
	if user.Status == "deleted" {
		return nil, errors.ErrUserNotFound.WithDetails("用户已被删除")
	}
	
	return &user, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(username, email string) (*DemoUser, error) {
	// 使用安全执行包装，防止panic
	result, err := errors.SafeExecuteWithResult(func() (interface{}, error) {
		// 检查用户名是否已存在
		var count int64
		err := s.db.Model(&DemoUser{}).Where("username = ?", username).Count(&count).Error
		if err != nil {
			return nil, errors.ConvertGormError(err)
		}
		
		if count > 0 {
			return nil, errors.ErrUsernameTaken.WithDetailsf("用户名 '%s' 已被占用", username)
		}
		
		// 检查邮箱是否已存在
		err = s.db.Model(&DemoUser{}).Where("email = ?", email).Count(&count).Error
		if err != nil {
			return nil, errors.ConvertGormError(err)
		}
		
		if count > 0 {
			return nil, errors.ErrEmailTaken.WithDetailsf("邮箱 '%s' 已被占用", email)
		}
		
		// 创建用户
		user := &DemoUser{
			Username: username,
			Email:    email,
			Status:   "active",
		}
		
		err = s.db.Create(user).Error
		if err != nil {
			return nil, errors.ConvertGormError(err)
		}
		
		return user, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	return result.(*DemoUser), nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(userID string, req *UpdateUserRequest) (*DemoUser, error) {
	// 先获取用户
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	
	// 更新字段
	updateData := make(map[string]interface{})
	if req.Username != "" {
		updateData["username"] = req.Username
	}
	if req.Email != "" {
		updateData["email"] = req.Email
	}
	if req.Status != "" {
		updateData["status"] = req.Status
	}
	
	if len(updateData) == 0 {
		return user, nil // 没有更新内容
	}
	
	err = s.db.Model(user).Updates(updateData).Error
	if err != nil {
		return nil, errors.ConvertGormError(err)
	}
	
	return user, nil
}

// DeleteUser 删除用户（软删除）
func (s *UserService) DeleteUser(userID string) error {
	// 检查用户是否存在
	_, err := s.GetUserByID(userID)
	if err != nil {
		return err
	}
	
	// 软删除：更新状态为deleted
	err = s.db.Model(&DemoUser{}).Where("id = ?", userID).Update("status", "deleted").Error
	if err != nil {
		return errors.ConvertGormError(err)
	}
	
	return nil
}

// GetUsers 获取用户列表
func (s *UserService) GetUsers(page, size int) ([]*DemoUser, int64, error) {
	var users []*DemoUser
	var total int64
	
	// 使用链式错误处理
	err := errors.Chain(
		func() error {
			// 获取总数
			return s.db.Model(&DemoUser{}).Where("status != ?", "deleted").Count(&total).Error
		},
		func() error {
			// 获取分页数据
			offset := (page - 1) * size
			return s.db.Where("status != ?", "deleted").Offset(offset).Limit(size).Find(&users).Error
		},
	)
	
	if err != nil {
		return nil, 0, errors.ConvertGormError(err)
	}
	
	return users, total, nil
}

// ===== 中间件使用示例 =====

// AuthMiddleware 认证中间件示例
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			// 直接中断并返回错误
			errors.AbortWithError(c, errors.ErrUnauthorized.WithDetails("缺少认证令牌"))
			return
		}
		
		// 验证token（这里简化处理）
		if token != "Bearer valid-token" {
			errors.AbortWithError(c, errors.ErrTokenInvalid.WithDetails("令牌无效或已过期"))
			return
		}
		
		c.Set("user_id", "123")
		c.Next()
	}
}

// PermissionMiddleware 权限检查中间件示例
func PermissionMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetString("user_id")
		if userID == "" {
			errors.AbortWithError(c, errors.ErrUnauthorized)
			return
		}
		
		// 检查权限（这里简化处理）
		if !hasPermission(userID, requiredPermission) {
			errors.AbortWithError(c, errors.ErrPermissionDenied.WithDetailsf("缺少权限: %s", requiredPermission))
			return
		}
		
		c.Next()
	}
}

// hasPermission 检查用户是否有指定权限（示例实现）
func hasPermission(userID, permission string) bool {
	// 这里应该查询数据库或缓存来检查权限
	return userID == "123" && permission == "user:read"
}

// ===== 路由注册示例 =====

// RegisterUserRoutes 注册用户相关路由
func RegisterUserRoutes(r *gin.Engine, userController *UserController) {
	// 使用错误处理中间件
	r.Use(errors.ErrorHandlerMiddleware())
	
	api := r.Group("/api/v1")
	{
		// 公开接口
		api.POST("/users", userController.CreateUser)
		
		// 需要认证的接口
		auth := api.Group("/users", AuthMiddleware())
		{
			auth.GET("/:id", userController.GetUser)
			auth.GET("", PermissionMiddleware("user:read"), userController.GetUsers)
			auth.PUT("/:id", PermissionMiddleware("user:write"), userController.UpdateUser)
			auth.DELETE("/:id", PermissionMiddleware("user:delete"), userController.DeleteUser)
		}
	}
}

// ===== 外部服务调用示例 =====

// ExternalAPIService 外部API服务示例
type ExternalAPIService struct {
	client interface{} // HTTP客户端
}

// CallExternalAPI 调用外部API（带重试机制）
func (s *ExternalAPIService) CallExternalAPI(userID string) (interface{}, error) {
	return errors.SafeExecuteWithResult(func() (interface{}, error) {
		// 使用重试机制调用外部API
		return nil, errors.RetryOnError(func() error {
			// 模拟API调用
			// result, err := s.client.Get(fmt.Sprintf("/api/users/%s", userID))
			// if err != nil {
			//     return errors.ErrNetworkError.WithDetails(err.Error())
			// }
			// return nil
			return nil
		}, 3) // 最多重试3次
	})
}