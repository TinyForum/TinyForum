package service

import (
	"errors"
	"fmt"
	"time"

	apperrors "tiny-forum/internal/errors"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
	jwtpkg "tiny-forum/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo        *repository.UserRepository
	jwtMgr      *jwtpkg.Manager
	notifSvc    *NotificationService
	roleChecker RoleChangeChecker
}

func NewUserService(
	repo *repository.UserRepository,
	jwtMgr *jwtpkg.Manager,
	notifSvc *NotificationService,
) *UserService {
	return &UserService{
		repo:        repo,
		jwtMgr:      jwtMgr,
		notifSvc:    notifSvc,
		roleChecker: RoleChangeChecker{},
	}
}

// ── Auth ─────────────────────────────────────────────────────────────────────

type RegisterInput struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResult struct {
	Token string      `json:"token"`
	User  *model.User `json:"user"`
}

func (s *UserService) Register(input RegisterInput) (*AuthResult, error) {
	if _, err := s.repo.FindByUsername(input.Username); err == nil {
		return nil, errors.New("用户名已被占用")
	}
	if _, err := s.repo.FindByEmail(input.Email); err == nil {
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
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	s.notifSvc.Create(user.ID, nil, model.NotifySystem, "欢迎加入 Tiny Forum！", nil, "")

	token, err := s.jwtMgr.Generate(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &AuthResult{Token: token, User: user}, nil
}

func (s *UserService) Login(input LoginInput) (*AuthResult, error) {
	user, err := s.repo.FindByEmail(input.Email)
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
	_ = s.repo.Update(user)

	token, err := s.jwtMgr.Generate(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}
	return &AuthResult{Token: token, User: user}, nil
}

// ── Profile ──────────────────────────────────────────────────────────────────

func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	return s.repo.FindByID(userID)
}

type UpdateProfileInput struct {
	Bio    string `json:"bio" binding:"max=500"`
	Avatar string `json:"avatar"`
}

func (s *UserService) UpdateProfile(userID uint, input UpdateProfileInput) error {
	fields := map[string]interface{}{}
	if input.Bio != "" {
		fields["bio"] = input.Bio
	}
	if input.Avatar != "" {
		fields["avatar"] = input.Avatar
	}
	if len(fields) == 0 {
		return nil
	}
	return s.repo.UpdateFields(userID, fields)
}

type UserProfileResponse struct {
	*model.User
	FollowerCount  int64 `json:"follower_count"`
	FollowingCount int64 `json:"following_count"`
	IsFollowing    bool  `json:"is_following"`
}

func (s *UserService) GetUserProfile(targetID, viewerID uint) (*UserProfileResponse, error) {
	user, err := s.repo.FindByID(targetID)
	if err != nil {
		return nil, err
	}
	resp := &UserProfileResponse{
		User:           user,
		FollowerCount:  s.repo.GetFollowerCount(targetID),
		FollowingCount: s.repo.GetFollowingCount(targetID),
	}
	if viewerID > 0 {
		resp.IsFollowing = s.repo.IsFollowing(viewerID, targetID)
	}
	return resp, nil
}

// ── Follow ───────────────────────────────────────────────────────────────────

func (s *UserService) Follow(followerID, followingID uint) error {
	if followerID == followingID {
		return errors.New("不能关注自己")
	}
	if err := s.repo.Follow(followerID, followingID); err != nil {
		return err
	}
	following, _ := s.repo.FindByID(followingID)
	if following != nil {
		s.notifSvc.Create(followingID, &followerID, model.NotifyFollow,
			following.Username+" 关注了你", nil, "")
	}
	return nil
}

func (s *UserService) Unfollow(followerID, followingID uint) error {
	return s.repo.Unfollow(followerID, followingID)
}

// 获取关注者列表
func (s *UserService) GetFollowers(userID uint, page, pageSize int) ([]model.User, int64, error) {
	return s.repo.GetFollowers(userID, page, pageSize)
}

// 获取关注列表
func (s *UserService) GetFollowing(userID uint, page, pageSize int) ([]model.User, int64, error) {
	return s.repo.GetFollowing(userID, page, pageSize)
}

// 查询积分
func (s *UserService) GetScoreById(userID uint) (int, error) {
	return s.repo.GetScoreById(userID)
}

// ── Admin ────────────────────────────────────────────────────────────────────

func (s *UserService) GetLeaderboard(limit int) ([]model.User, error) {
	return s.repo.GetTopUsers(limit)
}

func (s *UserService) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *UserService) SetActive(userID uint, active bool) error {
	return s.repo.UpdateFields(userID, map[string]interface{}{"is_active": active})
}

func (s *UserService) SetBlocked(userID uint, blocked bool) error {
	return s.repo.UpdateFields(userID, map[string]interface{}{"is_blocked": blocked})
}

func (s *UserService) SetScoreById(userID uint, score int) error {
	// 1. 参数验证
	if userID == 0 {
		return errors.New("用户ID不能为空")
	}

	err := s.repo.SetScoreById(userID, score)
	if err != nil {
		return fmt.Errorf("设置积分失败: %w", err)
	}

	// 4. 可选：触发积分变更事件（如发送通知、更新缓存等）
	go s.onScoreChanged(userID, score)

	return nil
}

// 积分变更后的回调处理
func (s *UserService) onScoreChanged(userID uint, newScore int) {
	// 可以在这里添加：
	// - 发送系统通知
	// - 更新Redis缓存
	// - 检查是否触发等级变更
	// - 记录日志到消息队列等
}

// ── Role Management ──────────────────────────────────────────────────────────

// SetRole 变更目标用户角色（含细粒度权限校验）。
// operatorID 来自 JWT，确保操作者身份可信。
func (s *UserService) SetRole(operatorID, targetID uint, newRole string) error {
	// 1. 校验新角色字面值
	targetRole := model.UserRole(newRole)
	if !model.IsValidRole(targetRole) {
		return fmt.Errorf("%w: %s", apperrors.ErrInvalidRole, newRole)
	}

	// 2. 加载操作者与目标用户
	operator, err := s.repo.FindByID(operatorID)
	if err != nil {
		return fmt.Errorf("操作者不存在: %w", err)
	}
	target, err := s.repo.FindByID(targetID)
	if err != nil {
		return err // 保留 gorm.ErrRecordNotFound 供 handler 判断
	}

	// 3. 幂等：角色无变化直接返回
	if target.Role == targetRole {
		return nil
	}

	// 4. 细粒度权限校验（职责交由 RoleChangeChecker）
	if err := s.roleChecker.Check(RoleChangeRequest{
		Operator: operator,
		Target:   target,
		NewRole:  targetRole,
	}); err != nil {
		return err
	}

	// 5. 执行更新
	return s.repo.UpdateFields(targetID, map[string]interface{}{"role": newRole})
}

// ── 内部工具 ─────────────────────────────────────────────────────────────────

func avatarURL(username string) string {
	return "https://api.dicebear.com/8.x/lorelei/svg?seed=" + username
}

type LoginResult struct {
	Token string    `json:"-"` // json:"-" 防止意外序列化到响应
	User  *UserInfo `json:"user"`
}

type UserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
