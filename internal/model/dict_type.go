package model

import (
	"time"

	"gorm.io/gorm"
)

// DictType 字典类型模型
type DictType struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Type      string         `gorm:"size:50;not null;uniqueIndex" json:"type"`
	Status    int8           `gorm:"default:1" json:"status"` // 0:禁用 1:启用
	Remark    string         `gorm:"size:255" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (DictType) TableName() string {
	return "kp_dict_type"
}
