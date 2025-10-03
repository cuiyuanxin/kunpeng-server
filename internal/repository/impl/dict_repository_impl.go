package impl

import (
	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// DictRepositoryImpl 字典仓储实现
type DictRepositoryImpl struct {
	BaseRepository
}

// NewDictRepository 创建字典仓储
func NewDictRepository() repository.DictRepository {
	return &DictRepositoryImpl{
		BaseRepository: NewBaseRepository(),
	}
}

// FindDictTypeByID 根据ID获取字典类型
func (r *DictRepositoryImpl) FindDictTypeByID(id uint) (*model.DictType, error) {
	var dictType model.DictType
	err := r.db.First(&dictType, id).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}
	return &dictType, nil
}

// FindAllDictTypes 获取所有字典类型
func (r *DictRepositoryImpl) FindAllDictTypes() ([]*model.DictType, error) {
	var dictTypes []*model.DictType
	err := r.db.Where("status = ?", 1).Find(&dictTypes).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return dictTypes, nil
}

// FindDictTypeList 获取字典类型列表
func (r *DictRepositoryImpl) FindDictTypeList(req *dto.DictTypePageReq) ([]*model.DictType, int64, error) {
	var dictTypes []*model.DictType
	var total int64

	db := r.db.Model(&model.DictType{})

	// 构建查询条件
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Type != "" {
		db = db.Where("type LIKE ?", "%"+req.Type+"%")
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	// 查询总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 分页查询
	err = db.Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&dictTypes).Error
	if err != nil {
		return nil, 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	return dictTypes, total, nil
}

// CreateDictType 创建字典类型
func (r *DictRepositoryImpl) CreateDictType(dictType *model.DictType) error {
	err := r.db.Create(dictType).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// UpdateDictType 更新字典类型
func (r *DictRepositoryImpl) UpdateDictType(dictType *model.DictType) error {
	err := r.db.Save(dictType).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// DeleteDictType 删除字典类型
func (r *DictRepositoryImpl) DeleteDictType(id uint) error {
	// 开启事务
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除字典类型
	if err := tx.Delete(&model.DictType{}, id).Error; err != nil {
		tx.Rollback()
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 删除关联的字典数据
	if err := tx.Where("dict_type_id = ?", id).Delete(&model.DictData{}).Error; err != nil {
		tx.Rollback()
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// BatchDeleteDictType 批量删除字典类型
func (r *DictRepositoryImpl) BatchDeleteDictType(ids []uint) error {
	// 开启事务
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 删除字典类型
	if err := tx.Delete(&model.DictType{}, ids).Error; err != nil {
		tx.Rollback()
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 删除关联的字典数据
	if err := tx.Where("dict_type_id IN ?", ids).Delete(&model.DictData{}).Error; err != nil {
		tx.Rollback()
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// FindDictDataByID 根据ID获取字典数据
func (r *DictRepositoryImpl) FindDictDataByID(id uint) (*model.DictData, error) {
	var dictData model.DictData
	err := r.db.First(&dictData, id).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}
	return &dictData, nil
}

// FindDictDataByType 根据字典类型获取字典数据
func (r *DictRepositoryImpl) FindDictDataByType(dictType string) ([]*model.DictData, error) {
	var dictData []*model.DictData
	err := r.db.
		Joins("JOIN dict_types ON dict_data.dict_type_id = dict_types.id").
		Where("dict_types.type = ? AND dict_data.status = ?", dictType, 1).
		Order("dict_data.sort ASC").
		Find(&dictData).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return dictData, nil
}

// FindDictDataList 获取字典数据列表
func (r *DictRepositoryImpl) FindDictDataList(req *dto.DictDataPageReq) ([]*model.DictData, int64, error) {
	var dictData []*model.DictData
	var total int64

	db := r.db.Model(&model.DictData{})

	// 构建查询条件
	if req.DictTypeID != 0 {
		db = db.Where("dict_type_id = ?", req.DictTypeID)
	}
	if req.Label != "" {
		db = db.Where("label LIKE ?", "%"+req.Label+"%")
	}
	if req.Value != "" {
		db = db.Where("value LIKE ?", "%"+req.Value+"%")
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	// 查询总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 分页查询
	err = db.Order("sort ASC").
		Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&dictData).Error
	if err != nil {
		return nil, 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	return dictData, total, nil
}

// CreateDictData 创建字典数据
func (r *DictRepositoryImpl) CreateDictData(dictData *model.DictData) error {
	err := r.db.Create(dictData).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// UpdateDictData 更新字典数据
func (r *DictRepositoryImpl) UpdateDictData(dictData *model.DictData) error {
	err := r.db.Save(dictData).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// DeleteDictData 删除字典数据
func (r *DictRepositoryImpl) DeleteDictData(id uint) error {
	err := r.db.Delete(&model.DictData{}, id).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// BatchDeleteDictData 批量删除字典数据
func (r *DictRepositoryImpl) BatchDeleteDictData(ids []uint) error {
	err := r.db.Delete(&model.DictData{}, ids).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}
