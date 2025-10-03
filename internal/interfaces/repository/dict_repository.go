package repository

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// DictRepository 字典仓储接口
type DictRepository interface {
	// 根据ID获取字典类型
	FindDictTypeByID(id uint) (*model.DictType, error)

	// 获取所有字典类型
	FindAllDictTypes() ([]*model.DictType, error)

	// 获取字典类型列表
	FindDictTypeList(req *dto.DictTypePageReq) ([]*model.DictType, int64, error)

	// 创建字典类型
	CreateDictType(dictType *model.DictType) error

	// 更新字典类型
	UpdateDictType(dictType *model.DictType) error

	// 删除字典类型
	DeleteDictType(id uint) error

	// 批量删除字典类型
	BatchDeleteDictType(ids []uint) error

	// 根据ID获取字典数据
	FindDictDataByID(id uint) (*model.DictData, error)

	// 根据字典类型获取字典数据
	FindDictDataByType(dictType string) ([]*model.DictData, error)

	// 获取字典数据列表
	FindDictDataList(req *dto.DictDataPageReq) ([]*model.DictData, int64, error)

	// 创建字典数据
	CreateDictData(dictData *model.DictData) error

	// 更新字典数据
	UpdateDictData(dictData *model.DictData) error

	// 删除字典数据
	DeleteDictData(id uint) error

	// 批量删除字典数据
	BatchDeleteDictData(ids []uint) error
}
