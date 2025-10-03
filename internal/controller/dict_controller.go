package controller

import (
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/internal/service"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/cuiyuanxin/kunpeng/pkg/validator"
	"github.com/gin-gonic/gin"
)

// DictController 字典控制器
type DictController struct{}

// GetDictTypeList 获取字典类型列表
// @Summary 获取字典类型列表
// @Description 获取字典类型列表
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page_num query int true "页码" default(1)
// @Param page_size query int true "每页数量" default(10)
// @Param name query string false "字典名称"
// @Param type query string false "字典类型"
// @Param status query int false "状态" Enums(0, 1)
// @Param begin_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} response.Response{data=dto.PageResp} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/types [get]
func (c *DictController) GetDictTypeList(ctx *gin.Context) {
	var req dto.DictTypePageReq
	if err := validator.BindAndValidateQueryI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetDictService().GetDictTypeList(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetAllDictTypes 获取所有字典类型
// @Summary 获取所有字典类型
// @Description 获取所有字典类型
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]model.DictType} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/types/all [get]
func (c *DictController) GetAllDictTypes(ctx *gin.Context) {
	// 调用服务
	resp, err := service.GetDictService().GetAllDictTypes()
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetDictTypeByID 根据ID获取字典类型
// @Summary 根据ID获取字典类型
// @Description 根据ID获取字典类型
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "字典类型ID"
// @Success 200 {object} response.Response{data=model.DictType} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/types/{id} [get]
func (c *DictController) GetDictTypeByID(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetDictService().GetDictTypeByID(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// CreateDictType 创建字典类型
// @Summary 创建字典类型
// @Description 创建字典类型
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.DictTypeCreateReq true "创建字典类型请求"
// @Success 200 {object} response.Response{data=uint} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/types [post]
func (c *DictController) CreateDictType(ctx *gin.Context) {
	var req dto.DictTypeCreateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	id, err := service.GetDictService().CreateDictType(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, id)
}

// UpdateDictType 更新字典类型
// @Summary 更新字典类型
// @Description 更新字典类型
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.DictTypeUpdateReq true "更新字典类型请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/types [put]
func (c *DictController) UpdateDictType(ctx *gin.Context) {
	var req dto.DictTypeUpdateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetDictService().UpdateDictType(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// DeleteDictType 删除字典类型
// @Summary 删除字典类型
// @Description 删除字典类型
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "字典类型ID"
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/types/{id} [delete]
func (c *DictController) DeleteDictType(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetDictService().DeleteDictType(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// BatchDeleteDictType 批量删除字典类型
// @Summary 批量删除字典类型
// @Description 批量删除字典类型
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.IDsReq true "批量删除字典类型请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/types/batch [delete]
func (c *DictController) BatchDeleteDictType(ctx *gin.Context) {
	var req dto.IDsReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetDictService().BatchDeleteDictType(req.IDs)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// GetDictDataList 获取字典数据列表
// @Summary 获取字典数据列表
// @Description 获取字典数据列表
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page_num query int true "页码" default(1)
// @Param page_size query int true "每页数量" default(10)
// @Param dict_type query string true "字典类型"
// @Param label query string false "字典标签"
// @Param status query int false "状态" Enums(0, 1)
// @Param begin_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} response.Response{data=dto.PageResp} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/data [get]
func (c *DictController) GetDictDataList(ctx *gin.Context) {
	var req dto.DictDataPageReq
	if err := validator.BindAndValidateQueryI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetDictService().GetDictDataList(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetDictDataByType 根据字典类型获取字典数据
// @Summary 根据字典类型获取字典数据
// @Description 根据字典类型获取字典数据
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param dict_type path string true "字典类型"
// @Success 200 {object} response.Response{data=[]model.DictData} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/data/type/{dict_type} [get]
func (c *DictController) GetDictDataByType(ctx *gin.Context) {
	var req dto.DictTypeReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetDictService().GetDictDataByType(req.DictType)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetDictDataByID 根据ID获取字典数据
// @Summary 根据ID获取字典数据
// @Description 根据ID获取字典数据
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "字典数据ID"
// @Success 200 {object} response.Response{data=model.DictData} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/data/{id} [get]
func (c *DictController) GetDictDataByID(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetDictService().GetDictDataByID(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// CreateDictData 创建字典数据
// @Summary 创建字典数据
// @Description 创建字典数据
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.DictDataCreateReq true "创建字典数据请求"
// @Success 200 {object} response.Response{data=uint} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/data [post]
func (c *DictController) CreateDictData(ctx *gin.Context) {
	var req dto.DictDataCreateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	id, err := service.GetDictService().CreateDictData(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, id)
}

// UpdateDictData 更新字典数据
// @Summary 更新字典数据
// @Description 更新字典数据
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.DictDataUpdateReq true "更新字典数据请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/data [put]
func (c *DictController) UpdateDictData(ctx *gin.Context) {
	var req dto.DictDataUpdateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetDictService().UpdateDictData(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// DeleteDictData 删除字典数据
// @Summary 删除字典数据
// @Description 删除字典数据
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "字典数据ID"
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/data/{id} [delete]
func (c *DictController) DeleteDictData(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务层删除字典数据
	err := service.GetDictService().DeleteDictData(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// BatchDeleteDictData 批量删除字典数据
// @Summary 批量删除字典数据
// @Description 批量删除字典数据
// @Tags 字典管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.IDsReq true "批量删除字典数据请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/dict/data/batch [delete]
func (c *DictController) BatchDeleteDictData(ctx *gin.Context) {
	var req dto.IDsReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetDictService().BatchDeleteDictData(req.IDs)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}
