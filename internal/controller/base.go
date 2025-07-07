package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/cuiyuanxin/kunpeng/pkg/app"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
)

type BaseController struct {
	formValidator *app.FormValidator
}

func NewBaseController() *BaseController {
	return &BaseController{
		formValidator: app.NewFormValidator(),
	}
}

// ShouldBind 仅绑定参数，不进行验证
func (bc *BaseController) ShouldBind(c *gin.Context, req interface{}) error {
	return bc.formValidator.ShouldBind(c, req)
}

// ShouldBindJSON 仅绑定JSON参数，不进行验证
func (bc *BaseController) ShouldBindJSON(c *gin.Context, req interface{}) error {
	return bc.formValidator.ShouldBindJSON(c, req)
}

// ShouldBindQuery 仅绑定查询参数，不进行验证
func (bc *BaseController) ShouldBindQuery(c *gin.Context, req interface{}) error {
	return bc.formValidator.ShouldBindQuery(c, req)
}

// Validate 仅验证结构体，不进行绑定
func (bc *BaseController) Validate(req interface{}) app.ValidErrors {
	return bc.formValidator.Validate(req)
}

// BindAndValidate 绑定并验证请求参数
func (bc *BaseController) BindAndValidate(c *gin.Context, req interface{}) app.ValidErrors {
	return bc.formValidator.BindAndValidate(c, req)
}

// BindJSONAndValidate 绑定JSON并验证请求参数
func (bc *BaseController) BindJSONAndValidate(c *gin.Context, req interface{}) app.ValidErrors {
	return bc.formValidator.BindJSONAndValidate(c, req)
}

// BindQueryAndValidate 绑定查询参数并验证
func (bc *BaseController) BindQueryAndValidate(c *gin.Context, req interface{}) app.ValidErrors {
	return bc.formValidator.BindQueryAndValidate(c, req)
}

// MustBindAndValidate 绑定并验证，失败时自动返回错误响应
func (bc *BaseController) MustBindAndValidate(c *gin.Context, req interface{}) bool {
	return bc.formValidator.MustBindAndValidate(c, req)
}

// MustBindJSONAndValidate 绑定JSON并验证，失败时自动返回错误响应
func (bc *BaseController) MustBindJSONAndValidate(c *gin.Context, req interface{}) bool {
	return bc.formValidator.MustBindJSONAndValidate(c, req)
}

// MustBindQueryAndValidate 绑定查询参数并验证，失败时自动返回错误响应
func (bc *BaseController) MustBindQueryAndValidate(c *gin.Context, req interface{}) bool {
	return bc.formValidator.MustBindQueryAndValidate(c, req)
}

// HandleValidationErrors 处理验证错误
func (bc *BaseController) HandleValidationErrors(c *gin.Context, errs app.ValidErrors) {
	bc.formValidator.HandleValidationErrors(c, errs)
}

// SetLanguage 设置验证器语言
func (bc *BaseController) SetLanguage(lang string) error {
	return bc.formValidator.SetLanguage(lang)
}

// GetFormValidator 获取表单验证器实例
func (bc *BaseController) GetFormValidator() *app.FormValidator {
	return bc.formValidator
}

// GetUserID 从上下文中获取用户ID
func (bc *BaseController) GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	return userID.(uint), true
}

// GetUserRole 从上下文中获取用户角色
func (bc *BaseController) GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	return role.(string), true
}

// RequireAuth 检查用户认证
func (bc *BaseController) RequireAuth(c *gin.Context) (uint, bool) {
	userID, exists := bc.GetUserID(c)
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return 0, false
	}
	return userID, true
}

// RequireAdmin 检查管理员权限
func (bc *BaseController) RequireAdmin(c *gin.Context) bool {
	role, exists := bc.GetUserRole(c)
	if !exists || role != "admin" {
		response.Forbidden(c, "Admin permission required")
		return false
	}
	return true
}

// HandleValidationError 处理验证错误
func (bc *BaseController) HandleValidationError(c *gin.Context, err error) {
	response.BadRequestWithData(c, "Validation failed", err.Error())
}

// HandleBindError 处理绑定错误
func (bc *BaseController) HandleBindError(c *gin.Context, err error) {
	response.BadRequest(c, "Invalid request format")
}