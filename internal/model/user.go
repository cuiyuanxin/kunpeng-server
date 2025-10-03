package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Username  string         `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password  string         `gorm:"size:100;not null" json:"-"`
	Nickname  string         `gorm:"size:50" json:"nickname"`
	RealName  string         `gorm:"size:50" json:"real_name"`
	Avatar    string         `gorm:"size:255" json:"avatar"`
	Gender    int8           `gorm:"default:0" json:"gender"` // 0:未知 1:男 2:女
	Email     string         `gorm:"size:100;uniqueIndex" json:"email"`
	Mobile    string         `gorm:"size:20;uniqueIndex" json:"mobile"`
	DeptID    uint           `gorm:"index" json:"dept_id"`
	PostID    uint           `gorm:"index" json:"post_id"`
	RoleID    uint           `gorm:"index" json:"role_id"`
	Status    int8           `gorm:"default:1" json:"status"` // 0:禁用 1:启用
	LoginIP   string         `gorm:"size:50" json:"login_ip"`
	LoginTime *time.Time     `json:"login_time"`
	AppKey    string         `gorm:"size:100" json:"app_key"`
	AppSecret string         `gorm:"size:100" json:"-"`
	Remark    string         `gorm:"size:255" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (User) TableName() string {
	return "kp_user"
}
