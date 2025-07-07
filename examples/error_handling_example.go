package examples

// 这个文件展示了如何使用新的错误中心
// 注意：这是示例代码，不会被编译到最终程序中

import (
	"github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ExampleHandler 示例处理器
type ExampleHandler struct {
	db *gorm.DB
}

// 示例1: 使用预定义错误
func (h *ExampleHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")

	// 参数验证
	if userID == "" {
		// 使用预定义错误并添加详细信息
		errors.AbortWithError(c, errors.ErrInvalidParams.WithDetails("用户ID不能为空"))
		return
	}

	// 模拟数据库查询
	var user struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	err := h.db.First(&user, "id = ?", userID).Error
	if err != nil {
		// 自动转换GORM错误为业务错误
		errors.AbortWithError(c, errors.ConvertGormError(err))
		return
	}

	// 成功响应
	response.Success(c, user)
}

// 示例2: 在中间件中设置错误
func (h *ExampleHandler) CreateUser(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
	}

	// 参数绑定和验证
	if err := c.ShouldBindJSON(&req); err != nil {
		// 设置验证错误，由中间件处理
		errors.SetError(c, errors.ErrValidationFailed.WithDetails(err.Error()))
		return
	}

	// 业务逻辑验证
	if len(req.Username) < 3 {
		// 设置业务错误
		errors.SetError(c, errors.ErrValidationFailed.WithDetails("用户名长度不能少于3个字符"))
		return
	}

	// 模拟数据库操作
	user := struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	}{
		ID:       "123",
		Username: req.Username,
		Email:    req.Email,
	}

	err := h.db.Create(&user).Error
	if err != nil {
		// 设置数据库错误
		errors.SetError(c, errors.ConvertGormError(err))
		return
	}

	response.Success(c, user)
}

// 示例3: 服务层错误处理
type ExampleService struct {
	db *gorm.DB
}

func (s *ExampleService) GetUserByID(id string) (interface{}, error) {
	// 参数验证
	if id == "" {
		return nil, errors.ErrInvalidParams.WithDetails("用户ID不能为空")
	}

	var user struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Status string `json:"status"`
	}

	// 数据库查询
	err := s.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, errors.ConvertGormError(err)
	}

	// 业务逻辑验证
	if user.Status == "deleted" {
		return nil, errors.ErrUserNotFound.WithDetails("用户已被删除")
	}

	return user, nil
}

func (s *ExampleService) CreateUser(username, email string) (interface{}, error) {
	// 使用安全执行包装可能panic的代码
	return errors.SafeExecuteWithResult(func() (interface{}, error) {
		// 验证用户名是否已存在
		var count int64
		err := s.db.Model(&struct{ Username string }{}).Where("username = ?", username).Count(&count).Error
		if err != nil {
			return nil, errors.ConvertGormError(err)
		}

		if count > 0 {
			return nil, errors.ErrUserExists.WithDetailsf("用户名 '%s' 已被占用", username)
		}

		// 创建用户
		user := struct {
			ID       string `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
		}{
			ID:       "new-id",
			Username: username,
			Email:    email,
		}

		err = s.db.Create(&user).Error
		if err != nil {
			return nil, errors.ConvertGormError(err)
		}

		return user, nil
	})
}

// 示例4: 链式错误处理
func (s *ExampleService) ComplexUserOperation(userID string) error {
	return errors.Chain(
		func() error {
			// 步骤1: 验证用户存在
			_, err := s.GetUserByID(userID)
			return err
		},
		func() error {
			// 步骤2: 检查权限
			return s.checkUserPermissions(userID)
		},
		func() error {
			// 步骤3: 执行操作
			return s.performUserOperation(userID)
		},
		func() error {
			// 步骤4: 更新缓存
			return s.updateUserCache(userID)
		},
	)
}

func (s *ExampleService) checkUserPermissions(userID string) error {
	// 模拟权限检查
	if userID == "forbidden" {
		return errors.ErrPermissionDenied.WithDetails("用户无权限执行此操作")
	}
	return nil
}

func (s *ExampleService) performUserOperation(userID string) error {
	// 模拟操作
	if userID == "error" {
		return errors.ErrInternalError.WithDetails("操作执行失败")
	}
	return nil
}

func (s *ExampleService) updateUserCache(userID string) error {
	// 模拟缓存更新
	return nil
}

// 示例5: 重试机制
func (s *ExampleService) CallExternalAPI(userID string) error {
	return errors.RetryOnError(func() error {
		// 模拟外部API调用
		if userID == "retry" {
			return errors.ErrExternalServiceError.WithDetails("外部服务暂时不可用")
		}
		return nil
	}, 3) // 最多重试3次
}

// 示例6: 并行操作错误收集
func (s *ExampleService) ParallelUserOperations(userIDs []string) error {
	var operations []func() error

	for _, userID := range userIDs {
		userID := userID // 避免闭包问题
		operations = append(operations, func() error {
			return s.performUserOperation(userID)
		})
	}

	errs := errors.Parallel(operations...)
	if len(errs) > 0 {
		return errors.CombineErrors(errs...)
	}

	return nil
}

// 示例7: 自定义错误
const (
	ErrCodeCustomValidation errors.ErrorCode = 9001
	ErrCodeBusinessRule     errors.ErrorCode = 9002
)

var (
	ErrCustomValidation = errors.NewBusinessError(
		ErrCodeCustomValidation,
		"自定义验证失败",
		400,
	)

	ErrBusinessRule = errors.NewBusinessError(
		ErrCodeBusinessRule,
		"业务规则违反",
		422,
	)
)

func (s *ExampleService) CustomValidation(data string) error {
	if len(data) > 100 {
		return ErrCustomValidation.WithDetailsf("数据长度 %d 超过最大限制 100", len(data))
	}

	if data == "forbidden" {
		return ErrBusinessRule.WithDetails("该数据值不被允许")
	}

	return nil
}

// 示例8: 中间件中的错误处理
func ExampleAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			// 直接中断并返回错误
			errors.AbortWithError(c, errors.ErrUnauthorized)
			return
		}

		// 验证token
		if token != "valid-token" {
			errors.AbortWithError(c, errors.ErrTokenInvalid.WithDetails("令牌格式无效"))
			return
		}

		c.Set("user_id", "123")
		c.Next()
	}
}

// 示例9: 错误包装
func (s *ExampleService) WrapperExample(userID string) error {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return errors.WrapErrorf(err, "获取用户信息失败, ID: %s", userID)
	}

	err = s.updateUserStatus(user)
	if err != nil {
		return errors.WrapError(err, "更新用户状态失败")
	}

	return nil
}

func (s *ExampleService) updateUserStatus(user interface{}) error {
	// 模拟更新操作
	return nil
}
