package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/cuiyuanxin/kunpeng/internal/auth"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/response"
	"github.com/cuiyuanxin/kunpeng/internal/service"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService *service.UserService
	jwtManager  *auth.JWTManager
	validator   *validator.Validate
}

// NewUserHandler 创建用户处理器
func NewUserHandler(userService *service.UserService, jwtManager *auth.JWTManager) *UserHandler {
	return &UserHandler{
		userService: userService,
		jwtManager:  jwtManager,
		validator:   validator.New(),
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
func (h *UserHandler) Register(c *gin.Context) {
	var req model.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		response.BadRequestWithData(c, "Validation failed", err.Error())
		return
	}

	// 创建用户
	user, err := h.userService.Create(&req)
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
func (h *UserHandler) Login(c *gin.Context) {
	var req model.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		response.BadRequestWithData(c, "Validation failed", err.Error())
		return
	}

	// 获取用户
	user, err := h.userService.GetByUsername(req.Username)
	if err != nil {
		response.Unauthorized(c, "Invalid username or password")
		return
	}

	// 验证密码
	if !h.userService.VerifyPassword(user, req.Password) {
		response.Unauthorized(c, "Invalid username or password")
		return
	}

	// 检查用户状态
	if user.Status != 1 {
		response.Forbidden(c, "User account is disabled")
		return
	}

	// 生成JWT令牌
	token, err := h.jwtManager.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		response.ServerError(c, "Failed to generate token")
		return
	}

	// 更新最后登录时间
	h.userService.UpdateLastLogin(user.ID)

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
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	user, err := h.userService.GetByID(userID.(uint))
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
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req model.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		response.BadRequestWithData(c, "Validation failed", err.Error())
		return
	}

	// 普通用户不能修改角色和状态
	role, _ := c.Get("role")
	if role != "admin" {
		req.Role = ""
		req.Status = nil
	}

	// 更新用户
	user, err := h.userService.Update(userID.(uint), &req)
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
// @Param request body model.UserChangePasswordRequest true "密码信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/user/password [put]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req model.UserChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		response.BadRequestWithData(c, "Validation failed", err.Error())
		return
	}

	// 修改密码
	if err := h.userService.ChangePassword(userID.(uint), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Password changed successfully", nil)
}

// GetUsers 获取用户列表（管理员）
// @Summary 获取用户列表
// @Description 获取用户列表（仅管理员）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param keyword query string false "搜索关键词"
// @Param role query string false "角色筛选"
// @Param status query int false "状态筛选"
// @Success 200 {object} response.PageResponse{data=[]model.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/admin/users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	var req model.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "Invalid query parameters")
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		response.BadRequestWithData(c, "Validation failed", err.Error())
		return
	}

	// 获取用户列表
	users, total, err := h.userService.List(&req)
	if err != nil {
		response.ServerError(c, "Failed to get users")
		return
	}

	// 转换为响应格式
	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	response.SuccessPage(c, userResponses, total, req.Page, req.PageSize)
}

// GetUser 获取指定用户（管理员）
// @Summary 获取指定用户
// @Description 获取指定用户信息（仅管理员）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/admin/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}

	response.Success(c, user.ToResponse())
}

// UpdateUser 更新指定用户（管理员）
// @Summary 更新指定用户
// @Description 更新指定用户信息（仅管理员）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Param request body model.UserUpdateRequest true "更新信息"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req model.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request format")
		return
	}

	// 验证请求参数
	if err := h.validator.Struct(&req); err != nil {
		response.BadRequestWithData(c, "Validation failed", err.Error())
		return
	}

	// 更新用户
	user, err := h.userService.Update(uint(id), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "User updated successfully", user.ToResponse())
}

// DeleteUser 删除指定用户（管理员）
// @Summary 删除指定用户
// @Description 删除指定用户（仅管理员）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	// 删除用户
	if err := h.userService.Delete(uint(id)); err != nil {
		response.ServerError(c, "Failed to delete user")
		return
	}

	response.SuccessWithMessage(c, "User deleted successfully", nil)
}
