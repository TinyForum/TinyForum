package user

import (
	"context"
	"time"
	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository/token" // 假设 token 包路径

	"gorm.io/gorm"
)

type userRepository struct {
	db        *gorm.DB
	tokenRepo token.TokenRepository
}

type UserRepository interface {
	// admin
	UpdateBlocked(ctx context.Context, userID uint, isBlocked bool) error
	UpdateActive(ctx context.Context, userID uint, isActive bool) error
	UpdateRole(ctx context.Context, userID uint, role string) error
	SoftDelete(ctx context.Context, userID uint) error
	HardDelete(ctx context.Context, userID uint) error
	UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error
	InvalidateUserTokens(ctx context.Context, userID uint) error
	RestoreDeleted(ctx context.Context, userID uint) error
	SetTempPasswordFlag(ctx context.Context, userID uint, isTemp bool, expireAt time.Time) error
	ClearTempPasswordFlag(ctx context.Context, userID uint) error
	BatchUpdateBlocked(ctx context.Context, userIDs []uint, isBlocked bool) (int64, error)
	//  crud
	Create(user *model.User) error
	FindByID(id uint) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByResetToken(ctx context.Context, token string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdateFields(id uint, fields map[string]interface{}) error
	List(page, pageSize int, keyword string) ([]model.User, int64, error)
	FindByIDs(ids []uint) ([]model.User, error)
	GetUserBasicInfo(id uint) (*model.User, error)
	GetUserBasicInfoById(userID uint) (*model.User, error)
	GetUserRoleById(userID uint) (string, error)
	FindByEmailUnscoped(ctx context.Context, email string) (*model.User, error)
	// follow
	Follow(followerID, followingID uint) error
	Unfollow(followerID, followingID uint) error
	IsFollowing(followerID, followingID uint) bool
	GetFollowerCount(userID uint) int64
	GetFollowingCount(userID uint) int64
	GetFollowers(userID uint, page, pageSize int) ([]model.User, int64, error)
	GetFollowing(userID uint, page, pageSize int) ([]model.User, int64, error)
	// rank
	GetTopScoreUsersSimple(ctx context.Context, limit int, excludeBlocked bool) ([]dto.LeaderboardUserSimple, error)
	GetTopScoreUsersDetail(ctx context.Context, limit int, excludeBlocked bool) ([]dto.LeaderboardUserDetail, error)
	// score
	GetScoreById(userID uint) (int, error)
	GetUsersScoreTotal() (int, error)
	GetEveryoneUsersScore() ([]model.User, error)
	AddScore(userID uint, score int) error
	DeductScore(tx *gorm.DB, userID uint, score int) error
	SetScoreById(id uint, score int) error
	// stats
	Count(ctx context.Context) (int64, error)
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
	CountActiveByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)
	GetActiveUsersByDateRange(
		ctx context.Context,
		startDate, endDate time.Time,
		limit int,
	) ([]*ActiveUserRow, error)
}

func NewUserRepository(db *gorm.DB, tokenRepo token.TokenRepository) UserRepository {
	return &userRepository{db: db, tokenRepo: tokenRepo}
}
