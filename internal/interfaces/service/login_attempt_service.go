package service

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// LoginAttemptService 登录尝试服务接口
type LoginAttemptService interface {
	// CheckAndRecordAttempt 检查并记录登录尝试
	// 返回是否被拉黑和错误信息
	CheckAndRecordAttempt(account, ip string, success bool) (blocked bool, err error)

	// IsBlocked 检查账号是否被拉黑
	IsBlocked(account, ip string) (bool, error)

	// ResetAttempts 重置账号的登录尝试记录
	ResetAttempts(account string) error

	// CleanupExpiredRecords 清理过期的登录尝试记录
	CleanupExpiredRecords() error

	// GetAttempt 获取登录尝试记录
	GetAttempt(account, ip string) (*model.LoginAttempt, error)
}
