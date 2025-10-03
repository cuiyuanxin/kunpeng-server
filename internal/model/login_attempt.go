package model

import (
	"time"

	"gorm.io/gorm"
)

// LoginAttempt 登录尝试记录模型
type LoginAttempt struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Account   string         `gorm:"size:100;not null;index" json:"account"` // 登录账号（用户名或手机号）
	IP        string         `gorm:"size:50;not null;index" json:"ip"`       // 登录IP
	Attempts  int            `gorm:"default:1" json:"attempts"`              // 失败次数
	LastTry   time.Time      `gorm:"not null" json:"last_try"`               // 最后尝试时间
	BlockedAt *time.Time     `json:"blocked_at"`                             // 拉黑时间
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (LoginAttempt) TableName() string {
	return "kp_login_attempt"
}

// IsBlocked 检查是否被拉黑
func (la *LoginAttempt) IsBlocked() bool {
	if la.BlockedAt == nil {
		return false
	}
	// 检查拉黑是否已过期（2小时）
	return time.Since(*la.BlockedAt) < 2*time.Hour
}

// ShouldBlock 检查是否应该拉黑
func (la *LoginAttempt) ShouldBlock() bool {
	return la.Attempts >= 5
}

// Block 拉黑账号
func (la *LoginAttempt) Block() {
	now := time.Now()
	la.BlockedAt = &now
}

// Reset 重置登录尝试记录
func (la *LoginAttempt) Reset() {
	la.Attempts = 0
	la.BlockedAt = nil
}
