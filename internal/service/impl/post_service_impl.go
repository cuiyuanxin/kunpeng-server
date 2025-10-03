package impl

import (
	"errors"

	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/pkg/database"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"gorm.io/gorm"
)

// PostServiceImpl 岗位服务实现
type PostServiceImpl struct{}

// GetPostList 获取岗位列表
func (s *PostServiceImpl) GetPostList(req *dto.PostPageReq) (*dto.PageResp, error) {
	var posts []model.Post
	var total int64

	db := database.GetDB().Model(&model.Post{})

	// 构建查询条件
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Code != "" {
		db = db.Where("code LIKE ?", "%"+req.Code+"%")
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
	err = db.Order("sort ASC").Offset((req.PageNum - 1) * req.PageSize).Limit(req.PageSize).Find(&posts).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 计算总页数
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		totalPages++
	}

	return &dto.PageResp{
		List:       posts,
		Total:      total,
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetAllPosts 获取所有岗位
func (s *PostServiceImpl) GetAllPosts() ([]*model.Post, error) {
	var posts []*model.Post
	err := database.GetDB().Where("status = ?", 1).Order("sort ASC").Find(&posts).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return posts, nil
}

// GetPostByID 根据ID获取岗位
func (s *PostServiceImpl) GetPostByID(id uint) (*model.Post, error) {
	var post model.Post
	err := database.GetDB().First(&post, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, kperrors.New(kperrors.ErrNotFound, err)
		}
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return &post, nil
}

// CreatePost 创建岗位
func (s *PostServiceImpl) CreatePost(req *dto.PostCreateReq) (uint, error) {
	// 检查岗位编码是否存在
	var count int64
	err := database.GetDB().Model(&model.Post{}).Where("code = ?", req.Code).Count(&count).Error
	if err != nil {
		return 0, kperrors.New(kperrors.ErrDatabase, err)
	}
	if count > 0 {
		return 0, kperrors.New(kperrors.ErrParam, nil).WithMessage("岗位编码已存在")
	}

	// 创建岗位
	post := model.Post{
		Name:   req.Name,
		Code:   req.Code,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	}

	err = database.GetDB().Create(&post).Error
	if err != nil {
		return 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	return post.ID, nil
}

// UpdatePost 更新岗位
func (s *PostServiceImpl) UpdatePost(req *dto.PostUpdateReq) error {
	// 检查岗位是否存在
	var post model.Post
	err := database.GetDB().First(&post, req.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrNotFound, err)
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 检查岗位编码是否存在
	if req.Code != post.Code {
		var count int64
		err = database.GetDB().Model(&model.Post{}).Where("code = ? AND id != ?", req.Code, req.ID).Count(&count).Error
		if err != nil {
			return kperrors.New(kperrors.ErrDatabase, err)
		}
		if count > 0 {
			return kperrors.New(kperrors.ErrParam, nil).WithMessage("岗位编码已存在")
		}
	}

	// 更新岗位
	updates := map[string]interface{}{
		"name":   req.Name,
		"code":   req.Code,
		"sort":   req.Sort,
		"status": req.Status,
		"remark": req.Remark,
	}

	err = database.GetDB().Model(&post).Updates(updates).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// DeletePost 删除岗位
func (s *PostServiceImpl) DeletePost(id uint) error {
	// 检查岗位是否存在
	var post model.Post
	err := database.GetDB().First(&post, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return kperrors.New(kperrors.ErrNotFound, err)
		}
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	// 检查岗位是否有关联用户
	var userCount int64
	err = database.GetDB().Model(&model.User{}).Where("post_id = ?", id).Count(&userCount).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	if userCount > 0 {
		return kperrors.New(kperrors.ErrParam, nil).WithMessage("岗位下存在用户，无法删除")
	}

	// 删除岗位
	err = database.GetDB().Delete(&post).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}

// BatchDeletePost 批量删除岗位
func (s *PostServiceImpl) BatchDeletePost(ids []uint) error {
	// 检查岗位是否有关联用户
	var userCount int64
	err := database.GetDB().Model(&model.User{}).Where("post_id IN ?", ids).Count(&userCount).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	if userCount > 0 {
		return kperrors.New(kperrors.ErrParam, nil).WithMessage("选中的岗位中有岗位下存在用户，无法删除")
	}

	// 删除岗位
	err = database.GetDB().Where("id IN ?", ids).Delete(&model.Post{}).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	return nil
}
