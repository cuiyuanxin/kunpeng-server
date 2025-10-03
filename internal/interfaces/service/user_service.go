package service

import (
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
)

// UserService 用户服务接口
type UserService interface {
	// Login 用户登录
	Login(req *dto.UserLoginReq, clientIP string) (*dto.UserLoginResp, error)

	// RefreshToken 刷新token
	RefreshToken(req *dto.RefreshTokenReq) (*dto.UserLoginResp, error)

	// GetUserInfo 获取用户信息
	GetUserInfo(userID uint) (*dto.UserInfoResp, error)

	// GetUserList 获取用户列表
	GetUserList(req *dto.UserPageReq) (*dto.PageResp, error)

	// GetUserByID 根据ID获取用户
	GetUserByID(id uint) (*model.User, error)

	// CreateUser 创建用户
	CreateUser(req *dto.UserCreateReq) (uint, error)

	// UpdateUser 更新用户
	UpdateUser(req *dto.UserUpdateReq) error

	// DeleteUser 删除用户
	DeleteUser(id uint) error

	// BatchDeleteUser 批量删除用户
	BatchDeleteUser(ids []uint) error

	// ChangeUserStatus 修改用户状态
	ChangeUserStatus(req *dto.StatusReq) error

	// ResetUserPassword 重置用户密码
	ResetUserPassword(id uint) error

	// ChangePassword 修改密码
	ChangePassword(userID uint, req *dto.UserChangePasswordReq) error
}
