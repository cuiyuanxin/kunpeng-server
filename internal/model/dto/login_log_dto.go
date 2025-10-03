package dto

import (
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// LoginLogListReq 登录日志列表请求
type LoginLogListReq struct {
	PageNum   int    `form:"page_num" binding:"required,min=1" example:"1"`
	PageSize  int    `form:"page_size" binding:"required,min=1,max=100" example:"10"`
	UserID    *uint  `form:"user_id" json:"user_id"`   // 用户ID
	Username  string `form:"username" json:"username"` // 用户名
	Status    *int8  `form:"status" json:"status"`     // 状态
	IP        string `form:"ip" json:"ip"`             // IP地址
	BeginTime string `form:"begin_time" example:"2023-01-01 00:00:00"`
	EndTime   string `form:"end_time" example:"2023-12-31 23:59:59"`
}

// LoginLogResp 登录日志响应
type LoginLogResp struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`     // 用户ID
	Username   string    `json:"username"`    // 用户名
	IP         string    `json:"ip"`          // IP地址
	Location   string    `json:"location"`    // 地理位置
	Browser    string    `json:"browser"`     // 浏览器
	OS         string    `json:"os"`          // 操作系统
	Device     string    `json:"device"`      // 设备
	Status     int8      `json:"status"`      // 状态
	StatusName string    `json:"status_name"` // 状态名称
	Message    string    `json:"message"`     // 消息描述
	CreatedAt  time.Time `json:"created_at"`
}

// ToLoginLogResp 转换为登录日志响应
func ToLoginLogResp(log *model.LoginLog) *LoginLogResp {
	resp := &LoginLogResp{
		ID:        log.ID,
		UserID:    log.UserID,
		Username:  log.Username,
		IP:        log.IP,
		Location:  log.Location,
		Browser:   log.Browser,
		OS:        log.OS,
		Device:    log.Device,
		Status:    log.Status,
		Message:   log.Message,
		CreatedAt: log.CreatedAt,
	}

	// 设置状态名称
	if log.Status == 1 {
		resp.StatusName = "成功"
	} else {
		resp.StatusName = "失败"
	}

	return resp
}

// ToLoginLogRespList 转换为登录日志响应列表
func ToLoginLogRespList(logs []*model.LoginLog) []*LoginLogResp {
	result := make([]*LoginLogResp, 0, len(logs))
	for _, log := range logs {
		result = append(result, ToLoginLogResp(log))
	}
	return result
}
