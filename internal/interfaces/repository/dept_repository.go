package repository

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// DeptRepository 部门仓储接口
type DeptRepository interface {
	// 根据ID获取部门
	FindByID(id uint) (*model.Dept, error)

	// 获取部门列表
	FindList(req *dto.DeptListReq) ([]*model.Dept, error)

	// 获取部门树
	FindTree() ([]*dto.DeptTreeResp, error)

	// 创建部门
	Create(dept *model.Dept) error

	// 更新部门
	Update(dept *model.Dept) error

	// 删除部门
	Delete(id uint) error

	// 检查部门是否有子部门
	HasChildren(id uint) (bool, error)

	// 检查部门是否有用户
	HasUsers(id uint) (bool, error)
}
