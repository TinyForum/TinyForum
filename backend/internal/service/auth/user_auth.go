package auth

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"
	"tiny-forum/internal/model"
	userSvc "tiny-forum/internal/service/user"
	apperrors "tiny-forum/pkg/errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register 用户注册
func (s *authService) Register(ctx context.Context, input userSvc.RegisterInput) (*userSvc.AuthResult, error) {
	if _, err := s.userRepo.FindByUsername(input.Username); err == nil {
		return nil, errors.New("用户名已被占用")
	}
	if _, err := s.userRepo.FindByEmail(ctx, input.Email); err == nil {
		return nil, errors.New("邮箱已被注册")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashed),
		Role:     model.RoleUser,
		Avatar:   avatarURL(input.Username),
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	s.notifSvc.Create(user.ID, nil, model.NotifySystem, "欢迎加入 Tiny Forum！", nil, "")

	token, err := s.jwtMgr.Generate(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &userSvc.AuthResult{Token: token, User: user}, nil
}

// Login 用户登录
func (s *authService) Login(ctx context.Context, input userSvc.LoginInput) (*userSvc.AuthResult, error) {
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("邮箱或密码错误")
		}
		return nil, err
	}
	if user.IsBlocked {
		return nil, errors.New("账户已被禁用")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("邮箱或密码错误")
	}

	now := time.Now()
	user.LastLogin = &now
	_ = s.userRepo.Update(ctx, user)

	token, err := s.jwtMgr.Generate(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &userSvc.AuthResult{Token: token, User: user}, nil
}

// ChangePassword 修改密码
func (s *authService) ChangePassword(userID uint, oldPassword, newPassword string) (string, error) {
	ctx := context.Background()

	targetUser, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", userID)
		}
		return "", fmt.Errorf("查询目标用户失败: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(targetUser.Password), []byte(oldPassword)); err != nil {
		return "", apperrors.ErrInvalidPassword
	}
	if oldPassword == newPassword {
		return "", apperrors.ErrPasswordSameAsOld
	}
	if len(newPassword) < 6 {
		return "", apperrors.ErrPasswordTooShort
	}

	strength := s.checkPasswordStrength(newPassword)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return "", fmt.Errorf("更新密码失败: %w", err)
	}

	if strength == "weak" {
		return "密码修改成功，但密码强度较弱，建议使用更复杂的密码", nil
	}
	return "密码修改成功", nil
}

// checkPasswordStrength 检查密码强度（内部辅助）
func (s *authService) checkPasswordStrength(password string) string {
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)

	score := 0
	if len(password) >= 8 {
		score++
	}
	if len(password) >= 12 {
		score++
	}
	if hasUpper && hasLower {
		score++
	}
	if hasNumber {
		score++
	}
	if hasSpecial {
		score++
	}

	if score <= 2 {
		return "weak"
	}
	if score <= 4 {
		return "medium"
	}
	return "strong"
}

// avatarURL 生成默认头像URL
func avatarURL(username string) string {
	return "https://api.dicebear.com/8.x/lorelei/svg?seed=" + username
}
