package impl

import (
	"math"
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/internal/repository"
)

// OperationLogServiceImpl 操作日志服务实现
type OperationLogServiceImpl struct{}

// GetOperationLogList 获取操作日志列表
func (s *OperationLogServiceImpl) GetOperationLogList(req *dto.OperationLogListReq) (*dto.PageResp, error) {
	logs, total, err := repository.GetOperationLogRepository().FindList(req)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	logList := dto.ToOperationLogRespList(logs)

	// 计算总页数
	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &dto.PageResp{
		List:       logList,
		Total:      total,
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetOperationLogByID 根据ID获取操作日志
func (s *OperationLogServiceImpl) GetOperationLogByID(id uint) (*dto.OperationLogResp, error) {
	log, err := repository.GetOperationLogRepository().FindByID(id)
	if err != nil {
		return nil, err
	}
	return dto.ToOperationLogResp(log), nil
}

// DeleteOperationLog 删除操作日志
func (s *OperationLogServiceImpl) DeleteOperationLog(id uint) error {
	return repository.GetOperationLogRepository().Delete(id)
}

// BatchDeleteOperationLog 批量删除操作日志
func (s *OperationLogServiceImpl) BatchDeleteOperationLog(req *dto.IDsReq) error {
	return repository.GetOperationLogRepository().DeleteBatch(req.IDs)
}

// CleanOldLogs 清理旧日志
func (s *OperationLogServiceImpl) CleanOldLogs(days int) error {
	return repository.GetOperationLogRepository().Clean(days)
}

// RecordOperationLog 记录操作日志
func (s *OperationLogServiceImpl) RecordOperationLog(userID uint, username, module, action, method, path, ip, userAgent, request, response string, status int8, errorMessage string, duration int) error {
	log := &model.OperationLog{
		UserID:       userID,
		Username:     username,
		Module:       module,
		Action:       action,
		Method:       method,
		Path:         path,
		IP:           ip,
		UserAgent:    userAgent,
		Request:      request,
		Response:     response,
		Status:       status,
		ErrorMessage: errorMessage,
		Duration:     duration,
		CreatedAt:    time.Now(),
	}

	return repository.GetOperationLogRepository().Create(log)
}
