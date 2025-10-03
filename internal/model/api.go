package model

import (
	"time"

	"gorm.io/gorm"
)

// API API模型
type API struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Group     string         `gorm:"size:50;not null" json:"group"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Method    string         `gorm:"size:10;not null" json:"method"`
	Path      string         `gorm:"size:100;not null" json:"path"`
	Status    int8           `gorm:"default:1" json:"status"` // 0:禁用 1:启用
	Remark    string         `gorm:"size:255" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (API) TableName() string {
	return "kp_api"
}
