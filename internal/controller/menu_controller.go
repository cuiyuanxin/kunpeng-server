package controller

import (
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/internal/service"
	"github.com/cuiyuanxin/kunpeng/pkg/jwt"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/cuiyuanxin/kunpeng/pkg/validator"
	"github.com/gin-gonic/gin"
)

// MenuController 菜单控制器
type MenuController struct{}

// GetMenuList 获取菜单列表
// @Summary 获取菜单列表
// @Description 获取菜单列表
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "菜单名称"
// @Param status query int false "状态" Enums(0, 1)
// @Success 200 {object} response.Response{data=[]model.Menu} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/menus [get]
func (c *MenuController) GetMenuList(ctx *gin.Context) {
	var req dto.MenuQueryReq
	if err := validator.BindAndValidateQueryI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetMenuService().GetMenuList(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetMenuTree 获取菜单树
// @Summary 获取菜单树
// @Description 获取菜单树
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "菜单名称"
// @Param status query int false "状态" Enums(0, 1)
// @Success 200 {object} response.Response{data=[]dto.MenuTreeResp} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/menus/tree [get]
func (c *MenuController) GetMenuTree(ctx *gin.Context) {
	var req dto.MenuQueryReq
	if err := validator.BindAndValidateQuery(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetMenuService().GetMenuTree(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetUserMenuTree 获取用户菜单树
// @Summary 获取用户菜单树
// @Description 获取用户菜单树
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]dto.MenuTreeResp} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/menus/user [get]
func (c *MenuController) GetUserMenuTree(ctx *gin.Context) {
	// 获取当前用户ID
	userID := jwt.GetUserID(ctx)

	// 调用服务
	resp, err := service.GetMenuService().GetUserMenuTree(userID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetMenuByID 根据ID获取菜单
// @Summary 根据ID获取菜单
// @Description 根据ID获取菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response{data=model.Menu} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/menus/{id} [get]
func (c *MenuController) GetMenuByID(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetMenuService().GetMenuByID(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// CreateMenu 创建菜单
// @Summary 创建菜单
// @Description 创建菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.MenuCreateReq true "创建菜单请求"
// @Success 200 {object} response.Response{data=uint} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/menus [post]
func (c *MenuController) CreateMenu(ctx *gin.Context) {
	var req dto.MenuCreateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	id, err := service.GetMenuService().CreateMenu(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, id)
}

// UpdateMenu 更新菜单
// @Summary 更新菜单
// @Description 更新菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.MenuUpdateReq true "更新菜单请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/menus [put]
func (c *MenuController) UpdateMenu(ctx *gin.Context) {
	var req dto.MenuUpdateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetMenuService().UpdateMenu(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// DeleteMenu 删除菜单
// @Summary 删除菜单
// @Description 删除菜单
// @Tags 菜单管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "菜单ID"
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/menus/{id} [delete]
func (c *MenuController) DeleteMenu(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetMenuService().DeleteMenu(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}
