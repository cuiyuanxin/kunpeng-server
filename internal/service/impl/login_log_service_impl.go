package impl

import (
	"math"
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/internal/repository"
)

// LoginLogServiceImpl 登录日志服务实现
type LoginLogServiceImpl struct{}

// GetLoginLogList 获取登录日志列表
func (s *LoginLogServiceImpl) GetLoginLogList(req *dto.LoginLogListReq) (*dto.PageResp, error) {
	logs, total, err := repository.GetLoginLogRepository().FindList(req)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	logList := dto.ToLoginLogRespList(logs)

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

// GetLoginLogByID 根据ID获取登录日志
func (s *LoginLogServiceImpl) GetLoginLogByID(id uint) (*dto.LoginLogResp, error) {
	log, err := repository.GetLoginLogRepository().FindByID(id)
	if err != nil {
		return nil, err
	}

	return dto.ToLoginLogResp(log), nil
}

// DeleteLoginLog 删除登录日志
func (s *LoginLogServiceImpl) DeleteLoginLog(id uint) error {
	return repository.GetLoginLogRepository().Delete(id)
}

// BatchDeleteLoginLog 批量删除登录日志
func (s *LoginLogServiceImpl) BatchDeleteLoginLog(req *dto.IDsReq) error {
	return repository.GetLoginLogRepository().BatchDelete(req.IDs)
}

// CleanOldLogs 清理旧日志
func (s *LoginLogServiceImpl) CleanOldLogs(days int) error {
	return repository.GetLoginLogRepository().CleanOldLogs(days)
}

// RecordLoginLog 记录登录日志
func (s *LoginLogServiceImpl) RecordLoginLog(userID uint, username, ip, browser, os, device, location string, status int8, message string) error {
	log := &model.LoginLog{
		UserID:    userID,
		Username:  username,
		IP:        ip,
		Location:  location,
		Browser:   browser,
		OS:        os,
		Device:    device,
		Status:    status,
		Message:   message,
		CreatedAt: time.Now(),
	}

	return repository.GetLoginLogRepository().Create(log)
}
