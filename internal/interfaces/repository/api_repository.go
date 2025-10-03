package repository

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// APIRepository API仓储接口
type APIRepository interface {
	// 根据ID获取API
	FindByID(id uint) (*model.API, error)

	// 获取所有API
	FindAll() ([]*model.API, error)

	// 获取API列表
	FindList(req *dto.APIPageReq) ([]*model.API, int64, error)

	// 创建API
	Create(api *model.API) error

	// 更新API
	Update(api *model.API) error

	// 删除API
	Delete(id uint) error

	// 批量删除API
	BatchDelete(ids []uint) error

	// 根据角色ID获取API ID列表
	FindAPIIDsByRoleID(roleID uint) ([]uint, error)
}
