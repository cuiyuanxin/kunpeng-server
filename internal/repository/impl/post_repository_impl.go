package impl

import (
	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// PostRepositoryImpl 岗位仓储实现
type PostRepositoryImpl struct {
	BaseRepository
}

// NewPostRepository 创建岗位仓储
func NewPostRepository() repository.PostRepository {
	return &PostRepositoryImpl{
		BaseRepository: NewBaseRepository(),
	}
}

// FindByID 根据ID获取岗位
func (r *PostRepositoryImpl) FindByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}
	return &post, nil
}

// FindAll 获取所有岗位
func (r *PostRepositoryImpl) FindAll() ([]*model.Post, error) {
	var posts []*model.Post
	err := r.db.Where("status = ?", 1).Order("sort ASC").Find(&posts).Error
	if err != nil {
		return nil, kperrors.New(kperrors.ErrDatabase, err)
	}
	return posts, nil
}

// FindList 获取岗位列表
func (r *PostRepositoryImpl) FindList(req *dto.PostPageReq) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	db := r.db.Model(&model.Post{})

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

	// 查询总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 分页查询
	err = db.Order("sort ASC").
		Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&posts).Error
	if err != nil {
		return nil, 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	return posts, total, nil
}

// Create 创建岗位
func (r *PostRepositoryImpl) Create(post *model.Post) error {
	err := r.db.Create(post).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// Update 更新岗位
func (r *PostRepositoryImpl) Update(post *model.Post) error {
	err := r.db.Save(post).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// Delete 删除岗位
func (r *PostRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&model.Post{}, id).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// BatchDelete 批量删除岗位
func (r *PostRepositoryImpl) BatchDelete(ids []uint) error {
	err := r.db.Delete(&model.Post{}, ids).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}
