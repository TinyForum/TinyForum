package auth

import (
	"context"
	"errors"
	"strings"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"

	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册
func (s *authService) Register(ctx context.Context, input request.RegisterRequest) (*vo.AuthResultVO, error) {
	//  校验用户名格式（只允许字母数字下划线连字符）
	if !usernameRegex.MatchString(input.Username) {
		return nil, errors.New("用户名只能包含字母、数字、下划线和连字符，长度 3-30 位")
	}

	// 检查保留用户名
	if reservedUsernames[strings.ToLower(input.Username)] {
		return nil, apperrors.ErrInvalidUsername
	}

	// 邮箱格式严格校验（防止换行符等邮件头注入）
	if !isValidEmail(input.Email) {
		logger.Debugf("无效的邮箱格式: %s", input.Email)
		return nil, apperrors.ErrInvalidEmail
	}

	if _, err := s.userRepo.FindByUsername(input.Username); err == nil {
		logger.Debugf("用户名已被注册: %s", input.Username)
		return nil, apperrors.ErrUserExist
	}
	if _, err := s.userRepo.FindByEmail(ctx, input.Email); err == nil {
		logger.Debugf("邮箱已被注册: %s", input.Email)
		return nil, apperrors.ErrInvalidEmail
	}

	// 注册时强制密码强度校验
	if err := validatePasswordStrength(input.Password); err != nil {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("密码加密失败: %v", err)
		return nil, apperrors.ErrInternalError
	}

	user := &do.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashed),
		Role:     do.RoleUser,
		AvatarUrl:   avatarURL(input.Username),
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	s.notifSvc.Create(user.ID, nil, do.NotifySystem, "欢迎加入 Tiny Forum！", nil, "")

	token, err := s.jwtMgr.Generate(user.ID, user.Username, string(user.Role))
	if err != nil {
		logger.Errorf("生成 jwt 失败: %v", err)
		return nil, apperrors.ErrInternalError
	}
	return &vo.AuthResultVO{
		Token: token,
		User: &vo.UserPrivateVO{
			ID:       user.ID,
			Username: user.Username,
			Avatar:   user.AvatarUrl,
			Email:    user.Email,
		}}, nil
}
