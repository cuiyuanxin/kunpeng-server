package dto

// DeptCreateReq 创建部门请求
type DeptCreateReq struct {
	ParentID uint   `json:"parent_id" example:"0"`
	Name     string `json:"name" binding:"required" example:"技术部"`
	Leader   string `json:"leader" example:"张三"`
	Phone    string `json:"phone" example:"13800138000"`
	Email    string `json:"email" example:"tech@example.com"`
	Sort     int    `json:"sort" example:"0"`
	Status   int8   `json:"status" example:"1"`
	Remark   string `json:"remark" example:"技术部门"`
}

// DeptUpdateReq 更新部门请求
type DeptUpdateReq struct {
	ID       uint   `json:"id" binding:"required" example:"1"`
	ParentID uint   `json:"parent_id" example:"0"`
	Name     string `json:"name" binding:"required" example:"技术部"`
	Leader   string `json:"leader" example:"张三"`
	Phone    string `json:"phone" example:"13800138000"`
	Email    string `json:"email" example:"tech@example.com"`
	Sort     int    `json:"sort" example:"0"`
	Status   int8   `json:"status" example:"1"`
	Remark   string `json:"remark" example:"技术部门"`
}

// DeptQueryReq 部门查询请求
type DeptQueryReq struct {
	Name   string `form:"name" example:"技术部"`
	Status int8   `form:"status" example:"1"`
}

// DeptListReq 部门列表请求
type DeptListReq struct {
	Name   string `form:"name" example:"技术部"`
	Status int8   `form:"status" example:"1"`
}

// DeptTreeResp 部门树响应
type DeptTreeResp struct {
	ID       uint            `json:"id"`
	ParentID uint            `json:"parent_id"`
	Name     string          `json:"name"`
	Leader   string          `json:"leader"`
	Phone    string          `json:"phone"`
	Email    string          `json:"email"`
	Sort     int             `json:"sort"`
	Status   int8            `json:"status"`
	Children []*DeptTreeResp `json:"children"`
}
