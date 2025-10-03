package impl

import (
	"errors"
	"time"

	"github.com/cuiyuanxin/kunpeng/internal/interfaces/service"
	"github.com/cuiyuanxin/kunpeng/internal/model"
	"github.com/cuiyuanxin/kunpeng/internal/model/dto"
	"github.com/cuiyuanxin/kunpeng/internal/repository"
	kperrors "github.com/cuiyuanxin/kunpeng/pkg/errors"
	"github.com/cuiyuanxin/kunpeng/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	loginAttemptService service.LoginAttemptService
}

// NewUserService 创建用户服务实例
func NewUserService(loginAttemptService service.LoginAttemptService) *UserServiceImpl {
	return &UserServiceImpl{
		loginAttemptService: loginAttemptService,
	}
}

// Login 用户登录
func (s *UserServiceImpl) Login(req *dto.UserLoginReq, clientIP string) (*dto.UserLoginResp, error) {
	// 检查是否被拉黑
	if s.loginAttemptService != nil {
		blocked, err := s.loginAttemptService.IsBlocked(req.Account, clientIP)
		if err != nil {
			return nil, kperrors.New(kperrors.ErrSystem, err)
		}
		if blocked {
			return nil, kperrors.New(kperrors.ErrAuthLocked, nil).WithMessage("登录失败次数过多，账号已被锁定2小时")
		}
	}

	var user *model.User
	var err error

	// 根据登录类型查询用户
	switch req.LoginType {
	case "username":
		user, err = repository.GetUserRepository().FindByUsername(req.Account)
	case "mobile":
		user, err = repository.GetUserRepository().FindByMobile(req.Account)
	default:
		return nil, kperrors.New(kperrors.ErrParam, nil)
	}

	if err != nil {
		if kperrors.IsCode(err, kperrors.ErrDBNotFound) {
			return nil, kperrors.New(kperrors.ErrUserNotFound, err)
		}
		return nil, err
	}

	// 检查用户状态
	if user.Status == 0 {
		return nil, kperrors.New(kperrors.ErrUserDisabled, nil)
	}
	if user.Status == 2 {
		return nil, kperrors.New(kperrors.ErrUserLocked, nil)
	}

	// 根据登录类型进行不同的验证
	var loginSuccess bool
	switch req.LoginType {
	case "username":
		// 账号登录：验证密码
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			// 记录登录失败
			if s.loginAttemptService != nil {
				s.loginAttemptService.CheckAndRecordAttempt(req.Account, clientIP, false)
			}
			return nil, kperrors.New(kperrors.ErrUserPassword, err)
		}
		loginSuccess = true
	case "mobile":
		// 手机号登录：验证验证码
		// TODO: 这里应该调用验证码服务验证验证码的有效性
		// 暂时跳过验证码验证，实际项目中需要实现验证码验证逻辑
		if req.Captcha == "" {
			// 记录登录失败
			if s.loginAttemptService != nil {
				s.loginAttemptService.CheckAndRecordAttempt(req.Account, clientIP, false)
			}
			return nil, kperrors.New(kperrors.ErrParam, errors.New("验证码不能为空"))
		}
		loginSuccess = true
	}

	// 记录登录成功
	if s.loginAttemptService != nil && loginSuccess {
		s.loginAttemptService.CheckAndRecordAttempt(req.Account, clientIP, true)
	}

	// 生成Token对（支持记住我功能）
	tokenPair, err := jwt.GenerateTokenPair(user.ID, user.Username, user.RoleID, user.AppKey, user.AppSecret, req.RememberMe)
	if err != nil {
		return nil, kperrors.New(kperrors.ErrSystem, err)
	}

	// 更新登录信息
	now := time.Now()
	user.LoginTime = &now
	user.LoginIP = clientIP
	repository.GetUserRepository().Update(user)

	// 返回结果
	return &dto.UserLoginResp{
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		ExpiresIn:        tokenPair.ExpiresIn,
		RefreshExpiresIn: tokenPair.RefreshExpiresIn,
	}, nil
}

// RefreshToken 刷新token
func (s *UserServiceImpl) RefreshToken(req *dto.RefreshTokenReq) (*dto.UserLoginResp, error) {
	// 解析refresh token获取用户信息
	claims, err := jwt.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, kperrors.New(kperrors.ErrInvalidToken, err)
	}

	// 验证token类型必须是refresh token
	if claims.TokenType != jwt.RefreshTokenType {
		return nil, kperrors.New(kperrors.ErrInvalidToken, nil)
	}

	// 检查用户是否存在且状态正常
	user, err := repository.GetUserRepository().FindByID(claims.UserID)
	if err != nil {
		if kperrors.IsCode(err, kperrors.ErrDBNotFound) {
			return nil, kperrors.New(kperrors.ErrUserNotFound, err)
		}
		return nil, err
	}

	// 检查用户状态
	if user.Status == 0 {
		return nil, kperrors.New(kperrors.ErrUserDisabled, nil)
	}
	if user.Status == 2 {
		return nil, kperrors.New(kperrors.ErrUserLocked, nil)
	}

	// 生成新的token对
	tokenPair, err := jwt.GenerateTokenPair(user.ID, user.Username, user.RoleID, user.AppKey, user.AppSecret, claims.RememberMe)
	if err != nil {
		return nil, kperrors.New(kperrors.ErrSystem, err)
	}

	// 返回结果
	return &dto.UserLoginResp{
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		ExpiresIn:        tokenPair.ExpiresIn,
		RefreshExpiresIn: tokenPair.RefreshExpiresIn,
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserServiceImpl) GetUserInfo(userID uint) (*dto.UserInfoResp, error) {
	// 查询用户
	user, err := repository.GetUserRepository().FindByID(userID)
	if err != nil {
		return nil, err
	}

	// 查询角色、部门、岗位信息
	role, _ := repository.GetRoleRepository().FindByID(user.RoleID)
	dept, _ := repository.GetDeptRepository().FindByID(user.DeptID)
	post, _ := repository.GetPostRepository().FindByID(user.PostID)

	roleName := ""
	deptName := ""
	postName := ""

	if role != nil {
		roleName = role.Name
	}
	if dept != nil {
		deptName = dept.Name
	}
	if post != nil {
		postName = post.Name
	}

	return &dto.UserInfoResp{
		ID:       user.ID,
		Username: user.Username,
		Nickname: user.Nickname,
		RealName: user.RealName,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		RoleName: roleName,
		DeptName: deptName,
		PostName: postName,
	}, nil
}

// GetUserList 获取用户列表
func (s *UserServiceImpl) GetUserList(req *dto.UserPageReq) (*dto.PageResp, error) {
	// 查询用户列表
	users, total, err := repository.GetUserRepository().FindList(req)
	if err != nil {
		return nil, err
	}

	// 计算总页数
	totalPages := int(total) / req.PageSize
	if int(total)%req.PageSize != 0 {
		totalPages++
	}

	return &dto.PageResp{
		List:       users,
		Total:      total,
		PageNum:    req.PageNum,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetUserByID 根据ID获取用户
func (s *UserServiceImpl) GetUserByID(id uint) (*model.User, error) {
	return repository.GetUserRepository().FindByID(id)
}

// CreateUser 创建用户
func (s *UserServiceImpl) CreateUser(req *dto.UserCreateReq) (uint, error) {
	// 检查用户名是否存在
	existUser, err := repository.GetUserRepository().FindByUsername(req.Username)
	if err == nil && existUser != nil {
		return 0, kperrors.New(kperrors.ErrUserNameExists, nil)
	} else if err != nil {
		// 如果错误不是"记录未找到"，则返回错误
		var kpErr *kperrors.Error
		if !errors.As(err, &kpErr) || kpErr.Code != kperrors.ErrUserNotFound {
			return 0, err
		}
	}

	// 生成密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, kperrors.New(kperrors.ErrSystem, err)
	}

	// 生成AppKey和AppSecret
	appKey, appSecret := jwt.GenerateAppKeyAndSecret(req.Username)

	// 创建用户
	user := model.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Nickname:  req.Nickname,
		RealName:  req.RealName,
		Avatar:    req.Avatar,
		Gender:    req.Gender,
		Email:     req.Email,
		Mobile:    req.Mobile,
		DeptID:    req.DeptID,
		PostID:    req.PostID,
		RoleID:    req.RoleID,
		Status:    req.Status,
		AppKey:    appKey,
		AppSecret: appSecret,
		Remark:    req.Remark,
	}

	err = repository.GetUserRepository().Create(&user)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// UpdateUser 更新用户
func (s *UserServiceImpl) UpdateUser(req *dto.UserUpdateReq) error {
	// 检查用户是否存在
	user, err := repository.GetUserRepository().FindByID(req.ID)
	if err != nil {
		return err
	}

	// 更新用户信息
	user.Nickname = req.Nickname
	user.RealName = req.RealName
	user.Avatar = req.Avatar
	user.Gender = req.Gender
	user.Email = req.Email
	user.Mobile = req.Mobile
	user.DeptID = req.DeptID
	user.PostID = req.PostID
	user.RoleID = req.RoleID
	user.Status = req.Status
	user.Remark = req.Remark

	return repository.GetUserRepository().Update(user)
}

// DeleteUser 删除用户
func (s *UserServiceImpl) DeleteUser(id uint) error {
	return repository.GetUserRepository().Delete(id)
}

// BatchDeleteUser 批量删除用户
func (s *UserServiceImpl) BatchDeleteUser(ids []uint) error {
	return repository.GetUserRepository().BatchDelete(ids)
}

// ChangeUserStatus 修改用户状态
func (s *UserServiceImpl) ChangeUserStatus(req *dto.StatusReq) error {
	return repository.GetUserRepository().UpdateStatus(req.ID, int(req.Status))
}

// ResetUserPassword 重置用户密码
func (s *UserServiceImpl) ResetUserPassword(id uint) error {
	// 生成新密码哈希(默认密码为123456)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	if err != nil {
		return kperrors.New(kperrors.ErrSystem, err)
	}

	return repository.GetUserRepository().ResetPassword(id, string(hashedPassword))
}

// ChangePassword 修改密码
func (s *UserServiceImpl) ChangePassword(userID uint, req *dto.UserChangePasswordReq) error {
	// 检查用户是否存在
	user, err := repository.GetUserRepository().FindByID(userID)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return kperrors.New(kperrors.ErrUserOldPassword, err)
	}

	// 生成新密码哈希
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return kperrors.New(kperrors.ErrSystem, err)
	}

	return repository.GetUserRepository().UpdatePassword(userID, string(hashedPassword))
}
