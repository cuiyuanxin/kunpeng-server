package impl

import (
	"errors"
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// LoginLogRepositoryImpl 登录日志仓储实现
type LoginLogRepositoryImpl struct {
	*BaseRepository
}

// NewLoginLogRepository 创建登录日志仓储实例
func NewLoginLogRepository() repository.LoginLogRepository {
	baseRepo := NewBaseRepository()
	return &LoginLogRepositoryImpl{
		BaseRepository: &baseRepo,
	}
}

// Create 创建登录日志
func (r *LoginLogRepositoryImpl) Create(log *model.LoginLog) error {
	if err := r.db.Create(log).Error; err != nil {
		return r.HandleDBError(err)
	}
	return nil
}

// FindByID 根据ID查询登录日志
func (r *LoginLogRepositoryImpl) FindByID(id uint) (*model.LoginLog, error) {
	var log model.LoginLog
	if err := r.db.First(&log, id).Error; err != nil {
		return nil, r.HandleDBError(err)
	}
	return &log, nil
}

// FindList 查询登录日志列表
func (r *LoginLogRepositoryImpl) FindList(req *dto.LoginLogListReq) ([]*model.LoginLog, int64, error) {
	var logs []*model.LoginLog
	var total int64

	db := r.db.Model(&model.LoginLog{})

	// 构建查询条件
	if req.UserID != nil {
		db = db.Where("user_id = ?", *req.UserID)
	}
	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.IP != "" {
		db = db.Where("ip LIKE ?", "%"+req.IP+"%")
	}
	if req.BeginTime != "" && req.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.BeginTime, req.EndTime)
	}

	// 查询总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, r.HandleDBError(err)
	}

	// 分页查询
	if err := db.Order("created_at DESC").
		Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&logs).Error; err != nil {
		return nil, 0, r.HandleDBError(err)
	}

	return logs, total, nil
}

// Delete 删除登录日志
func (r *LoginLogRepositoryImpl) Delete(id uint) error {
	if err := r.db.Delete(&model.LoginLog{}, id).Error; err != nil {
		return r.HandleDBError(err)
	}
	return nil
}

// BatchDelete 批量删除登录日志
func (r *LoginLogRepositoryImpl) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return kperrors.New(kperrors.ErrParam, errors.New("删除ID列表不能为空"))
	}

	if err := r.db.Where("id IN ?", ids).Delete(&model.LoginLog{}).Error; err != nil {
		return r.HandleDBError(err)
	}
	return nil
}

// DeleteByUserID 根据用户ID删除日志
func (r *LoginLogRepositoryImpl) DeleteByUserID(userID uint) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&model.LoginLog{}).Error; err != nil {
		return r.HandleDBError(err)
	}
	return nil
}

// CleanOldLogs 清理旧日志（保留指定天数）
func (r *LoginLogRepositoryImpl) CleanOldLogs(days int) error {
	if days <= 0 {
		return kperrors.New(kperrors.ErrParam, errors.New("保留天数必须大于0"))
	}

	// 计算截止时间
	cutoffTime := time.Now().AddDate(0, 0, -days)

	if err := r.db.Where("created_at < ?", cutoffTime).Delete(&model.LoginLog{}).Error; err != nil {
		return r.HandleDBError(err)
	}
	return nil
}
