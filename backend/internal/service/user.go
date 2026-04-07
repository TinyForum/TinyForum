package service

import (
	"errors"
	"time"

	"bbs-forum/internal/model"
	"bbs-forum/internal/repository"
	jwtpkg "bbs-forum/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	repo     *repository.UserRepository
	jwtMgr   *jwtpkg.Manager
	notifSvc *NotificationService
}

func NewUserService(repo *repository.UserRepository, jwtMgr *jwtpkg.Manager, notifSvc *NotificationService) *UserService {
	return &UserService{repo: repo, jwtMgr: jwtMgr, notifSvc: notifSvc}
}

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
	Token string     `json:"token"`
	User  *model.User `json:"user"`
}

func (s *UserService) Register(input RegisterInput) (*AuthResult, error) {
	// Check username
	if _, err := s.repo.FindByUsername(input.Username); err == nil {
		return nil, errors.New("用户名已被占用")
	}
	// Check email
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

	// Send welcome notification
	s.notifSvc.Create(user.ID, nil, model.NotifySystem, "欢迎加入 BBS Forum！", nil, "")

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

	if !user.IsActive {
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

func (s *UserService) Follow(followerID, followingID uint) error {
	if followerID == followingID {
		return errors.New("不能关注自己")
	}
	if err := s.repo.Follow(followerID, followingID); err != nil {
		return err
	}
	// Notification
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

func (s *UserService) GetLeaderboard(limit int) ([]model.User, error) {
	return s.repo.GetTopUsers(limit)
}

func (s *UserService) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *UserService) SetActive(userID uint, active bool) error {
	return s.repo.UpdateFields(userID, map[string]interface{}{"is_active": active})
}

func avatarURL(username string) string {
	return "https://api.dicebear.com/8.x/initials/svg?seed=" + username
}
