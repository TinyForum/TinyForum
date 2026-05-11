package auth

import (
	"context"
	"errors"
	"strings"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"

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
		return nil, errors.New("该用户名不可用，请换一个")
	}

	// 邮箱格式严格校验（防止换行符等邮件头注入）
	if !isValidEmail(input.Email) {
		return nil, errors.New("邮箱格式无效")
	}

	if _, err := s.userRepo.FindByUsername(input.Username); err == nil {
		return nil, errors.New("用户名已被占用")
	}
	if _, err := s.userRepo.FindByEmail(ctx, input.Email); err == nil {
		return nil, errors.New("邮箱已被注册")
	}

	// 注册时强制密码强度校验
	if err := validatePasswordStrength(input.Password); err != nil {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &do.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashed),
		Role:     do.RoleUser,
		Avatar:   avatarURL(input.Username),
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	s.notifSvc.Create(user.ID, nil, do.NotifySystem, "欢迎加入 Tiny Forum！", nil, "")

	token, err := s.jwtMgr.Generate(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &vo.AuthResultVO{Token: token, User: &vo.UserVO{
		ID:       user.ID,
		Username: user.Username,
		Avatar:   user.Avatar,
	}}, nil
	// return &userSvc.AuthResult{User: user}, nil
}
