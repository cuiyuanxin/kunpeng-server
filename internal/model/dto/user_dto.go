package dto

import (
	"errors"
	"regexp"

	"github.com/cuiyuanxin/kunpeng/pkg/constants"
	"github.com/cuiyuanxin/kunpeng/pkg/i18n"
)

// UserLoginReq 用户登录请求
type UserLoginReq struct {
	LoginType  string `json:"login_type" binding:"required,oneof=username mobile" example:"username"` // username: 用户名登录, mobile: 手机号登录
	Account    string `json:"account" binding:"required" example:"admin"`                             // 账号（用户名或手机号）
	Password   string `json:"password" example:"Abc123!@#"`                                           // 密码（账号登录时必填）
	Captcha    string `json:"captcha" example:"123456"`                                               // 验证码（手机号登录时必填）
	CaptchaID  string `json:"captcha_id" example:"abcd1234"`
	RememberMe bool   `json:"remember_me"`
}

// Validate 自定义验证方法
func (req *UserLoginReq) Validate() error {
	return req.validateInternal()
}

// validateInternal 内部验证方法
func (req *UserLoginReq) validateInternal() error {
	// 根据登录类型验证账号格式和必填字段
	if req.LoginType == "username" {
		// 验证用户名格式
		matched, _ := regexp.MatchString(constants.UsernameRegex, req.Account)
		if !matched {
			return errors.New(i18n.TWithField("validator.username", "username"))
		}
		// 账号登录必须有密码
		if req.Password == "" {
			return errors.New(i18n.TWithField("validator.required", "password"))
		}
		// 验证密码格式
		if !constants.ValidatePassword(req.Password) {
			return errors.New(i18n.TWithField("validator.password", "password"))
		}
	} else if req.LoginType == "mobile" {
		// 验证手机号格式
		matched, _ := regexp.MatchString(constants.MobileRegex, req.Account)
		if !matched {
			return errors.New(i18n.TWithField("validator.mobile", "mobile"))
		}
		// 手机号登录必须有验证码
		if req.Captcha == "" {
			return errors.New(i18n.TWithField("validator.required", "captcha"))
		}
		// 验证码必须是6位数字
		matched, _ = regexp.MatchString(constants.CaptchaRegex, req.Captcha)
		if !matched {
			return errors.New(i18n.TWithField("validator.captcha", "captcha"))
		}
	}
	return nil
}

// UserLoginResp 用户登录响应
type UserLoginResp struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`         // access token过期时间（秒）
	RefreshExpiresIn int64  `json:"refresh_expires_in"` // refresh token过期时间（秒）
}

// RefreshTokenReq 刷新token请求
type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// UserInfoResp 用户信息响应
type UserInfoResp struct {
	ID       uint   `json:"userId"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	RealName string `json:"realname"`
	Avatar   string `json:"avatar"`
	Gender   int8   `json:"gender"`
	RoleName string `json:"role_name"`
	DeptName string `json:"dept_name"`
	PostName string `json:"post_name"`
}

// UserCreateReq 创建用户请求
type UserCreateReq struct {
	Username string `json:"username" binding:"required,username" example:"zhangsan"`
	Password string `json:"password" binding:"required,password" example:"Abc123!@#"`
	Nickname string `json:"nickname" example:"张三"`
	RealName string `json:"real_name" example:"张三"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar.png"`
	Gender   int8   `json:"gender" example:"1"`
	Email    string `json:"email" binding:"email" example:"zhangsan@example.com"`
	Mobile   string `json:"mobile" example:"13800138000"`
	DeptID   uint   `json:"dept_id" example:"1"`
	PostID   uint   `json:"post_id" example:"1"`
	RoleID   uint   `json:"role_id" example:"1"`
	Status   int8   `json:"status" example:"1"`
	Remark   string `json:"remark" example:"测试用户"`
}

// UserUpdateReq 更新用户请求
type UserUpdateReq struct {
	ID       uint   `json:"id" binding:"required" example:"1"`
	Nickname string `json:"nickname" example:"张三"`
	RealName string `json:"real_name" example:"张三"`
	Avatar   string `json:"avatar" example:"https://example.com/avatar.png"`
	Gender   int8   `json:"gender" example:"1"`
	Email    string `json:"email" binding:"email" example:"zhangsan@example.com"`
	Mobile   string `json:"mobile" example:"13800138000"`
	DeptID   uint   `json:"dept_id" example:"1"`
	PostID   uint   `json:"post_id" example:"1"`
	RoleID   uint   `json:"role_id" example:"1"`
	Status   int8   `json:"status" example:"1"`
	Remark   string `json:"remark" example:"测试用户"`
}

// UserChangePasswordReq 修改密码请求
type UserChangePasswordReq struct {
	OldPassword string `json:"old_password" binding:"required,password" example:"Abc123!@#"`
	NewPassword string `json:"new_password" binding:"required,password" example:"Def456$%^"`
}

// UserPageReq 用户分页请求
type UserPageReq struct {
	PageNum   int    `form:"page_num" binding:"required,min=1" example:"1"`
	PageSize  int    `form:"page_size" binding:"required,min=1,max=100" example:"10"`
	Username  string `form:"username" example:"admin"`
	Nickname  string `form:"nickname" example:"管理员"`
	Mobile    string `form:"mobile" example:"13800138000"`
	Status    int8   `form:"status" example:"1"`
	DeptID    uint   `form:"dept_id" example:"1"`
	BeginTime string `form:"begin_time" example:"2023-01-01 00:00:00"`
	EndTime   string `form:"end_time" example:"2023-12-31 23:59:59"`
}
