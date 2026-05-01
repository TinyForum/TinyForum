package auth

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
	"tiny-forum/internal/model/po"
	userSvc "tiny-forum/internal/service/user"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 保留用户名黑名单（防止注册 admin、root 等特权名）
var reservedUsernames = map[string]bool{
	"admin": true, "root": true, "system": true, "moderator": true,
	"support": true, "help": true, "info": true, "contact": true,
	"superadmin": true, "administrator": true, "staff": true, "official": true,
}

// usernameRegex: 只允许字母、数字、下划线、连字符，长度 3-30
//
//	防止用户名包含 XSS/注入特殊字符
var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`)

// Register 用户注册
func (s *authService) Register(ctx context.Context, input userSvc.RegisterInput) (*userSvc.AuthResult, error) {
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

	user := &po.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashed),
		Role:     po.RoleUser,
		Avatar:   avatarURL(input.Username),
	}
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	s.notifSvc.Create(user.ID, nil, po.NotifySystem, "欢迎加入 Tiny Forum！", nil, "")

	token, err := s.jwtMgr.Generate(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &userSvc.AuthResult{Token: token, User: user}, nil
}

type AuthResult struct {
	Token          string          `json:"token"`
	User           *po.User     `json:"user"`
	DeletionStatus *DeletionStatus `json:"deletion_status,omitempty"`
}

// Login 用户登录
func (s *authService) Login(ctx context.Context, input userSvc.LoginInput) (*AuthResult, error) {
	user, err := s.userRepo.FindByEmailUnscoped(ctx, input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("邮箱或密码错误") // 不区分邮箱/密码错误
		}
		return nil, err
	}

	if user.IsBlocked {
		return nil, errors.New("账户已被禁用")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("邮箱或密码错误")
	}

	var deletionStatus *DeletionStatus
	if user.DeletedAt.Valid {
		remainingDays := 30 - int(time.Since(user.DeletedAt.Time).Hours()/24)
		deletionStatus = &DeletionStatus{
			IsDeleted:     true,
			DeletedAt:     &user.DeletedAt.Time,
			CanRestore:    remainingDays > 0,
			RemainingDays: remainingDays,
		}
		if remainingDays <= 0 {
			return nil, errors.New("账户已永久删除，无法登录")
		}
	}

	now := time.Now()
	user.LastLogin = &now
	_ = s.userRepo.Update(ctx, user)

	token, err := s.jwtMgr.Generate(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &AuthResult{
		Token:          token,
		User:           user,
		DeletionStatus: deletionStatus,
	}, nil
}

// RevokeToken  将 Token 加入黑名单，注销后无法再使用
// 写入 tokenRepo 的黑名单，中间件验证时检查黑名单
func (s *authService) RevokeToken(ctx context.Context, rawToken string) error {
	// 解析 token 获取 JTI 和过期时间，精确存储黑名单
	claims, err := s.jwtMgr.Parse(rawToken)
	if err != nil {
		// token 已无效，无需加入黑名单
		return nil
	}
	// 将 JTI 存入黑名单，TTL 与 token 剩余有效期一致，节省存储
	return s.tokenRepo.RevokeToken(ctx, claims.ID, claims.ExpiresAt.Time)
}
func (s *authService) FinduUserEmailByID(userID uint) (string, error) {
	user := &po.User{}
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}
	return user.Email, nil

}

// ChangePassword 修改密码
func (s *authService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) (string, error) {
	// 1. 获取用户信息
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperrors.ErrUserNotFound
		}
		return "", fmt.Errorf("查询用户失败: %w", err)
	}

	// 2. 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return "", apperrors.ErrInvalidCurrentPassword
	}

	// 3. 禁止与旧密码相同
	if oldPassword == newPassword {
		return "", apperrors.ErrPasswordSameAsOld
	}

	// 4. 密码强度校验
	if err := validatePasswordStrength(newPassword); err != nil {
		return "", err
	}

	// 5. 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("密码加密失败: %w", err)
	}

	// 6. 更新密码
	if err := s.userRepo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
		return "", fmt.Errorf("更新密码失败: %w", err)
	}

	// 7. 吊销所有 refresh token（强制其他设备下线）
	if err := s.tokenRepo.RevokeAllUserTokens(ctx, userID); err != nil {
		// 记录日志但不影响主流程
		logger.Warnf("警告: 吊销 token 失败: %v", err)
	}

	// 8. 返回消息（根据密码强度）
	strength := s.checkPasswordStrength(newPassword)
	if strength == "weak" {
		return "密码修改成功，但密码强度较弱，建议使用更复杂的密码", nil
	}
	return "密码修改成功", nil
}

// validatePasswordStrength 统一密码强度校验（注册、修改、重置均使用）
// 规则：至少 8 位，包含字母和数字
func validatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("密码长度不得少于 8 位")
	}
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasLetter || !hasDigit {
		return apperrors.ErrInvalidPassword.WithMessage("密码必须同时包含字母和数字")
	}
	return nil
}

// isValidEmail 严格邮箱校验，过滤换行符等注入字符
func isValidEmail(email string) bool {
	// 拒绝包含换行、回车、空格等危险字符（防邮件头注入）
	if strings.ContainsAny(email, "\r\n\t ") {
		return false
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email) && len(email) <= 254
}

// checkPasswordStrength 密码强度评级（用于提示，不强制）
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
