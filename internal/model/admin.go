package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// AdminUser 管理员用户模型
type AdminUser struct {
	ID                uint           `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Username          string         `json:"username" gorm:"type:varchar(50);uniqueIndex;not null;comment:用户名/账号" validate:"required,min=3,max=50"`
	Phone             *string        `json:"phone" gorm:"type:varchar(20);uniqueIndex;comment:手机号" validate:"omitempty,len=11"`
	Email             *string        `json:"email" gorm:"type:varchar(100);uniqueIndex;comment:邮箱" validate:"omitempty,email"`
	Password          string         `json:"-" gorm:"type:varchar(255);not null;comment:密码(bcrypt加密)" validate:"required,min=6"`
	RealName          *string        `json:"real_name" gorm:"type:varchar(50);comment:真实姓名"`
	Nickname          *string        `json:"nickname" gorm:"type:varchar(50);comment:昵称"`
	Avatar            *string        `json:"avatar" gorm:"type:varchar(255);comment:头像URL"`
	Gender            int8           `json:"gender" gorm:"type:tinyint(1);default:0;comment:性别: 0-未知, 1-男, 2-女"`
	Birthday          *time.Time     `json:"birthday" gorm:"type:date;comment:生日"`
	DepartmentID      *uint          `json:"department_id" gorm:"comment:部门ID"`
	Position          *string        `json:"position" gorm:"type:varchar(100);comment:职位"`
	Status            int8           `json:"status" gorm:"type:tinyint(1);not null;default:1;comment:状态: 0-禁用, 1-启用"`
	IsSuperAdmin      bool           `json:"is_super_admin" gorm:"type:tinyint(1);not null;default:0;comment:是否超级管理员: 0-否, 1-是"`
	LastLoginTime     *time.Time     `json:"last_login_time" gorm:"comment:最后登录时间"`
	LastLoginIP       *string        `json:"last_login_ip" gorm:"type:varchar(45);comment:最后登录IP"`
	LoginCount        int            `json:"login_count" gorm:"type:int;not null;default:0;comment:登录次数"`
	PasswordChangedAt *time.Time     `json:"password_changed_at" gorm:"comment:密码修改时间"`
	Remark            *string        `json:"remark" gorm:"type:text;comment:备注"`
	CreatedAt         time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt         gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间(软删除)"`

	// 关联关系
	Department *AdminDepartment `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	Roles      []AdminRole      `json:"roles,omitempty" gorm:"many2many:admin_user_roles;"`
}

// GetCasbinSubject 获取Casbin主体标识
func (u *AdminUser) GetCasbinSubject() string {
	return fmt.Sprintf("user:%d", u.ID)
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "admin_users"
}

// ToResponse 转换为响应格式
func (u *AdminUser) ToResponse() *AdminUserResponse {
	return &AdminUserResponse{
		ID:            u.ID,
		Username:      u.Username,
		Phone:         u.Phone,
		Email:         u.Email,
		RealName:      u.RealName,
		Nickname:      u.Nickname,
		Avatar:        u.Avatar,
		Gender:        u.Gender,
		Birthday:      u.Birthday,
		DepartmentID:  u.DepartmentID,
		Position:      u.Position,
		Status:        u.Status,
		IsSuperAdmin:  u.IsSuperAdmin,
		LastLoginTime: u.LastLoginTime,
		LastLoginIP:   u.LastLoginIP,
		LoginCount:    u.LoginCount,
		Remark:        u.Remark,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		Department:    u.Department,
		Roles:         u.Roles,
	}
}

// AdminRole 管理员角色模型
type AdminRole struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Name        string         `json:"name" gorm:"type:varchar(50);not null;comment:角色名称" validate:"required,max=50"`
	Code        string         `json:"code" gorm:"type:varchar(50);uniqueIndex;not null;comment:角色编码" validate:"required,max=50"`
	Description *string        `json:"description" gorm:"type:text;comment:角色描述"`
	Level       int            `json:"level" gorm:"type:int;not null;default:1;comment:角色级别(数字越小权限越高)"`
	Status      int8           `json:"status" gorm:"type:tinyint(1);not null;default:1;comment:状态: 0-禁用, 1-启用"`
	IsSystem    bool           `json:"is_system" gorm:"type:tinyint(1);not null;default:0;comment:是否系统角色: 0-否, 1-是"`
	SortOrder   int            `json:"sort_order" gorm:"type:int;not null;default:0;comment:排序"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间(软删除)"`

	// 关联关系
	Users       []AdminUser       `json:"users,omitempty" gorm:"many2many:admin_user_roles;"`
	Permissions []AdminPermission `json:"permissions,omitempty" gorm:"many2many:admin_role_permissions;"`
}

// GetCasbinRole 获取Casbin角色标识
func (r *AdminRole) GetCasbinRole() string {
	return fmt.Sprintf("role:%d", r.ID)
}

// TableName 指定表名
func (AdminRole) TableName() string {
	return "admin_roles"
}

// AdminPermission 管理员权限模型
type AdminPermission struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	ParentID    uint           `json:"parent_id" gorm:"not null;default:0;comment:父级权限ID"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null;comment:权限名称" validate:"required,max=100"`
	Code        string         `json:"code" gorm:"type:varchar(100);uniqueIndex;not null;comment:权限编码" validate:"required,max=100"`
	Type        int8           `json:"type" gorm:"type:tinyint(1);not null;default:1;comment:权限类型: 1-菜单, 2-按钮, 3-接口"`
	Path        *string        `json:"path" gorm:"type:varchar(255);comment:路由路径"`
	Component   *string        `json:"component" gorm:"type:varchar(255);comment:组件路径"`
	Icon        *string        `json:"icon" gorm:"type:varchar(100);comment:图标"`
	Method      *string        `json:"method" gorm:"type:varchar(10);comment:HTTP方法(GET,POST,PUT,DELETE等)"`
	URL         *string        `json:"url" gorm:"type:varchar(255);comment:API接口地址"`
	Level       int            `json:"level" gorm:"type:int;not null;default:1;comment:层级"`
	SortOrder   int            `json:"sort_order" gorm:"type:int;not null;default:0;comment:排序"`
	Status      int8           `json:"status" gorm:"type:tinyint(1);not null;default:1;comment:状态: 0-禁用, 1-启用"`
	IsHidden    bool           `json:"is_hidden" gorm:"type:tinyint(1);not null;default:0;comment:是否隐藏: 0-否, 1-是"`
	Description *string        `json:"description" gorm:"type:text;comment:权限描述"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间(软删除)"`

	// 关联关系
	Children []AdminPermission `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Parent   *AdminPermission  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Roles    []AdminRole       `json:"roles,omitempty" gorm:"many2many:admin_role_permissions;"`
}

// GetCasbinObject 获取Casbin对象标识
func (p *AdminPermission) GetCasbinObject() string {
	return p.Code
}

// GetCasbinAction 获取Casbin动作标识
func (p *AdminPermission) GetCasbinAction() string {
	if p.Method != nil {
		return *p.Method
	}
	return "*"
}

// TableName 指定表名
func (AdminPermission) TableName() string {
	return "admin_permissions"
}

// AdminDepartment 部门模型
type AdminDepartment struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	ParentID    uint           `json:"parent_id" gorm:"not null;default:0;comment:父级部门ID"`
	Name        string         `json:"name" gorm:"type:varchar(100);not null;comment:部门名称" validate:"required,max=100"`
	Code        *string        `json:"code" gorm:"type:varchar(50);comment:部门编码"`
	Level       int            `json:"level" gorm:"type:int;not null;default:1;comment:层级"`
	Path        *string        `json:"path" gorm:"type:varchar(500);comment:层级路径(如: 1,2,3)"`
	ManagerID   *uint          `json:"manager_id" gorm:"comment:部门负责人ID"`
	Phone       *string        `json:"phone" gorm:"type:varchar(20);comment:部门电话"`
	Email       *string        `json:"email" gorm:"type:varchar(100);comment:部门邮箱" validate:"omitempty,email"`
	Address     *string        `json:"address" gorm:"type:varchar(255);comment:部门地址"`
	SortOrder   int            `json:"sort_order" gorm:"type:int;not null;default:0;comment:排序"`
	Status      int8           `json:"status" gorm:"type:tinyint(1);not null;default:1;comment:状态: 0-禁用, 1-启用"`
	Description *string        `json:"description" gorm:"type:text;comment:部门描述"`
	CreatedAt   time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间(软删除)"`

	// 关联关系
	Children []AdminDepartment `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Parent   *AdminDepartment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Manager  *AdminUser        `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`
	Users    []AdminUser       `json:"users,omitempty" gorm:"foreignKey:DepartmentID"`
}

// TableName 指定表名
func (AdminDepartment) TableName() string {
	return "admin_departments"
}

// AdminLoginLog 登录日志模型
type AdminLoginLog struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID        *uint     `json:"user_id" gorm:"comment:用户ID"`
	Username      *string   `json:"username" gorm:"type:varchar(50);comment:用户名"`
	LoginType     int8      `json:"login_type" gorm:"type:tinyint(1);not null;default:1;comment:登录类型: 1-用户名, 2-手机号, 3-邮箱"`
	LoginMethod   string    `json:"login_method" gorm:"type:varchar(20);not null;default:'password';comment:登录方式: password-密码, sms-短信, qrcode-二维码"`
	IPAddress     string    `json:"ip_address" gorm:"type:varchar(45);not null;comment:IP地址"`
	UserAgent     *string   `json:"user_agent" gorm:"type:text;comment:用户代理"`
	DeviceType    *string   `json:"device_type" gorm:"type:varchar(20);comment:设备类型: web, mobile, tablet"`
	Browser       *string   `json:"browser" gorm:"type:varchar(50);comment:浏览器"`
	OS            *string   `json:"os" gorm:"type:varchar(50);comment:操作系统"`
	Location      *string   `json:"location" gorm:"type:varchar(100);comment:登录地点"`
	Status        int8      `json:"status" gorm:"type:tinyint(1);not null;default:1;comment:登录状态: 0-失败, 1-成功"`
	FailureReason *string   `json:"failure_reason" gorm:"type:varchar(255);comment:失败原因"`
	LoginTime     time.Time `json:"login_time" gorm:"not null;default:CURRENT_TIMESTAMP;comment:登录时间"`
	LogoutTime    *time.Time `json:"logout_time" gorm:"comment:退出时间"`

	// 关联关系
	User *AdminUser `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (AdminLoginLog) TableName() string {
	return "admin_login_logs"
}

// AdminOperationLog 操作日志模型
type AdminOperationLog struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID        *uint     `json:"user_id" gorm:"comment:操作用户ID"`
	Username      *string   `json:"username" gorm:"type:varchar(50);comment:操作用户名"`
	Module        string    `json:"module" gorm:"type:varchar(50);not null;comment:操作模块"`
	Action        string    `json:"action" gorm:"type:varchar(50);not null;comment:操作动作"`
	Description   *string   `json:"description" gorm:"type:varchar(255);comment:操作描述"`
	Method        string    `json:"method" gorm:"type:varchar(10);not null;comment:HTTP方法"`
	URL           string    `json:"url" gorm:"type:varchar(255);not null;comment:请求URL"`
	Params        *string   `json:"params" gorm:"type:text;comment:请求参数"`
	Result        *string   `json:"result" gorm:"type:text;comment:操作结果"`
	IPAddress     string    `json:"ip_address" gorm:"type:varchar(45);not null;comment:IP地址"`
	UserAgent     *string   `json:"user_agent" gorm:"type:text;comment:用户代理"`
	Status        int8      `json:"status" gorm:"type:tinyint(1);not null;default:1;comment:操作状态: 0-失败, 1-成功"`
	ErrorMessage  *string   `json:"error_message" gorm:"type:text;comment:错误信息"`
	ExecutionTime *int      `json:"execution_time" gorm:"type:int;comment:执行时间(毫秒)"`
	CreatedAt     time.Time `json:"created_at" gorm:"comment:操作时间"`

	// 关联关系
	User *AdminUser `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (AdminOperationLog) TableName() string {
	return "admin_operation_logs"
}

// AdminConfig 系统配置模型
type AdminConfig struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	GroupName   string    `json:"group_name" gorm:"type:varchar(50);not null;comment:配置分组" validate:"required,max=50"`
	ConfigKey   string    `json:"config_key" gorm:"type:varchar(100);not null;comment:配置键" validate:"required,max=100"`
	ConfigValue *string   `json:"config_value" gorm:"type:text;comment:配置值"`
	ConfigType  string    `json:"config_type" gorm:"type:varchar(20);not null;default:'string';comment:配置类型: string, int, bool, json"`
	Description *string   `json:"description" gorm:"type:varchar(255);comment:配置描述"`
	IsSystem    bool      `json:"is_system" gorm:"type:tinyint(1);not null;default:0;comment:是否系统配置: 0-否, 1-是"`
	SortOrder   int       `json:"sort_order" gorm:"type:int;not null;default:0;comment:排序"`
	CreatedAt   time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

// TableName 指定表名
func (AdminConfig) TableName() string {
	return "admin_configs"
}

// ===== 请求和响应结构体 =====

// AdminUserCreateRequest 创建管理员用户请求
type AdminUserCreateRequest struct {
	Username     string  `json:"username" validate:"required,min=3,max=50"`
	Phone        *string `json:"phone" validate:"omitempty,len=11"`
	Email        *string `json:"email" validate:"omitempty,email"`
	Password     string  `json:"password" validate:"required,min=6"`
	RealName     *string `json:"real_name" validate:"omitempty,max=50"`
	Nickname     *string `json:"nickname" validate:"omitempty,max=50"`
	Avatar       *string `json:"avatar"`
	Gender       int8    `json:"gender" validate:"omitempty,oneof=0 1 2"`
	Birthday     *string `json:"birthday"` // 前端传入字符串格式
	DepartmentID *uint   `json:"department_id"`
	Position     *string `json:"position" validate:"omitempty,max=100"`
	Status       int8    `json:"status" validate:"omitempty,oneof=0 1"`
	RoleIDs      []uint  `json:"role_ids"` // 角色ID列表
	Remark       *string `json:"remark"`
}

// AdminUserUpdateRequest 更新管理员用户请求
type AdminUserUpdateRequest struct {
	Phone        *string `json:"phone" validate:"omitempty,len=11"`
	Email        *string `json:"email" validate:"omitempty,email"`
	RealName     *string `json:"real_name" validate:"omitempty,max=50"`
	Nickname     *string `json:"nickname" validate:"omitempty,max=50"`
	Avatar       *string `json:"avatar"`
	Gender       *int8   `json:"gender" validate:"omitempty,oneof=0 1 2"`
	Birthday     *string `json:"birthday"` // 前端传入字符串格式
	DepartmentID *uint   `json:"department_id"`
	Position     *string `json:"position" validate:"omitempty,max=100"`
	Status       *int8   `json:"status" validate:"omitempty,oneof=0 1"`
	RoleIDs      []uint  `json:"role_ids"` // 角色ID列表
	Remark       *string `json:"remark"`
}

// AdminUserLoginRequest 管理员登录请求
type AdminUserLoginRequest struct {
	Account   string `json:"account" validate:"required"` // 用户名、手机号或邮箱
	Password  string `json:"password" validate:"required,min=6"`
	Captcha   string `json:"captcha"` // 验证码
	CaptchaID string `json:"captcha_id"` // 验证码ID
	Remember  bool   `json:"remember"` // 记住登录
}

// AdminUserChangePasswordRequest 修改密码请求
type AdminUserChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required,min=6"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// AdminProfileUpdateRequest 更新个人信息请求
type AdminProfileUpdateRequest struct {
	Nickname string `json:"nickname" validate:"omitempty,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"omitempty,len=11"`
	Avatar   string `json:"avatar" validate:"omitempty,url"`
}

// AdminUserListRequest 管理员用户列表请求
type AdminUserListRequest struct {
	Page         int     `json:"page" form:"page" validate:"omitempty,min=1"`
	PageSize     int     `json:"page_size" form:"page_size" validate:"omitempty,min=1,max=100"`
	Keyword      string  `json:"keyword" form:"keyword"`
	DepartmentID *uint   `json:"department_id" form:"department_id"`
	Status       *int8   `json:"status" form:"status" validate:"omitempty,oneof=0 1"`
	RoleID       *uint   `json:"role_id" form:"role_id"`
	StartTime    *string `json:"start_time" form:"start_time"`
	EndTime      *string `json:"end_time" form:"end_time"`
}

// AdminUserResponse 管理员用户响应
type AdminUserResponse struct {
	ID            uint              `json:"id"`
	Username      string            `json:"username"`
	Phone         *string           `json:"phone"`
	Email         *string           `json:"email"`
	RealName      *string           `json:"real_name"`
	Nickname      *string           `json:"nickname"`
	Avatar        *string           `json:"avatar"`
	Gender        int8              `json:"gender"`
	Birthday      *time.Time        `json:"birthday"`
	DepartmentID  *uint             `json:"department_id"`
	Position      *string           `json:"position"`
	Status        int8              `json:"status"`
	IsSuperAdmin  bool              `json:"is_super_admin"`
	LastLoginTime *time.Time        `json:"last_login_time"`
	LastLoginIP   *string           `json:"last_login_ip"`
	LoginCount    int               `json:"login_count"`
	Remark        *string           `json:"remark"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Department    *AdminDepartment  `json:"department,omitempty"`
	Roles         []AdminRole       `json:"roles,omitempty"`
}

// AdminLoginResponse 管理员登录响应
type AdminLoginResponse struct {
	Token     string              `json:"token"`
	ExpiresAt int64               `json:"expires_at"`
	User      *AdminUserResponse  `json:"user"`
}

// AdminRoleCreateRequest 创建角色请求
type AdminRoleCreateRequest struct {
	Name          string  `json:"name" validate:"required,max=50"`
	Code          string  `json:"code" validate:"required,max=50"`
	Description   *string `json:"description"`
	Level         int     `json:"level" validate:"omitempty,min=1"`
	Status        int8    `json:"status" validate:"omitempty,oneof=0 1"`
	SortOrder     int     `json:"sort_order"`
	PermissionIDs []uint  `json:"permission_ids"` // 权限ID列表
}

// AdminRoleUpdateRequest 更新角色请求
type AdminRoleUpdateRequest struct {
	Name          *string `json:"name" validate:"omitempty,max=50"`
	Description   *string `json:"description"`
	Level         *int    `json:"level" validate:"omitempty,min=1"`
	Status        *int8   `json:"status" validate:"omitempty,oneof=0 1"`
	SortOrder     *int    `json:"sort_order"`
	PermissionIDs []uint  `json:"permission_ids"` // 权限ID列表
}

// AdminPermissionCreateRequest 创建权限请求
type AdminPermissionCreateRequest struct {
	ParentID    uint    `json:"parent_id"`
	Name        string  `json:"name" validate:"required,max=100"`
	Code        string  `json:"code" validate:"required,max=100"`
	Type        int8    `json:"type" validate:"required,oneof=1 2 3"`
	Path        *string `json:"path"`
	Component   *string `json:"component"`
	Icon        *string `json:"icon"`
	Method      *string `json:"method"`
	URL         *string `json:"url"`
	SortOrder   int     `json:"sort_order"`
	Status      int8    `json:"status" validate:"omitempty,oneof=0 1"`
	IsHidden    bool    `json:"is_hidden"`
	Description *string `json:"description"`
}

// AdminDepartmentCreateRequest 创建部门请求
type AdminDepartmentCreateRequest struct {
	ParentID    uint    `json:"parent_id"`
	Name        string  `json:"name" validate:"required,max=100"`
	Code        *string `json:"code" validate:"omitempty,max=50"`
	ManagerID   *uint   `json:"manager_id"`
	Phone       *string `json:"phone"`
	Email       *string `json:"email" validate:"omitempty,email"`
	Address     *string `json:"address"`
	SortOrder   int     `json:"sort_order"`
	Status      int8    `json:"status" validate:"omitempty,oneof=0 1"`
	Description *string `json:"description"`
}



// AdminDepartmentUpdateRequest 更新部门请求
type AdminDepartmentUpdateRequest struct {
	Name        *string `json:"name" validate:"omitempty,max=50" comment:"部门名称"`
	Code        *string `json:"code" validate:"omitempty,max=50" comment:"部门编码"`
	ParentID    *uint   `json:"parent_id" comment:"父级部门ID"`
	Description *string `json:"description" validate:"omitempty,max=255" comment:"部门描述"`
	SortOrder   *int    `json:"sort_order" comment:"排序"`
	Status      *int8   `json:"status" validate:"omitempty,oneof=0 1" comment:"状态: 0-禁用, 1-启用"`
}

// FileInfo 文件信息模型
type FileInfo struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	OriginalName string         `json:"original_name" gorm:"type:varchar(255);not null;comment:原始文件名"`
	Filename     string         `json:"filename" gorm:"type:varchar(255);not null;comment:存储文件名"`
	FilePath     string         `json:"file_path" gorm:"type:varchar(500);not null;comment:文件路径"`
	FileSize     int64          `json:"file_size" gorm:"type:bigint;not null;comment:文件大小(字节)"`
	MimeType     string         `json:"mime_type" gorm:"type:varchar(100);comment:MIME类型"`
	FileExt      string         `json:"file_ext" gorm:"type:varchar(20);comment:文件扩展名"`
	URL          string         `json:"url" gorm:"type:varchar(500);comment:访问URL"`
	UploadedBy   *uint          `json:"uploaded_by" gorm:"comment:上传用户ID"`
	CreatedAt    time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index;comment:删除时间(软删除)"`

	// 关联关系
	Uploader *AdminUser `json:"uploader,omitempty" gorm:"foreignKey:UploadedBy"`
}

// TableName 指定表名
func (FileInfo) TableName() string {
	return "file_infos"
}



// 为了兼容性，添加简化的别名
type Department = AdminDepartment
type DepartmentCreateRequest = AdminDepartmentCreateRequest
type DepartmentUpdateRequest = AdminDepartmentUpdateRequest
type Role = AdminRole
type Permission = AdminPermission
type RoleCreateRequest = AdminRoleCreateRequest
type RoleUpdateRequest = AdminRoleUpdateRequest