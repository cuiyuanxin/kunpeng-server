package dto

// RoleCreateReq 创建角色请求
type RoleCreateReq struct {
	Name   string `json:"name" binding:"required" example:"测试角色"`
	Code   string `json:"code" binding:"required" example:"test"`
	Sort   int    `json:"sort" example:"0"`
	Status int8   `json:"status" example:"1"`
	Remark string `json:"remark" example:"测试角色"`
}

// RoleUpdateReq 更新角色请求
type RoleUpdateReq struct {
	ID     uint   `json:"id" binding:"required" example:"1"`
	Name   string `json:"name" binding:"required" example:"测试角色"`
	Code   string `json:"code" binding:"required" example:"test"`
	Sort   int    `json:"sort" example:"0"`
	Status int8   `json:"status" example:"1"`
	Remark string `json:"remark" example:"测试角色"`
}

// RolePageReq 角色分页请求
type RolePageReq struct {
	PageNum   int    `form:"page_num" binding:"required,min=1" example:"1"`
	PageSize  int    `form:"page_size" binding:"required,min=1,max=100" example:"10"`
	Name      string `form:"name" example:"测试角色"`
	Code      string `form:"code" example:"test"`
	Status    int8   `form:"status" example:"1"`
	BeginTime string `form:"begin_time" example:"2023-01-01 00:00:00"`
	EndTime   string `form:"end_time" example:"2023-12-31 23:59:59"`
}

// RoleMenuReq 角色菜单请求
type RoleMenuReq struct {
	RoleID  uint   `json:"role_id" binding:"required" example:"1"`
	MenuIDs []uint `json:"menu_ids" binding:"required" example:"[1,2,3]"`
}

// RoleAPIReq 角色API请求
type RoleAPIReq struct {
	RoleID uint   `json:"role_id" binding:"required" example:"1"`
	APIIDs []uint `json:"api_ids" binding:"required" example:"[1,2,3]"`
}
