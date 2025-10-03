package impl

import (
	"errors"

	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/pkg/database"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"gorm.io/gorm"
)

// DictServiceImpl 字典服务实现
type DictServiceImpl struct{}

// GetDictTypeList 获取字典类型列表
func (s *DictServiceImpl) GetDictTypeList(req *dto.DictTypePageReq) (*dto.PageResp, error) {
	var dictTypes []model.DictType
	var total int64

	db := database.GetDB().Model(&model.DictType{})

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
	if req.BeginTime != "" && req.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.BeginTime, req.EndTime)
	}

	// 查询总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 分页查询
	err = db.Order("id DESC").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&dictTypes).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 计算总页数
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		totalPages++
	}

	return &dto.PageResp{
		List:       dictTypes,
		Total:      total,
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAllDictTypes 获取所有字典类型
func (s *DictServiceImpl) GetAllDictTypes() ([]*model.DictType, error) {
	var dictTypes []*model.DictType
	err := database.GetDB().Where("status = ?", 1).Find(&dictTypes).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return dictTypes, nil
}

// GetDictTypeByID 根据ID获取字典类型
func (s *DictServiceImpl) GetDictTypeByID(id uint) (*model.DictType, error) {
	var dictType model.DictType
	err := database.GetDB().First(&dictType, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kperrors.New(kperrors.ErrNotFound, err)
		}
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return &dictType, nil
}

// CreateDictType 创建字典类型
func (s *DictServiceImpl) CreateDictType(req *dto.DictTypeCreateReq) (uint, error) {
	// 检查字典类型是否存在
	var count int64
	err := database.GetDB().Model(&model.DictType{}).Where("type = ?", req.Type).Count(&count).Error
	if err != nil {
		return 0, kperrors.New(kperrors.ErrDatabase, err)
	}
	if count > 0 {
		return 0, kperrors.New(kperrors.ErrParam, nil).WithMessage("字典类型已存在")
	}

	// 创建字典类型
	dictType := model.DictType{
		Name:   req.Name,
		Type:   req.Type,
		Status: req.Status,
		Remark: req.Remark,
	}

	err = database.GetDB().Create(&dictType).Error
	if err != nil {
		return 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	return dictType.ID, nil
}

// UpdateDictType 更新字典类型
func (s *DictServiceImpl) UpdateDictType(req *dto.DictTypeUpdateReq) error {
	// 检查字典类型是否存在
	var dictType model.DictType
	err := database.GetDB().First(&dictType, req.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrNotFound, err)
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 检查字典类型是否存在
	if req.Type != dictType.Type {
		var count int64
		err = database.GetDB().Model(&model.DictType{}).Where("type = ? AND id != ?", req.Type, req.ID).Count(&count).Error
		if err != nil {
			return kperrors.New(kperrors.ErrDatabase, err)
		}
		if count > 0 {
			return kperrors.New(kperrors.ErrParam, nil).WithMessage("字典类型已存在")
		}
	}

	// 更新字典类型
	updates := map[string]interface{}{
		"name":   req.Name,
		"type":   req.Type,
		"status": req.Status,
		"remark": req.Remark,
	}

	err = database.GetDB().Model(&dictType).Updates(updates).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 如果字典类型发生变化，更新字典数据的类型
	if req.Type != dictType.Type {
		err = database.GetDB().Model(&model.DictData{}).Where("dict_type = ?", dictType.Type).Update("dict_type", req.Type).Error
		if err != nil {
			return kperrors.New(kperrors.ErrDatabase, err)
		}
	}

	return nil
}

// DeleteDictType 删除字典类型
func (s *DictServiceImpl) DeleteDictType(id uint) error {
	// 检查字典类型是否存在
	var dictType model.DictType
	err := database.GetDB().First(&dictType, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrNotFound, err)
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 检查是否有关联的字典数据
	var count int64
	err = database.GetDB().Model(&model.DictData{}).Where("dict_type = ?", dictType.Type).Count(&count).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	if count > 0 {
		return kperrors.New(kperrors.ErrParam, nil).WithMessage("字典类型下存在字典数据，无法删除")
	}

	// 删除字典类型
	err = database.GetDB().Delete(&dictType).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// BatchDeleteDictType 批量删除字典类型
func (s *DictServiceImpl) BatchDeleteDictType(ids []uint) error {
	// 获取字典类型
	var dictTypes []model.DictType
	err := database.GetDB().Where("id IN ?", ids).Find(&dictTypes).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 提取字典类型
	types := make([]string, 0, len(dictTypes))
	for _, dictType := range dictTypes {
		types = append(types, dictType.Type)
	}

	// 检查是否有关联的字典数据
	var count int64
	err = database.GetDB().Model(&model.DictData{}).Where("dict_type IN ?", types).Count(&count).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	if count > 0 {
		return kperrors.New(kperrors.ErrParam, nil).WithMessage("选中的字典类型中有字典类型下存在字典数据，无法删除")
	}

	// 删除字典类型
	err = database.GetDB().Where("id IN ?", ids).Delete(&model.DictType{}).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// GetDictDataList 获取字典数据列表
func (s *DictServiceImpl) GetDictDataList(req *dto.DictDataPageReq) (*dto.PageResp, error) {
	var dictDatas []model.DictData
	var total int64

	db := database.GetDB().Model(&model.DictData{})

	// 构建查询条件
	db = db.Where("dict_type = ?", req.DictType)
	if req.Label != "" {
		db = db.Where("label LIKE ?", "%"+req.Label+"%")
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
	err = db.Order("sort ASC").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&dictDatas).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 计算总页数
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		totalPages++
	}

	return &dto.PageResp{
		List:       dictDatas,
		Total:      total,
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetDictDataByType 根据字典类型获取字典数据
func (s *DictServiceImpl) GetDictDataByType(dictType string) ([]*model.DictData, error) {
	var dictDatas []*model.DictData
	err := database.GetDB().Where("dict_type = ? AND status = ?", dictType, 1).Order("sort ASC").Find(&dictDatas).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return dictDatas, nil
}

// GetDictDataByID 根据ID获取字典数据
func (s *DictServiceImpl) GetDictDataByID(id uint) (*model.DictData, error) {
	var dictData model.DictData
	err := database.GetDB().First(&dictData, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kperrors.New(kperrors.ErrNotFound, err)
		}
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return &dictData, nil
}

// CreateDictData 创建字典数据
func (s *DictServiceImpl) CreateDictData(req *dto.DictDataCreateReq) (uint, error) {
	// 检查字典类型是否存在
	var dictType model.DictType
	err := database.GetDB().Where("type = ?", req.DictType).First(&dictType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, kperrors.New(kperrors.ErrParam, err).WithMessage("字典类型不存在")
		}
		return 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 检查字典数据是否存在
	var count int64
	err = database.GetDB().Model(&model.DictData{}).Where("dict_type = ? AND value = ?", req.DictType, req.Value).Count(&count).Error
	if err != nil {
		return 0, kperrors.New(kperrors.ErrDatabase, err)
	}
	if count > 0 {
		return 0, kperrors.New(kperrors.ErrParam, nil).WithMessage("字典数据值已存在")
	}

	// 创建字典数据
	dictData := model.DictData{
		DictType: req.DictType,
		Label:    req.Label,
		Value:    req.Value,
		Sort:     req.Sort,
		Status:   req.Status,
		Remark:   req.Remark,
	}

	err = database.GetDB().Create(&dictData).Error
	if err != nil {
		return 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	return dictData.ID, nil
}

// UpdateDictData 更新字典数据
func (s *DictServiceImpl) UpdateDictData(req *dto.DictDataUpdateReq) error {
	// 检查字典数据是否存在
	var dictData model.DictData
	err := database.GetDB().First(&dictData, req.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrNotFound, err)
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 检查字典类型是否存在
	var dictType model.DictType
	err = database.GetDB().Where("type = ?", req.DictType).First(&dictType).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrParam, err).WithMessage("字典类型不存在")
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 检查字典数据是否存在
	if req.Value != dictData.Value || req.DictType != dictData.DictType {
		var count int64
		err = database.GetDB().Model(&model.DictData{}).Where("dict_type = ? AND value = ? AND id != ?", req.DictType, req.Value, req.ID).Count(&count).Error
		if err != nil {
			return kperrors.New(kperrors.ErrDatabase, err)
		}
		if count > 0 {
			return kperrors.New(kperrors.ErrParam, nil).WithMessage("字典数据值已存在")
		}
	}

	// 更新字典数据
	updates := map[string]interface{}{
		"dict_type": req.DictType,
		"label":     req.Label,
		"value":     req.Value,
		"sort":      req.Sort,
		"status":    req.Status,
		"remark":    req.Remark,
	}

	err = database.GetDB().Model(&dictData).Updates(updates).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// DeleteDictData 删除字典数据
func (s *DictServiceImpl) DeleteDictData(id uint) error {
	// 检查字典数据是否存在
	var dictData model.DictData
	err := database.GetDB().First(&dictData, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrNotFound, err)
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 删除字典数据
	err = database.GetDB().Delete(&dictData).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// BatchDeleteDictData 批量删除字典数据
func (s *DictServiceImpl) BatchDeleteDictData(ids []uint) error {
	// 删除字典数据
	err := database.GetDB().Where("id IN ?", ids).Delete(&model.DictData{}).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}
