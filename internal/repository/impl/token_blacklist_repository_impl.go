package impl

import (
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// TokenBlacklistRepositoryImpl token黑名单仓储实现
type TokenBlacklistRepositoryImpl struct {
	BaseRepository
}

// NewTokenBlacklistRepository 创建token黑名单仓储
func NewTokenBlacklistRepository() repository.TokenBlacklistRepository {
	return &TokenBlacklistRepositoryImpl{
		BaseRepository: NewBaseRepository(),
	}
}

// Create 添加token到黑名单
func (r *TokenBlacklistRepositoryImpl) Create(blacklist *model.TokenBlacklist) error {
	err := r.db.Create(blacklist).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// IsTokenBlacklisted 检查token是否在黑名单中
func (r *TokenBlacklistRepositoryImpl) IsTokenBlacklisted(token string) (bool, error) {
	var count int64
	err := r.db.Model(&model.TokenBlacklist{}).Where("token = ?", token).Count(&count).Error
	if err != nil {
		return false, kperrors.New(kperrors.ErrDatabase, err)
	}
	return count > 0, nil
}

// FindByToken 根据token查找黑名单记录
func (r *TokenBlacklistRepositoryImpl) FindByToken(token string) (*model.TokenBlacklist, error) {
	var blacklist model.TokenBlacklist
	err := r.db.Where("token = ?", token).First(&blacklist).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}
	return &blacklist, nil
}

// CleanExpired 清理过期的黑名单记录
func (r *TokenBlacklistRepositoryImpl) CleanExpired() error {
	now := time.Now()
	err := r.db.Where("expires_at < ?", now).Delete(&model.TokenBlacklist{}).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// FindByUserID 根据用户ID获取黑名单记录
func (r *TokenBlacklistRepositoryImpl) FindByUserID(userID uint) ([]*model.TokenBlacklist, error) {
	var blacklists []*model.TokenBlacklist
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&blacklists).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return blacklists, nil
}

// Delete 删除黑名单记录
func (r *TokenBlacklistRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&model.TokenBlacklist{}, id).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}
