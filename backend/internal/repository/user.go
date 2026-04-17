// package repository

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"time"
// 	"tiny-forum/internal/model"

// 	"gorm.io/gorm"
// )

// type UserRepository struct {
// 	db        *gorm.DB
// 	tokenRepo *TokenRepository // 新增
// }

// func NewUserRepository(db *gorm.DB, tokenRepo *TokenRepository) *UserRepository {
// 	return &UserRepository{db: db, tokenRepo: tokenRepo}
// }

// // MARK: 用户信息
// func (r *UserRepository) Create(user *model.User) error {
// 	return r.db.Create(user).Error
// }

// // MARK: - dfd
// func (r *UserRepository) FindByID(id uint) (*model.User, error) {
// 	var user model.User
// 	err := r.db.First(&user, id).Error
// 	return &user, err
// }

// // MARK: - 获取用户的基本信息
// func (r *UserRepository) GetUserBasicInfo(id uint) (*model.User, error) {
// 	var user model.User
// 	err := r.db.Select("id, username, avatar").First(&user, id).Error
// 	return &user, err
// }

// // MARK: - 获取用户的基本信息
// func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
// 	var user model.User
// 	err := r.db.Where("email = ?", email).First(&user).Error
// 	return &user, err
// }

// // MARK: - 获取用户的基本信息
// func (r *UserRepository) FindByUsername(username string) (*model.User, error) {
// 	var user model.User
// 	err := r.db.Where("username = ?", username).First(&user).Error
// 	return &user, err
// }

// // MARK: - 更新
// func (r *UserRepository) Update(user *model.User) error {
// 	return r.db.Save(user).Error
// }

// // 更新字段
// func (r *UserRepository) UpdateFields(id uint, fields map[string]interface{}) error {
// 	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(fields).Error
// }

// // 列出用户
// func (r *UserRepository) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
// 	var users []model.User
// 	var total int64
// 	query := r.db.Model(&model.User{})
// 	if keyword != "" {
// 		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
// 	}
// 	if err := query.Count(&total).Error; err != nil {
// 		return nil, 0, err
// 	}
// 	offset := (page - 1) * pageSize
// 	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error
// 	return users, total, err
// }

// // MARK: Score

// // TopUsersQuery 排行榜查询参数
// type TopUsersQuery struct {
// 	Limit          int      // 数量限制，默认20，最大100
// 	ExcludeBlocked bool     // 是否排除封禁用户，默认true
// 	Fields         []string // 需要返回的字段，空则返回全部（谨慎使用）
// }

// func (r *UserRepository) GetTopUsers(ctx context.Context, query TopUsersQuery) ([]model.User, error) {
// 	// 默认值
// 	if query.Limit <= 0 {
// 		query.Limit = 20
// 	}
// 	if query.Limit > 100 {
// 		query.Limit = 100
// 	}

// 	dbQuery := r.db.WithContext(ctx).Model(&model.User{})

// 	// 字段筛选
// 	if len(query.Fields) > 0 {
// 		dbQuery = dbQuery.Select(query.Fields)
// 	}

// 	// 排除封禁用户
// 	if query.ExcludeBlocked {
// 		dbQuery = dbQuery.Where("is_blocked = ?", false)
// 	}

// 	var users []model.User
// 	err := dbQuery.
// 		Order("COALESCE(score, 0) DESC").
// 		Limit(query.Limit).
// 		Find(&users).Error

// 	dbQuery = dbQuery.Debug()
// 	if err != nil {
// 		return nil, fmt.Errorf("获取排行榜失败: %w", err)
// 	}
// 	if users == nil {
// 		return []model.User{}, nil
// 	}
// 	return users, nil
// }

// // 关注排名
// func (r *UserRepository) GetTopFollowers(limit int) ([]model.User, error) {
// 	var users []model.User
// 	err := r.db.Table("users").
// 		Select("users.*, COUNT(follows.follower_id) as follower_count").
// 		Joins("LEFT JOIN follows ON users.id = follows.following_id").
// 		Group("users.id").Order("follower_count DESC").Limit(limit).Find(&users).Error
// 	return users, err
// }

// // GetFollowing 获取用户关注的列表
// func (r *UserRepository) GetFollowing(userID uint, page, pageSize int) ([]model.User, int64, error) {
// 	var users []model.User
// 	var total int64

// 	offset := (page - 1) * pageSize

// 	// 获取关注总数
// 	err := r.db.Model(&model.Follow{}).
// 		Where("follower_id = ?", userID).
// 		Count(&total).Error
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// 获取关注的用户列表
// 	err = r.db.Model(&model.Follow{}).
// 		Select("users.*").
// 		Joins("JOIN users ON follows.following_id = users.id").
// 		Where("follows.follower_id = ?", userID).
// 		Offset(offset).
// 		Limit(pageSize).
// 		Find(&users).Error

// 	return users, total, err
// }

// // 用户查询自己的积分
// func (r *UserRepository) GetScoreById(userID uint) (int, error) {
// 	var user model.User
// 	err := r.db.Model(&model.User{}).Select("score").Where("id = ?", userID).First(&user).Error
// 	return user.Score, err
// }

// // GetAllUsersScore 获取所有用户的总积分
// func (r *UserRepository) GetUsersScoreTotal() (int, error) {
// 	var totalScore int64
// 	err := r.db.Model(&model.User{}).Select("COALESCE(SUM(score), 0)").Scan(&totalScore).Error
// 	return int(totalScore), err
// }

// // GetUsersScore 获取每个用户的积分
// func (r *UserRepository) GetEveryoneUsersScore() ([]model.User, error) {
// 	var users []model.User
// 	err := r.db.Model(&model.User{}).Select("id, score").Find(&users).Error
// 	return users, err
// }

// // MARK: Follow
// // Follow/unfollow
// func (r *UserRepository) Follow(followerID, followingID uint) error {
// 	follow := model.Follow{FollowerID: followerID, FollowingID: followingID}
// 	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).
// 		FirstOrCreate(&follow).Error
// }

// func (r *UserRepository) Unfollow(followerID, followingID uint) error {
// 	return r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).
// 		Delete(&model.Follow{}).Error
// }

// func (r *UserRepository) IsFollowing(followerID, followingID uint) bool {
// 	var count int64
// 	r.db.Model(&model.Follow{}).
// 		Where("follower_id = ? AND following_id = ?", followerID, followingID).
// 		Count(&count)
// 	return count > 0
// }

// func (r *UserRepository) GetFollowerCount(userID uint) int64 {
// 	var count int64
// 	r.db.Model(&model.Follow{}).Where("following_id = ?", userID).Count(&count)
// 	return count
// }

// func (r *UserRepository) GetFollowingCount(userID uint) int64 {
// 	var count int64
// 	r.db.Model(&model.Follow{}).Where("follower_id = ?", userID).Count(&count)
// 	return count
// }

// // GetFollowers 获取用户的粉丝列表（谁关注了该用户）
// func (r *UserRepository) GetFollowers(userID uint, page, pageSize int) ([]model.User, int64, error) {
// 	var users []model.User
// 	var total int64

// 	offset := (page - 1) * pageSize

// 	// 获取粉丝总数
// 	err := r.db.Model(&model.Follow{}).
// 		Where("following_id = ?", userID).
// 		Count(&total).Error
// 	if err != nil {
// 		return nil, 0, err
// 	}

// 	// 获取粉丝用户列表
// 	err = r.db.Model(&model.Follow{}).
// 		Select("users.*").
// 		Joins("JOIN users ON follows.follower_id = users.id").
// 		Where("follows.following_id = ?", userID).
// 		Offset(offset).
// 		Limit(pageSize).
// 		Find(&users).Error

// 	return users, total, err
// }

// // MARK: Score
// //
// //	扣减用户积分（使用事务）
// func (r *UserRepository) DeductScore(tx *gorm.DB, userID uint, score int) error {
// 	if score <= 0 {
// 		return nil
// 	}

// 	var user model.User
// 	if err := tx.Where("id = ?", userID).First(&user).Error; err != nil {
// 		return errors.New("用户不存在")
// 	}

// 	if user.Score < score {
// 		return errors.New("积分不足")
// 	}

// 	return tx.Model(&user).Update("score", gorm.Expr("score - ?", score)).Error
// }

// // 增加用户积分
// func (r *UserRepository) AddScore(userID uint, score int) error {
// 	return r.db.Model(&model.User{}).Where("id = ?", userID).
// 		UpdateColumn("score", gorm.Expr("score + ?", score)).Error
// }

// // 设置用户积分
// func (r *UserRepository) SetScoreById(id uint, score int) error {
// 	// 1. 验证ID有效性
// 	if id == 0 {
// 		return errors.New("无效的用户ID")
// 	}

// 	// 2. 验证积分范围（根据业务需求调整）
// 	if score < 0 {
// 		return errors.New("积分不能为负数")
// 	}
// 	if score > 999999 {
// 		return errors.New("积分超出最大限制")
// 	}

// 	// 3. 执行更新操作
// 	result := r.db.Model(&model.User{}).
// 		Where("id = ?", id).
// 		Update("score", score)

// 	// 4. 检查更新结果
// 	if result.Error != nil {
// 		return fmt.Errorf("更新积分失败: %w", result.Error)
// 	}

// 	// 5. 检查是否更新了记录
// 	if result.RowsAffected == 0 {
// 		return errors.New("用户不存在")
// 	}

// 	return nil
// }

// // FindByIDs 批量查询用户
// func (r *UserRepository) FindByIDs(ids []uint) ([]model.User, error) {
// 	if len(ids) == 0 {
// 		return []model.User{}, nil
// 	}
// 	var users []model.User
// 	err := r.db.Where("id IN ?", ids).Find(&users).Error
// 	return users, err
// }

// // Count 返回用户总数
// func (r *UserRepository) Count(ctx context.Context) (int64, error) {
// 	var count int64
// 	err := r.db.WithContext(ctx).Model(&model.User{}).Count(&count).Error
// 	return count, err
// }

// // CountByDateRange 统计指定时间段内新增用户数
// func (r *UserRepository) CountByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
// 	var count int64
// 	err := r.db.WithContext(ctx).
// 		Model(&model.User{}).
// 		Where("created_at BETWEEN ? AND ?", startDate, endDate).
// 		Count(&count).Error
// 	return count, err
// }

// // CountActiveByDateRange 统计指定时间段内活跃用户数（有发帖或评论行为）
// func (r *UserRepository) CountActiveByDateRange(ctx context.Context, startDate, endDate time.Time) (int64, error) {
// 	var count int64

// 	// 活跃用户 = 该时间段内有发帖或发评论的去重用户数
// 	err := r.db.WithContext(ctx).
// 		Table("users u").
// 		Where(`u.deleted_at IS NULL AND EXISTS (
// 			SELECT 1 FROM posts p
// 			WHERE p.author_id = u.id AND p.deleted_at IS NULL
// 			  AND p.created_at BETWEEN ? AND ?
// 		) OR EXISTS (
// 			SELECT 1 FROM comments c
// 			WHERE c.author_id = u.id AND c.deleted_at IS NULL
// 			  AND c.created_at BETWEEN ? AND ?
// 		)`, startDate, endDate, startDate, endDate).
// 		Count(&count).Error

// 	return count, err
// }

// // ActiveUserRow 活跃用户查询结果行
// type ActiveUserRow struct {
// 	ID       uint
// 	Username string
// 	Avatar   string
// }

// // GetActiveUsersByDateRange 获取指定时间段内活跃用户列表
// func (r *UserRepository) GetActiveUsersByDateRange(
// 	ctx context.Context,
// 	startDate, endDate time.Time,
// 	limit int,
// ) ([]*ActiveUserRow, error) {
// 	var rows []*ActiveUserRow

// 	err := r.db.WithContext(ctx).
// 		Table("users u").
// 		Select("u.id, u.username, u.avatar").
// 		Where(`u.deleted_at IS NULL AND (
// 			EXISTS (
// 				SELECT 1 FROM posts p
// 				WHERE p.author_id = u.id AND p.deleted_at IS NULL
// 				  AND p.created_at BETWEEN ? AND ?
// 			) OR EXISTS (
// 				SELECT 1 FROM comments c
// 				WHERE c.author_id = u.id AND c.deleted_at IS NULL
// 				  AND c.created_at BETWEEN ? AND ?
// 			)
// 		)`, startDate, endDate, startDate, endDate).
// 		Order("u.score DESC").
// 		Limit(limit).
// 		Scan(&rows).Error

// 	return rows, err
// }

// // MARK: Blocked 封禁用户
// // UpdateBlocked 更新用户封禁状态
// func (r *UserRepository) UpdateBlocked(ctx context.Context, userID uint, isBlocked bool) error {
// 	result := r.db.WithContext(ctx).
// 		Model(&model.User{}).
// 		Where("id = ?", userID).
// 		Update("is_blocked", isBlocked)

// 	if result.Error != nil {
// 		return fmt.Errorf("更新用户封禁状态失败: %w", result.Error)
// 	}

// 	if result.RowsAffected == 0 {
// 		return errors.New("用户不存在")
// 	}

// 	// 如果是封禁操作，清理用户的所有 Token
// 	if isBlocked {
// 		_ = r.tokenRepo.DeleteByUserID(ctx, userID)
// 	}

// 	return nil
// }

// // UpdateActive 更新用户激活状态（如果你有 is_active 字段）
// func (r *UserRepository) UpdateActive(ctx context.Context, userID uint, isActive bool) error {
// 	result := r.db.WithContext(ctx).
// 		Model(&model.User{}).
// 		Where("id = ?", userID).
// 		Update("is_active", isActive)

// 	if result.Error != nil {
// 		return fmt.Errorf("更新用户激活状态失败: %w", result.Error)
// 	}

// 	if result.RowsAffected == 0 {
// 		return errors.New("用户不存在")
// 	}

// 	return nil
// }

// // UpdateRole 更新用户角色
// func (r *UserRepository) UpdateRole(ctx context.Context, userID uint, role string) error {
// 	result := r.db.WithContext(ctx).
// 		Model(&model.User{}).
// 		Where("id = ?", userID).
// 		Update("role", role)

// 	if result.Error != nil {
// 		return fmt.Errorf("更新用户角色失败: %w", result.Error)
// 	}

// 	if result.RowsAffected == 0 {
// 		return errors.New("用户不存在")
// 	}

// 	return nil
// }

// // GetUserRoleById 获取用户角色
// func (r *UserRepository) GetUserRoleById(userID uint) (string, error) {
// 	var role string
// 	err := r.db.Model(&model.User{}).
// 		Select("role").
// 		Where("id = ?", userID).
// 		Scan(&role).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return "", err
// 		}
// 		return "", fmt.Errorf("查询用户角色失败: %w", err)
// 	}
// 	return role, nil
// }

// // GetUserBasicInfoById 获取用户基本信息（ID、用户名、角色）
// func (r *UserRepository) GetUserBasicInfoById(userID uint) (*model.User, error) {
// 	var user model.User
// 	err := r.db.Model(&model.User{}).
// 		Select("id, username, role").
// 		Where("id = ?", userID).
// 		First(&user).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &user, nil
// }

// // BatchUpdateBlocked 批量更新封禁状态
// func (r *UserRepository) BatchUpdateBlocked(ctx context.Context, userIDs []uint, isBlocked bool) (int64, error) {
// 	if len(userIDs) == 0 {
// 		return 0, nil
// 	}

// 	result := r.db.WithContext(ctx).
// 		Model(&model.User{}).
// 		Where("id IN ?", userIDs).
// 		Update("is_blocked", isBlocked)

// 	return result.RowsAffected, result.Error
// }

// // MARK: - 管理员操作相关方法

// // SoftDelete 软删除用户
// func (r *UserRepository) SoftDelete(ctx context.Context, userID uint) error {
// 	// 1. 先清理用户的所有 Token
// 	_ = r.tokenRepo.DeleteByUserID(ctx, userID)

// 	// 2. 软删除用户
// 	result := r.db.WithContext(ctx).
// 		Model(&model.User{}).
// 		Where("id = ?", userID).
// 		Update("deleted_at", time.Now())

// 	if result.Error != nil {
// 		return fmt.Errorf("软删除用户失败: %w", result.Error)
// 	}

// 	if result.RowsAffected == 0 {
// 		return errors.New("用户不存在或已被删除")
// 	}

// 	return nil
// }

// // HardDelete 硬删除用户（慎用，仅用于测试或数据清理）
// func (r *UserRepository) HardDelete(ctx context.Context, userID uint) error {
// 	result := r.db.WithContext(ctx).
// 		Unscoped().
// 		Delete(&model.User{}, userID)

// 	if result.Error != nil {
// 		return fmt.Errorf("硬删除用户失败: %w", result.Error)
// 	}

// 	if result.RowsAffected == 0 {
// 		return errors.New("用户不存在")
// 	}

// 	return nil
// }

// // UpdatePassword 更新用户密码
// func (r *UserRepository) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
// 	// 1. 更新密码
// 	result := r.db.WithContext(ctx).
// 		Model(&model.User{}).
// 		Where("id = ?", userID).
// 		Update("password", hashedPassword)

// 	if result.Error != nil {
// 		return fmt.Errorf("更新密码失败: %w", result.Error)
// 	}

// 	if result.RowsAffected == 0 {
// 		return errors.New("用户不存在")
// 	}

// 	// 2. 清理用户的所有 Token（强制重新登录）
// 	_ = r.tokenRepo.DeleteByUserID(ctx, userID)

// 	return nil
// }

// // InvalidateUserTokens 使用户的所有 Token 失效
// func (r *UserRepository) InvalidateUserTokens(ctx context.Context, userID uint) error {
// 	return r.tokenRepo.DeleteByUserID(ctx, userID)
// }

// // RestoreDeleted 恢复已删除的用户
// func (r *UserRepository) RestoreDeleted(ctx context.Context, userID uint) error {
// 	result := r.db.WithContext(ctx).
// 		Model(&model.User{}).
// 		Where("id = ?", userID).
// 		Update("deleted_at", nil)

// 	if result.Error != nil {
// 		return fmt.Errorf("恢复用户失败: %w", result.Error)
// 	}

// 	if result.RowsAffected == 0 {
// 		return errors.New("用户不存在或未被删除")
// 	}

// 	return nil
// }

// // SetTempPasswordFlag 设置临时密码标记
// func (r *UserRepository) SetTempPasswordFlag(ctx context.Context, userID uint, isTemp bool, expireAt time.Time) error {
// 	updates := map[string]interface{}{
// 		"is_temp_password":     isTemp,
// 		"temp_password_expire": expireAt,
// 	}

// 	result := r.db.WithContext(ctx).
// 		Model(&model.User{}).
// 		Where("id = ?", userID).
// 		Updates(updates)

// 	if result.Error != nil {
// 		return fmt.Errorf("设置临时密码标记失败: %w", result.Error)
// 	}

// 	return nil
// }

// // ClearTempPasswordFlag 清除临时密码标记
// func (r *UserRepository) ClearTempPasswordFlag(ctx context.Context, userID uint) error {
// 	result := r.db.WithContext(ctx).
// 		Model(&model.User{}).
// 		Where("id = ?", userID).
// 		Updates(map[string]interface{}{
// 			"is_temp_password":     false,
// 			"temp_password_expire": nil,
// 		})

//		return result.Error
//	}
package repository
