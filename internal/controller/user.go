package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/service"
	"github.com/cuiyuanxin/kunpeng/pkg/auth"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
)

// UserController 用户控制器
type UserController struct {
	*BaseController
	userService *service.UserService
	jwtManager  *auth.JWTManager
}

// NewUserController 创建用户控制器
func NewUserController(userService *service.UserService, jwtManager *auth.JWTManager) *UserController {
	return &UserController{
		BaseController: NewBaseController(),
		userService:    userService,
		jwtManager:     jwtManager,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body model.UserCreateRequest true "用户注册信息"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/register [post]
func (uc *UserController) Register(c *gin.Context) {
	var req model.UserCreateRequest
	if err := uc.BindAndValidate(c, &req); err != nil {
		uc.HandleValidationError(c, err)
		return
	}

	// 创建用户
	user, err := uc.userService.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "User registered successfully", user.ToResponse())
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取JWT令牌
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body model.UserLoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=model.LoginResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var req model.UserLoginRequest
	if err := uc.BindAndValidate(c, &req); err != nil {
		uc.HandleValidationError(c, err)
		return
	}

	// 获取用户
	user, err := uc.userService.GetByUsername(req.Username)
	if err != nil {
		response.Unauthorized(c, "Invalid username or password")
		return
	}

	// 验证密码
	if !uc.userService.VerifyPassword(user, req.Password) {
		response.Unauthorized(c, "Invalid username or password")
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		response.Forbidden(c, "User account is disabled")
		return
	}

	// 生成JWT令牌
	token, err := uc.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.ServerError(c, "Failed to generate token")
		return
	}

	// 更新最后登录时间
	uc.userService.UpdateLastLogin(user.ID)

	// 返回登录响应
	loginResp := model.LoginResponse{
		Token: token,
		User:  user.ToResponse(),
	}

	response.SuccessWithMessage(c, "Login successful", loginResp)
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户的资料信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/user/profile [get]
func (uc *UserController) GetProfile(c *gin.Context) {
	userID, ok := uc.RequireAuth(c)
	if !ok {
		return
	}

	user, err := uc.userService.GetByID(userID)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, user.ToResponse())
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前登录用户的资料信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body model.UserUpdateRequest true "更新信息"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/user/profile [put]
func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID, ok := uc.RequireAuth(c)
	if !ok {
		return
	}

	var req model.UserUpdateRequest
	if err := uc.BindAndValidate(c, &req); err != nil {
		uc.HandleValidationError(c, err)
		return
	}

	// 普通用户不能修改角色和状态
	role, _ := uc.GetUserRole(c)
	if role != "admin" {
		req.Role = ""
		req.Status = nil
	}

	// 更新用户
	user, err := uc.userService.Update(userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Profile updated successfully", user.ToResponse())
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户的密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body model.ChangePasswordRequest true "修改密码信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/user/change-password [post]
func (uc *UserController) ChangePassword(c *gin.Context) {
	userID, ok := uc.RequireAuth(c)
	if !ok {
		return
	}

	var req model.ChangePasswordRequest
	if err := uc.BindAndValidate(c, &req); err != nil {
		uc.HandleValidationError(c, err)
		return
	}

	// 获取用户
	user, err := uc.userService.GetByID(userID)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	// 验证旧密码
	if !uc.userService.VerifyPassword(user, req.OldPassword) {
		response.BadRequest(c, "Invalid old password")
		return
	}

	// 更新密码
	if err := uc.userService.ChangePassword(userID, &req); err != nil {
		response.ServerError(c, "Failed to change password")
		return
	}

	response.SuccessWithMessage(c, "Password changed successfully", nil)
}