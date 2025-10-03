package repository

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// LoginLogRepository 登录日志仓储接口
type LoginLogRepository interface {
	// Create 创建登录日志
	Create(log *model.LoginLog) error

	// FindByID 根据ID查询登录日志
	FindByID(id uint) (*model.LoginLog, error)

	// FindList 查询登录日志列表
	FindList(req *dto.LoginLogListReq) ([]*model.LoginLog, int64, error)

	// Delete 删除登录日志
	Delete(id uint) error

	// BatchDelete 批量删除登录日志
	BatchDelete(ids []uint) error

	// DeleteByUserID 根据用户ID删除日志
	DeleteByUserID(userID uint) error

	// CleanOldLogs 清理旧日志（保留指定天数）
	CleanOldLogs(days int) error
}
