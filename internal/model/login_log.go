package model

import (
	"time"
)

// LoginLog 登录日志模型
type LoginLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"size:50" json:"username"`
	IP        string    `gorm:"size:50" json:"ip"`
	Location  string    `gorm:"size:100" json:"location"`
	Browser   string    `gorm:"size:50" json:"browser"`
	OS        string    `gorm:"size:50" json:"os"`
	Device    string    `gorm:"size:50" json:"device"`
	Status    int8      `gorm:"default:1" json:"status"` // 0:失败 1:成功
	Message   string    `gorm:"size:255" json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (LoginLog) TableName() string {
	return "kp_login_log"
}
