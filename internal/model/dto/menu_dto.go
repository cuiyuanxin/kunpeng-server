package dto

// MenuCreateReq 创建菜单请求
type MenuCreateReq struct {
	ParentID   uint   `json:"parent_id" example:"0"`
	Name       string `json:"name" binding:"required" example:"系统管理"`
	Type       int8   `json:"type" binding:"required" example:"0"`
	Path       string `json:"path" example:"/system"`
	Component  string `json:"component" example:"Layout"`
	Permission string `json:"permission" example:"system"`
	Icon       string `json:"icon" example:"setting"`
	Sort       int    `json:"sort" example:"0"`
	Visible    int8   `json:"visible" example:"1"`
	Status     int8   `json:"status" example:"1"`
	Remark     string `json:"remark" example:"系统管理菜单"`
}

// MenuUpdateReq 更新菜单请求
type MenuUpdateReq struct {
	ID         uint   `json:"id" binding:"required" example:"1"`
	ParentID   uint   `json:"parent_id" example:"0"`
	Name       string `json:"name" binding:"required" example:"系统管理"`
	Type       int8   `json:"type" binding:"required" example:"0"`
	Path       string `json:"path" example:"/system"`
	Component  string `json:"component" example:"Layout"`
	Permission string `json:"permission" example:"system"`
	Icon       string `json:"icon" example:"setting"`
	Sort       int    `json:"sort" example:"0"`
	Visible    int8   `json:"visible" example:"1"`
	Status     int8   `json:"status" example:"1"`
	Remark     string `json:"remark" example:"系统管理菜单"`
}

// MenuQueryReq 菜单查询请求
type MenuQueryReq struct {
	Name   string `form:"name" example:"系统管理"`
	Status int8   `form:"status" example:"1"`
}

// MenuListReq 菜单列表请求
type MenuListReq struct {
	Name   string `form:"name" example:"系统管理"`
	Status int8   `form:"status" example:"1"`
}

// MenuTreeResp 菜单树响应
type MenuTreeResp struct {
	ID         uint            `json:"id"`
	ParentID   uint            `json:"parent_id"`
	Name       string          `json:"name"`
	Type       int8            `json:"type"`
	Path       string          `json:"path"`
	Component  string          `json:"component"`
	Permission string          `json:"permission"`
	Icon       string          `json:"icon"`
	Sort       int             `json:"sort"`
	Visible    int8            `json:"visible"`
	Status     int8            `json:"status"`
	IsCache    int8            `json:"is_cache"`
	IsFrame    int8            `json:"is_frame"`
	Children   []*MenuTreeResp `json:"children"`
}

// MenuRoleResp 角色菜单响应
type MenuRoleResp struct {
	MenuIDs []uint `json:"menu_ids"`
}
