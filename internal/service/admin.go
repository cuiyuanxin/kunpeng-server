package service

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// AdminUserService 管理员用户服务
type AdminUserService struct {
	db *gorm.DB
}

// NewAdminUserService 创建管理员用户服务实例
func NewAdminUserService(db *gorm.DB) *AdminUserService {
	return &AdminUserService{
		db: db,
	}
}

// Create 创建管理员用户
func (s *AdminUserService) Create(req *model.AdminUserCreateRequest) (*model.AdminUser, error) {
	// 检查用户名是否已存在
	var existingUser model.AdminUser
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查手机号是否已存在
	if req.Phone != nil && *req.Phone != "" {
		if err := s.db.Where("phone = ?", *req.Phone).First(&existingUser).Error; err == nil {
			return nil, errors.New("手机号已存在")
		}
	}

	// 检查邮箱是否已存在
	if req.Email != nil && *req.Email != "" {
		if err := s.db.Where("email = ?", *req.Email).First(&existingUser).Error; err == nil {
			return nil, errors.New("邮箱已存在")
		}
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}

	// 处理生日
	var birthday *time.Time
	if req.Birthday != nil && *req.Birthday != "" {
		if parsedTime, err := time.Parse("2006-01-02", *req.Birthday); err == nil {
			birthday = &parsedTime
		}
	}

	// 创建用户
	user := &model.AdminUser{
		Username:          req.Username,
		Phone:             req.Phone,
		Email:             req.Email,
		Password:          string(hashedPassword),
		RealName:          req.RealName,
		Nickname:          req.Nickname,
		Avatar:            req.Avatar,
		Gender:            req.Gender,
		Birthday:          birthday,
		DepartmentID:      req.DepartmentID,
		Position:          req.Position,
		Status:            req.Status,
		PasswordChangedAt: &[]time.Time{time.Now()}[0],
		Remark:            req.Remark,
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建用户
	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建用户失败: %v", err)
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		var roles []model.AdminRole
		if err := tx.Where("id IN ?", req.RoleIDs).Find(&roles).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("查询角色失败: %v", err)
		}
		if err := tx.Model(user).Association("Roles").Append(roles); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("分配角色失败: %v", err)
		}
	}

	tx.Commit()
	return user, nil
}

// GetByID 根据ID获取管理员用户
func (s *AdminUserService) GetByID(id uint) (*model.AdminUser, error) {
	var user model.AdminUser
	err := s.db.Preload("Department").Preload("Roles").Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取管理员用户
func (s *AdminUserService) GetByUsername(username string) (*model.AdminUser, error) {
	var user model.AdminUser
	err := s.db.Preload("Department").Preload("Roles").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号获取管理员用户
func (s *AdminUserService) GetByPhone(phone string) (*model.AdminUser, error) {
	var user model.AdminUser
	err := s.db.Preload("Department").Preload("Roles").Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取管理员用户
func (s *AdminUserService) GetByEmail(email string) (*model.AdminUser, error) {
	var user model.AdminUser
	err := s.db.Preload("Department").Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// GetByAccount 根据账号获取管理员用户（支持用户名、手机号、邮箱）
func (s *AdminUserService) GetByAccount(account string) (*model.AdminUser, error) {
	// 判断账号类型
	if isPhone(account) {
		return s.GetByPhone(account)
	} else if isEmail(account) {
		return s.GetByEmail(account)
	} else {
		return s.GetByUsername(account)
	}
}

// Update 更新管理员用户
func (s *AdminUserService) Update(id uint, req *model.AdminUserUpdateRequest) (*model.AdminUser, error) {
	// 获取用户
	user, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 检查手机号是否已被其他用户使用
	if req.Phone != nil && *req.Phone != "" {
		var existingUser model.AdminUser
		if err := s.db.Where("phone = ? AND id != ?", *req.Phone, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("手机号已被其他用户使用")
		}
	}

	// 检查邮箱是否已被其他用户使用
	if req.Email != nil && *req.Email != "" {
		var existingUser model.AdminUser
		if err := s.db.Where("email = ? AND id != ?", *req.Email, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("邮箱已被其他用户使用")
		}
	}

	// 更新字段
	updates := make(map[string]interface{})
	if req.Phone != nil {
		updates["phone"] = req.Phone
	}
	if req.Email != nil {
		updates["email"] = req.Email
	}
	if req.RealName != nil {
		updates["real_name"] = req.RealName
	}
	if req.Nickname != nil {
		updates["nickname"] = req.Nickname
	}
	if req.Avatar != nil {
		updates["avatar"] = req.Avatar
	}
	if req.Gender != nil {
		updates["gender"] = *req.Gender
	}
	if req.Birthday != nil && *req.Birthday != "" {
		if parsedTime, err := time.Parse("2006-01-02", *req.Birthday); err == nil {
			updates["birthday"] = parsedTime
		}
	}
	if req.DepartmentID != nil {
		updates["department_id"] = req.DepartmentID
	}
	if req.Position != nil {
		updates["position"] = req.Position
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != nil {
		updates["remark"] = req.Remark
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 更新用户信息
	if len(updates) > 0 {
		if err := tx.Model(user).Updates(updates).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("更新用户失败: %v", err)
		}
	}

	// 更新角色
	if req.RoleIDs != nil {
		// 清除现有角色
		if err := tx.Model(user).Association("Roles").Clear(); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("清除角色失败: %v", err)
		}

		// 分配新角色
		if len(req.RoleIDs) > 0 {
			var roles []model.AdminRole
			if err := tx.Where("id IN ?", req.RoleIDs).Find(&roles).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("查询角色失败: %v", err)
			}
			if err := tx.Model(user).Association("Roles").Append(roles); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("分配角色失败: %v", err)
			}
		}
	}

	tx.Commit()

	// 重新获取更新后的用户信息
	return s.GetByID(id)
}

// Delete 删除管理员用户（软删除）
func (s *AdminUserService) Delete(id uint) error {
	user, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// 检查是否为超级管理员
	if user.IsSuperAdmin {
		return errors.New("不能删除超级管理员")
	}

	return s.db.Delete(&model.AdminUser{}, id).Error
}

// List 获取管理员用户列表
func (s *AdminUserService) List(req *model.AdminUserListRequest) ([]model.AdminUser, int64, error) {
	var users []model.AdminUser
	var total int64

	// 构建查询
	query := s.db.Model(&model.AdminUser{})

	// 关键词搜索
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("username LIKE ? OR real_name LIKE ? OR nickname LIKE ? OR phone LIKE ? OR email LIKE ?",
			keyword, keyword, keyword, keyword, keyword)
	}

	// 部门筛选
	if req.DepartmentID != nil {
		query = query.Where("department_id = ?", *req.DepartmentID)
	}

	// 状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 角色筛选
	if req.RoleID != nil {
		query = query.Joins("JOIN admin_user_roles ON admin_users.id = admin_user_roles.admin_user_id").
			Where("admin_user_roles.admin_role_id = ?", *req.RoleID)
	}

	// 时间范围筛选
	if req.StartTime != nil && *req.StartTime != "" {
		if startTime, err := time.Parse("2006-01-02", *req.StartTime); err == nil {
			query = query.Where("created_at >= ?", startTime)
		}
	}
	if req.EndTime != nil && *req.EndTime != "" {
		if endTime, err := time.Parse("2006-01-02", *req.EndTime); err == nil {
			endTime = endTime.Add(24 * time.Hour) // 包含当天
			query = query.Where("created_at < ?", endTime)
		}
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 查询数据
	err := query.Preload("Department").Preload("Roles").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&users).Error

	return users, total, err
}

// VerifyPassword 验证密码
func (s *AdminUserService) VerifyPassword(user *model.AdminUser, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// ValidateUser 验证用户登录
func (s *AdminUserService) ValidateUser(account, password string) (*model.AdminUser, error) {
	var user model.AdminUser

	// 根据账号类型查询用户
	query := s.db.Where("status = ?", 1)

	// 判断账号类型并查询
	if isEmail(account) {
		query = query.Where("email = ?", account)
	} else if isPhone(account) {
		query = query.Where("phone = ?", account)
	} else {
		query = query.Where("username = ?", account)
	}

	if err := query.First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在或已被禁用")
		}
		return nil, err
	}

	// 验证密码
	if !checkPassword(password, user.Password) {
		return nil, errors.New("密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	return &user, nil
}

// hashPassword 加密密码
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}



// ChangePassword 修改密码
func (s *AdminUserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	// 获取用户
	var user model.AdminUser
	if err := s.db.First(&user, id).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if !checkPassword(oldPassword, user.Password) {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// 更新密码
	now := time.Now()
	updates := map[string]interface{}{
		"password":            hashedPassword,
		"password_changed_at": &now,
	}

	return s.db.Model(&user).Updates(updates).Error
}

// RecordLoginLog 记录登录日志
func (s *AdminUserService) RecordLoginLog(userID uint, ip, userAgent string) {
	// 这里可以实现登录日志记录逻辑
	// 暂时为空实现
}

// UpdateProfile 更新个人信息
func (s *AdminUserService) UpdateProfile(id uint, req *model.AdminProfileUpdateRequest) (*model.AdminUser, error) {
	// 获取用户
	var user model.AdminUser
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 更新字段
	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}

	// 执行更新
	if err := s.db.Model(&user).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新获取更新后的用户信息
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateLoginInfo 更新登录信息
func (s *AdminUserService) UpdateLoginInfo(id uint, ip string) error {
	updates := map[string]interface{}{
		"last_login_time": time.Now(),
		"last_login_ip":   ip,
		"login_count":     gorm.Expr("login_count + 1"),
	}
	return s.db.Model(&model.AdminUser{}).Where("id = ?", id).Updates(updates).Error
}

// Login 管理员登录
func (s *AdminUserService) Login(req *model.AdminUserLoginRequest, ip string) (*model.AdminUser, error) {
	// 根据账号获取用户
	user, err := s.GetByAccount(req.Account)
	if err != nil {
		return nil, errors.New("账号或密码错误")
	}

	// 检查用户状态
	if user.Status == 0 {
		return nil, errors.New("账号已被禁用")
	}

	// 验证密码
	if !s.VerifyPassword(user, req.Password) {
		return nil, errors.New("账号或密码错误")
	}

	// 更新登录信息
	if err := s.UpdateLoginInfo(user.ID, ip); err != nil {
		// 记录日志，但不影响登录
		fmt.Printf("更新登录信息失败: %v\n", err)
	}

	return user, nil
}

// ===== 角色服务 =====

// AdminRoleService 管理员角色服务
type AdminRoleService struct {
	db *gorm.DB
}

// NewAdminRoleService 创建管理员角色服务实例
func NewAdminRoleService(db *gorm.DB) *AdminRoleService {
	return &AdminRoleService{
		db: db,
	}
}

// Create 创建角色
func (s *AdminRoleService) Create(req *model.AdminRoleCreateRequest) (*model.AdminRole, error) {
	// 检查角色编码是否已存在
	var existingRole model.AdminRole
	if err := s.db.Where("code = ?", req.Code).First(&existingRole).Error; err == nil {
		return nil, errors.New("角色编码已存在")
	}

	// 创建角色
	role := &model.AdminRole{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Level:       req.Level,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
	}

	// 开启事务
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 创建角色
	if err := tx.Create(role).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("创建角色失败: %v", err)
	}

	// 分配权限
	if len(req.PermissionIDs) > 0 {
		var permissions []model.AdminPermission
		if err := tx.Where("id IN ?", req.PermissionIDs).Find(&permissions).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("查询权限失败: %v", err)
		}
		if err := tx.Model(role).Association("Permissions").Append(permissions); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("分配权限失败: %v", err)
		}
	}

	tx.Commit()
	return role, nil
}

// GetByID 根据ID获取角色
func (s *AdminRoleService) GetByID(id uint) (*model.AdminRole, error) {
	var role model.AdminRole
	err := s.db.Preload("Permissions").Where("id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, err
	}
	return &role, nil
}

// List 获取角色列表
func (s *AdminRoleService) List() ([]model.AdminRole, error) {
	var roles []model.AdminRole
	err := s.db.Preload("Permissions").Order("sort_order ASC, created_at DESC").Find(&roles).Error
	return roles, err
}

// ===== 权限服务 =====

// AdminPermissionService 管理员权限服务
type AdminPermissionService struct {
	db *gorm.DB
}

// NewAdminPermissionService 创建管理员权限服务实例
func NewAdminPermissionService(db *gorm.DB) *AdminPermissionService {
	return &AdminPermissionService{
		db: db,
	}
}

// GetUserPermissions 获取用户权限
func (s *AdminPermissionService) GetUserPermissions(userID uint) ([]model.AdminPermission, error) {
	var permissions []model.AdminPermission

	// 通过用户角色获取权限
	err := s.db.Table("admin_permissions").
		Joins("JOIN admin_role_permissions ON admin_permissions.id = admin_role_permissions.admin_permission_id").
		Joins("JOIN admin_user_roles ON admin_role_permissions.admin_role_id = admin_user_roles.admin_role_id").
		Where("admin_user_roles.admin_user_id = ? AND admin_permissions.status = 1", userID).
		Group("admin_permissions.id").
		Order("admin_permissions.sort_order ASC").
		Find(&permissions).Error

	return permissions, err
}

// GetMenuTree 获取菜单树
func (s *AdminPermissionService) GetMenuTree(userID uint) ([]model.AdminPermission, error) {
	// 获取用户权限
	permissions, err := s.GetUserPermissions(userID)
	if err != nil {
		return nil, err
	}

	// 过滤菜单权限
	var menuPermissions []model.AdminPermission
	for _, permission := range permissions {
		if permission.Type == 1 && !permission.IsHidden { // 菜单类型且不隐藏
			menuPermissions = append(menuPermissions, permission)
		}
	}

	// 构建树形结构
	return buildPermissionTree(menuPermissions, 0), nil
}

// List 获取权限列表
func (s *AdminPermissionService) List() ([]model.AdminPermission, error) {
	var permissions []model.AdminPermission
	err := s.db.Order("sort_order ASC, created_at DESC").Find(&permissions).Error
	return permissions, err
}

// GetTree 获取权限树
func (s *AdminPermissionService) GetTree() ([]model.AdminPermission, error) {
	permissions, err := s.List()
	if err != nil {
		return nil, err
	}
	return buildPermissionTree(permissions, 0), nil
}

// ===== 部门服务 =====

// AdminDepartmentService 管理员部门服务
type AdminDepartmentService struct {
	db *gorm.DB
}

// NewAdminDepartmentService 创建管理员部门服务实例
func NewAdminDepartmentService(db *gorm.DB) *AdminDepartmentService {
	return &AdminDepartmentService{
		db: db,
	}
}

// List 获取部门列表
func (s *AdminDepartmentService) List() ([]model.AdminDepartment, error) {
	var departments []model.AdminDepartment
	err := s.db.Preload("Manager").Order("sort_order ASC, created_at DESC").Find(&departments).Error
	return departments, err
}

// GetTree 获取部门树
func (s *AdminDepartmentService) GetTree() ([]model.AdminDepartment, error) {
	departments, err := s.List()
	if err != nil {
		return nil, err
	}
	return buildDepartmentTree(departments, 0), nil
}

// ===== 工具函数 =====

// checkPassword 验证密码
func checkPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// isEmail 检查是否为邮箱格式
func isEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// isPhone 检查是否为手机号格式
func isPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	return phoneRegex.MatchString(phone)
}



// buildPermissionTree 构建权限树
func buildPermissionTree(permissions []model.AdminPermission, parentID uint) []model.AdminPermission {
	var tree []model.AdminPermission
	for _, permission := range permissions {
		if permission.ParentID == parentID {
			children := buildPermissionTree(permissions, permission.ID)
			if len(children) > 0 {
				permission.Children = children
			}
			tree = append(tree, permission)
		}
	}
	return tree
}

// buildDepartmentTree 构建部门树
func buildDepartmentTree(departments []model.AdminDepartment, parentID uint) []model.AdminDepartment {
	var tree []model.AdminDepartment
	for _, department := range departments {
		if department.ParentID == parentID {
			children := buildDepartmentTree(departments, department.ID)
			if len(children) > 0 {
				department.Children = children
			}
			tree = append(tree, department)
		}
	}
	return tree
}