package model

import (
	"time"
)

// OperationLog 操作日志模型
type OperationLog struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	Username     string    `gorm:"size:50" json:"username"`
	Module       string    `gorm:"size:50" json:"module"`
	Action       string    `gorm:"size:50" json:"action"`
	Method       string    `gorm:"size:10" json:"method"`
	Path         string    `gorm:"size:100" json:"path"`
	IP           string    `gorm:"size:50" json:"ip"`
	UserAgent    string    `gorm:"size:255" json:"user_agent"`
	Request      string    `gorm:"type:text" json:"request"`
	Response     string    `gorm:"type:text" json:"response"`
	Status       int8      `gorm:"default:1" json:"status"` // 0:失败 1:成功
	ErrorMessage string    `gorm:"type:text" json:"error_message"`
	Duration     int       `gorm:"default:0" json:"duration"` // 执行时长(ms)
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 表名
func (OperationLog) TableName() string {
	return "kp_operation_log"
}
