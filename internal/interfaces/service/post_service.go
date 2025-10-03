package service

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// PostService 岗位服务接口
type PostService interface {
	// GetPostList 获取岗位列表
	GetPostList(req *dto.PostPageReq) (*dto.PageResp, error)

	// GetAllPosts 获取所有岗位
	GetAllPosts() ([]*model.Post, error)

	// GetPostByID 根据ID获取岗位
	GetPostByID(id uint) (*model.Post, error)

	// CreatePost 创建岗位
	CreatePost(req *dto.PostCreateReq) (uint, error)

	// UpdatePost 更新岗位
	UpdatePost(req *dto.PostUpdateReq) error

	// DeletePost 删除岗位
	DeletePost(id uint) error

	// BatchDeletePost 批量删除岗位
	BatchDeletePost(ids []uint) error
}
