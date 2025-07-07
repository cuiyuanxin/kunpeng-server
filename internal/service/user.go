package service

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/internal/database"

	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务
func NewUserService() *UserService {
	return &UserService{
		db: database.GetDB(),
	}
}

// NewUserServiceWithDB 使用指定数据库连接创建用户服务
func NewUserServiceWithDB(db *gorm.DB) *UserService {
	return &UserService{
		db: db,
	}
}

// Create 创建用户
func (s *UserService) Create(req *model.UserCreateRequest) (*model.User, error) {
	// 检查用户名是否已存在
	var existUser model.User
	if err := s.db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existUser).Error; err == nil {
		if existUser.Username == req.Username {
			return nil, errors.New("username already exists")
		}
		if existUser.Email == req.Email {
			return nil, errors.New("email already exists")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Role:     req.Role,
		Status:   1, // 默认启用
	}

	if user.Role == "" {
		user.Role = "user" // 默认角色
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (s *UserService) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (s *UserService) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// Update 更新用户
func (s *UserService) Update(id uint, req *model.UserUpdateRequest) (*model.User, error) {
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.db.Save(user).Error; err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// Delete 删除用户（软删除）
func (s *UserService) Delete(id uint) error {
	if err := s.db.Delete(&model.User{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// List 获取用户列表
func (s *UserService) List(req *model.UserListRequest) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// 构建查询
	query := s.db.Model(&model.User{})

	// 关键词搜索
	if req.Keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			"%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 角色筛选
	if req.Role != "" {
		query = query.Where("role = ?", req.Role)
	}

	// 状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// 分页查询
	if err := query.Scopes(database.Paginate(req.Page, req.PageSize)).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	return users, total, nil
}

// VerifyPassword 验证密码
func (s *UserService) VerifyPassword(user *model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id uint, req *model.UserChangePasswordRequest) error {
	user, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !s.VerifyPassword(user, req.OldPassword) {
		return errors.New("old password is incorrect")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	// 更新密码
	if err := s.db.Model(user).Update("password", string(hashedPassword)).Error; err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

// UpdateLastLogin 更新最后登录时间
func (s *UserService) UpdateLastLogin(id uint) error {
	now := time.Now()
	if err := s.db.Model(&model.User{}).Where("id = ?", id).Update("last_login", &now).Error; err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}
