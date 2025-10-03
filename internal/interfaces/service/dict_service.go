package service

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// DictService 字典服务接口
type DictService interface {
	// GetDictTypeList 获取字典类型列表
	GetDictTypeList(req *dto.DictTypePageReq) (*dto.PageResp, error)

	// GetAllDictTypes 获取所有字典类型
	GetAllDictTypes() ([]*model.DictType, error)

	// GetDictTypeByID 根据ID获取字典类型
	GetDictTypeByID(id uint) (*model.DictType, error)

	// CreateDictType 创建字典类型
	CreateDictType(req *dto.DictTypeCreateReq) (uint, error)

	// UpdateDictType 更新字典类型
	UpdateDictType(req *dto.DictTypeUpdateReq) error

	// DeleteDictType 删除字典类型
	DeleteDictType(id uint) error

	// BatchDeleteDictType 批量删除字典类型
	BatchDeleteDictType(ids []uint) error

	// GetDictDataList 获取字典数据列表
	GetDictDataList(req *dto.DictDataPageReq) (*dto.PageResp, error)

	// GetDictDataByType 根据字典类型获取字典数据
	GetDictDataByType(dictType string) ([]*model.DictData, error)

	// GetDictDataByID 根据ID获取字典数据
	GetDictDataByID(id uint) (*model.DictData, error)

	// CreateDictData 创建字典数据
	CreateDictData(req *dto.DictDataCreateReq) (uint, error)

	// UpdateDictData 更新字典数据
	UpdateDictData(req *dto.DictDataUpdateReq) error

	// DeleteDictData 删除字典数据
	DeleteDictData(id uint) error

	// BatchDeleteDictData 批量删除字典数据
	BatchDeleteDictData(ids []uint) error
}
