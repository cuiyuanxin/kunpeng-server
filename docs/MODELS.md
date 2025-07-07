# 数据模型系统文档

## 概述

本项目采用 GORM 作为 ORM 框架，提供了完整的数据模型系统。模型系统支持多种数据库类型，包含用户管理、权限控制、内容管理等核心功能模块。所有模型都遵循统一的设计规范，支持软删除、时间戳自动管理、数据验证等特性。

## 系统架构

### 模型系统组件

```
Model System
├── Core Models              # 核心模型
│   ├── User Model          # 用户模型
│   ├── Role Model          # 角色模型
│   └── Permission Model    # 权限模型
├── Content Models          # 内容模型
│   ├── Article Model       # 文章模型
│   ├── Category Model      # 分类模型
│   ├── Tag Model           # 标签模型
│   └── Comment Model       # 评论模型
├── System Models           # 系统模型
│   ├── File Model          # 文件模型
│   ├── Log Model           # 日志模型
│   └── Config Model        # 配置模型
├── Request/Response DTOs   # 数据传输对象
│   ├── Create Requests     # 创建请求
│   ├── Update Requests     # 更新请求
│   ├── List Requests       # 列表请求
│   └── Response Objects    # 响应对象
└── Model Utilities         # 模型工具
    ├── Validation          # 数据验证
    ├── Serialization       # 序列化
    └── Relationships       # 关联关系
```

## 核心模型定义

### 用户模型 (User)

<mcfile name="user.go" path="internal/model/user.go"></mcfile> 文件定义了完整的用户模型系统：

```go
// User 用户模型
type User struct {
    ID        uint           `json:"id" gorm:"primarykey"`
    Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null" validate:"required,min=3,max=50"`
    Email     string         `json:"email" gorm:"uniqueIndex;size:100;not null" validate:"required,email"`
    Password  string         `json:"-" gorm:"size:255;not null" validate:"required,min=6"`
    Nickname  string         `json:"nickname" gorm:"size:50"`
    Avatar    string         `json:"avatar" gorm:"size:255"`
    Phone     string         `json:"phone" gorm:"size:20"`
    Role      string         `json:"role" gorm:"size:20;default:user" validate:"oneof=admin user guest"`
    Status    int            `json:"status" gorm:"default:1;comment:1-正常 0-禁用"`
    LastLogin *time.Time     `json:"last_login"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName 指定表名
func (User) TableName() string {
    return "users"
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() UserResponse {
    return UserResponse{
        ID:        u.ID,
        Username:  u.Username,
        Email:     u.Email,
        Nickname:  u.Nickname,
        Avatar:    u.Avatar,
        Phone:     u.Phone,
        Role:      u.Role,
        Status:    u.Status,
        LastLogin: u.LastLogin,
        CreatedAt: u.CreatedAt,
        UpdatedAt: u.UpdatedAt,
    }
}
```

**字段说明**:
- `ID`: 主键，自增ID
- `Username`: 用户名，唯一索引，3-50字符
- `Email`: 邮箱，唯一索引，必须符合邮箱格式
- `Password`: 密码，JSON序列化时隐藏，最少6字符
- `Nickname`: 昵称，可选，最多50字符
- `Avatar`: 头像URL，可选
- `Phone`: 手机号，可选
- `Role`: 用户角色，枚举值：admin/user/guest
- `Status`: 状态，1-正常，0-禁用
- `LastLogin`: 最后登录时间
- `CreatedAt/UpdatedAt`: GORM自动管理的时间戳
- `DeletedAt`: 软删除时间戳

### 角色模型 (Role) - 建议实现

```go
// Role 角色模型
type Role struct {
    ID          uint           `json:"id" gorm:"primarykey"`
    Name        string         `json:"name" gorm:"uniqueIndex;size:50;not null" validate:"required,min=2,max=50"`
    DisplayName string         `json:"display_name" gorm:"size:100;not null" validate:"required,max=100"`
    Description string         `json:"description" gorm:"size:255"`
    Status      int            `json:"status" gorm:"default:1;comment:1-启用 0-禁用"`
    IsSystem    bool           `json:"is_system" gorm:"default:false;comment:是否为系统角色"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    
    // 关联关系
    Users       []User       `json:"users,omitempty" gorm:"many2many:user_roles;"`
    Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;"`
}

func (Role) TableName() string {
    return "roles"
}
```

### 权限模型 (Permission) - 建议实现

```go
// Permission 权限模型
type Permission struct {
    ID          uint           `json:"id" gorm:"primarykey"`
    Name        string         `json:"name" gorm:"uniqueIndex;size:100;not null" validate:"required,min=2,max=100"`
    DisplayName string         `json:"display_name" gorm:"size:100;not null" validate:"required,max=100"`
    Description string         `json:"description" gorm:"size:255"`
    Resource    string         `json:"resource" gorm:"size:50;not null" validate:"required,max=50"`
    Action      string         `json:"action" gorm:"size:50;not null" validate:"required,max=50"`
    Status      int            `json:"status" gorm:"default:1;comment:1-启用 0-禁用"`
    IsSystem    bool           `json:"is_system" gorm:"default:false;comment:是否为系统权限"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    
    // 关联关系
    Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permissions;"`
}

func (Permission) TableName() string {
    return "permissions"
}
```

## 内容管理模型

### 文章模型 (Article) - 建议实现

```go
// Article 文章模型
type Article struct {
    ID          uint           `json:"id" gorm:"primarykey"`
    Title       string         `json:"title" gorm:"size:200;not null" validate:"required,min=1,max=200"`
    Slug        string         `json:"slug" gorm:"uniqueIndex;size:200;not null"`
    Summary     string         `json:"summary" gorm:"size:500"`
    Content     string         `json:"content" gorm:"type:longtext;not null" validate:"required"`
    CoverImage  string         `json:"cover_image" gorm:"size:255"`
    Status      int            `json:"status" gorm:"default:1;comment:1-发布 2-草稿 3-下线"`
    ViewCount   int64          `json:"view_count" gorm:"default:0"`
    LikeCount   int64          `json:"like_count" gorm:"default:0"`
    CommentCount int64         `json:"comment_count" gorm:"default:0"`
    IsTop       bool           `json:"is_top" gorm:"default:false;comment:是否置顶"`
    IsRecommend bool           `json:"is_recommend" gorm:"default:false;comment:是否推荐"`
    PublishedAt *time.Time     `json:"published_at"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    
    // 外键
    AuthorID   uint `json:"author_id" gorm:"not null;index"`
    CategoryID uint `json:"category_id" gorm:"not null;index"`
    
    // 关联关系
    Author   User     `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
    Category Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
    Tags     []Tag    `json:"tags,omitempty" gorm:"many2many:article_tags;"`
    Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:ArticleID"`
}

func (Article) TableName() string {
    return "articles"
}
```

### 分类模型 (Category) - 建议实现

```go
// Category 分类模型
type Category struct {
    ID          uint           `json:"id" gorm:"primarykey"`
    Name        string         `json:"name" gorm:"uniqueIndex;size:50;not null" validate:"required,min=1,max=50"`
    Slug        string         `json:"slug" gorm:"uniqueIndex;size:50;not null"`
    Description string         `json:"description" gorm:"size:255"`
    Icon        string         `json:"icon" gorm:"size:100"`
    Color       string         `json:"color" gorm:"size:20"`
    Sort        int            `json:"sort" gorm:"default:0;comment:排序权重"`
    Status      int            `json:"status" gorm:"default:1;comment:1-启用 0-禁用"`
    ArticleCount int64         `json:"article_count" gorm:"default:0"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    
    // 自关联（支持父子分类）
    ParentID uint       `json:"parent_id" gorm:"default:0;index"`
    Parent   *Category  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
    Children []Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
    
    // 关联关系
    Articles []Article `json:"articles,omitempty" gorm:"foreignKey:CategoryID"`
}

func (Category) TableName() string {
    return "categories"
}
```

### 标签模型 (Tag) - 建议实现

```go
// Tag 标签模型
type Tag struct {
    ID          uint           `json:"id" gorm:"primarykey"`
    Name        string         `json:"name" gorm:"uniqueIndex;size:50;not null" validate:"required,min=1,max=50"`
    Slug        string         `json:"slug" gorm:"uniqueIndex;size:50;not null"`
    Description string         `json:"description" gorm:"size:255"`
    Color       string         `json:"color" gorm:"size:20"`
    Status      int            `json:"status" gorm:"default:1;comment:1-启用 0-禁用"`
    ArticleCount int64         `json:"article_count" gorm:"default:0"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    
    // 关联关系
    Articles []Article `json:"articles,omitempty" gorm:"many2many:article_tags;"`
}

func (Tag) TableName() string {
    return "tags"
}
```

### 评论模型 (Comment) - 建议实现

```go
// Comment 评论模型
type Comment struct {
    ID        uint           `json:"id" gorm:"primarykey"`
    Content   string         `json:"content" gorm:"type:text;not null" validate:"required,min=1,max=1000"`
    Status    int            `json:"status" gorm:"default:1;comment:1-正常 2-待审核 3-已删除"`
    LikeCount int64          `json:"like_count" gorm:"default:0"`
    IPAddress string         `json:"ip_address" gorm:"size:45"`
    UserAgent string         `json:"user_agent" gorm:"size:255"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
    
    // 外键
    ArticleID uint `json:"article_id" gorm:"not null;index"`
    UserID    uint `json:"user_id" gorm:"not null;index"`
    
    // 自关联（支持回复评论）
    ParentID uint      `json:"parent_id" gorm:"default:0;index"`
    Parent   *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
    Replies  []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID"`
    
    // 关联关系
    Article Article `json:"article,omitempty" gorm:"foreignKey:ArticleID"`
    User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Comment) TableName() string {
    return "comments"
}
```

## 系统管理模型

### 文件模型 (File) - 建议实现

```go
// File 文件模型
type File struct {
    ID          uint           `json:"id" gorm:"primarykey"`
    OriginalName string        `json:"original_name" gorm:"size:255;not null"`
    FileName    string         `json:"file_name" gorm:"size:255;not null"`
    FilePath    string         `json:"file_path" gorm:"size:500;not null"`
    FileSize    int64          `json:"file_size" gorm:"not null"`
    MimeType    string         `json:"mime_type" gorm:"size:100;not null"`
    Extension   string         `json:"extension" gorm:"size:20;not null"`
    Hash        string         `json:"hash" gorm:"size:64;uniqueIndex"`
    StorageType string         `json:"storage_type" gorm:"size:20;default:local;comment:存储类型:local,oss,cos等"`
    Status      int            `json:"status" gorm:"default:1;comment:1-正常 0-已删除"`
    DownloadCount int64        `json:"download_count" gorm:"default:0"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
    
    // 外键
    UploaderID uint `json:"uploader_id" gorm:"not null;index"`
    
    // 关联关系
    Uploader User `json:"uploader,omitempty" gorm:"foreignKey:UploaderID"`
}

func (File) TableName() string {
    return "files"
}
```

### 日志模型 (Log) - 建议实现

```go
// Log 操作日志模型
type Log struct {
    ID        uint      `json:"id" gorm:"primarykey"`
    Level     string    `json:"level" gorm:"size:20;not null;index"`
    Module    string    `json:"module" gorm:"size:50;not null;index"`
    Action    string    `json:"action" gorm:"size:100;not null"`
    Message   string    `json:"message" gorm:"type:text;not null"`
    Context   string    `json:"context" gorm:"type:json"`
    IPAddress string    `json:"ip_address" gorm:"size:45;index"`
    UserAgent string    `json:"user_agent" gorm:"size:255"`
    RequestID string    `json:"request_id" gorm:"size:100;index"`
    CreatedAt time.Time `json:"created_at"`
    
    // 外键（可选，匿名操作时为0）
    UserID uint `json:"user_id" gorm:"default:0;index"`
    
    // 关联关系
    User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (Log) TableName() string {
    return "logs"
}
```

### 配置模型 (Config) - 建议实现

```go
// Config 系统配置模型
type Config struct {
    ID          uint           `json:"id" gorm:"primarykey"`
    Key         string         `json:"key" gorm:"uniqueIndex;size:100;not null" validate:"required,min=1,max=100"`
    Value       string         `json:"value" gorm:"type:text"`
    Type        string         `json:"type" gorm:"size:20;default:string;comment:配置类型:string,int,bool,json"`
    Group       string         `json:"group" gorm:"size:50;default:system;index"`
    Description string         `json:"description" gorm:"size:255"`
    IsPublic    bool           `json:"is_public" gorm:"default:false;comment:是否为公开配置"`
    IsSystem    bool           `json:"is_system" gorm:"default:false;comment:是否为系统配置"`
    Sort        int            `json:"sort" gorm:"default:0"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

func (Config) TableName() string {
    return "configs"
}
```

## 数据传输对象 (DTOs)

### 用户相关 DTOs

当前项目已实现的用户相关 DTOs：

```go
// UserCreateRequest 用户创建请求
type UserCreateRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
    Nickname string `json:"nickname" validate:"max=50"`
    Phone    string `json:"phone" validate:"max=20"`
    Role     string `json:"role" validate:"oneof=admin user guest"`
}

// UserUpdateRequest 用户更新请求
type UserUpdateRequest struct {
    Nickname string `json:"nickname" validate:"max=50"`
    Avatar   string `json:"avatar" validate:"max=255"`
    Phone    string `json:"phone" validate:"max=20"`
    Role     string `json:"role" validate:"oneof=admin user guest"`
    Status   *int   `json:"status" validate:"oneof=0 1"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
    Username string `json:"username" validate:"required"`
    Password string `json:"password" validate:"required"`
}

// UserChangePasswordRequest 用户修改密码请求
type UserChangePasswordRequest struct {
    OldPassword string `json:"old_password" validate:"required"`
    NewPassword string `json:"new_password" validate:"required,min=6"`
}

// UserListRequest 用户列表请求
type UserListRequest struct {
    Page     int    `form:"page" validate:"min=1"`
    PageSize int    `form:"page_size" validate:"min=1,max=100"`
    Keyword  string `form:"keyword"`
    Role     string `form:"role" validate:"oneof=admin user guest ''"`
    Status   *int   `form:"status" validate:"oneof=0 1"`
}

// UserResponse 用户响应
type UserResponse struct {
    ID        uint       `json:"id"`
    Username  string     `json:"username"`
    Email     string     `json:"email"`
    Nickname  string     `json:"nickname"`
    Avatar    string     `json:"avatar"`
    Phone     string     `json:"phone"`
    Role      string     `json:"role"`
    Status    int        `json:"status"`
    LastLogin *time.Time `json:"last_login"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
}

// LoginResponse 登录响应
type LoginResponse struct {
    Token string       `json:"token"`
    User  UserResponse `json:"user"`
}
```

### 其他模型 DTOs - 建议实现

```go
// 角色相关 DTOs
type RoleCreateRequest struct {
    Name        string `json:"name" validate:"required,min=2,max=50"`
    DisplayName string `json:"display_name" validate:"required,max=100"`
    Description string `json:"description" validate:"max=255"`
    Status      int    `json:"status" validate:"oneof=0 1"`
}

type RoleUpdateRequest struct {
    DisplayName string `json:"display_name" validate:"max=100"`
    Description string `json:"description" validate:"max=255"`
    Status      *int   `json:"status" validate:"oneof=0 1"`
}

type RoleResponse struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`
    DisplayName string    `json:"display_name"`
    Description string    `json:"description"`
    Status      int       `json:"status"`
    IsSystem    bool      `json:"is_system"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

// 文章相关 DTOs
type ArticleCreateRequest struct {
    Title      string `json:"title" validate:"required,min=1,max=200"`
    Summary    string `json:"summary" validate:"max=500"`
    Content    string `json:"content" validate:"required"`
    CoverImage string `json:"cover_image" validate:"max=255"`
    CategoryID uint   `json:"category_id" validate:"required"`
    TagIDs     []uint `json:"tag_ids"`
    Status     int    `json:"status" validate:"oneof=1 2 3"`
    IsTop      bool   `json:"is_top"`
    IsRecommend bool  `json:"is_recommend"`
}

type ArticleUpdateRequest struct {
    Title       string `json:"title" validate:"min=1,max=200"`
    Summary     string `json:"summary" validate:"max=500"`
    Content     string `json:"content"`
    CoverImage  string `json:"cover_image" validate:"max=255"`
    CategoryID  uint   `json:"category_id"`
    TagIDs      []uint `json:"tag_ids"`
    Status      *int   `json:"status" validate:"oneof=1 2 3"`
    IsTop       *bool  `json:"is_top"`
    IsRecommend *bool  `json:"is_recommend"`
}

type ArticleListRequest struct {
    Page       int    `form:"page" validate:"min=1"`
    PageSize   int    `form:"page_size" validate:"min=1,max=100"`
    Keyword    string `form:"keyword"`
    CategoryID uint   `form:"category_id"`
    TagID      uint   `form:"tag_id"`
    AuthorID   uint   `form:"author_id"`
    Status     *int   `form:"status" validate:"oneof=1 2 3"`
    IsTop      *bool  `form:"is_top"`
    IsRecommend *bool `form:"is_recommend"`
}

type ArticleResponse struct {
    ID           uint         `json:"id"`
    Title        string       `json:"title"`
    Slug         string       `json:"slug"`
    Summary      string       `json:"summary"`
    Content      string       `json:"content,omitempty"`
    CoverImage   string       `json:"cover_image"`
    Status       int          `json:"status"`
    ViewCount    int64        `json:"view_count"`
    LikeCount    int64        `json:"like_count"`
    CommentCount int64        `json:"comment_count"`
    IsTop        bool         `json:"is_top"`
    IsRecommend  bool         `json:"is_recommend"`
    PublishedAt  *time.Time   `json:"published_at"`
    CreatedAt    time.Time    `json:"created_at"`
    UpdatedAt    time.Time    `json:"updated_at"`
    Author       UserResponse `json:"author,omitempty"`
    Category     CategoryResponse `json:"category,omitempty"`
    Tags         []TagResponse `json:"tags,omitempty"`
}
```

## 模型设计规范

### 1. 命名规范

```go
// 模型命名：使用大驼峰命名法
type User struct {}
type UserRole struct {}

// 字段命名：使用大驼峰命名法
type User struct {
    ID        uint   `json:"id"`
    Username  string `json:"username"`
    CreatedAt time.Time `json:"created_at"`
}

// 表名：使用复数形式，小写+下划线
func (User) TableName() string {
    return "users"
}

func (UserRole) TableName() string {
    return "user_roles"
}
```

### 2. 标签规范

```go
type User struct {
    // GORM 标签
    ID       uint   `gorm:"primarykey"`                    // 主键
    Username string `gorm:"uniqueIndex;size:50;not null"` // 唯一索引，长度限制，非空
    Email    string `gorm:"uniqueIndex;size:100;not null"` // 唯一索引
    Status   int    `gorm:"default:1;comment:状态说明"`      // 默认值，注释
    
    // JSON 标签
    Password  string `json:"-"`          // 隐藏字段
    CreatedAt time.Time `json:"created_at"` // 自定义JSON字段名
    
    // 验证标签
    Username string `validate:"required,min=3,max=50"`     // 必填，长度限制
    Email    string `validate:"required,email"`            // 必填，邮箱格式
    Role     string `validate:"oneof=admin user guest"`    // 枚举值
}
```

### 3. 时间戳管理

```go
type BaseModel struct {
    ID        uint           `json:"id" gorm:"primarykey"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // 软删除
}

// 继承基础模型
type User struct {
    BaseModel
    Username string `json:"username" gorm:"uniqueIndex;size:50;not null"`
    // 其他字段...
}
```

### 4. 关联关系

```go
// 一对一关系
type User struct {
    ID      uint    `gorm:"primarykey"`
    Profile Profile `gorm:"foreignKey:UserID"`
}

type Profile struct {
    ID     uint `gorm:"primarykey"`
    UserID uint `gorm:"uniqueIndex"`
    Bio    string
}

// 一对多关系
type User struct {
    ID       uint      `gorm:"primarykey"`
    Articles []Article `gorm:"foreignKey:AuthorID"`
}

type Article struct {
    ID       uint `gorm:"primarykey"`
    AuthorID uint `gorm:"not null;index"`
    Author   User `gorm:"foreignKey:AuthorID"`
}

// 多对多关系
type User struct {
    ID    uint   `gorm:"primarykey"`
    Roles []Role `gorm:"many2many:user_roles;"`
}

type Role struct {
    ID    uint   `gorm:"primarykey"`
    Users []User `gorm:"many2many:user_roles;"`
}
```

### 5. 索引设计

```go
type User struct {
    ID       uint   `gorm:"primarykey"`                    // 主键索引
    Username string `gorm:"uniqueIndex;size:50;not null"` // 唯一索引
    Email    string `gorm:"uniqueIndex;size:100;not null"` // 唯一索引
    Status   int    `gorm:"index"`                         // 普通索引
    Role     string `gorm:"index"`                         // 普通索引
    
    // 复合索引
    CreatedAt time.Time      `gorm:"index:idx_created_status"`
    Status    int            `gorm:"index:idx_created_status"`
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

## 模型工具方法

### 1. 响应转换

```go
// 单个模型转换
func (u *User) ToResponse() UserResponse {
    return UserResponse{
        ID:        u.ID,
        Username:  u.Username,
        Email:     u.Email,
        // 其他字段...
    }
}

// 批量转换
func UsersToResponse(users []User) []UserResponse {
    responses := make([]UserResponse, len(users))
    for i, user := range users {
        responses[i] = user.ToResponse()
    }
    return responses
}

// 带关联关系的转换
func (u *User) ToDetailResponse() UserDetailResponse {
    response := UserDetailResponse{
        UserResponse: u.ToResponse(),
    }
    
    // 转换关联数据
    if len(u.Articles) > 0 {
        response.Articles = ArticlesToResponse(u.Articles)
    }
    
    return response
}
```

### 2. 数据验证

```go
// 创建前验证
func (u *User) BeforeCreate(tx *gorm.DB) error {
    // 密码加密
    if u.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.Password = string(hashedPassword)
    }
    
    // 生成默认昵称
    if u.Nickname == "" {
        u.Nickname = u.Username
    }
    
    return nil
}

// 更新前验证
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    // 如果密码被修改，重新加密
    if tx.Statement.Changed("Password") && u.Password != "" {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        u.Password = string(hashedPassword)
    }
    
    return nil
}

// 自定义验证方法
func (u *User) ValidatePassword(password string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
    return err == nil
}
```

### 3. 查询作用域

```go
// 状态过滤
func (u *User) ScopeActive(db *gorm.DB) *gorm.DB {
    return db.Where("status = ?", 1)
}

func (u *User) ScopeByRole(role string) func(db *gorm.DB) *gorm.DB {
    return func(db *gorm.DB) *gorm.DB {
        if role != "" {
            return db.Where("role = ?", role)
        }
        return db
    }
}

// 使用作用域
var users []User
db.Scopes(user.ScopeActive, user.ScopeByRole("admin")).Find(&users)
```

## 数据库迁移

### 模型注册

在 <mcfile name="migrate.go" path="internal/migrate/migrate.go"></mcfile> 中注册所有模型：

```go
// AllModels 所有需要迁移的模型
var AllModels = []interface{}{
    &model.User{},
    &model.Role{},
    &model.Permission{},
    &model.Article{},
    &model.Category{},
    &model.Tag{},
    &model.Comment{},
    &model.File{},
    &model.Log{},
    &model.Config{},
}
```

### 自动迁移

```go
// AutoMigrate 自动迁移所有数据库表
func AutoMigrate() error {
    if database.DB == nil {
        return fmt.Errorf("database not initialized")
    }

    for _, model := range AllModels {
        if err := database.DB.AutoMigrate(model); err != nil {
            modelName := reflect.TypeOf(model).Elem().Name()
            return fmt.Errorf("failed to migrate model %s: %w", modelName, err)
        }
    }

    klogger.Logger.Info("Database migration completed successfully")
    return nil
}
```

## 使用示例

### 基础 CRUD 操作

```go
package main

import (
    "github.com/your-project/internal/model"
    "github.com/your-project/internal/database"
)

func main() {
    // 创建用户
    user := &model.User{
        Username: "john_doe",
        Email:    "john@example.com",
        Password: "password123",
        Role:     "user",
    }
    
    result := database.DB.Create(user)
    if result.Error != nil {
        log.Printf("Failed to create user: %v", result.Error)
        return
    }
    
    // 查询用户
    var foundUser model.User
    database.DB.Where("username = ?", "john_doe").First(&foundUser)
    
    // 更新用户
    database.DB.Model(&foundUser).Update("nickname", "John")
    
    // 删除用户（软删除）
    database.DB.Delete(&foundUser)
    
    // 查询未删除的用户
    var activeUsers []model.User
    database.DB.Where("status = ?", 1).Find(&activeUsers)
}
```

### 关联查询

```go
// 预加载关联数据
var articles []model.Article
database.DB.Preload("Author").Preload("Category").Preload("Tags").Find(&articles)

// 条件预加载
database.DB.Preload("Comments", "status = ?", 1).Find(&articles)

// 嵌套预加载
database.DB.Preload("Comments.User").Find(&articles)

// 自定义预加载
database.DB.Preload("Comments", func(db *gorm.DB) *gorm.DB {
    return db.Order("created_at DESC").Limit(10)
}).Find(&articles)
```

### 复杂查询

```go
// 分页查询
func GetUserList(req model.UserListRequest) ([]model.User, int64, error) {
    var users []model.User
    var total int64
    
    query := database.DB.Model(&model.User{})
    
    // 条件过滤
    if req.Keyword != "" {
        query = query.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?", 
            "%"+req.Keyword+"%", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
    }
    
    if req.Role != "" {
        query = query.Where("role = ?", req.Role)
    }
    
    if req.Status != nil {
        query = query.Where("status = ?", *req.Status)
    }
    
    // 获取总数
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // 分页查询
    offset := (req.Page - 1) * req.PageSize
    if err := query.Offset(offset).Limit(req.PageSize).Order("created_at DESC").Find(&users).Error; err != nil {
        return nil, 0, err
    }
    
    return users, total, nil
}
```

### 事务操作

```go
// 创建文章并更新分类计数
func CreateArticle(req model.ArticleCreateRequest, authorID uint) error {
    return database.DB.Transaction(func(tx *gorm.DB) error {
        // 创建文章
        article := &model.Article{
            Title:      req.Title,
            Content:    req.Content,
            AuthorID:   authorID,
            CategoryID: req.CategoryID,
            Status:     req.Status,
        }
        
        if err := tx.Create(article).Error; err != nil {
            return err
        }
        
        // 更新分类文章计数
        if err := tx.Model(&model.Category{}).Where("id = ?", req.CategoryID).
            UpdateColumn("article_count", gorm.Expr("article_count + ?", 1)).Error; err != nil {
            return err
        }
        
        // 关联标签
        if len(req.TagIDs) > 0 {
            var tags []model.Tag
            if err := tx.Where("id IN ?", req.TagIDs).Find(&tags).Error; err != nil {
                return err
            }
            
            if err := tx.Model(article).Association("Tags").Append(tags); err != nil {
                return err
            }
        }
        
        return nil
    })
}
```

## 性能优化

### 1. 查询优化

```go
// 使用索引
db.Where("username = ?", username).First(&user)  // 利用唯一索引
db.Where("status = ? AND role = ?", 1, "admin")   // 利用复合索引

// 选择特定字段
db.Select("id, username, email").Find(&users)

// 避免 N+1 查询
db.Preload("Articles").Find(&users)  // 预加载

// 批量操作
db.CreateInBatches(users, 100)  // 批量创建
```

### 2. 连接池优化

```go
// 配置连接池
sqlDB, err := db.DB()
if err != nil {
    return err
}

sqlDB.SetMaxOpenConns(100)                // 最大打开连接数
sqlDB.SetMaxIdleConns(10)                 // 最大空闲连接数
sqlDB.SetConnMaxLifetime(time.Hour)       // 连接最大生命周期
sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 连接最大空闲时间
```

### 3. 缓存策略

```go
// Redis 缓存用户信息
func GetUserByID(id uint) (*model.User, error) {
    cacheKey := fmt.Sprintf("user:%d", id)
    
    // 尝试从缓存获取
    cached, err := redis.GetClient().Get(context.Background(), cacheKey).Result()
    if err == nil {
        var user model.User
        if err := json.Unmarshal([]byte(cached), &user); err == nil {
            return &user, nil
        }
    }
    
    // 从数据库查询
    var user model.User
    if err := database.DB.First(&user, id).Error; err != nil {
        return nil, err
    }
    
    // 写入缓存
    userJSON, _ := json.Marshal(user)
    redis.GetClient().Set(context.Background(), cacheKey, userJSON, time.Hour)
    
    return &user, nil
}
```

## 测试

### 模型测试示例

```go
package model

import (
    "testing"
    "time"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("failed to connect database")
    }
    
    // 自动迁移
    db.AutoMigrate(&User{}, &Role{}, &Article{})
    
    return db
}

func TestUserModel(t *testing.T) {
    db := setupTestDB()
    
    t.Run("Create User", func(t *testing.T) {
        user := &User{
            Username: "testuser",
            Email:    "test@example.com",
            Password: "password123",
            Role:     "user",
        }
        
        err := db.Create(user).Error
        require.NoError(t, err)
        assert.NotZero(t, user.ID)
        assert.NotZero(t, user.CreatedAt)
    })
    
    t.Run("User Validation", func(t *testing.T) {
        user := &User{
            Username: "test",
            Email:    "invalid-email",
            Password: "123", // 太短
        }
        
        // 这里需要集成验证器进行测试
        // err := validator.Validate(user)
        // assert.Error(t, err)
    })
    
    t.Run("User Response Conversion", func(t *testing.T) {
        user := &User{
            ID:       1,
            Username: "testuser",
            Email:    "test@example.com",
            Password: "hashed_password",
            Role:     "user",
            Status:   1,
            CreatedAt: time.Now(),
            UpdatedAt: time.Now(),
        }
        
        response := user.ToResponse()
        
        assert.Equal(t, user.ID, response.ID)
        assert.Equal(t, user.Username, response.Username)
        assert.Equal(t, user.Email, response.Email)
        assert.Equal(t, user.Role, response.Role)
        // 密码不应该在响应中
        assert.Empty(t, response.Password)
    })
}

func TestUserRelationships(t *testing.T) {
    db := setupTestDB()
    
    // 创建用户
    user := &User{
        Username: "author",
        Email:    "author@example.com",
        Password: "password123",
        Role:     "user",
    }
    db.Create(user)
    
    // 创建文章
    article := &Article{
        Title:    "Test Article",
        Content:  "This is a test article",
        AuthorID: user.ID,
        Status:   1,
    }
    db.Create(article)
    
    // 测试关联查询
    var foundUser User
    err := db.Preload("Articles").First(&foundUser, user.ID).Error
    require.NoError(t, err)
    
    assert.Len(t, foundUser.Articles, 1)
    assert.Equal(t, "Test Article", foundUser.Articles[0].Title)
}
```

## 最佳实践

### 1. 模型设计原则

- **单一职责**: 每个模型只负责一个业务实体
- **数据完整性**: 使用数据库约束和验证确保数据完整性
- **性能考虑**: 合理设计索引，避免过度关联
- **扩展性**: 预留扩展字段，支持业务发展

### 2. 安全考虑

- **敏感数据**: 密码等敏感字段使用 `json:"-"` 隐藏
- **SQL注入**: 使用参数化查询，避免字符串拼接
- **权限控制**: 在模型层实现基础权限检查
- **数据验证**: 严格验证输入数据

### 3. 维护性

- **文档完整**: 为每个模型和字段添加注释
- **版本控制**: 使用数据库迁移管理模型变更
- **测试覆盖**: 为模型编写完整的测试用例
- **代码规范**: 遵循统一的命名和编码规范

## 相关文档

- [数据库系统完整指南](DATABASE_GUIDE.md)
- [数据库迁移文档](MIGRATION.md)
- [配置管理文档](CONFIG.md)
- [统一响应格式文档](RESPONSE.md)
- [GORM 官方文档](https://gorm.io/docs/)
- [Go Validator 文档](https://github.com/go-playground/validator)

---

**最佳实践**: 遵循 RESTful 设计原则；使用 DTO 分离内部模型和外部接口；实现完整的数据验证；合理设计数据库索引；使用事务确保数据一致性；编写完整的测试用例；定期进行性能优化。