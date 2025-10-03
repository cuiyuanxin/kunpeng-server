package service

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// TokenBlacklistService token黑名单服务接口
type TokenBlacklistService interface {
	// 将token加入黑名单
	BlacklistToken(token string, userID uint, username string, reason string) error

	// 检查token是否在黑名单中
	IsTokenBlacklisted(token string) (bool, error)

	// 根据用户ID获取黑名单记录
	GetUserBlacklistTokens(userID uint) ([]*model.TokenBlacklist, error)

	// 清理过期的黑名单记录
	CleanExpiredTokens() error

	// 删除黑名单记录
	DeleteBlacklistToken(id uint) error
}
