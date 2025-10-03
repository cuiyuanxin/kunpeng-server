package controller

import (
	"strconv"

	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/internal/service"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/cuiyuanxin/kunpeng/pkg/validator"
	"github.com/gin-gonic/gin"
)

// LoginLogController 登录日志控制器
type LoginLogController struct{}

// GetLoginLogList 获取登录日志列表
// @Summary 获取登录日志列表
// @Description 获取登录日志列表
// @Tags 登录日志管理
// @Accept json
// @Produce json
// @Param pageNum query int false "页码" default(1)
// @Param pageSize query int false "每页数量" default(10)
// @Param userID query int false "用户ID"
// @Param username query string false "用户名"
// @Param status query int false "状态"
// @Param ip query string false "IP地址"
// @Param beginTime query string false "开始时间"
// @Param endTime query string false "结束时间"
// @Success 200 {object} response.Response{data=dto.PageResp}
// @Router /api/v1/login-logs [get]
// @Security ApiKeyAuth
func (c *LoginLogController) GetLoginLogList(ctx *gin.Context) {
	req := &dto.LoginLogListReq{}

	// 绑定查询参数
	if err := validator.BindAndValidateQueryI18n(ctx, req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 设置默认值
	if req.PageNum <= 0 {
		req.PageNum = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取登录日志列表
	result, err := service.GetLoginLogService().GetLoginLogList(req)
	if err != nil {
		response.FailWithCode(ctx, kperrors.ErrLoginLogGetList)
		return
	}

	response.OkWithData(ctx, result)
}

// GetLoginLogByID 根据ID获取登录日志
// @Summary 根据ID获取登录日志
// @Description 根据ID获取登录日志详情
// @Tags 登录日志管理
// @Accept json
// @Produce json
// @Param id path int true "登录日志ID"
// @Success 200 {object} response.Response{data=dto.LoginLogResp}
// @Router /api/v1/login-logs/{id} [get]
// @Security ApiKeyAuth
func (c *LoginLogController) GetLoginLogByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithCode(ctx, kperrors.ErrParam)
		return
	}

	result, err := service.GetLoginLogService().GetLoginLogByID(uint(id))
	if err != nil {
		response.FailWithCode(ctx, kperrors.ErrLoginLogGetByID)
		return
	}

	response.OkWithData(ctx, result)
}

// DeleteLoginLog 删除登录日志
// @Summary 删除登录日志
// @Description 根据ID删除登录日志
// @Tags 登录日志管理
// @Accept json
// @Produce json
// @Param id path int true "登录日志ID"
// @Success 200 {object} response.Response
// @Router /api/v1/login-logs/{id} [delete]
// @Security ApiKeyAuth
func (c *LoginLogController) DeleteLoginLog(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithCode(ctx, kperrors.ErrParam)
		return
	}

	err = service.GetLoginLogService().DeleteLoginLog(uint(id))
	if err != nil {
		response.FailWithCode(ctx, kperrors.ErrLoginLogDelete)
		return
	}

	response.OkWithMessage(ctx, "删除成功")
}

// BatchDeleteLoginLog 批量删除登录日志
// @Summary 批量删除登录日志
// @Description 根据ID列表批量删除登录日志
// @Tags 登录日志管理
// @Accept json
// @Produce json
// @Param req body dto.IDsReq true "ID列表"
// @Success 200 {object} response.Response
// @Router /api/v1/login-logs/batch [delete]
// @Security ApiKeyAuth
func (c *LoginLogController) BatchDeleteLoginLog(ctx *gin.Context) {
	req := &dto.IDsReq{}
	if err := validator.BindAndValidateJSONI18n(ctx, req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	err := service.GetLoginLogService().BatchDeleteLoginLog(req)
	if err != nil {
		response.FailWithCode(ctx, kperrors.ErrLoginLogBatchDelete)
		return
	}

	response.OkWithMessage(ctx, "删除成功")
}

// CleanOldLogs 清理旧日志
// @Summary 清理旧日志
// @Description 清理指定天数之前的登录日志
// @Tags 登录日志管理
// @Accept json
// @Produce json
// @Param days query int true "保留天数" default(30)
// @Success 200 {object} response.Response
// @Router /api/v1/login-logs/clean [post]
// @Security ApiKeyAuth
func (c *LoginLogController) CleanOldLogs(ctx *gin.Context) {
	daysStr := ctx.Query("days")
	if daysStr == "" {
		daysStr = "30" // 默认保留30天
	}

	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		response.FailWithCode(ctx, kperrors.ErrParam)
		return
	}

	err = service.GetLoginLogService().CleanOldLogs(days)
	if err != nil {
		response.FailWithCode(ctx, kperrors.ErrLoginLogCleanOld)
		return
	}

	response.OkWithMessage(ctx, "清理成功")
}
