package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/service"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
)

// AdminController 管理员控制器
type AdminController struct {
	*BaseController
	userService       *service.UserService
	roleService       *service.RoleService
	permissionService *service.PermissionService
	departmentService *service.DepartmentService
	fileService       *service.FileService
	statService       *service.StatService
}

// NewAdminController 创建管理员控制器
func NewAdminController(
	userService *service.UserService,
	roleService *service.RoleService,
	permissionService *service.PermissionService,
	departmentService *service.DepartmentService,
	fileService *service.FileService,
	statService *service.StatService,
) *AdminController {
	return &AdminController{
		BaseController:    NewBaseController(),
		userService:       userService,
		roleService:       roleService,
		permissionService: permissionService,
		departmentService: departmentService,
		fileService:       fileService,
		statService:       statService,
	}
}

// ===== 用户管理 =====

// GetUsers 获取用户列表
// @Summary 获取用户列表
// @Description 获取系统中所有用户的列表
// @Tags 管理员-用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param limit query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} response.Response{data=response.PageResponse}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/admin/users [get]
func (ac *AdminController) GetUsers(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	// 解析查询参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	// search := c.Query("search")

	// 获取用户列表
	req := &model.UserListRequest{
		Page:     page,
		PageSize: limit,
	}
	users, total, err := ac.userService.List(req)
	if err != nil {
		response.ServerError(c, "Failed to get users")
		return
	}

	// 转换为响应格式
	userResponses := make([]model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	response.Success(c, response.PageResponse{
		Data:  userResponses,
		Total:   total,
		Page:    page,
		Size:    limit,
	})
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 管理员创建新用户
// @Tags 管理员-用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body model.UserCreateRequest true "用户信息"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/admin/users [post]
func (ac *AdminController) CreateUser(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	var req model.UserCreateRequest
	if err := ac.BindAndValidate(c, &req); err != nil {
		ac.HandleValidationError(c, err)
		return
	}

	// 创建用户
	user, err := ac.userService.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "User created successfully", user.ToResponse())
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 管理员更新用户信息
// @Tags 管理员-用户管理
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
func (ac *AdminController) UpdateUser(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	// 获取用户ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	var req model.UserUpdateRequest
	if err := ac.BindAndValidate(c, &req); err != nil {
		ac.HandleValidationError(c, err)
		return
	}

	// 更新用户
	user, err := ac.userService.Update(uint(userID), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "User updated successfully", user.ToResponse())
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 管理员删除用户
// @Tags 管理员-用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /api/v1/admin/users/{id} [delete]
func (ac *AdminController) DeleteUser(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	// 获取用户ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	// 删除用户
	if err := ac.userService.Delete(uint(userID)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "User deleted successfully", nil)
}

// ===== 角色管理 =====

// GetRoles 获取角色列表
func (ac *AdminController) GetRoles(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	roles, err := ac.roleService.GetAll()
	if err != nil {
		response.ServerError(c, "Failed to get roles")
		return
	}

	response.Success(c, roles)
}

// CreateRole 创建角色
func (ac *AdminController) CreateRole(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	var req model.RoleCreateRequest
	if err := ac.BindAndValidate(c, &req); err != nil {
		ac.HandleValidationError(c, err)
		return
	}

	role, err := ac.roleService.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Role created successfully", role)
}

// ===== 权限管理 =====

// GetPermissions 获取权限列表
func (ac *AdminController) GetPermissions(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	permissions, err := ac.permissionService.GetAll()
	if err != nil {
		response.ServerError(c, "Failed to get permissions")
		return
	}

	response.Success(c, permissions)
}

// ===== 部门管理 =====

// GetDepartments 获取部门列表
func (ac *AdminController) GetDepartments(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	departments, err := ac.departmentService.GetAll()
	if err != nil {
		response.ServerError(c, "Failed to get departments")
		return
	}

	response.Success(c, departments)
}

// CreateDepartment 创建部门
func (ac *AdminController) CreateDepartment(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	var req model.DepartmentCreateRequest
	if err := ac.BindAndValidate(c, &req); err != nil {
		ac.HandleValidationError(c, err)
		return
	}

	department, err := ac.departmentService.Create(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "Department created successfully", department)
}

// ===== 文件管理 =====

// UploadFile 上传文件
func (ac *AdminController) UploadFile(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "No file uploaded")
		return
	}

	fileInfo, err := ac.fileService.Upload(file)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessWithMessage(c, "File uploaded successfully", fileInfo)
}

// GetFiles 获取文件列表
func (ac *AdminController) GetFiles(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	files, total, err := ac.fileService.GetList(page, limit)
	if err != nil {
		response.ServerError(c, "Failed to get files")
		return
	}

	response.Success(c, response.PageResponse{
		Data:  files,
		Total: total,
		Page:  page,
		Size:  limit,
	})
}

// ===== 统计分析 =====

// GetDashboardStats 获取仪表板统计
func (ac *AdminController) GetDashboardStats(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	stats, err := ac.statService.GetDashboardStats()
	if err != nil {
		response.ServerError(c, "Failed to get dashboard stats")
		return
	}

	response.Success(c, stats)
}

// GetUserStats 获取用户统计
func (ac *AdminController) GetUserStats(c *gin.Context) {
	if !ac.RequireAdmin(c) {
		return
	}

	stats, err := ac.statService.GetUserStats()
	if err != nil {
		response.ServerError(c, "Failed to get user stats")
		return
	}

	response.Success(c, stats)
}