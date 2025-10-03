package dto

// PageReq 分页请求
type PageReq struct {
	PageNum   int    `form:"page_num" binding:"required,min=1" example:"1"`
	PageSize  int    `form:"page_size" binding:"required,min=1,max=100" example:"10"`
	BeginTime string `form:"begin_time" example:"2023-01-01 00:00:00"`
	EndTime   string `form:"end_time" example:"2023-12-31 23:59:59"`
}

// PageResp 分页响应
type PageResp struct {
	List       interface{} `json:"list"`        // 数据列表
	Total      int64       `json:"total"`       // 总数
	PageNum    int         `json:"page_num"`    // 当前页码
	PageSize   int         `json:"page_size"`   // 每页数量
	TotalPages int         `json:"total_pages"` // 总页数
}

// CaptchaResp 验证码响应
type CaptchaResp struct {
	CaptchaID     string `json:"captcha_id"`     // 验证码ID
	CaptchaBase64 string `json:"captcha_base64"` // 验证码Base64
}

// IDReq ID请求
type IDReq struct {
	ID uint `uri:"id" binding:"required" example:"1"` // ID
}

// IDsReq IDs请求
type IDsReq struct {
	IDs []uint `json:"ids" binding:"required" example:"[1,2,3]"` // ID列表
}

// StatusReq 状态请求
type StatusReq struct {
	ID     uint `json:"id" binding:"required" example:"1"`     // ID
	Status int8 `json:"status" binding:"required" example:"1"` // 状态
}
