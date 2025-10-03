package impl

import (
	"github.com/cuiyuanxin/kunpeng/internal/interfaces/repository"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
)

// UserRepositoryImpl 用户仓储实现
type UserRepositoryImpl struct {
	BaseRepository
}

// NewUserRepository 创建用户仓储
func NewUserRepository() repository.UserRepository {
	return &UserRepositoryImpl{
		BaseRepository: NewBaseRepository(),
	}
}

// FindByID 根据ID获取用户
func (r *UserRepositoryImpl) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}
	return &user, nil
}

// FindByUsername 根据用户名获取用户
func (r *UserRepositoryImpl) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}
	return &user, nil
}

// FindByMobile 根据手机号获取用户
func (r *UserRepositoryImpl) FindByMobile(mobile string) (*model.User, error) {
	var user model.User
	err := r.db.Where("mobile = ?", mobile).First(&user).Error
	if err != nil {
		return nil, r.HandleDBError(err)
	}
	return &user, nil
}

// FindList 获取用户列表
func (r *UserRepositoryImpl) FindList(req *dto.UserPageReq) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	db := r.db.Model(&model.User{})

	// 构建查询条件
	if req.Username != "" {
		db = db.Where("username LIKE ?", "%"+req.Username+"%")
	}
	if req.Nickname != "" {
		db = db.Where("nickname LIKE ?", "%"+req.Nickname+"%")
	}
	if req.Mobile != "" {
		db = db.Where("mobile LIKE ?", "%"+req.Mobile+"%")
	}
	if req.Status != 0 {
		db = db.Where("status = ?", req.Status)
	}
	if req.DeptID != 0 {
		db = db.Where("dept_id = ?", req.DeptID)
	}
	if req.BeginTime != "" && req.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.BeginTime, req.EndTime)
	}

	// 查询总数
	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	// 分页查询
	err = db.Preload("Role").Preload("Dept").Preload("Post").
		Offset((req.PageNum - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&users).Error
	if err != nil {
		return nil, 0, kperrors.New(kperrors.ErrDatabase, err)
	}

	return users, total, nil
}

// Create 创建用户
func (r *UserRepositoryImpl) Create(user *model.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// Update 更新用户
func (r *UserRepositoryImpl) Update(user *model.User) error {
	err := r.db.Save(user).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// Delete 删除用户
func (r *UserRepositoryImpl) Delete(id uint) error {
	err := r.db.Delete(&model.User{}, id).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// BatchDelete 批量删除用户
func (r *UserRepositoryImpl) BatchDelete(ids []uint) error {
	err := r.db.Delete(&model.User{}, ids).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// UpdateStatus 更新用户状态
func (r *UserRepositoryImpl) UpdateStatus(id uint, status int) error {
	err := r.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// ResetPassword 重置用户密码
func (r *UserRepositoryImpl) ResetPassword(id uint, password string) error {
	err := r.db.Model(&model.User{}).Where("id = ?", id).Update("password", password).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// UpdatePassword 更新用户密码
func (r *UserRepositoryImpl) UpdatePassword(id uint, password string) error {
	err := r.db.Model(&model.User{}).Where("id = ?", id).Update("password", password).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}
	return nil
}

// FindByRoleID 根据角色ID查找用户
func (r *UserRepositoryImpl) FindByRoleID(roleID uint) error {
	var count int64
	err := r.db.Model(&model.User{}).Where("role_id = ?", roleID).Count(&count).Error
	if err != nil {
		return kperrors.New(kperrors.ErrDatabase, err)
	}

	if count > 0 {
		return nil // 找到了用户，返回nil表示存在
	}

	return kperrors.New(kperrors.ErrUserNotFound, nil) // 没找到用户，返回错误
}
