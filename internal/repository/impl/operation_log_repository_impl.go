package impl

import (
	"errors"
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// OperationLogRepositoryImpl 操作日志仓储实现
type OperationLogRepositoryImpl struct {
	*BaseRepository
}

// NewOperationLogRepository 创建操作日志仓储实例
func NewOperationLogRepository() repository.OperationLogRepository {
	baseRepo := NewBaseRepository()
	return &OperationLogRepositoryImpl{
		BaseRepository: &baseRepo,
	}
}

// Create 创建操作日志
func (r *OperationLogRepositoryImpl) Create(log *model.OperationLog) error {
	if err := r.db.Create(log).Error; err != nil {
		return r.HandleDBError(err)
	}
	return nil
}

// FindByID 根据ID查找操作日志
func (r *OperationLogRepositoryImpl) FindByID(id uint) (*model.OperationLog, error) {
	var log model.OperationLog
	err := r.db.First(&log, id).Error
	if err != nil {
		if errors.Is(err, r.db.Error) {
			return nil, kperrors.New(kperrors.ErrNotFound, err)
		}
		return nil, r.HandleDBError(err)
	}
	return &log, nil
}

// FindList 分页查询操作日志列表
func (r *OperationLogRepositoryImpl) FindList(req *dto.OperationLogListReq) ([]*model.OperationLog, int64, error) {
	var logs []*model.OperationLog
	var total int64

	db := r.db.Model(&model.OperationLog{})

	// 构建查询条件
	if req.UserID != nil {
		db = db.Where("user_id = ?", *req.UserID)
	}
	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Module != "" {
		db = db.Where("module LIKE ?", "%"+req.Module+"%")
	}
	if req.Action != "" {
		db = db.Where("action LIKE ?", "%"+req.Action+"%")
	}
	if req.Method != "" {
		db = db.Where("method = ?", req.Method)
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
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, r.HandleDBError(err)
	}

	// 分页查询
	err = db.Order("created_at DESC").
		Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&logs).Error
	if err != nil {
		return nil, 0, r.HandleDBError(err)
	}

	return logs, total, nil
}

// Delete 删除操作日志
func (r *OperationLogRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&model.OperationLog{}, id).Error
	if err != nil {
		return r.HandleDBError(err)
	}
	return nil
}

// DeleteBatch 批量删除操作日志
func (r *OperationLogRepositoryImpl) DeleteBatch(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	err := r.db.Where("id IN ?", ids).Delete(&model.OperationLog{}).Error
	if err != nil {
		return r.HandleDBError(err)
	}
	return nil
}

// Clean 清理指定天数前的操作日志
func (r *OperationLogRepositoryImpl) Clean(days int) error {
	if days <= 0 {
		return errors.New("days must be greater than 0")
	}

	cutoffTime := time.Now().AddDate(0, 0, -days)
	err := r.db.Where("created_at < ?", cutoffTime).Delete(&model.OperationLog{}).Error
	if err != nil {
		return r.HandleDBError(err)
	}
	return nil
}
