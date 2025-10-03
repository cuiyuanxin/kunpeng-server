package dto

// APICreateReq 创建API请求
type APICreateReq struct {
	Group  string `json:"group" binding:"required" example:"用户管理"`
	Name   string `json:"name" binding:"required" example:"获取用户列表"`
	Method string `json:"method" binding:"required" example:"GET"`
	Path   string `json:"path" binding:"required" example:"/api/v1/users"`
	Status int8   `json:"status" example:"1"`
	Remark string `json:"remark" example:"获取用户列表API"`
}

// APIUpdateReq 更新API请求
type APIUpdateReq struct {
	ID     uint   `json:"id" binding:"required" example:"1"`
	Group  string `json:"group" binding:"required" example:"用户管理"`
	Name   string `json:"name" binding:"required" example:"获取用户列表"`
	Method string `json:"method" binding:"required" example:"GET"`
	Path   string `json:"path" binding:"required" example:"/api/v1/users"`
	Status int8   `json:"status" example:"1"`
	Remark string `json:"remark" example:"获取用户列表API"`
}

// APIPageReq API分页请求
type APIPageReq struct {
	PageNum   int    `form:"page_num" binding:"required,min=1" example:"1"`
	PageSize  int    `form:"page_size" binding:"required,min=1,max=100" example:"10"`
	Group     string `form:"group" example:"用户管理"`
	Name      string `form:"name" example:"获取用户列表"`
	Method    string `form:"method" example:"GET"`
	Path      string `form:"path" example:"/api/v1/users"`
	Status    int8   `form:"status" example:"1"`
	BeginTime string `form:"begin_time" example:"2023-01-01 00:00:00"`
	EndTime   string `form:"end_time" example:"2023-12-31 23:59:59"`
}

// APIRoleResp 角色API响应
type APIRoleResp struct {
	APIIDs []uint `json:"api_ids"`
}
