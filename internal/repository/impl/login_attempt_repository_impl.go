package impl

import (
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// LoginAttemptRepositoryImpl 登录尝试仓储实现
type LoginAttemptRepositoryImpl struct {
	BaseRepository
}

// NewLoginAttemptRepository 创建登录尝试仓储实例
func NewLoginAttemptRepository() repository.LoginAttemptRepository {
	return &LoginAttemptRepositoryImpl{
		BaseRepository: NewBaseRepository(),
	}
}

// FindByAccountAndIP 根据账号和IP获取登录尝试记录
func (r *LoginAttemptRepositoryImpl) FindByAccountAndIP(account, ip string) (*model.LoginAttempt, error) {
	var attempt model.LoginAttempt
	err := r.db.Where("account = ? AND ip = ?", account, ip).First(&attempt).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}
	return &attempt, nil
}

// FindByAccount 根据账号获取登录尝试记录
func (r *LoginAttemptRepositoryImpl) FindByAccount(account string) (*model.LoginAttempt, error) {
	var attempt model.LoginAttempt
	err := r.db.Where("account = ?", account).First(&attempt).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}
	return &attempt, nil
}

// Create 创建登录尝试记录
func (r *LoginAttemptRepositoryImpl) Create(attempt *model.LoginAttempt) error {
	err := r.db.Create(attempt).Error
	if err != nil {
		return r.HandleDBError(err)
	}
	return nil
}

// Update 更新登录尝试记录
func (r *LoginAttemptRepositoryImpl) Update(attempt *model.LoginAttempt) error {
	err := r.db.Save(attempt).Error
	if err != nil {
		return r.HandleDBError(err)
	}
	return nil
}

// Delete 删除登录尝试记录
func (r *LoginAttemptRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&model.LoginAttempt{}, id).Error
	if err != nil {
		return r.HandleDBError(err)
	}
	return nil
}

// CleanExpired 清理过期的登录尝试记录
func (r *LoginAttemptRepositoryImpl) CleanExpired() error {
	// 清理30天前的记录
	expiredTime := time.Now().AddDate(0, 0, -30)
	return r.db.Where("created_at < ?", expiredTime).Delete(&model.LoginAttempt{}).Error
}

// ResetByAccount 重置账号的登录尝试记录
func (r *LoginAttemptRepositoryImpl) ResetByAccount(account string) error {
	err := r.db.Where("account = ?", account).Delete(&model.LoginAttempt{}).Error
	if err != nil {
		return r.HandleDBError(err)
	}
	return nil
}
