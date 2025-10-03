package service

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// DeptService 部门服务接口
type DeptService interface {
	// GetDeptList 获取部门列表
	GetDeptList(req *dto.DeptQueryReq) ([]*model.Dept, error)

	// GetDeptTree 获取部门树
	GetDeptTree(req *dto.DeptQueryReq) ([]*dto.DeptTreeResp, error)

	// GetDeptByID 根据ID获取部门
	GetDeptByID(id uint) (*model.Dept, error)

	// CreateDept 创建部门
	CreateDept(req *dto.DeptCreateReq) (uint, error)

	// UpdateDept 更新部门
	UpdateDept(req *dto.DeptUpdateReq) error

	// DeleteDept 删除部门
	DeleteDept(id uint) error
}
