package errors

import (
	"fmt"
	"github.com/cuiyuanxin/kunpeng/pkg/i18n"
)

// getI18nMessage 获取国际化消息，如果翻译失败则返回回退消息
func getI18nMessage(messageID, fallback string) string {
	result := i18n.T(messageID)
	// 如果翻译失败（返回的是 messageID 本身），则使用回退消息
	if result == messageID {
		return fallback
	}
	return result
}

// 错误码规则：
// 1. 错误码为5位数字
// 2. 第1位表示错误级别：1-系统级错误，2-业务级错误
// 3. 第2-3位表示模块：00-公共，01-用户，02-权限，03-配置，...
// 4. 第4-5位表示具体错误码

// 系统级错误码 (10000-19999)
const (
	// 公共系统错误 (10000-10099)
	ErrSystem           = 10000 // 系统内部错误
	ErrUnknown          = 10001 // 未知错误
	ErrParam            = 10002 // 参数错误
	ErrSignature        = 10003 // 签名错误
	ErrUnauthorized     = 10004 // 未授权
	ErrForbidden        = 10005 // 禁止访问
	ErrNotFound         = 10006 // 资源不存在
	ErrMethodNotAllowed = 10007 // 方法不允许
	ErrTimeout          = 10008 // 超时
	ErrTooManyRequests  = 10009 // 请求过多
	ErrInvalidToken     = 10010 // 无效的令牌
	ErrTokenExpired     = 10011 // 令牌已过期
	ErrInvalidSign      = 10012 // 无效的签名
	ErrDatabase         = 10013 // 数据库错误
	ErrRedis            = 10014 // Redis错误
	ErrNetwork          = 10015 // 网络错误
	ErrThirdParty       = 10016 // 第三方服务错误
	ErrUpload           = 10017 // 上传错误
	ErrDownload         = 10018 // 下载错误
	ErrExport           = 10019 // 导出错误
	ErrImport           = 10020 // 导入错误
	ErrValidation       = 10021 // 参数验证失败

	// 数据库错误(10100-10199)
	ErrDBConnection   = 10100 // 数据库连接错误
	ErrDBExecution    = 10101 // 数据库执行错误
	ErrDBTransaction  = 10102 // 数据库事务错误
	ErrDBDuplicate    = 10103 // 数据库重复记录
	ErrDBNotFound     = 10104 // 数据库记录不存在
	ErrDBForeignKey   = 10105 // 数据库外键约束错误
	ErrDBUniqueKey    = 10106 // 数据库唯一约束错误
	ErrDBPrimaryKey   = 10107 // 数据库主键约束错误
	ErrDBDeadlock     = 10108 // 数据库死锁错误
	ErrDBTimeout      = 10109 // 数据库超时错误
	ErrDBSyntax       = 10110 // 数据库语法错误
	ErrDBConstraint   = 10111 // 数据库约束错误
	ErrDBDataTooLong  = 10112 // 数据库数据过长错误
	ErrDBDataInvalid  = 10113 // 数据库数据无效错误
	ErrDBDataTruncate = 10114 // 数据库数据截断错误

	// 缓存错误 (10200-10299)
	ErrCacheConnection = 10200 // 缓存连接错误
	ErrCacheExecution  = 10201 // 缓存执行错误
	ErrCacheNotFound   = 10202 // 缓存记录不存在
	ErrCacheExpired    = 10203 // 缓存已过期
	ErrCacheFull       = 10204 // 缓存已满
	ErrCacheLock       = 10205 // 缓存锁错误

	// 认证错误 (10300-10399)
	ErrAuthLogin        = 10300 // 登录失败
	ErrAuthPassword     = 10301 // 密码错误
	ErrAuthAccount      = 10302 // 账号不存在
	ErrAuthDisabled     = 10303 // 账号已禁用
	ErrAuthLocked       = 10304 // 账号已锁定
	ErrAuthExpired      = 10305 // 账号已过期
	ErrAuthLogout       = 10306 // 登出失败
	ErrAuthCaptcha      = 10307 // 验证码错误
	ErrAuthCaptchaEmpty = 10308 // 验证码为空
	ErrAuthCaptchaUsed  = 10309 // 验证码已使用
	ErrAuthCaptchaFreq  = 10310 // 验证码发送频率过高
	ErrAuthTwoFactor    = 10311 // 两因素认证失败
	ErrAuthToken        = 10312 // 令牌错误
	ErrAuthSession      = 10313 // 会话错误
	ErrAuthSSO          = 10314 // 单点登录错误

	// 权限错误 (10400-10499)
	ErrPermDenied      = 10400 // 权限不足
	ErrPermRole        = 10401 // 角色不存在
	ErrPermRoleDisable = 10402 // 角色已禁用
	ErrPermPolicy      = 10403 // 策略不存在
	ErrPermAction      = 10404 // 操作不允许
	ErrPermResource    = 10405 // 资源不可访问
	ErrPermMenu        = 10406 // 菜单不可访问
	ErrPermAPI         = 10407 // API不可访问
	ErrPermData        = 10408 // 数据不可访问
	ErrPermButton      = 10409 // 按钮不可访问
	ErrPermField       = 10410 // 字段不可访问
	ErrPermFile        = 10411 // 文件不可访问
	ErrPermExport      = 10412 // 导出不允许
	ErrPermImport      = 10413 // 导入不允许
	ErrPermDelete      = 10414 // 删除不允许
	ErrPermUpdate      = 10415 // 更新不允许
	ErrPermCreate      = 10416 // 创建不允许
	ErrPermRead        = 10417 // 读取不允许

	// 文件错误 (10500-10599)
	ErrFileUpload      = 10500 // 文件上传失败
	ErrFileDownload    = 10501 // 文件下载失败
	ErrFileDelete      = 10502 // 文件删除失败
	ErrFileNotFound    = 10503 // 文件不存在
	ErrFileSize        = 10504 // 文件大小超限
	ErrFileType        = 10505 // 文件类型不支持
	ErrFileFormat      = 10506 // 文件格式错误
	ErrFileCorrupt     = 10507 // 文件已损坏
	ErrFilePermission  = 10508 // 文件权限不足
	ErrFileDuplicate   = 10509 // 文件已存在
	ErrFileStorage     = 10510 // 文件存储错误
	ErrFileRead        = 10511 // 文件读取错误
	ErrFileWrite       = 10512 // 文件写入错误
	ErrFilePath        = 10513 // 文件路径错误
	ErrFileNameTooLong = 10514 // 文件名过长
	ErrFileNameInvalid = 10515 // 文件名无效
)

// 业务级错误码 (20000-29999)
const (
	// 用户模块错误 (20100-20199)
	ErrUserNotFound      = 20100 // 用户不存在
	ErrUserDisabled      = 20101 // 用户已禁用
	ErrUserLocked        = 20102 // 用户已锁定
	ErrUserExpired       = 20103 // 用户已过期
	ErrUserPassword      = 20104 // 用户密码错误
	ErrUserOldPassword   = 20105 // 用户旧密码错误
	ErrUserExists        = 20106 // 用户已存在
	ErrUserNameExists    = 20107 // 用户名已存在
	ErrUserEmailExists   = 20108 // 用户邮箱已存在
	ErrUserPhoneExists   = 20109 // 用户手机号已存在
	ErrUserNameInvalid   = 20110 // 用户名无效
	ErrUserEmailInvalid  = 20111 // 用户邮箱无效
	ErrUserPhoneInvalid  = 20112 // 用户手机号无效
	ErrUserRoleInvalid   = 20113 // 用户角色无效
	ErrUserDeptInvalid   = 20114 // 用户部门无效
	ErrUserPostInvalid   = 20115 // 用户岗位无效
	ErrUserStatusInvalid = 20116 // 用户状态无效

	// 角色模块错误 (20200-20299)
	ErrRoleNotFound      = 20200 // 角色不存在
	ErrRoleDisabled      = 20201 // 角色已禁用
	ErrRoleExists        = 20202 // 角色已存在
	ErrRoleNameExists    = 20203 // 角色名已存在
	ErrRoleCodeExists    = 20204 // 角色编码已存在
	ErrRoleNameInvalid   = 20205 // 角色名无效
	ErrRoleCodeInvalid   = 20206 // 角色编码无效
	ErrRoleStatusInvalid = 20207 // 角色状态无效
	ErrRoleHasUsers      = 20208 // 角色下有用户
	ErrRoleHasChildren   = 20209 // 角色下有子角色
	ErrRoleHasMenus      = 20210 // 角色下有菜单
	ErrRoleHasPerms      = 20211 // 角色下有权限
	ErrRoleHasDepts      = 20212 // 角色下有部门
	ErrRoleHasPosts      = 20213 // 角色下有岗位
	ErrRoleHasApis       = 20214 // 角色下有API
	ErrRoleHasButtons    = 20215 // 角色下有按钮

	// 菜单模块错误 (20300-20399)
	ErrMenuNotFound         = 20300 // 菜单不存在
	ErrMenuDisabled         = 20301 // 菜单已禁用
	ErrMenuExists           = 20302 // 菜单已存在
	ErrMenuNameExists       = 20303 // 菜单名已存在
	ErrMenuCodeExists       = 20304 // 菜单编码已存在
	ErrMenuNameInvalid      = 20305 // 菜单名无效
	ErrMenuCodeInvalid      = 20306 // 菜单编码无效
	ErrMenuStatusInvalid    = 20307 // 菜单状态无效
	ErrMenuHasChildren      = 20308 // 菜单下有子菜单
	ErrMenuHasButtons       = 20309 // 菜单下有按钮
	ErrMenuHasApis          = 20310 // 菜单下有API
	ErrMenuParentNotFound   = 20311 // 父菜单不存在
	ErrMenuParentDisabled   = 20312 // 父菜单已禁用
	ErrMenuParentInvalid    = 20313 // 父菜单无效
	ErrMenuTypeInvalid      = 20314 // 菜单类型无效
	ErrMenuPathInvalid      = 20315 // 菜单路径无效
	ErrMenuComponentInvalid = 20316 // 菜单组件无效
	ErrMenuPermInvalid      = 20317 // 菜单权限无效
	ErrMenuIconInvalid      = 20318 // 菜单图标无效
	ErrMenuSortInvalid      = 20319 // 菜单排序无效

	// 登录日志模块错误码 (20400-20499)
	ErrLoginLogNotFound    = 20400 // 登录日志不存在
	ErrLoginLogGetList     = 20401 // 获取登录日志列表失败
	ErrLoginLogGetByID     = 20402 // 根据ID获取登录日志失败
	ErrLoginLogDelete      = 20403 // 删除登录日志失败
	ErrLoginLogBatchDelete = 20404 // 批量删除登录日志失败
	ErrLoginLogCleanOld    = 20405 // 清理旧登录日志失败

	// 操作日志模块错误码 (20500-20599)
	ErrOperationLogNotFound    = 20500 // 操作日志不存在
	ErrOperationLogGetList     = 20501 // 获取操作日志列表失败
	ErrOperationLogGetByID     = 20502 // 根据ID获取操作日志失败
	ErrOperationLogDelete      = 20503 // 删除操作日志失败
	ErrOperationLogBatchDelete = 20504 // 批量删除操作日志失败
	ErrOperationLogCleanOld    = 20505 // 清理旧操作日志失败
	ErrOperationLogInvalidID   = 20506 // 操作日志ID无效
)

// 错误码映射表 - 使用 i18n 获取国际化错误消息
func getCodeMessageMap() map[int]string {
	return map[int]string{
		// 系统级错误
		ErrSystem:           getI18nMessage(fmt.Sprintf("error.%d", ErrSystem), "系统内部错误"),
		ErrUnknown:          getI18nMessage(fmt.Sprintf("error.%d", ErrUnknown), "未知错误"),
		ErrParam:            getI18nMessage(fmt.Sprintf("error.%d", ErrParam), "参数错误"),
		ErrSignature:        getI18nMessage(fmt.Sprintf("error.%d", ErrSignature), "签名错误"),
		ErrUnauthorized:     getI18nMessage(fmt.Sprintf("error.%d", ErrUnauthorized), "未授权"),
		ErrForbidden:        getI18nMessage(fmt.Sprintf("error.%d", ErrForbidden), "禁止访问"),
		ErrNotFound:         getI18nMessage(fmt.Sprintf("error.%d", ErrNotFound), "资源不存在"),
		ErrMethodNotAllowed: getI18nMessage(fmt.Sprintf("error.%d", ErrMethodNotAllowed), "方法不允许"),
		ErrTimeout:          getI18nMessage(fmt.Sprintf("error.%d", ErrTimeout), "超时"),
		ErrTooManyRequests:  getI18nMessage(fmt.Sprintf("error.%d", ErrTooManyRequests), "请求过多"),
		ErrInvalidToken:     getI18nMessage(fmt.Sprintf("error.%d", ErrInvalidToken), "无效的令牌"),
		ErrTokenExpired:     getI18nMessage(fmt.Sprintf("error.%d", ErrTokenExpired), "令牌已过期"),
		ErrInvalidSign:      getI18nMessage(fmt.Sprintf("error.%d", ErrInvalidSign), "无效的签名"),
		ErrDatabase:         getI18nMessage(fmt.Sprintf("error.%d", ErrDatabase), "数据库错误"),
		ErrRedis:            getI18nMessage(fmt.Sprintf("error.%d", ErrRedis), "Redis错误"),
		ErrNetwork:          getI18nMessage(fmt.Sprintf("error.%d", ErrNetwork), "网络错误"),
		ErrThirdParty:       getI18nMessage(fmt.Sprintf("error.%d", ErrThirdParty), "第三方服务错误"),
		ErrUpload:           getI18nMessage(fmt.Sprintf("error.%d", ErrUpload), "上传错误"),
		ErrDownload:         getI18nMessage(fmt.Sprintf("error.%d", ErrDownload), "下载错误"),
		ErrExport:           getI18nMessage(fmt.Sprintf("error.%d", ErrExport), "导出错误"),
		ErrImport:           getI18nMessage(fmt.Sprintf("error.%d", ErrImport), "导入错误"),
		ErrValidation:       getI18nMessage(fmt.Sprintf("error.%d", ErrValidation), "参数验证失败"),

		// 数据库错误
		ErrDBConnection:   getI18nMessage(fmt.Sprintf("error.%d", ErrDBConnection), "数据库连接错误"),
		ErrDBExecution:    getI18nMessage(fmt.Sprintf("error.%d", ErrDBExecution), "数据库执行错误"),
		ErrDBTransaction:  getI18nMessage(fmt.Sprintf("error.%d", ErrDBTransaction), "数据库事务错误"),
		ErrDBDuplicate:    getI18nMessage(fmt.Sprintf("error.%d", ErrDBDuplicate), "数据库重复记录"),
		ErrDBNotFound:     getI18nMessage(fmt.Sprintf("error.%d", ErrDBNotFound), "数据库记录不存在"),
		ErrDBForeignKey:   getI18nMessage(fmt.Sprintf("error.%d", ErrDBForeignKey), "数据库外键约束错误"),
		ErrDBUniqueKey:    getI18nMessage(fmt.Sprintf("error.%d", ErrDBUniqueKey), "数据库唯一约束错误"),
		ErrDBPrimaryKey:   getI18nMessage(fmt.Sprintf("error.%d", ErrDBPrimaryKey), "数据库主键约束错误"),
		ErrDBDeadlock:     getI18nMessage(fmt.Sprintf("error.%d", ErrDBDeadlock), "数据库死锁错误"),
		ErrDBTimeout:      getI18nMessage(fmt.Sprintf("error.%d", ErrDBTimeout), "数据库超时错误"),
		ErrDBSyntax:       getI18nMessage(fmt.Sprintf("error.%d", ErrDBSyntax), "数据库语法错误"),
		ErrDBConstraint:   getI18nMessage(fmt.Sprintf("error.%d", ErrDBConstraint), "数据库约束错误"),
		ErrDBDataTooLong:  getI18nMessage(fmt.Sprintf("error.%d", ErrDBDataTooLong), "数据库数据过长错误"),
		ErrDBDataInvalid:  getI18nMessage(fmt.Sprintf("error.%d", ErrDBDataInvalid), "数据库数据无效错误"),
		ErrDBDataTruncate: getI18nMessage(fmt.Sprintf("error.%d", ErrDBDataTruncate), "数据库数据截断错误"),

		// 缓存错误
		ErrCacheConnection: getI18nMessage(fmt.Sprintf("error.%d", ErrCacheConnection), "缓存连接错误"),
		ErrCacheExecution:  getI18nMessage(fmt.Sprintf("error.%d", ErrCacheExecution), "缓存执行错误"),
		ErrCacheNotFound:   getI18nMessage(fmt.Sprintf("error.%d", ErrCacheNotFound), "缓存记录不存在"),
		ErrCacheExpired:    getI18nMessage(fmt.Sprintf("error.%d", ErrCacheExpired), "缓存已过期"),
		ErrCacheFull:       getI18nMessage(fmt.Sprintf("error.%d", ErrCacheFull), "缓存已满"),
		ErrCacheLock:       getI18nMessage(fmt.Sprintf("error.%d", ErrCacheLock), "缓存锁错误"),

		// 认证错误
		ErrAuthLogin:        getI18nMessage(fmt.Sprintf("error.%d", ErrAuthLogin), "登录失败"),
		ErrAuthPassword:     getI18nMessage(fmt.Sprintf("error.%d", ErrAuthPassword), "密码错误"),
		ErrAuthAccount:      getI18nMessage(fmt.Sprintf("error.%d", ErrAuthAccount), "账号不存在"),
		ErrAuthDisabled:     getI18nMessage(fmt.Sprintf("error.%d", ErrAuthDisabled), "账号已禁用"),
		ErrAuthLocked:       getI18nMessage(fmt.Sprintf("error.%d", ErrAuthLocked), "账号已锁定"),
		ErrAuthExpired:      getI18nMessage(fmt.Sprintf("error.%d", ErrAuthExpired), "账号已过期"),
		ErrAuthLogout:       getI18nMessage(fmt.Sprintf("error.%d", ErrAuthLogout), "登出失败"),
		ErrAuthCaptcha:      getI18nMessage(fmt.Sprintf("error.%d", ErrAuthCaptcha), "验证码错误"),
		ErrAuthCaptchaEmpty: getI18nMessage(fmt.Sprintf("error.%d", ErrAuthCaptchaEmpty), "验证码为空"),
		ErrAuthCaptchaUsed:  getI18nMessage(fmt.Sprintf("error.%d", ErrAuthCaptchaUsed), "验证码已使用"),
		ErrAuthCaptchaFreq:  getI18nMessage(fmt.Sprintf("error.%d", ErrAuthCaptchaFreq), "验证码发送频率过高"),
		ErrAuthTwoFactor:    getI18nMessage(fmt.Sprintf("error.%d", ErrAuthTwoFactor), "两因素认证失败"),
		ErrAuthToken:        getI18nMessage(fmt.Sprintf("error.%d", ErrAuthToken), "令牌错误"),
		ErrAuthSession:      getI18nMessage(fmt.Sprintf("error.%d", ErrAuthSession), "会话错误"),
		ErrAuthSSO:          getI18nMessage(fmt.Sprintf("error.%d", ErrAuthSSO), "单点登录错误"),

		// 权限错误
		ErrPermDenied:      getI18nMessage(fmt.Sprintf("error.%d", ErrPermDenied), "权限不足"),
		ErrPermRole:        getI18nMessage(fmt.Sprintf("error.%d", ErrPermRole), "角色不存在"),
		ErrPermRoleDisable: getI18nMessage(fmt.Sprintf("error.%d", ErrPermRoleDisable), "角色已禁用"),
		ErrPermPolicy:      getI18nMessage(fmt.Sprintf("error.%d", ErrPermPolicy), "策略不存在"),
		ErrPermAction:      getI18nMessage(fmt.Sprintf("error.%d", ErrPermAction), "操作不允许"),
		ErrPermResource:    getI18nMessage(fmt.Sprintf("error.%d", ErrPermResource), "资源不可访问"),
		ErrPermMenu:        getI18nMessage(fmt.Sprintf("error.%d", ErrPermMenu), "菜单不可访问"),
		ErrPermAPI:         getI18nMessage(fmt.Sprintf("error.%d", ErrPermAPI), "API不可访问"),
		ErrPermData:        getI18nMessage(fmt.Sprintf("error.%d", ErrPermData), "数据不可访问"),
		ErrPermButton:      getI18nMessage(fmt.Sprintf("error.%d", ErrPermButton), "按钮不可访问"),
		ErrPermField:       getI18nMessage(fmt.Sprintf("error.%d", ErrPermField), "字段不可访问"),
		ErrPermFile:        getI18nMessage(fmt.Sprintf("error.%d", ErrPermFile), "文件不可访问"),
		ErrPermExport:      getI18nMessage(fmt.Sprintf("error.%d", ErrPermExport), "导出不允许"),
		ErrPermImport:      getI18nMessage(fmt.Sprintf("error.%d", ErrPermImport), "导入不允许"),
		ErrPermDelete:      getI18nMessage(fmt.Sprintf("error.%d", ErrPermDelete), "删除不允许"),
		ErrPermUpdate:      getI18nMessage(fmt.Sprintf("error.%d", ErrPermUpdate), "更新不允许"),
		ErrPermCreate:      getI18nMessage(fmt.Sprintf("error.%d", ErrPermCreate), "创建不允许"),
		ErrPermRead:        getI18nMessage(fmt.Sprintf("error.%d", ErrPermRead), "读取不允许"),

		// 文件错误
		ErrFileUpload:      getI18nMessage(fmt.Sprintf("error.%d", ErrFileUpload), "文件上传失败"),
		ErrFileDownload:    getI18nMessage(fmt.Sprintf("error.%d", ErrFileDownload), "文件下载失败"),
		ErrFileDelete:      getI18nMessage(fmt.Sprintf("error.%d", ErrFileDelete), "文件删除失败"),
		ErrFileNotFound:    getI18nMessage(fmt.Sprintf("error.%d", ErrFileNotFound), "文件不存在"),
		ErrFileSize:        getI18nMessage(fmt.Sprintf("error.%d", ErrFileSize), "文件大小超限"),
		ErrFileType:        getI18nMessage(fmt.Sprintf("error.%d", ErrFileType), "文件类型不支持"),
		ErrFileFormat:      getI18nMessage(fmt.Sprintf("error.%d", ErrFileFormat), "文件格式错误"),
		ErrFileCorrupt:     getI18nMessage(fmt.Sprintf("error.%d", ErrFileCorrupt), "文件已损坏"),
		ErrFilePermission:  getI18nMessage(fmt.Sprintf("error.%d", ErrFilePermission), "文件权限不足"),
		ErrFileDuplicate:   getI18nMessage(fmt.Sprintf("error.%d", ErrFileDuplicate), "文件已存在"),
		ErrFileStorage:     getI18nMessage(fmt.Sprintf("error.%d", ErrFileStorage), "文件存储错误"),
		ErrFileRead:        getI18nMessage(fmt.Sprintf("error.%d", ErrFileRead), "文件读取错误"),
		ErrFileWrite:       getI18nMessage(fmt.Sprintf("error.%d", ErrFileWrite), "文件写入错误"),
		ErrFilePath:        getI18nMessage(fmt.Sprintf("error.%d", ErrFilePath), "文件路径错误"),
		ErrFileNameTooLong: getI18nMessage(fmt.Sprintf("error.%d", ErrFileNameTooLong), "文件名过长"),
		ErrFileNameInvalid: getI18nMessage(fmt.Sprintf("error.%d", ErrFileNameInvalid), "文件名无效"),

		// 用户模块错误
		ErrUserNotFound:      getI18nMessage(fmt.Sprintf("error.%d", ErrUserNotFound), "用户不存在"),
		ErrUserDisabled:      getI18nMessage(fmt.Sprintf("error.%d", ErrUserDisabled), "用户已禁用"),
		ErrUserLocked:        getI18nMessage(fmt.Sprintf("error.%d", ErrUserLocked), "用户已锁定"),
		ErrUserExpired:       getI18nMessage(fmt.Sprintf("error.%d", ErrUserExpired), "用户已过期"),
		ErrUserPassword:      getI18nMessage(fmt.Sprintf("error.%d", ErrUserPassword), "用户密码错误"),
		ErrUserOldPassword:   getI18nMessage(fmt.Sprintf("error.%d", ErrUserOldPassword), "用户旧密码错误"),
		ErrUserExists:        getI18nMessage(fmt.Sprintf("error.%d", ErrUserExists), "用户已存在"),
		ErrUserNameExists:    getI18nMessage(fmt.Sprintf("error.%d", ErrUserNameExists), "用户名已存在"),
		ErrUserEmailExists:   getI18nMessage(fmt.Sprintf("error.%d", ErrUserEmailExists), "用户邮箱已存在"),
		ErrUserPhoneExists:   getI18nMessage(fmt.Sprintf("error.%d", ErrUserPhoneExists), "用户手机号已存在"),
		ErrUserNameInvalid:   getI18nMessage(fmt.Sprintf("error.%d", ErrUserNameInvalid), "用户名无效"),
		ErrUserEmailInvalid:  getI18nMessage(fmt.Sprintf("error.%d", ErrUserEmailInvalid), "用户邮箱无效"),
		ErrUserPhoneInvalid:  getI18nMessage(fmt.Sprintf("error.%d", ErrUserPhoneInvalid), "用户手机号无效"),
		ErrUserRoleInvalid:   getI18nMessage(fmt.Sprintf("error.%d", ErrUserRoleInvalid), "用户角色无效"),
		ErrUserDeptInvalid:   getI18nMessage(fmt.Sprintf("error.%d", ErrUserDeptInvalid), "用户部门无效"),
		ErrUserPostInvalid:   getI18nMessage(fmt.Sprintf("error.%d", ErrUserPostInvalid), "用户岗位无效"),
		ErrUserStatusInvalid: getI18nMessage(fmt.Sprintf("error.%d", ErrUserStatusInvalid), "用户状态无效"),

		// 角色模块错误
		ErrRoleNotFound:      getI18nMessage(fmt.Sprintf("error.%d", ErrRoleNotFound), "角色不存在"),
		ErrRoleDisabled:      getI18nMessage(fmt.Sprintf("error.%d", ErrRoleDisabled), "角色已禁用"),
		ErrRoleExists:        getI18nMessage(fmt.Sprintf("error.%d", ErrRoleExists), "角色已存在"),
		ErrRoleNameExists:    getI18nMessage(fmt.Sprintf("error.%d", ErrRoleNameExists), "角色名已存在"),
		ErrRoleCodeExists:    getI18nMessage(fmt.Sprintf("error.%d", ErrRoleCodeExists), "角色编码已存在"),
		ErrRoleNameInvalid:   getI18nMessage(fmt.Sprintf("error.%d", ErrRoleNameInvalid), "角色名无效"),
		ErrRoleCodeInvalid:   getI18nMessage(fmt.Sprintf("error.%d", ErrRoleCodeInvalid), "角色编码无效"),
		ErrRoleStatusInvalid: getI18nMessage(fmt.Sprintf("error.%d", ErrRoleStatusInvalid), "角色状态无效"),
		ErrRoleHasUsers:      getI18nMessage(fmt.Sprintf("error.%d", ErrRoleHasUsers), "角色下有用户"),
		ErrRoleHasChildren:   getI18nMessage(fmt.Sprintf("error.%d", ErrRoleHasChildren), "角色下有子角色"),
		ErrRoleHasMenus:      getI18nMessage(fmt.Sprintf("error.%d", ErrRoleHasMenus), "角色下有菜单"),
		ErrRoleHasPerms:      getI18nMessage(fmt.Sprintf("error.%d", ErrRoleHasPerms), "角色下有权限"),
		ErrRoleHasDepts:      getI18nMessage(fmt.Sprintf("error.%d", ErrRoleHasDepts), "角色下有部门"),
		ErrRoleHasPosts:      getI18nMessage(fmt.Sprintf("error.%d", ErrRoleHasPosts), "角色下有岗位"),
		ErrRoleHasApis:       getI18nMessage(fmt.Sprintf("error.%d", ErrRoleHasApis), "角色下有API"),
		ErrRoleHasButtons:    getI18nMessage(fmt.Sprintf("error.%d", ErrRoleHasButtons), "角色下有按钮"),

		// 菜单模块错误
		ErrMenuNotFound:         getI18nMessage(fmt.Sprintf("error.%d", ErrMenuNotFound), "菜单不存在"),
		ErrMenuDisabled:         getI18nMessage(fmt.Sprintf("error.%d", ErrMenuDisabled), "菜单已禁用"),
		ErrMenuExists:           getI18nMessage(fmt.Sprintf("error.%d", ErrMenuExists), "菜单已存在"),
		ErrMenuNameExists:       getI18nMessage(fmt.Sprintf("error.%d", ErrMenuNameExists), "菜单名已存在"),
		ErrMenuCodeExists:       getI18nMessage(fmt.Sprintf("error.%d", ErrMenuCodeExists), "菜单编码已存在"),
		ErrMenuNameInvalid:      getI18nMessage(fmt.Sprintf("error.%d", ErrMenuNameInvalid), "菜单名无效"),
		ErrMenuCodeInvalid:      getI18nMessage(fmt.Sprintf("error.%d", ErrMenuCodeInvalid), "菜单编码无效"),
		ErrMenuStatusInvalid:    getI18nMessage(fmt.Sprintf("error.%d", ErrMenuStatusInvalid), "菜单状态无效"),
		ErrMenuHasChildren:      getI18nMessage(fmt.Sprintf("error.%d", ErrMenuHasChildren), "菜单下有子菜单"),
		ErrMenuHasButtons:       getI18nMessage(fmt.Sprintf("error.%d", ErrMenuHasButtons), "菜单下有按钮"),
		ErrMenuHasApis:          getI18nMessage(fmt.Sprintf("error.%d", ErrMenuHasApis), "菜单下有API"),
		ErrMenuParentNotFound:   getI18nMessage(fmt.Sprintf("error.%d", ErrMenuParentNotFound), "父菜单不存在"),
		ErrMenuParentDisabled:   getI18nMessage(fmt.Sprintf("error.%d", ErrMenuParentDisabled), "父菜单已禁用"),
		ErrMenuParentInvalid:    getI18nMessage(fmt.Sprintf("error.%d", ErrMenuParentInvalid), "父菜单无效"),
		ErrMenuTypeInvalid:      getI18nMessage(fmt.Sprintf("error.%d", ErrMenuTypeInvalid), "菜单类型无效"),
		ErrMenuPathInvalid:      getI18nMessage(fmt.Sprintf("error.%d", ErrMenuPathInvalid), "菜单路径无效"),
		ErrMenuComponentInvalid: getI18nMessage(fmt.Sprintf("error.%d", ErrMenuComponentInvalid), "菜单组件无效"),
		ErrMenuPermInvalid:      getI18nMessage(fmt.Sprintf("error.%d", ErrMenuPermInvalid), "菜单权限无效"),
		ErrMenuIconInvalid:      getI18nMessage(fmt.Sprintf("error.%d", ErrMenuIconInvalid), "菜单图标无效"),
		ErrMenuSortInvalid:      getI18nMessage(fmt.Sprintf("error.%d", ErrMenuSortInvalid), "菜单排序无效"),

		// 登录日志模块错误
		ErrLoginLogNotFound:    getI18nMessage(fmt.Sprintf("error.%d", ErrLoginLogNotFound), "登录日志不存在"),
		ErrLoginLogGetList:     getI18nMessage(fmt.Sprintf("error.%d", ErrLoginLogGetList), "获取登录日志列表失败"),
		ErrLoginLogGetByID:     getI18nMessage(fmt.Sprintf("error.%d", ErrLoginLogGetByID), "根据ID获取登录日志失败"),
		ErrLoginLogDelete:      getI18nMessage(fmt.Sprintf("error.%d", ErrLoginLogDelete), "删除登录日志失败"),
		ErrLoginLogBatchDelete: getI18nMessage(fmt.Sprintf("error.%d", ErrLoginLogBatchDelete), "批量删除登录日志失败"),
		ErrLoginLogCleanOld:    getI18nMessage(fmt.Sprintf("error.%d", ErrLoginLogCleanOld), "清理旧登录日志失败"),

		// 操作日志模块错误
		ErrOperationLogNotFound:    getI18nMessage(fmt.Sprintf("error.%d", ErrOperationLogNotFound), "操作日志不存在"),
		ErrOperationLogGetList:     getI18nMessage(fmt.Sprintf("error.%d", ErrOperationLogGetList), "获取操作日志列表失败"),
		ErrOperationLogGetByID:     getI18nMessage(fmt.Sprintf("error.%d", ErrOperationLogGetByID), "根据ID获取操作日志失败"),
		ErrOperationLogDelete:      getI18nMessage(fmt.Sprintf("error.%d", ErrOperationLogDelete), "删除操作日志失败"),
		ErrOperationLogBatchDelete: getI18nMessage(fmt.Sprintf("error.%d", ErrOperationLogBatchDelete), "批量删除操作日志失败"),
		ErrOperationLogCleanOld:    getI18nMessage(fmt.Sprintf("error.%d", ErrOperationLogCleanOld), "清理旧操作日志失败"),
		ErrOperationLogInvalidID:   getI18nMessage(fmt.Sprintf("error.%d", ErrOperationLogInvalidID), "操作日志ID无效"),
	}
}

// GetMessage 获取错误码对应的消息（支持i18n）
func GetMessage(code int) string {
	// 尝试获取国际化消息
	messageKey := fmt.Sprintf("error.%d", code)
	localizedMessage := i18n.T(messageKey)

	// 如果找到翻译，使用翻译消息
	if localizedMessage != messageKey {
		return localizedMessage
	}

	// 否则使用默认消息
	codeMessageMap := getCodeMessageMap()
	if msg, ok := codeMessageMap[code]; ok {
		return msg
	}
	return "未知错误"
}

// IsSystemError 判断是否为系统级错误
func IsSystemError(code int) bool {
	return code >= 10000 && code < 20000
}

// IsBusinessError 判断是否为业务级错误
func IsBusinessError(code int) bool {
	return code >= 20000 && code < 30000
}
