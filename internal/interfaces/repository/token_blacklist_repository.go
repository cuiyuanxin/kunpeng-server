package repository

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// TokenBlacklistRepository token黑名单仓储接口
type TokenBlacklistRepository interface {
	// 添加token到黑名单
	Create(blacklist *model.TokenBlacklist) error

	// 检查token是否在黑名单中
	IsTokenBlacklisted(token string) (bool, error)

	// 根据token查找黑名单记录
	FindByToken(token string) (*model.TokenBlacklist, error)

	// 清理过期的黑名单记录
	CleanExpired() error

	// 根据用户ID获取黑名单记录
	FindByUserID(userID uint) ([]*model.TokenBlacklist, error)

	// 删除黑名单记录
	Delete(id uint) error
}
