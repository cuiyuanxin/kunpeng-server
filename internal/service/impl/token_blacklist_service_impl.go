package impl

import (
	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/interfaces/service"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/cuiyuanxin/kunpeng/pkg/jwt"
)

// TokenBlacklistServiceImpl token黑名单服务实现
type TokenBlacklistServiceImpl struct {
	tokenBlacklistRepository repository.TokenBlacklistRepository
}

// NewTokenBlacklistService 创建token黑名单服务
func NewTokenBlacklistService(tokenBlacklistRepository repository.TokenBlacklistRepository) service.TokenBlacklistService {
	return &TokenBlacklistServiceImpl{
		tokenBlacklistRepository: tokenBlacklistRepository,
	}
}

// BlacklistToken 将token加入黑名单
func (s *TokenBlacklistServiceImpl) BlacklistToken(token string, userID uint, username string, reason string) error {
	// 解析token获取过期时间
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return kperrors.New(kperrors.ErrInvalidToken, err)
	}

	// 创建黑名单记录
	blacklist := &model.TokenBlacklist{
		Token:     token,
		UserID:    userID,
		Username:  username,
		Reason:    reason,
		ExpiresAt: claims.ExpiresAt.Time,
	}

	return s.tokenBlacklistRepository.Create(blacklist)
}

// IsTokenBlacklisted 检查token是否在黑名单中
func (s *TokenBlacklistServiceImpl) IsTokenBlacklisted(token string) (bool, error) {
	return s.tokenBlacklistRepository.IsTokenBlacklisted(token)
}

// GetUserBlacklistTokens 根据用户ID获取黑名单记录
func (s *TokenBlacklistServiceImpl) GetUserBlacklistTokens(userID uint) ([]*model.TokenBlacklist, error) {
	return s.tokenBlacklistRepository.FindByUserID(userID)
}

// CleanExpiredTokens 清理过期的黑名单记录
func (s *TokenBlacklistServiceImpl) CleanExpiredTokens() error {
	return s.tokenBlacklistRepository.CleanExpired()
}

// DeleteBlacklistToken 删除黑名单记录
func (s *TokenBlacklistServiceImpl) DeleteBlacklistToken(id uint) error {
	return s.tokenBlacklistRepository.Delete(id)
}
