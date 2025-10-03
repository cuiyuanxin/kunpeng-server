package model

import (
	"time"

	"gorm.io/gorm"
)

// DictData 字典数据模型
type DictData struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	DictType  string         `gorm:"size:50;not null;index" json:"dict_type"`
	Label     string         `gorm:"size:100;not null" json:"label"`
	Value     string         `gorm:"size:100;not null" json:"value"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Status    int8           `gorm:"default:1" json:"status"` // 0:禁用 1:启用
	Remark    string         `gorm:"size:255" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (DictData) TableName() string {
	return "kp_dict_data"
}
