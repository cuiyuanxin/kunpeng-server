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

// OperationLogController 操作日志控制器
type OperationLogController struct{}

// GetOperationLogList 获取操作日志列表
func (c *OperationLogController) GetOperationLogList(ctx *gin.Context) {
	var req dto.OperationLogListReq
	if err := validator.BindAndValidateQueryI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	operationLogService := service.GetOperationLogService()
	result, err := operationLogService.GetOperationLogList(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, result)
}

// GetOperationLogByID 根据ID获取操作日志
func (c *OperationLogController) GetOperationLogByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.FailWithCode(ctx, kperrors.ErrOperationLogInvalidID)
		return
	}

	operationLogService := service.GetOperationLogService()
	result, err := operationLogService.GetOperationLogByID(uint(id))
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, result)
}

// DeleteOperationLog 删除操作日志
func (c *OperationLogController) DeleteOperationLog(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		response.FailWithCode(ctx, kperrors.ErrOperationLogInvalidID)
		return
	}

	operationLogService := service.GetOperationLogService()
	err = operationLogService.DeleteOperationLog(uint(id))
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// BatchDeleteOperationLog 批量删除操作日志
func (c *OperationLogController) BatchDeleteOperationLog(ctx *gin.Context) {
	var req dto.IDsReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	operationLogService := service.GetOperationLogService()
	err := operationLogService.BatchDeleteOperationLog(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// CleanOldLogs 清理旧的操作日志
func (c *OperationLogController) CleanOldLogs(ctx *gin.Context) {
	daysStr := ctx.Query("days")
	days := 30 // 默认清理30天前的日志
	if daysStr != "" {
		if d, err := strconv.Atoi(daysStr); err == nil && d > 0 {
			days = d
		}
	}

	operationLogService := service.GetOperationLogService()
	err := operationLogService.CleanOldLogs(days)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}
