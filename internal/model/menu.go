package model

import (
	"time"

	"gorm.io/gorm"
)

// Menu 菜单模型
type Menu struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	ParentID   uint           `gorm:"default:0" json:"parent_id"`
	Name       string         `gorm:"size:50;not null" json:"name"`
	Type       int8           `gorm:"default:0" json:"type"` // 0:目录 1:菜单 2:按钮
	Path       string         `gorm:"size:100" json:"path"`
	Component  string         `gorm:"size:100" json:"component"`
	Permission string         `gorm:"size:100" json:"permission"`
	Icon       string         `gorm:"size:100" json:"icon"`
	Sort       int            `gorm:"default:0" json:"sort"`
	Visible    int8           `gorm:"default:1" json:"visible"`  // 0:隐藏 1:显示
	Status     int8           `gorm:"default:1" json:"status"`   // 0:禁用 1:启用
	IsCache    int8           `gorm:"default:0" json:"is_cache"` // 0:不缓存 1:缓存
	IsFrame    int8           `gorm:"default:1" json:"is_frame"` // 0:外链 1:内部页面
	Remark     string         `gorm:"size:255" json:"remark"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Menu) TableName() string {
	return "kp_menu"
}
