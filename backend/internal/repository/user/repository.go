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
	UpdateBlocked(ctx context.Context, userID uint, isBlocked bool) error                        // 更新用户状态：禁用
	UpdateActive(ctx context.Context, userID uint, isActive bool) error                          // 更新用户状态：激活
	UpdateRole(ctx context.Context, userID uint, role string) error                              // 更新用户角色
	SoftDelete(ctx context.Context, userID uint) error                                           // 软删除用户
	HardDelete(ctx context.Context, userID uint) error                                           // 硬删除用户
	UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error                // 更新用户密码
	InvalidateUserTokens(ctx context.Context, userID uint) error                                 // 使用户的所有令牌失效
	RestoreDeleted(ctx context.Context, userID uint) error                                       // 恢复软删除的用户
	SetTempPasswordFlag(ctx context.Context, userID uint, isTemp bool, expireAt time.Time) error // 设置临时密码标志
	ClearTempPasswordFlag(ctx context.Context, userID uint) error                                // 清除临时密码标志
	BatchUpdateBlocked(ctx context.Context, userIDs []uint, isBlocked bool) (int64, error)       // 批量更新用户状态：禁用
	//  crud
	Create(user *model.User) error                                              // 创建用户
	FindByID(id uint) (*model.User, error)                                      // 根据ID查找用户
	FindByEmail(ctx context.Context, email string) (*model.User, error)         // 根据邮箱查找用户
	FindByUsername(username string) (*model.User, error)                        // 根据用户名查找用户
	Update(ctx context.Context, user *model.User) error                         // 更新用户
	UpdateFields(id uint, fields map[string]any) error                          // 更新用户字段
	List(page, pageSize int, keyword string) ([]model.User, int64, error)       // 分页查询用户
	FindByIDs(ids []uint) ([]model.User, error)                                 // 根据ID列表查找用户
	GetUserBasicInfo(id uint) (*model.User, error)                              // 获取用户基本信息
	GetUserBasicInfoById(userID uint) (*model.User, error)                      // 根据ID获取用户基本信息
	GetUserRoleById(userID uint) (string, error)                                // 根据ID获取用户角色
	FindByEmailUnscoped(ctx context.Context, email string) (*model.User, error) // 根据邮箱查找用户（不使用软删除）
	// follow
	Follow(followerID, followingID uint) error                                 // 关注用户
	Unfollow(followerID, followingID uint) error                               // 取消关注用户
	IsFollowing(followerID, followingID uint) bool                             // 检查用户是否关注
	GetFollowerCount(userID uint) int64                                        // 获取关注者数量
	GetFollowingCount(userID uint) int64                                       // 获取关注数量
	GetFollowers(userID uint, page, pageSize int) ([]model.User, int64, error) // 获取关注者列表
	GetFollowing(userID uint, page, pageSize int) ([]model.User, int64, error) // 获取关注列表
	// rank
	GetTopScoreUsersSimple(ctx context.Context, limit int, excludeBlocked bool) ([]dto.LeaderboardUserSimple, error) // 获取积分排名前N的用户
	GetTopScoreUsersDetail(ctx context.Context, limit int, excludeBlocked bool) ([]dto.LeaderboardUserDetail, error) // 获取积分排名前N的用户详细信息
	// score
	GetScoreById(userID uint) (int, error)                 // 根据ID获取用户积分
	GetUsersScoreTotal() (int, error)                      // 获取所有用户的总积分
	GetEveryoneUsersScore() ([]model.User, error)          // 获取所有用户的积分
	AddScore(userID uint, score int) error                 // 增加用户积分
	DeductScore(tx *gorm.DB, userID uint, score int) error // 扣除用户积分
	SetScoreById(id uint, score int) error                 // 设置用户积分
	// stats
	Count(ctx context.Context) (int64, error)                                                                         // 获取用户总数
	CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)                                // 获取指定日期范围内的用户总数
	CountActiveByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error)                          // 获取指定日期范围内活跃用户总数
	GetActiveUsersByDateRange(ctx context.Context, startDate, endDate time.Time, limit int) ([]*ActiveUserRow, error) // 获取指定日期范围内活跃用户
	IsUserExistsByEmail(email string) (bool, error)                                                                   // 检查用户是否存在

}

func NewUserRepository(db *gorm.DB, tokenRepo token.TokenRepository) UserRepository {
	return &userRepository{db: db, tokenRepo: tokenRepo}
}
