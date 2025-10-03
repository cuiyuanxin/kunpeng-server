package model

import (
	"time"
)

// RoleAPI 角色API关联模型
type RoleAPI struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	RoleID    uint      `gorm:"not null;index:idx_role_api,unique" json:"role_id"`
	APIID     uint      `gorm:"not null;index:idx_role_api,unique" json:"api_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (RoleAPI) TableName() string {
	return "kp_role_api"
}
