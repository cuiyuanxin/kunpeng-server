package impl

import (
	"errors"

	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/pkg/database"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"gorm.io/gorm"
)

// DeptServiceImpl 部门服务实现
type DeptServiceImpl struct{}

// GetDeptList 获取部门列表
func (s *DeptServiceImpl) GetDeptList(req *dto.DeptQueryReq) ([]*model.Dept, error) {
	var depts []*model.Dept
	db := database.GetDB().Model(&model.Dept{})

	// 构建查询条件
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}

	// 查询部门
	err := db.Order("sort ASC").Find(&depts).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	return depts, nil
}

// GetDeptTree 获取部门树
func (s *DeptServiceImpl) GetDeptTree(req *dto.DeptQueryReq) ([]*dto.DeptTreeResp, error) {
	// 获取部门列表
	depts, err := s.GetDeptList(req)
	if err != nil {
		return nil, err
	}

	// 构建部门树
	return s.buildDeptTree(depts, 0), nil
}

// GetDeptByID 根据ID获取部门
func (s *DeptServiceImpl) GetDeptByID(id uint) (*model.Dept, error) {
	var dept model.Dept
	err := database.GetDB().First(&dept, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kperrors.New(kperrors.ErrNotFound, err)
		}
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return &dept, nil
}

// CreateDept 创建部门
func (s *DeptServiceImpl) CreateDept(req *dto.DeptCreateReq) (uint, error) {
	// 检查父部门是否存在
	if req.ParentID != 0 {
		var parentDept model.Dept
		err := database.GetDB().First(&parentDept, req.ParentID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return 0, kperrors.New(kperrors.ErrParam, err).WithMessage("父部门不存在")
			}
			return 0, kperrors.New(kperrors.ErrDatabase, err)
		}

		// 检查父部门状态
		if parentDept.Status != 1 {
			return 0, kperrors.New(kperrors.ErrParam, nil).WithMessage("父部门已被禁用")
		}
	}

	// 创建部门
	dept := model.Dept{
		ParentID: req.ParentID,
		Name:     req.Name,
		Leader:   req.Leader,
		Phone:    req.Phone,
		Email:    req.Email,
		Sort:     req.Sort,
		Status:   req.Status,
		Remark:   req.Remark,
	}

	err := database.GetDB().Create(&dept).Error
	if err != nil {
		return 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	return dept.ID, nil
}

// UpdateDept 更新部门
func (s *DeptServiceImpl) UpdateDept(req *dto.DeptUpdateReq) error {
	// 检查部门是否存在
	var dept model.Dept
	err := database.GetDB().First(&dept, req.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrNotFound, err)
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 检查父部门是否存在
	if req.ParentID != 0 {
		var parentDept model.Dept
		err := database.GetDB().First(&parentDept, req.ParentID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return kperrors.New(kperrors.ErrParam, err).WithMessage("父部门不存在")
			}
			return kperrors.New(kperrors.ErrDatabase, err)
		}

		// 检查父部门状态
		if parentDept.Status != 1 {
			return kperrors.New(kperrors.ErrParam, nil).WithMessage("父部门已被禁用")
		}

		// 检查是否将部门设置为自己的子部门
		if req.ParentID == req.ID {
			return kperrors.New(kperrors.ErrParam, nil).WithMessage("不能将部门设置为自己的子部门")
		}

		// 检查是否将部门设置为自己子部门的子部门
		var childDepts []*model.Dept
		err = database.GetDB().Where("parent_id = ?", req.ID).Find(&childDepts).Error
		if err != nil {
			return kperrors.New(kperrors.ErrDatabase, err)
		}
		childIDs := make([]uint, 0, len(childDepts))
		for _, childDept := range childDepts {
			childIDs = append(childIDs, childDept.ID)
		}
		if len(childIDs) > 0 && contains(childIDs, req.ParentID) {
			return kperrors.New(kperrors.ErrParam, nil).WithMessage("不能将部门设置为自己子部门的子部门")
		}
	}

	// 更新部门
	updates := map[string]interface{}{
		"parent_id": req.ParentID,
		"name":      req.Name,
		"leader":    req.Leader,
		"phone":     req.Phone,
		"email":     req.Email,
		"sort":      req.Sort,
		"status":    req.Status,
		"remark":    req.Remark,
	}

	err = database.GetDB().Model(&dept).Updates(updates).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// DeleteDept 删除部门
func (s *DeptServiceImpl) DeleteDept(id uint) error {
	// 检查部门是否存在
	var dept model.Dept
	err := database.GetDB().First(&dept, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrNotFound, err)
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 检查是否有子部门
	var childCount int64
	err = database.GetDB().Model(&model.Dept{}).Where("parent_id = ?", id).Count(&childCount).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	if childCount > 0 {
		return kperrors.New(kperrors.ErrParam, nil).WithMessage("存在子部门，无法删除")
	}

	// 检查是否有关联用户
	var userCount int64
	err = database.GetDB().Model(&model.User{}).Where("dept_id = ?", id).Count(&userCount).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	if userCount > 0 {
		return kperrors.New(kperrors.ErrParam, nil).WithMessage("部门下存在用户，无法删除")
	}

	// 删除部门
	err = database.GetDB().Delete(&dept).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// buildDeptTree 构建部门树
func (s *DeptServiceImpl) buildDeptTree(depts []*model.Dept, parentID uint) []*dto.DeptTreeResp {
	var tree []*dto.DeptTreeResp
	for _, dept := range depts {
		if dept.ParentID == parentID {
			node := &dto.DeptTreeResp{
				ID:       dept.ID,
				ParentID: dept.ParentID,
				Name:     dept.Name,
				Leader:   dept.Leader,
				Phone:    dept.Phone,
				Email:    dept.Email,
				Sort:     dept.Sort,
				Status:   dept.Status,
				Children: s.buildDeptTree(depts, dept.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
}
