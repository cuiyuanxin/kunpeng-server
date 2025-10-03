package dto

// PostCreateReq 创建岗位请求
type PostCreateReq struct {
	Name   string `json:"name" binding:"required" example:"技术总监"`
	Code   string `json:"code" binding:"required" example:"CTO"`
	Sort   int    `json:"sort" example:"0"`
	Status int8   `json:"status" example:"1"`
	Remark string `json:"remark" example:"技术总监岗位"`
}

// PostUpdateReq 更新岗位请求
type PostUpdateReq struct {
	ID     uint   `json:"id" binding:"required" example:"1"`
	Name   string `json:"name" binding:"required" example:"技术总监"`
	Code   string `json:"code" binding:"required" example:"CTO"`
	Sort   int    `json:"sort" example:"0"`
	Status int8   `json:"status" example:"1"`
	Remark string `json:"remark" example:"技术总监岗位"`
}

// PostPageReq 岗位分页请求
type PostPageReq struct {
	PageNum   int    `form:"page_num" binding:"required,min=1" example:"1"`
	PageSize  int    `form:"page_size" binding:"required,min=1,max=100" example:"10"`
	Name      string `form:"name" example:"技术总监"`
	Code      string `form:"code" example:"CTO"`
	Status    int8   `form:"status" example:"1"`
	BeginTime string `form:"begin_time" example:"2023-01-01 00:00:00"`
	EndTime   string `form:"end_time" example:"2023-12-31 23:59:59"`
}
