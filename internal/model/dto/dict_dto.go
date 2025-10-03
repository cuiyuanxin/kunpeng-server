package dto

// DictTypeCreateReq 创建字典类型请求
type DictTypeCreateReq struct {
	Name   string `json:"name" binding:"required" example:"性别"`
	Type   string `json:"type" binding:"required" example:"sys_gender"`
	Status int8   `json:"status" example:"1"`
	Remark string `json:"remark" example:"性别字典"`
}

// DictTypeUpdateReq 更新字典类型请求
type DictTypeUpdateReq struct {
	ID     uint   `json:"id" binding:"required" example:"1"`
	Name   string `json:"name" binding:"required" example:"性别"`
	Type   string `json:"type" binding:"required" example:"sys_gender"`
	Status int8   `json:"status" example:"1"`
	Remark string `json:"remark" example:"性别字典"`
}

// DictTypePageReq 字典类型分页请求
type DictTypePageReq struct {
	PageNum   int    `form:"page_num" binding:"required,min=1" example:"1"`
	PageSize  int    `form:"page_size" binding:"required,min=1,max=100" example:"10"`
	Name      string `form:"name" example:"性别"`
	Type      string `form:"type" example:"sys_gender"`
	Status    int8   `form:"status" example:"1"`
	BeginTime string `form:"begin_time" example:"2023-01-01 00:00:00"`
	EndTime   string `form:"end_time" example:"2023-12-31 23:59:59"`
}

// DictTypeReq 字典类型请求
type DictTypeReq struct {
	DictType string `uri:"dict_type" binding:"required" example:"sys_gender"`
}

// DictDataCreateReq 创建字典数据请求
type DictDataCreateReq struct {
	DictType string `json:"dict_type" binding:"required" example:"sys_gender"`
	Label    string `json:"label" binding:"required" example:"男"`
	Value    string `json:"value" binding:"required" example:"1"`
	Sort     int    `json:"sort" example:"0"`
	Status   int8   `json:"status" example:"1"`
	Remark   string `json:"remark" example:"男性"`
}

// DictDataUpdateReq 更新字典数据请求
type DictDataUpdateReq struct {
	ID       uint   `json:"id" binding:"required" example:"1"`
	DictType string `json:"dict_type" binding:"required" example:"sys_gender"`
	Label    string `json:"label" binding:"required" example:"男"`
	Value    string `json:"value" binding:"required" example:"1"`
	Sort     int    `json:"sort" example:"0"`
	Status   int8   `json:"status" example:"1"`
	Remark   string `json:"remark" example:"男性"`
}

// DictDataPageReq 字典数据分页请求
type DictDataPageReq struct {
	PageNum    int    `form:"page_num" binding:"required,min=1" example:"1"`
	PageSize   int    `form:"page_size" binding:"required,min=1,max=100" example:"10"`
	DictType   string `form:"dict_type" example:"sys_gender"`
	DictTypeID uint   `form:"dict_type_id" example:"1"`
	Label      string `form:"label" example:"男"`
	Value      string `form:"value" example:"1"`
	Status     int8   `form:"status" example:"1"`
	BeginTime  string `form:"begin_time" example:"2023-01-01 00:00:00"`
	EndTime    string `form:"end_time" example:"2023-12-31 23:59:59"`
}
