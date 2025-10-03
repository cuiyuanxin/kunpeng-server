package controller

import (
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/internal/service"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/cuiyuanxin/kunpeng/pkg/validator"
	"github.com/gin-gonic/gin"
)

// DeptController 部门控制器
type DeptController struct{}

// GetDeptList 获取部门列表
// @Summary 获取部门列表
// @Description 获取部门列表
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "部门名称"
// @Param status query int false "状态" Enums(0, 1)
// @Success 200 {object} response.Response{data=[]model.Dept} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/depts [get]
func (c *DeptController) GetDeptList(ctx *gin.Context) {
	var req dto.DeptQueryReq
	if err := validator.BindAndValidateQueryI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetDeptService().GetDeptList(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetDeptTree 获取部门树
// @Summary 获取部门树
// @Description 获取部门树
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param name query string false "部门名称"
// @Param status query int false "状态" Enums(0, 1)
// @Success 200 {object} response.Response{data=[]dto.DeptTreeResp} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/depts/tree [get]
func (c *DeptController) GetDeptTree(ctx *gin.Context) {
	var req dto.DeptQueryReq
	if err := validator.BindAndValidateQuery(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetDeptService().GetDeptTree(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetDeptByID 根据ID获取部门
// @Summary 根据ID获取部门
// @Description 根据ID获取部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "部门ID"
// @Success 200 {object} response.Response{data=model.Dept} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/depts/{id} [get]
func (c *DeptController) GetDeptByID(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetDeptService().GetDeptByID(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// CreateDept 创建部门
// @Summary 创建部门
// @Description 创建部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.DeptCreateReq true "创建部门请求"
// @Success 200 {object} response.Response{data=uint} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/depts [post]
func (c *DeptController) CreateDept(ctx *gin.Context) {
	var req dto.DeptCreateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	id, err := service.GetDeptService().CreateDept(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, id)
}

// UpdateDept 更新部门
// @Summary 更新部门
// @Description 更新部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.DeptUpdateReq true "更新部门请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/depts [put]
func (c *DeptController) UpdateDept(ctx *gin.Context) {
	var req dto.DeptUpdateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetDeptService().UpdateDept(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// DeleteDept 删除部门
// @Summary 删除部门
// @Description 删除部门
// @Tags 部门管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "部门ID"
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/depts/{id} [delete]
func (c *DeptController) DeleteDept(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetDeptService().DeleteDept(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}
