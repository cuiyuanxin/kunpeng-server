package dto

import (
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/model"
)

// OperationLogListReq 操作日志列表请求
type OperationLogListReq struct {
	PageReq
	UserID   *uint  `form:"user_id" json:"user_id"`   // 用户ID
	Username string `form:"username" json:"username"` // 用户名
	Module   string `form:"module" json:"module"`     // 模块名
	Action   string `form:"action" json:"action"`     // 操作动作
	Method   string `form:"method" json:"method"`     // 请求方法
	Status   *int8  `form:"status" json:"status"`     // 状态
	IP       string `form:"ip" json:"ip"`             // IP地址
}

// OperationLogResp 操作日志响应
type OperationLogResp struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`       // 用户ID
	Username     string    `json:"username"`      // 用户名
	Module       string    `json:"module"`        // 模块名
	Action       string    `json:"action"`        // 操作动作
	Method       string    `json:"method"`        // 请求方法
	Path         string    `json:"path"`          // 请求路径
	IP           string    `json:"ip"`            // IP地址
	UserAgent    string    `json:"user_agent"`    // 用户代理
	Request      string    `json:"request"`       // 请求参数
	Response     string    `json:"response"`      // 响应结果
	Status       int8      `json:"status"`        // 状态
	StatusName   string    `json:"status_name"`   // 状态名称
	ErrorMessage string    `json:"error_message"` // 错误信息
	Duration     int       `json:"duration"`      // 执行时长(ms)
	CreatedAt    time.Time `json:"created_at"`
}

// ToOperationLogResp 转换为操作日志响应
func ToOperationLogResp(log *model.OperationLog) *OperationLogResp {
	resp := &OperationLogResp{
		ID:           log.ID,
		UserID:       log.UserID,
		Username:     log.Username,
		Module:       log.Module,
		Action:       log.Action,
		Method:       log.Method,
		Path:         log.Path,
		IP:           log.IP,
		UserAgent:    log.UserAgent,
		Request:      log.Request,
		Response:     log.Response,
		Status:       log.Status,
		ErrorMessage: log.ErrorMessage,
		Duration:     log.Duration,
		CreatedAt:    log.CreatedAt,
	}

	// 设置状态名称
	if log.Status == 1 {
		resp.StatusName = "成功"
	} else {
		resp.StatusName = "失败"
	}

	return resp
}

// ToOperationLogRespList 转换为操作日志响应列表
func ToOperationLogRespList(logs []*model.OperationLog) []*OperationLogResp {
	result := make([]*OperationLogResp, 0, len(logs))
	for _, log := range logs {
		result = append(result, ToOperationLogResp(log))
	}
	return result
}
