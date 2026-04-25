package user

import (
	"context"
	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"
	userRepo "tiny-forum/internal/repository/user"
	"tiny-forum/internal/service/notification"
	jwtpkg "tiny-forum/pkg/jwt"
	"tiny-forum/pkg/validator"
)

type UserService interface {
	// admin crud
	List(page, pageSize int, keyword string) ([]model.User, int64, error)
	DeleteUser(operatorID uint, targetID uint) error
	// admin password
	ResetUserPasswordWithTemp(operatorID uint, targetID uint) (string, error)
	ResetUserPassword(operatorID uint, targetID uint, newPassword string) error
	// adin score
	GetScoreById(userID uint) (int, error)
	SetScoreById(userID uint, score int) error
	onScoreChanged(userID uint, newScore int) error // TODO: 未完成
	GetAllUsersWithScore() ([]UserScoreResponse, error)
	// user status
	SetBlocked(targetID uint, operatorID uint, isBlocked bool) error
	SetActive(targetID uint, operatorID uint, isActive bool) error
	SetRole(operatorID, targetID uint, newRole string) error
	// follow
	Follow(followerID, followingID uint) error
	Unfollow(followerID, followingID uint) error
	GetFollowers(userID uint, page, pageSize int) ([]model.User, int64, error)
	GetFollowing(userID uint, page, pageSize int) ([]model.User, int64, error)
	// leaderboard
	GetSimpleLeaderboardData(ctx context.Context, limit int) ([]dto.LeaderboardUserSimple, error)
	GetDetailLeaderboardData(ctx context.Context, limit int) ([]dto.LeaderboardUserDetail, error)
	// profile
	GetProfile(userID uint) (*model.User, error)
	UpdateProfile(userID uint, input model.UpdateProfileInput) error
	GetUserProfile(targetID, viewerID uint) (*UserProfileResponse, error)
	GetUserBasicInfo(userID uint) (*model.User, error)
	GetUserRoleById(userID uint) (string, error)
}
type userService struct {
	repo        userRepo.UserRepository
	jwtMgr      *jwtpkg.Manager
	notifSvc    *notification.NotificationService // 注意：NotificationService 定义在别的包，需正确导入
	roleChecker *validator.RoleChangeChecker      // 改为指针类型
	// roleChange  validator.RoleChangeRequest
}

func NewUserService(
	repo userRepo.UserRepository,
	jwtMgr *jwtpkg.Manager,
	notifSvc *notification.NotificationService,
	// roleChange validator.RoleChangeRequest,
) UserService {
	roleValidator := validator.NewRoleValidator()
	roleChecker := validator.NewRoleChangeChecker(roleValidator)
	return &userService{
		repo:     repo,
		jwtMgr:   jwtMgr,
		notifSvc: notifSvc,
		// roleChecker: validator.RoleChangeChecker{},
		// roleChange:  validator.RoleChangeRequest{},
		roleChecker: roleChecker,
	}
}
