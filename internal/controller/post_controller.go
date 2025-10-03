package controller

import (
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/internal/service"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/cuiyuanxin/kunpeng/pkg/validator"
	"github.com/gin-gonic/gin"
)

// PostController 岗位控制器
type PostController struct{}

// GetPostList 获取岗位列表
// @Summary 获取岗位列表
// @Description 获取岗位列表
// @Tags 岗位管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page_num query int true "页码" default(1)
// @Param page_size query int true "每页数量" default(10)
// @Param name query string false "岗位名称"
// @Param code query string false "岗位编码"
// @Param status query int false "状态" Enums(0, 1)
// @Param begin_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} response.Response{data=dto.PageResp} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/posts [get]
func (c *PostController) GetPostList(ctx *gin.Context) {
	var req dto.PostPageReq
	if err := validator.BindAndValidateQueryI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetPostService().GetPostList(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetAllPosts 获取所有岗位
// @Summary 获取所有岗位
// @Description 获取所有岗位
// @Tags 岗位管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]model.Post} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/posts/all [get]
func (c *PostController) GetAllPosts(ctx *gin.Context) {
	// 调用服务
	resp, err := service.GetPostService().GetAllPosts()
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetPostByID 根据ID获取岗位
// @Summary 根据ID获取岗位
// @Description 根据ID获取岗位
// @Tags 岗位管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "岗位ID"
// @Success 200 {object} response.Response{data=model.Post} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/posts/{id} [get]
func (c *PostController) GetPostByID(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetPostService().GetPostByID(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// CreatePost 创建岗位
// @Summary 创建岗位
// @Description 创建岗位
// @Tags 岗位管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.PostCreateReq true "创建岗位请求"
// @Success 200 {object} response.Response{data=uint} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/posts [post]
func (c *PostController) CreatePost(ctx *gin.Context) {
	var req dto.PostCreateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	id, err := service.GetPostService().CreatePost(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, id)
}

// UpdatePost 更新岗位
// @Summary 更新岗位
// @Description 更新岗位
// @Tags 岗位管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.PostUpdateReq true "更新岗位请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/posts [put]
func (c *PostController) UpdatePost(ctx *gin.Context) {
	var req dto.PostUpdateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetPostService().UpdatePost(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// DeletePost 删除岗位
// @Summary 删除岗位
// @Description 删除岗位
// @Tags 岗位管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "岗位ID"
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/posts/{id} [delete]
func (c *PostController) DeletePost(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetPostService().DeletePost(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// BatchDeletePost 批量删除岗位
// @Summary 批量删除岗位
// @Description 批量删除岗位
// @Tags 岗位管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.IDsReq true "批量删除岗位请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/posts/batch [delete]
func (c *PostController) BatchDeletePost(ctx *gin.Context) {
	var req dto.IDsReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetPostService().BatchDeletePost(req.IDs)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}
