package impl

import (
	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// DeptRepositoryImpl 部门仓储实现
type DeptRepositoryImpl struct {
	BaseRepository
}

// NewDeptRepository 创建部门仓储
func NewDeptRepository() repository.DeptRepository {
	return &DeptRepositoryImpl{
		BaseRepository: NewBaseRepository(),
	}
}

// FindByID 根据ID获取部门
func (r *DeptRepositoryImpl) FindByID(id uint) (*model.Dept, error) {
	var dept model.Dept
	err := r.db.First(&dept, id).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}
	return &dept, nil
}

// FindList 获取部门列表
func (r *DeptRepositoryImpl) FindList(req *dto.DeptListReq) ([]*model.Dept, error) {
	var depts []*model.Dept

	db := r.db.Model(&model.Dept{})

	// 构建查询条件
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	// 查询
	err := db.Order("sort ASC").Find(&depts).Error
	if err != nil {
		return nil, errors.New(errors.ErrDatabase, err)
	}

	return depts, nil
}

// FindTree 获取部门树
func (r *DeptRepositoryImpl) FindTree() ([]*dto.DeptTreeResp, error) {
	var depts []*model.Dept

	// 查询所有部门
	err := r.db.Order("sort ASC").Find(&depts).Error
	if err != nil {
		return nil, errors.New(errors.ErrDatabase, err)
	}

	// 构建部门树
	return buildDeptTree(depts, 0), nil
}

// Create 创建部门
func (r *DeptRepositoryImpl) Create(dept *model.Dept) error {
	err := r.db.Create(dept).Error
	if err != nil {
		return errors.New(errors.ErrDatabase, err)
	}
	return nil
}

// Update 更新部门
func (r *DeptRepositoryImpl) Update(dept *model.Dept) error {
	err := r.db.Save(dept).Error
	if err != nil {
		return errors.New(errors.ErrDatabase, err)
	}
	return nil
}

// Delete 删除部门
func (r *DeptRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&model.Dept{}, id).Error
	if err != nil {
		return errors.New(errors.ErrDatabase, err)
	}
	return nil
}

// HasChildren 检查部门是否有子部门
func (r *DeptRepositoryImpl) HasChildren(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Dept{}).Where("parent_id = ?", id).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ErrDatabase, err)
	}
	return count > 0, nil
}

// HasUsers 检查部门是否有用户
func (r *DeptRepositoryImpl) HasUsers(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("dept_id = ?", id).Count(&count).Error
	if err != nil {
		return false, errors.New(errors.ErrDatabase, err)
	}
	return count > 0, nil
}

// buildDeptTree 构建部门树
func buildDeptTree(depts []*model.Dept, parentID uint) []*dto.DeptTreeResp {
	var tree []*dto.DeptTreeResp

	for _, dept := range depts {
		if dept.ParentID == parentID {
			node := &dto.DeptTreeResp{
				ID:       dept.ID,
				ParentID: dept.ParentID,
				Name:     dept.Name,
				Sort:     dept.Sort,
				Leader:   dept.Leader,
				Phone:    dept.Phone,
				Email:    dept.Email,
				Status:   dept.Status,
				Children: buildDeptTree(depts, dept.ID),
			}
			tree = append(tree, node)
		}
	}

	return tree
}
