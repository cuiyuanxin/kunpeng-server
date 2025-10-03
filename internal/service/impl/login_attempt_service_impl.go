package impl

import (
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/interfaces/service"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
)

const (
	// MaxFailedAttempts 最大失败尝试次数
	MaxFailedAttempts = 5
	// BlockDuration 拉黑时长（2小时）
	BlockDuration = 2 * time.Hour
)

// LoginAttemptServiceImpl 登录尝试服务实现
type LoginAttemptServiceImpl struct {
	loginAttemptRepo repository.LoginAttemptRepository
}

// NewLoginAttemptService 创建登录尝试服务实例
func NewLoginAttemptService(loginAttemptRepo repository.LoginAttemptRepository) service.LoginAttemptService {
	return &LoginAttemptServiceImpl{
		loginAttemptRepo: loginAttemptRepo,
	}
}

// CheckAndRecordAttempt 检查并记录登录尝试
func (s *LoginAttemptServiceImpl) CheckAndRecordAttempt(account, ip string, success bool) (blocked bool, err error) {
	// 获取现有记录
	attempt, err := s.loginAttemptRepo.FindByAccountAndIP(account, ip)
	if err != nil {
		if kpErr, ok := err.(*kperrors.Error); ok && kpErr.Code == kperrors.ErrDBNotFound {
			// 记录不存在，继续创建新记录
		} else {
			return false, err
		}
	}

	// 如果没有记录，创建新记录
	if attempt == nil {
		attempt = &model.LoginAttempt{
			Account:  account,
			IP:       ip,
			Attempts: 0,
			LastTry:  time.Now(),
		}
	}

	// 检查是否已被拉黑
	if attempt.IsBlocked() {
		return true, nil
	}

	// 更新尝试记录
	attempt.LastTry = time.Now()

	if success {
		// 登录成功，重置失败次数
		attempt.Reset()
	} else {
		// 登录失败，增加失败次数
		attempt.Attempts++

		// 检查是否需要拉黑
		if attempt.ShouldBlock() {
			attempt.Block()
			blocked = true
		}
	}

	// 保存记录
	if attempt.ID == 0 {
		err = s.loginAttemptRepo.Create(attempt)
	} else {
		err = s.loginAttemptRepo.Update(attempt)
	}

	return blocked, err
}

// IsBlocked 检查账号是否被拉黑
func (s *LoginAttemptServiceImpl) IsBlocked(account, ip string) (bool, error) {
	attempt, err := s.loginAttemptRepo.FindByAccountAndIP(account, ip)
	if err != nil {
		if kperrors.IsCode(err, kperrors.ErrDBNotFound) {
			return false, nil
		}
		return false, err
	}

	return attempt.IsBlocked(), nil
}

// ResetAttempts 重置账号的登录尝试记录
func (s *LoginAttemptServiceImpl) ResetAttempts(account string) error {
	return s.loginAttemptRepo.ResetByAccount(account)
}

// CleanupExpiredRecords 清理过期的登录尝试记录
func (s *LoginAttemptServiceImpl) CleanupExpiredRecords() error {
	// 清理过期的记录
	return s.loginAttemptRepo.CleanExpired()
}

// GetAttempt 获取登录尝试记录
func (s *LoginAttemptServiceImpl) GetAttempt(account, ip string) (*model.LoginAttempt, error) {
	return s.loginAttemptRepo.FindByAccountAndIP(account, ip)
}
