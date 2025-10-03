package model

import (
	"time"

	"gorm.io/gorm"
)

// Dept 部门模型
type Dept struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	ParentID  uint           `gorm:"default:0" json:"parent_id"`
	Name      string         `gorm:"size:50;not null" json:"name"`
	Leader    string         `gorm:"size:50" json:"leader"`
	Phone     string         `gorm:"size:20" json:"phone"`
	Email     string         `gorm:"size:100" json:"email"`
	Sort      int            `gorm:"default:0" json:"sort"`
	Status    int8           `gorm:"default:1" json:"status"` // 0:禁用 1:启用
	Remark    string         `gorm:"size:255" json:"remark"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 表名
func (Dept) TableName() string {
	return "kp_dept"
}
