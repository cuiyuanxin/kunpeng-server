package model

import (
	"time"
)

// RoleMenu 角色菜单关联模型
type RoleMenu struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	RoleID    uint      `gorm:"not null;index:idx_role_menu,unique" json:"role_id"`
	MenuID    uint      `gorm:"not null;index:idx_role_menu,unique" json:"menu_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 表名
func (RoleMenu) TableName() string {
	return "kp_role_menu"
}
