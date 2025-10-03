package controller

import (
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/internal/service"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/cuiyuanxin/kunpeng/pkg/validator"
	"github.com/gin-gonic/gin"
)

// RoleController 角色控制器
type RoleController struct{}

// GetRoleList 获取角色列表
// @Summary 获取角色列表
// @Description 获取角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page_num query int true "页码" default(1)
// @Param page_size query int true "每页数量" default(10)
// @Param name query string false "角色名称"
// @Param code query string false "角色编码"
// @Param status query int false "状态" Enums(0, 1)
// @Param begin_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} response.Response{data=dto.PageResp} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/roles [get]
func (c *RoleController) GetRoleList(ctx *gin.Context) {
	var req dto.RolePageReq
	if err := validator.BindAndValidateQueryI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetRoleService().GetRoleList(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetAllRoles 获取所有角色
// @Summary 获取所有角色
// @Description 获取所有角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]model.Role} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/roles/all [get]
func (c *RoleController) GetAllRoles(ctx *gin.Context) {
	// 调用服务
	resp, err := service.GetRoleService().GetAllRoles()
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetRoleByID 根据ID获取角色
// @Summary 根据ID获取角色
// @Description 根据ID获取角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=model.Role} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/roles/{id} [get]
func (c *RoleController) GetRoleByID(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetRoleService().GetRoleByID(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.RoleCreateReq true "创建角色请求"
// @Success 200 {object} response.Response{data=uint} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/roles [post]
func (c *RoleController) CreateRole(ctx *gin.Context) {
	var req dto.RoleCreateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	id, err := service.GetRoleService().CreateRole(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, id)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.RoleUpdateReq true "更新角色请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/roles [put]
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	var req dto.RoleUpdateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetRoleService().UpdateRole(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/roles/{id} [delete]
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetRoleService().DeleteRole(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// BatchDeleteRole 批量删除角色
// @Summary 批量删除角色
// @Description 批量删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.IDsReq true "批量删除角色请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/roles/batch [delete]
func (c *RoleController) BatchDeleteRole(ctx *gin.Context) {
	var req dto.IDsReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetRoleService().BatchDeleteRole(req.IDs)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// ChangeRoleStatus 修改角色状态
// @Summary 修改角色状态
// @Description 修改角色状态
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.StatusReq true "修改角色状态请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/roles/status [put]
func (c *RoleController) ChangeRoleStatus(ctx *gin.Context) {
	var req dto.StatusReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetRoleService().ChangeRoleStatus(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// GetRoleMenus 获取角色菜单
// @Summary 获取角色菜单
// @Description 获取角色菜单
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=dto.MenuRoleResp} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/roles/{id}/menus [get]
func (c *RoleController) GetRoleMenus(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务层获取角色菜单
	resp, err := service.GetRoleService().GetRoleMenus(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// UpdateRoleMenus 更新角色菜单
// @Summary 更新角色菜单
// @Description 更新角色菜单
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.RoleMenuReq true "更新角色菜单请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/roles/menus [put]
func (c *RoleController) AssignMenus(ctx *gin.Context) {
	var req dto.RoleMenuReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetRoleService().UpdateRoleMenus(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// GetRoleAPIs 获取角色API
// @Summary 获取角色API
// @Description 获取角色API
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response{data=dto.APIRoleResp} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/roles/{id}/apis [get]
func (c *RoleController) GetRoleAPIs(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUri(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务层获取角色API
	resp, err := service.GetRoleService().GetRoleAPIs(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// UpdateRoleAPIs 更新角色API
// @Summary 更新角色API
// @Description 更新角色API
// @Tags 角色管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.RoleAPIReq true "更新角色API请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/roles/apis [put]
func (c *RoleController) UpdateRoleAPIs(ctx *gin.Context) {
	var req dto.RoleAPIReq
	if err := validator.BindAndValidateJSON(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetRoleService().UpdateRoleAPIs(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}
