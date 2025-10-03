package model

import (
	"time"

	"gorm.io/gorm"
)

// Config 系统配置模型
type Config struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Key       string         `gorm:"size:100;not null;uniqueIndex" json:"key"`
	Value     string         `gorm:"type:text;not null" json:"value"`
	Type      string         `gorm:"size:20;default:string" json:"type"`
	Group     string         `gorm:"size:50;default:default" json:"group"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Status    int8           `gorm:"default:1" json:"status"` // 0:禁用 1:启用
	Remark    string         `gorm:"size:255" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Config) TableName() string {
	return "kp_config"
}
