package impl

import (
	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// APIRepositoryImpl API仓储实现
type APIRepositoryImpl struct {
	BaseRepository
}

// NewAPIRepository 创建API仓储
func NewAPIRepository() repository.APIRepository {
	return &APIRepositoryImpl{
		BaseRepository: NewBaseRepository(),
	}
}

// FindByID 根据ID获取API
func (r *APIRepositoryImpl) FindByID(id uint) (*model.API, error) {
	var api model.API
	err := r.db.First(&api, id).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}
	return &api, nil
}

// FindAll 获取所有API
func (r *APIRepositoryImpl) FindAll() ([]*model.API, error) {
	var apis []*model.API
	err := r.db.Find(&apis).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return apis, nil
}

// FindList 获取API列表
func (r *APIRepositoryImpl) FindList(req *dto.APIPageReq) ([]*model.API, int64, error) {
	var apis []*model.API
	var total int64

	db := r.db.Model(&model.API{})

	// 构建查询条件
	if req.Path != "" {
		db = db.Where("path LIKE ?", "%"+req.Path+"%")
	}
	if req.Method != "" {
		db = db.Where("method = ?", req.Method)
	}
	if req.Group != "" {
		db = db.Where("group = ?", req.Group)
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}

	// 查询总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 分页查询
	err = db.Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&apis).Error
	if err != nil {
		return nil, 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	return apis, total, nil
}

// Create 创建API
func (r *APIRepositoryImpl) Create(api *model.API) error {
	err := r.db.Create(api).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// Update 更新API
func (r *APIRepositoryImpl) Update(api *model.API) error {
	err := r.db.Save(api).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// Delete 删除API
func (r *APIRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&model.API{}, id).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// BatchDelete 批量删除API
func (r *APIRepositoryImpl) BatchDelete(ids []uint) error {
	err := r.db.Delete(&model.API{}, ids).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// FindAPIIDsByRoleID 根据角色ID获取API ID列表
func (r *APIRepositoryImpl) FindAPIIDsByRoleID(roleID uint) ([]uint, error) {
	var apiIDs []uint
	err := r.db.Model(&model.RoleAPI{}).
		Select("api_id").
		Where("role_id = ?", roleID).
		Pluck("api_id", &apiIDs).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return apiIDs, nil
}
