package impl

import (
	"errors"

	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/pkg/database"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"gorm.io/gorm"
)

// APIServiceImpl API服务实现
type APIServiceImpl struct{}

// GetAPIList 获取API列表
func (s *APIServiceImpl) GetAPIList(req *dto.APIPageReq) (*dto.PageResp, error) {
	var apis []model.API
	var total int64

	db := database.GetDB().Model(&model.API{})

	// 构建查询条件
	if req.Group != "" {
		db = db.Where("group LIKE ?", "%"+req.Group+"%")
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Method != "" {
		db = db.Where("method = ?", req.Method)
	}
	if req.Path != "" {
		db = db.Where("path LIKE ?", "%"+req.Path+"%")
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.BeginTime != "" && req.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.BeginTime, req.EndTime)
	}

	// 查询总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 分页查询
	err = db.Order("id DESC").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&apis).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 计算总页数
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		totalPages++
	}

	return &dto.PageResp{
		List:       apis,
		Total:      total,
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAllAPIs 获取所有API
func (s *APIServiceImpl) GetAllAPIs() ([]*model.API, error) {
	var apis []*model.API
	err := database.GetDB().Where("status = ?", 1).Find(&apis).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return apis, nil
}

// GetAPIByID 根据ID获取API
func (s *APIServiceImpl) GetAPIByID(id uint) (*model.API, error) {
	var api model.API
	err := database.GetDB().First(&api, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kperrors.New(kperrors.ErrNotFound, err)
		}
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return &api, nil
}

// CreateAPI 创建API
func (s *APIServiceImpl) CreateAPI(req *dto.APICreateReq) (uint, error) {
	// 检查API是否已存在
	var count int64
	err := database.GetDB().Model(&model.API{}).Where("method = ? AND path = ?", req.Method, req.Path).Count(&count).Error
	if err != nil {
		return 0, kperrors.New(kperrors.ErrDatabase, err)
	}
	if count > 0 {
		return 0, kperrors.New(kperrors.ErrParam, nil).WithMessage("该API已存在")
	}

	// 创建API
	api := model.API{
		Group:  req.Group,
		Name:   req.Name,
		Method: req.Method,
		Path:   req.Path,
		Status: req.Status,
		Remark: req.Remark,
	}

	err = database.GetDB().Create(&api).Error
	if err != nil {
		return 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	return api.ID, nil
}

// UpdateAPI 更新API
func (s *APIServiceImpl) UpdateAPI(req *dto.APIUpdateReq) error {
	// 检查API是否存在
	var api model.API
	err := database.GetDB().First(&api, req.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrNotFound, err)
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 检查API是否已存在
	if req.Method != api.Method || req.Path != api.Path {
		var count int64
		err = database.GetDB().Model(&model.API{}).Where("method = ? AND path = ? AND id != ?", req.Method, req.Path, req.ID).Count(&count).Error
		if err != nil {
			return kperrors.New(kperrors.ErrDatabase, err)
		}
		if count > 0 {
			return kperrors.New(kperrors.ErrParam, nil).WithMessage("该API已存在")
		}
	}

	// 更新API
	updates := map[string]interface{}{
		"group":  req.Group,
		"name":   req.Name,
		"method": req.Method,
		"path":   req.Path,
		"status": req.Status,
		"remark": req.Remark,
	}

	err = database.GetDB().Model(&api).Updates(updates).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// DeleteAPI 删除API
func (s *APIServiceImpl) DeleteAPI(id uint) error {
	// 检查API是否存在
	var api model.API
	err := database.GetDB().First(&api, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrNotFound, err)
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 删除API
	err = database.GetDB().Delete(&api).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 删除角色API关联
	err = database.GetDB().Where("api_id = ?", id).Delete(&model.RoleAPI{}).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// BatchDeleteAPI 批量删除API
func (s *APIServiceImpl) BatchDeleteAPI(ids []uint) error {
	// 删除API
	err := database.GetDB().Where("id IN ?", ids).Delete(&model.API{}).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 删除角色API关联
	err = database.GetDB().Where("api_id IN ?", ids).Delete(&model.RoleAPI{}).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}
