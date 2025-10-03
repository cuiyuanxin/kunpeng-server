package repository

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// PostRepository 岗位仓储接口
type PostRepository interface {
	// 根据ID获取岗位
	FindByID(id uint) (*model.Post, error)

	// 获取所有岗位
	FindAll() ([]*model.Post, error)

	// 获取岗位列表
	FindList(req *dto.PostPageReq) ([]*model.Post, int64, error)

	// 创建岗位
	Create(post *model.Post) error

	// 更新岗位
	Update(post *model.Post) error

	// 删除岗位
	Delete(id uint) error

	// 批量删除岗位
	BatchDelete(ids []uint) error
}
