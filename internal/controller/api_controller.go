package controller

import (
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/internal/service"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/cuiyuanxin/kunpeng/pkg/validator"
	"github.com/gin-gonic/gin"
)

// APIController API控制器
type APIController struct{}

// GetAPIList 获取API列表
// @Summary 获取API列表
// @Description 获取API列表
// @Tags API管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page_num query int true "页码" default(1)
// @Param page_size query int true "每页数量" default(10)
// @Param group query string false "API分组"
// @Param name query string false "API名称"
// @Param method query string false "请求方法"
// @Param path query string false "请求路径"
// @Param status query int false "状态" Enums(0, 1)
// @Param begin_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} response.Response{data=dto.PageResp} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/apis [get]
func (c *APIController) GetAPIList(ctx *gin.Context) {
	var req dto.APIPageReq
	if err := validator.BindAndValidateQueryI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetAPIService().GetAPIList(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetAllAPIs 获取所有API
// @Summary 获取所有API
// @Description 获取所有API
// @Tags API管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]model.API} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/apis/all [get]
func (c *APIController) GetAllAPIs(ctx *gin.Context) {
	// 调用服务
	resp, err := service.GetAPIService().GetAllAPIs()
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetAPIByID 根据ID获取API
// @Summary 根据ID获取API
// @Description 根据ID获取API
// @Tags API管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "API ID"
// @Success 200 {object} response.Response{data=model.API} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/apis/{id} [get]
func (c *APIController) GetAPIByID(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetAPIService().GetAPIByID(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// CreateAPI 创建API
// @Summary 创建API
// @Description 创建API
// @Tags API管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.APICreateReq true "创建API请求"
// @Success 200 {object} response.Response{data=uint} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/apis [post]
func (c *APIController) CreateAPI(ctx *gin.Context) {
	var req dto.APICreateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	id, err := service.GetAPIService().CreateAPI(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, id)
}

// UpdateAPI 更新API
// @Summary 更新API
// @Description 更新API
// @Tags API管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.APIUpdateReq true "更新API请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/apis [put]
func (c *APIController) UpdateAPI(ctx *gin.Context) {
	var req dto.APIUpdateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetAPIService().UpdateAPI(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// DeleteAPI 删除API
// @Summary 删除API
// @Description 删除API
// @Tags API管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "API ID"
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/apis/{id} [delete]
func (c *APIController) DeleteAPI(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetAPIService().DeleteAPI(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// BatchDeleteAPI 批量删除API
// @Summary 批量删除API
// @Description 批量删除API
// @Tags API管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.IDsReq true "批量删除API请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/apis/batch [delete]
func (c *APIController) BatchDeleteAPI(ctx *gin.Context) {
	var req dto.IDsReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetAPIService().BatchDeleteAPI(req.IDs)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}
