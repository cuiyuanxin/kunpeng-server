package controller

import (
	"strconv"

	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/internal/service"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/cuiyuanxin/kunpeng/pkg/jwt"
	"github.com/cuiyuanxin/kunpeng/pkg/response"
	"github.com/cuiyuanxin/kunpeng/pkg/utils"
	"github.com/cuiyuanxin/kunpeng/pkg/validator"
	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct{}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.UserLoginReq true "用户登录请求"
// @Success 200 {object} response.Response{data=dto.UserLoginResp} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req dto.UserLoginReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		// 记录登录失败日志
		go c.recordLoginLog(0, req.Account, ctx, 0, "参数验证失败")
		return
	}

	// 自定义验证账号格式（支持多语言）
	if err := req.Validate(); err != nil {
		response.FailWithMessage(ctx, kperrors.ErrParam, err.Error())
		// 记录登录失败日志
		go c.recordLoginLog(0, req.Account, ctx, 0, "账号格式验证失败")
		return
	}

	// 调用服务
	realIP := utils.GetRealIP(ctx.Request)
	resp, err := service.GetUserService().Login(&req, realIP)
	if err != nil {
		response.FailWithError(ctx, err)
		// 记录登录失败日志
		go c.recordLoginLog(0, req.Account, ctx, 0, err.Error())
		return
	}

	// 从access token中解析用户ID记录登录成功日志
	claims, _ := jwt.ParseToken(resp.AccessToken)
	go c.recordLoginLog(claims.UserID, req.Account, ctx, 1, "登录成功")

	response.OkWithData(ctx, resp)
}

// RefreshToken 刷新token
// @Summary 刷新token
// @Description 通过refresh token刷新access token和refresh token
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.RefreshTokenReq true "刷新token请求"
// @Success 200 {object} response.Response{data=dto.UserLoginResp} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/refresh-token [post]
func (c *UserController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务刷新token
	resp, err := service.GetUserService().RefreshToken(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// Logout 用户退出登录
// @Summary 用户退出登录
// @Description 用户退出登录
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/logout [post]
func (c *UserController) Logout(ctx *gin.Context) {
	// 获取用户信息
	userID, username := c.getUserInfo(ctx)

	// 获取token并加入黑名单
	authorization := ctx.GetHeader("Authorization")
	token := jwt.ExtractTokenFromHeader(authorization)
	if token != "" {
		tokenBlacklistService := service.GetTokenBlacklistService()
		err := tokenBlacklistService.BlacklistToken(token, userID, username, "用户主动退出登录")
		if err != nil {
			// 记录错误但不影响退出登录流程
			go c.recordLoginLog(userID, username, ctx, 0, "退出登录成功，但token黑名单添加失败")
		} else {
			// 记录退出登录日志
			go c.recordLoginLog(userID, username, ctx, 1, "退出登录成功")
		}
	} else {
		// 没有token的情况下也记录退出登录日志
		go c.recordLoginLog(userID, username, ctx, 1, "退出登录成功")
	}

	response.OkWithMessage(ctx, "退出登录成功")
}

// getUserInfo 获取用户信息
func (c *UserController) getUserInfo(ctx *gin.Context) (uint, string) {
	userID := uint(0)
	username := "anonymous"

	// 从JWT中获取用户信息
	if userIDStr, exists := ctx.Get("userID"); exists {
		if id, ok := userIDStr.(string); ok {
			if parsedID, err := strconv.ParseUint(id, 10, 32); err == nil {
				userID = uint(parsedID)
			}
		}
	}

	if usernameStr, exists := ctx.Get("username"); exists {
		if name, ok := usernameStr.(string); ok {
			username = name
		}
	}

	return userID, username
}

// recordLoginLog 记录登录日志
func (c *UserController) recordLoginLog(userID uint, username string, ctx *gin.Context, status int8, message string) {
	userAgent := ctx.Request.UserAgent()

	// 简单解析User-Agent（实际项目中可以使用专门的库如ua-parser）
	browser := "Unknown"
	os := "Unknown"
	device := "Unknown"

	// 简单的User-Agent解析逻辑
	if userAgent != "" {
		// 这里可以使用更复杂的解析逻辑或第三方库
		browser = "Browser" // 简化处理
		os = "OS"           // 简化处理
		device = "Device"   // 简化处理
	}

	loginLogService := service.GetLoginLogService()
	loginLogService.RecordLoginLog(
		userID,
		username,
		ctx.ClientIP(),
		browser,
		os,
		device,
		"", // location 可以根据IP获取地理位置
		status,
		message,
	)
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=dto.UserInfoResp} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/user/info [get]
func (c *UserController) GetUserInfo(ctx *gin.Context) {
	// 获取当前用户ID
	userID := jwt.GetUserID(ctx)

	// 调用服务
	resp, err := service.GetUserService().GetUserInfo(userID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetUserList 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page_num query int true "页码" default(1)
// @Param page_size query int true "每页数量" default(10)
// @Param username query string false "用户名"
// @Param nickname query string false "昵称"
// @Param mobile query string false "手机号"
// @Param status query int false "状态" Enums(0, 1)
// @Param dept_id query int false "部门ID"
// @Param begin_time query string false "开始时间"
// @Param end_time query string false "结束时间"
// @Success 200 {object} response.Response{data=dto.PageResp} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/users [get]
func (c *UserController) GetUserList(ctx *gin.Context) {
	var req dto.UserPageReq
	if err := validator.BindAndValidateQueryI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetUserService().GetUserList(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// GetUserByID 根据ID获取用户
// @Summary 根据ID获取用户
// @Description 根据ID获取用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=model.User} "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetUserByID(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	resp, err := service.GetUserService().GetUserByID(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, resp)
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.UserCreateReq true "创建用户请求"
// @Success 200 {object} response.Response{data=uint} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/users [post]
func (c *UserController) CreateUser(ctx *gin.Context) {
	var req dto.UserCreateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	id, err := service.GetUserService().CreateUser(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.OkWithData(ctx, id)
}

// UpdateUser 更新用户
// @Summary 更新用户
// @Description 更新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.UserUpdateReq true "更新用户请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/users [put]
func (c *UserController) UpdateUser(ctx *gin.Context) {
	var req dto.UserUpdateReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetUserService().UpdateUser(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// DeleteUser 删除用户
// @Summary 删除用户
// @Description 删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/users/{id} [delete]
func (c *UserController) DeleteUser(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetUserService().DeleteUser(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// BatchDeleteUser 批量删除用户
// @Summary 批量删除用户
// @Description 批量删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.IDsReq true "批量删除用户请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/users/batch [delete]
func (c *UserController) BatchDeleteUser(ctx *gin.Context) {
	var req dto.IDsReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetUserService().BatchDeleteUser(req.IDs)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// ChangeUserStatus 修改用户状态
// @Summary 修改用户状态
// @Description 修改用户状态
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.StatusReq true "修改用户状态请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/users/status [put]
func (c *UserController) ChangeUserStatus(ctx *gin.Context) {
	var req dto.StatusReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetUserService().ChangeUserStatus(&req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// ResetUserPassword 重置用户密码
// @Summary 重置用户密码
// @Description 重置用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response "成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/users/{id}/password/reset [put]
func (c *UserController) ResetUserPassword(ctx *gin.Context) {
	var req dto.IDReq
	if err := validator.BindAndValidateUriI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 调用服务
	err := service.GetUserService().ResetUserPassword(req.ID)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param data body dto.UserChangePasswordReq true "修改密码请求"
// @Success 200 {object} response.Response "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/users/password [put]
func (c *UserController) ChangePassword(ctx *gin.Context) {
	var req dto.UserChangePasswordReq
	if err := validator.BindAndValidateJSONI18n(ctx, &req); err != nil {
		response.FailWithError(ctx, err)
		return
	}

	// 获取当前用户ID
	userID := jwt.GetUserID(ctx)

	// 调用服务
	err := service.GetUserService().ChangePassword(userID, &req)
	if err != nil {
		response.FailWithError(ctx, err)
		return
	}

	response.Ok(ctx)
}
