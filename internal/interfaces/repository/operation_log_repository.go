package repository

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// OperationLogRepository 操作日志仓储接口
type OperationLogRepository interface {
	// Create 创建操作日志
	Create(log *model.OperationLog) error

	// FindByID 根据ID查找操作日志
	FindByID(id uint) (*model.OperationLog, error)

	// FindList 分页查询操作日志列表
	FindList(req *dto.OperationLogListReq) ([]*model.OperationLog, int64, error)

	// Delete 删除操作日志
	Delete(id uint) error

	// DeleteBatch 批量删除操作日志
	DeleteBatch(ids []uint) error

	// Clean 清理指定天数前的操作日志
	Clean(days int) error
}
