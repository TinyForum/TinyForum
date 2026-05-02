package user

import (
	"context"
	"tiny-forum/internal/infra/validator"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	commentRepo "tiny-forum/internal/repository/comment"
	postRepo "tiny-forum/internal/repository/post"
	userRepo "tiny-forum/internal/repository/user"
	"tiny-forum/internal/service/notification"
	jwtpkg "tiny-forum/pkg/jwt"
)

type UserService interface {
	// admin crud
	List(page, pageSize int, keyword string) ([]do.User, int64, error)
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
	GetFollowers(userID uint, page, pageSize int) ([]do.User, int64, error)
	GetFollowing(userID uint, page, pageSize int) ([]do.User, int64, error)
	// leaderboard
	GetSimpleLeaderboardData(ctx context.Context, limit int) ([]dto.LeaderboardUserSimple, error)
	GetDetailLeaderboardData(ctx context.Context, limit int) ([]dto.LeaderboardUserDetail, error)
	// profile
	GetProfile(userID uint) (*do.User, error)
	UpdateProfile(userID uint, input do.UpdateProfileInput) error
	GetUserProfile(targetID, viewerID uint) (*UserProfileResponse, error)
	GetUserBasicInfo(userID uint) (*do.User, error)
	GetUserRoleById(userID uint) (string, error)
	// stats
	GetGlobalStatsCount(ctx context.Context, userID uint) (*dto.StatsInfo, error)
	// posts
	GetUserPosts(ctx context.Context, req request.GetUserPostsRequest, userID uint) (*vo.BasicPageData, error)
}
type userService struct {
	repo        userRepo.UserRepository
	jwtMgr      *jwtpkg.JWTManager
	notifSvc    notification.NotificationService
	roleChecker *validator.RoleChangeChecker
	postRepo    postRepo.PostRepository
	commentRepo commentRepo.CommentRepository
	// roleChange  validator.RoleChangeRequest
}

func NewUserService(
	repo userRepo.UserRepository,
	jwtMgr *jwtpkg.JWTManager,
	notifSvc notification.NotificationService,
	postRepo postRepo.PostRepository,
	commetnRepo commentRepo.CommentRepository,
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
		postRepo:    postRepo,
		commentRepo: commetnRepo,
	}
}
