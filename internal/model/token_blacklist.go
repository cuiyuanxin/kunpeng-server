package model

import (
	"time"

	"gorm.io/gorm"
)

// TokenBlacklist token黑名单模型
type TokenBlacklist struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Token     string         `json:"token" gorm:"type:text;not null;comment:被拉黑的token"`
	UserID    uint           `json:"user_id" gorm:"not null;comment:用户ID"`
	Username  string         `json:"username" gorm:"type:varchar(50);not null;comment:用户名"`
	Reason    string         `json:"reason" gorm:"type:varchar(100);comment:拉黑原因"`
	ExpiresAt time.Time      `json:"expires_at" gorm:"not null;comment:token原始过期时间"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// TableName 指定表名
func (TokenBlacklist) TableName() string {
	return "kp_token_blacklist"
}
