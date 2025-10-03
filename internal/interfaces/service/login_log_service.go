package service

import (
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// LoginLogService 登录日志服务接口
type LoginLogService interface {
	// GetLoginLogList 获取登录日志列表
	GetLoginLogList(req *dto.LoginLogListReq) (*dto.PageResp, error)

	// GetLoginLogByID 根据ID获取登录日志
	GetLoginLogByID(id uint) (*dto.LoginLogResp, error)

	// DeleteLoginLog 删除登录日志
	DeleteLoginLog(id uint) error

	// BatchDeleteLoginLog 批量删除登录日志
	BatchDeleteLoginLog(req *dto.IDsReq) error

	// CleanOldLogs 清理旧日志
	CleanOldLogs(days int) error

	// RecordLoginLog 记录登录日志
	RecordLoginLog(userID uint, username, ip, browser, os, device, location string, status int8, message string) error
}
