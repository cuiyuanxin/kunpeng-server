package service

import (
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// OperationLogService 操作日志服务接口
type OperationLogService interface {
	// GetOperationLogList 获取操作日志列表
	GetOperationLogList(req *dto.OperationLogListReq) (*dto.PageResp, error)

	// GetOperationLogByID 根据ID获取操作日志
	GetOperationLogByID(id uint) (*dto.OperationLogResp, error)

	// DeleteOperationLog 删除操作日志
	DeleteOperationLog(id uint) error

	// BatchDeleteOperationLog 批量删除操作日志
	BatchDeleteOperationLog(req *dto.IDsReq) error

	// CleanOldLogs 清理旧日志
	CleanOldLogs(days int) error

	// RecordOperationLog 记录操作日志
	RecordOperationLog(userID uint, username, module, action, method, path, ip, userAgent, request, response string, status int8, errorMessage string, duration int) error
}
