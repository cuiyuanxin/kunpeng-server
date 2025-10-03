package service

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// APIService API服务接口
type APIService interface {
	// GetAPIList 获取API列表
	GetAPIList(req *dto.APIPageReq) (*dto.PageResp, error)

	// GetAllAPIs 获取所有API
	GetAllAPIs() ([]*model.API, error)

	// GetAPIByID 根据ID获取API
	GetAPIByID(id uint) (*model.API, error)

	// CreateAPI 创建API
	CreateAPI(req *dto.APICreateReq) (uint, error)

	// UpdateAPI 更新API
	UpdateAPI(req *dto.APIUpdateReq) error

	// DeleteAPI 删除API
	DeleteAPI(id uint) error

	// BatchDeleteAPI 批量删除API
	BatchDeleteAPI(ids []uint) error
}
