package repository

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	// 根据ID获取用户
	FindByID(id uint) (*model.User, error)

	// 根据用户名获取用户
	FindByUsername(username string) (*model.User, error)

	// 根据手机号获取用户
	FindByMobile(mobile string) (*model.User, error)

	// 获取用户列表
	FindList(req *dto.UserPageReq) ([]*model.User, int64, error)

	// 创建用户
	Create(user *model.User) error

	// 更新用户
	Update(user *model.User) error

	// 删除用户
	Delete(id uint) error

	// 批量删除用户
	BatchDelete(ids []uint) error

	// 更新用户状态
	UpdateStatus(id uint, status int) error

	// 重置用户密码
	ResetPassword(id uint, password string) error

	// 更新用户密码
	UpdatePassword(id uint, password string) error

	// 根据角色ID查找用户
	FindByRoleID(roleID uint) error
}
