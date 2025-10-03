package repository

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// LoginAttemptRepository 登录尝试记录仓储接口
type LoginAttemptRepository interface {
	// 根据账号和IP获取登录尝试记录
	FindByAccountAndIP(account, ip string) (*model.LoginAttempt, error)

	// 根据账号获取登录尝试记录
	FindByAccount(account string) (*model.LoginAttempt, error)

	// 创建登录尝试记录
	Create(attempt *model.LoginAttempt) error

	// 更新登录尝试记录
	Update(attempt *model.LoginAttempt) error

	// 删除登录尝试记录
	Delete(id uint) error

	// 清理过期的登录尝试记录
	CleanExpired() error

	// 重置账号的登录尝试记录
	ResetByAccount(account string) error
}
