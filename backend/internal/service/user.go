// package service

// import (
// 	"context"
// 	"crypto/rand"
// 	"errors"
// 	"fmt"
// 	"math/big"
// 	"regexp"
// 	"time"
// 	"tiny-forum/internal/model"
// 	"tiny-forum/internal/repository"
// 	apperrors "tiny-forum/pkg/errors"
// 	"tiny-forum/pkg/fields"
// 	jwtpkg "tiny-forum/pkg/jwt"

// 	"golang.org/x/crypto/bcrypt"
// 	"gorm.io/gorm"
// )

// type UserService struct {
// 	repo        *repository.UserRepository
// 	jwtMgr      *jwtpkg.Manager
// 	notifSvc    *NotificationService
// 	roleChecker RoleChangeChecker
// }

// func NewUserService(
// 	repo *repository.UserRepository,
// 	jwtMgr *jwtpkg.Manager,
// 	notifSvc *NotificationService,
// ) *UserService {
// 	return &UserService{
// 		repo:        repo,
// 		jwtMgr:      jwtMgr,
// 		notifSvc:    notifSvc,
// 		roleChecker: RoleChangeChecker{},
// 	}
// }

// // ── Auth ─────────────────────────────────────────────────────────────────────

// type RegisterInput struct {
// 	Username string `json:"username" binding:"required,min=2,max=50"`
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required,min=6"`
// }

// type LoginInput struct {
// 	Email    string `json:"email" binding:"required,email"`
// 	Password string `json:"password" binding:"required"`
// }

// type AuthResult struct {
// 	Token string      `json:"token"`
// 	User  *model.User `json:"user"`
// }

// func (s *UserService) Register(input RegisterInput) (*AuthResult, error) {
// 	if _, err := s.repo.FindByUsername(input.Username); err == nil {
// 		return nil, errors.New("用户名已被占用")
// 	}
// 	if _, err := s.repo.FindByEmail(input.Email); err == nil {
// 		return nil, errors.New("邮箱已被注册")
// 	}

// 	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
// 	if err != nil {
// 		return nil, err
// 	}

// 	user := &model.User{
// 		Username: input.Username,
// 		Email:    input.Email,
// 		Password: string(hashed),
// 		Role:     model.RoleUser,
// 		Avatar:   avatarURL(input.Username),
// 	}
// 	if err := s.repo.Create(user); err != nil {
// 		return nil, err
// 	}

// 	s.notifSvc.Create(user.ID, nil, model.NotifySystem, "欢迎加入 Tiny Forum！", nil, "")

// 	token, err := s.jwtMgr.Generate(user.ID, user.Username, string(user.Role))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &AuthResult{Token: token, User: user}, nil
// }

// func (s *UserService) Login(input LoginInput) (*AuthResult, error) {
// 	user, err := s.repo.FindByEmail(input.Email)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("邮箱或密码错误")
// 		}
// 		return nil, err
// 	}
// 	if user.IsBlocked {
// 		return nil, errors.New("账户已被禁用")
// 	}
// 	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
// 		return nil, errors.New("邮箱或密码错误")
// 	}

// 	now := time.Now()
// 	user.LastLogin = &now
// 	_ = s.repo.Update(user)

// 	token, err := s.jwtMgr.Generate(user.ID, user.Username, string(user.Role))
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &AuthResult{Token: token, User: user}, nil
// }

// // ── Profile ──────────────────────────────────────────────────────────────────

// func (s *UserService) GetProfile(userID uint) (*model.User, error) {
// 	return s.repo.FindByID(userID)
// }

// func (s *UserService) UpdateProfile(userID uint, input model.UpdateProfileInput) error {
// 	fields := map[string]interface{}{}
// 	if input.Bio != "" {
// 		fields["bio"] = input.Bio
// 	}
// 	if input.Avatar != "" {
// 		fields["avatar"] = input.Avatar
// 	}
// 	if input.Email != "" {
// 		fields["email"] = input.Email
// 	}
// 	if len(fields) == 0 {
// 		return nil
// 	}
// 	return s.repo.UpdateFields(userID, fields)
// }

// func (s *UserService) ChangePassword(userID uint, oldPassword, newPassword string) (string, error) {
// 	ctx := context.Background()

// 	// 1. 查询用户
// 	targetUser, err := s.repo.FindByID(userID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return "", apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", userID)
// 		}
// 		return "", fmt.Errorf("查询目标用户失败: %w", err)
// 	}

// 	// 2. 验证旧密码 - 这里失败会返回 20005
// 	if err := bcrypt.CompareHashAndPassword([]byte(targetUser.Password), []byte(oldPassword)); err != nil {
// 		// 添加日志以便调试
// 		fmt.Printf("密码验证失败 - UserID: %d, Error: %v\n", userID, err)
// 		return "", apperrors.ErrInvalidPassword // 错误码 20005
// 	}

// 	// 3. 验证新旧密码不能相同
// 	if oldPassword == newPassword {
// 		return "", apperrors.ErrPasswordSameAsOld
// 	}

// 	// 4. 验证新密码长度（不强制强度）
// 	if len(newPassword) < 6 {
// 		return "", apperrors.ErrPasswordTooShort
// 	}

// 	// 5. 检查密码强度（仅用于提示）
// 	strength := s.checkPasswordStrength(newPassword)

// 	// 6. 加密新密码
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
// 	if err != nil {
// 		return "", fmt.Errorf("密码加密失败: %w", err)
// 	}

// 	// 7. 更新密码
// 	if err := s.repo.UpdatePassword(ctx, userID, string(hashedPassword)); err != nil {
// 		return "", fmt.Errorf("更新密码失败: %w", err)
// 	}

// 	// 8. 返回结果
// 	if strength == "weak" {
// 		return "密码修改成功，但密码强度较弱，建议使用更复杂的密码", nil
// 	}

// 	return "密码修改成功", nil
// }

// // checkPasswordStrength 检查密码强度（仅用于建议）
// func (s *UserService) checkPasswordStrength(password string) string {
// 	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
// 	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
// 	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
// 	hasSpecial := regexp.MustCompile(`[^A-Za-z0-9]`).MatchString(password)

// 	score := 0
// 	if len(password) >= 8 {
// 		score++
// 	}
// 	if len(password) >= 12 {
// 		score++
// 	}
// 	if hasUpper && hasLower {
// 		score++
// 	}
// 	if hasNumber {
// 		score++
// 	}
// 	if hasSpecial {
// 		score++
// 	}

// 	if score <= 2 {
// 		return "weak"
// 	}
// 	if score <= 4 {
// 		return "medium"
// 	}
// 	return "strong"
// }

// type UserProfileResponse struct {
// 	*model.User
// 	FollowerCount  int64 `json:"follower_count"`
// 	FollowingCount int64 `json:"following_count"`
// 	IsFollowing    bool  `json:"is_following"`
// }

// func (s *UserService) GetUserProfile(targetID, viewerID uint) (*UserProfileResponse, error) {
// 	user, err := s.repo.FindByID(targetID)
// 	if err != nil {
// 		// 转换数据库错误为业务错误
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", targetID)
// 		}
// 		return nil, fmt.Errorf("查询用户失败: %w", err)
// 	}

// 	resp := &UserProfileResponse{
// 		User:           user,
// 		FollowerCount:  s.repo.GetFollowerCount(targetID),
// 		FollowingCount: s.repo.GetFollowingCount(targetID),
// 	}

// 	if viewerID > 0 {
// 		resp.IsFollowing = s.repo.IsFollowing(viewerID, targetID)
// 	}

// 	return resp, nil
// }

// // ── Follow ───────────────────────────────────────────────────────────────────

// func (s *UserService) Follow(followerID, followingID uint) error {
// 	if followerID == followingID {
// 		return errors.New("不能关注自己")
// 	}
// 	if err := s.repo.Follow(followerID, followingID); err != nil {
// 		return err
// 	}
// 	following, _ := s.repo.FindByID(followingID)
// 	if following != nil {
// 		s.notifSvc.Create(followingID, &followerID, model.NotifyFollow,
// 			following.Username+" 关注了你", nil, "")
// 	}
// 	return nil
// }

// func (s *UserService) Unfollow(followerID, followingID uint) error {
// 	return s.repo.Unfollow(followerID, followingID)
// }

// // 获取关注者列表
// func (s *UserService) GetFollowers(userID uint, page, pageSize int) ([]model.User, int64, error) {
// 	return s.repo.GetFollowers(userID, page, pageSize)
// }

// // 获取关注列表
// func (s *UserService) GetFollowing(userID uint, page, pageSize int) ([]model.User, int64, error) {
// 	return s.repo.GetFollowing(userID, page, pageSize)
// }

// // 查询积分
// func (s *UserService) GetScoreById(userID uint) (int, error) {
// 	return s.repo.GetScoreById(userID)
// }

// // GetUserRoleById 获取用户角色
// func (s *UserService) GetUserRoleById(userID uint) (string, error) {
// 	return s.repo.GetUserRoleById(userID)
// }

// // GetUserBasicInfo 获取用户基本信息
// func (s *UserService) GetUserBasicInfo(userID uint) (*model.User, error) {
// 	return s.repo.GetUserBasicInfoById(userID)
// }

// type UserScoreResponse struct {
// 	ID       uint   `json:"id"`
// 	Username string `json:"username"`
// 	Avatar   string `json:"avatar_url"`
// 	Score    int    `json:"score"`
// }

// func (s *UserService) GetAllUsersWithScore() ([]UserScoreResponse, error) {
// 	users, err := s.repo.GetEveryoneUsersScore()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var result []UserScoreResponse
// 	for _, user := range users {
// 		// 获取用户基本信息（包含用户名和头像）
// 		basicInfo, err := s.repo.GetUserBasicInfo(user.ID)
// 		if err != nil {
// 			continue
// 		}

// 		result = append(result, UserScoreResponse{
// 			ID:       user.ID,
// 			Username: basicInfo.Username,
// 			Avatar:   basicInfo.Avatar,
// 			Score:    user.Score,
// 		})
// 	}

// 	return result, nil
// }

// // ── Admin ────────────────────────────────────────────────────────────────────
// type LeaderboardItem struct {
// 	ID       uint   `json:"id"`
// 	Username string `json:"username"`
// 	Avatar   string `json:"avatar"`
// 	Score    int    `json:"score"`
// 	Rank     int    `json:"rank"` // 额外添加排名
// }

// func (s *UserService) GetLeaderboard(ctx context.Context, limit int, fieldsParam string) ([]LeaderboardItem, error) {
// 	// 1. 参数校验
// 	if limit < 1 {
// 		limit = 20
// 	}
// 	if limit > 100 {
// 		limit = 100
// 	}

// 	// 2. 字段过滤（规范调用）
// 	selectedFields := fields.Filter(
// 		fieldsParam,
// 		model.UserPublicFields,
// 		model.UserDefaultFields,
// 	)

// 	// 3. 构造 Repository 查询
// 	query := repository.TopUsersQuery{
// 		Limit:          limit,
// 		ExcludeBlocked: true,
// 		Fields:         selectedFields,
// 	}

// 	users, err := s.repo.GetTopUsers(ctx, query)
// 	if err != nil {
// 		return nil, fmt.Errorf("查询排行榜失败: %w", err)
// 	}

// 	// 4. 转换为 DTO 并添加排名
// 	items := make([]LeaderboardItem, len(users))
// 	for i, u := range users {
// 		items[i] = LeaderboardItem{
// 			ID:       u.ID,
// 			Username: u.Username,
// 			Avatar:   u.Avatar,
// 			Score:    u.Score,
// 			Rank:     i + 1,
// 		}
// 	}
// 	return items, nil
// }

// func (s *UserService) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
// 	return s.repo.List(page, pageSize, keyword)
// }

// // MARK: Status
// // SetBlocked 设置用户封禁状态
// // func (s *UserService) SetBlocked(ctx context.Context, targetID uint, operatorID uint, isBlocked bool) error {
// // 	// 1. 查询目标用户
// // 	targetUser, err := s.repo.FindByID(targetID)
// // 	if err != nil {
// // 		if errors.Is(err, gorm.ErrRecordNotFound) {
// // 			return apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", targetID)
// // 		}
// // 		return fmt.Errorf("查询目标用户失败: %w", err)
// // 	}

// // 	// 2. 查询操作者信息
// // 	operator, err := s.repo.FindByID(operatorID)
// // 	if err != nil {
// // 		if errors.Is(err, gorm.ErrRecordNotFound) {
// // 			return apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", operatorID)
// // 		}
// // 		return fmt.Errorf("查询操作者信息失败: %w", err)
// // 	}

// // 	// 3. 安全检查：不能封禁自己
// // 	if targetID == operatorID {
// // 		return apperrors.ErrCannotModifySelf
// // 	}

// // 	// 4. 安全检查：不能封禁超级管理员
// // 	if targetUser.Role == model.RoleSuperAdmin {
// // 		return apperrors.ErrCannotChangeOwnerRole
// // 	}

// // 	// 5. 安全检查：普通管理员不能封禁其他管理员
// // 	if operator.Role != model.RoleSuperAdmin && targetUser.Role == model.RoleAdmin {
// // 		return apperrors.ErrInsufficientPermission
// // 	}

// // 	// 6. 幂等性检查
// // 	if targetUser.IsBlocked == isBlocked {
// // 		return nil
// // 	}

// // 	// 7. 更新封禁状态
// // 	if err := s.repo.UpdateBlocked(ctx, targetID, isBlocked); err != nil {
// // 		if errors.Is(err, gorm.ErrRecordNotFound) {
// // 			return apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", targetID)
// // 		}
// // 		return fmt.Errorf("更新用户封禁状态失败: %w", err)
// // 	}

// // 	return nil
// // }

// // SetActive 设置用户激活状态
// // func (s *UserService) SetActive(ctx context.Context, targetID uint, operatorID uint, isActive bool) error {
// // 	// 1. 查询目标用户
// // 	targetUser, err := s.repo.FindByID(targetID)
// // 	if err != nil {
// // 		if errors.Is(err, gorm.ErrRecordNotFound) {
// // 			return apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", targetID)
// // 		}
// // 		return fmt.Errorf("查询目标用户失败: %w", err)
// // 	}

// // 	// 2. 安全检查：不能禁用自己
// // 	if targetID == operatorID {
// // 		return apperrors.ErrCannotModifySelf
// // 	}

// // 	// 3. 安全检查：不能禁用超级管理员
// // 	if targetUser.Role == model.RoleSuperAdmin && !isActive {
// // 		return apperrors.ErrCannotChangeOwnerRole
// // 	}

// // 	// 4. 幂等性检查
// // 	if targetUser.IsActive == isActive {
// // 		return nil
// // 	}

// // 	// 5. 更新激活状态
// // 	if err := s.repo.UpdateActive(ctx, targetID, isActive); err != nil {
// // 		if errors.Is(err, gorm.ErrRecordNotFound) {
// // 			return apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", targetID)
// // 		}
// // 		return fmt.Errorf("更新用户激活状态失败: %w", err)
// // 	}

// // 	return nil
// // }

// // MARK: Scoer
// func (s *UserService) SetScoreById(userID uint, score int) error {
// 	// 1. 参数验证
// 	if userID == 0 {
// 		return errors.New("用户ID不能为空")
// 	}

// 	err := s.repo.SetScoreById(userID, score)
// 	if err != nil {
// 		return fmt.Errorf("设置积分失败: %w", err)
// 	}

// 	// 4. 可选：触发积分变更事件（如发送通知、更新缓存等）
// 	go s.onScoreChanged(userID, score)

// 	return nil
// }

// // 积分变更后的回调处理
// func (s *UserService) onScoreChanged(userID uint, newScore int) {
// 	// 可以在这里添加：
// 	// - 发送系统通知
// 	// - 更新Redis缓存
// 	// - 检查是否触发等级变更
// 	// - 记录日志到消息队列等
// }

// // ── Role Management ──────────────────────────────────────────────────────────

// // SetRole 变更目标用户角色（含细粒度权限校验）。
// // operatorID 来自 JWT，确保操作者身份可信。
// func (s *UserService) SetRole(operatorID, targetID uint, newRole string) error {
// 	// 1. 校验新角色字面值
// 	targetRole := model.UserRole(newRole)
// 	if !model.IsValidRole(targetRole) {
// 		return fmt.Errorf("%w: %s", apperrors.ErrInvalidRole, newRole)
// 	}

// 	// 2. 加载操作者与目标用户
// 	operator, err := s.repo.FindByID(operatorID)
// 	if err != nil {
// 		return fmt.Errorf("操作者不存在: %w", err)
// 	}
// 	target, err := s.repo.FindByID(targetID)
// 	if err != nil {
// 		return err // 保留 gorm.ErrRecordNotFound 供 handler 判断
// 	}

// 	// 3. 幂等：角色无变化直接返回
// 	if target.Role == targetRole {
// 		return nil
// 	}

// 	// 4. 细粒度权限校验（职责交由 RoleChangeChecker）
// 	if err := s.roleChecker.Check(RoleChangeRequest{
// 		Operator: operator,
// 		Target:   target,
// 		NewRole:  targetRole,
// 	}); err != nil {
// 		return err
// 	}

// 	// 5. 执行更新
// 	return s.repo.UpdateFields(targetID, map[string]interface{}{"role": newRole})
// }

// // DeleteUser 管理员删除用户
// func (s *UserService) DeleteUser(operatorID uint, targetID uint) error {
// 	ctx := context.Background()

// 	// 1. 查询目标用户
// 	targetUser, err := s.repo.FindByID(targetID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", targetID)
// 		}
// 		return fmt.Errorf("查询目标用户失败: %w", err)
// 	}

// 	// 2. 查询操作者
// 	operator, err := s.repo.FindByID(operatorID)
// 	if err != nil {
// 		return fmt.Errorf("查询操作者信息失败: %w", err)
// 	}

// 	// 3. 安全检查：不能删除自己
// 	if targetID == operatorID {
// 		return apperrors.Wrap(apperrors.ErrCannotModifySelf, "不能删除自己的账号")
// 	}

// 	// 4. 安全检查：不能删除超级管理员
// 	if targetUser.Role == model.RoleSuperAdmin {
// 		return apperrors.Wrap(apperrors.ErrCannotChangeOwnerRole, "不能删除超级管理员")
// 	}

// 	// 5. 安全检查：普通管理员不能删除其他管理员
// 	if operator.Role != model.RoleSuperAdmin && targetUser.Role == model.RoleAdmin {
// 		return apperrors.Wrap(apperrors.ErrInsufficientPermission, "只有超级管理员才能删除其他管理员")
// 	}

// 	// 6. 执行软删除（会自动清理 Token，因为 SoftDelete 内部调用了 tokenRepo.DeleteByUserID）
// 	if err := s.repo.SoftDelete(ctx, targetID); err != nil {
// 		return fmt.Errorf("删除用户失败: %w", err)
// 	}

// 	// 7. 记录审计日志（可选）
// 	s.logAudit(ctx, operatorID, targetID, "delete_user",
// 		fmt.Sprintf("管理员 %d 删除了用户 %d (%s)", operatorID, targetID, targetUser.Username))

// 	return nil
// }

// const (
// 	tempPasswordLength = 12
// 	tempPasswordTTL    = 30 * time.Minute
// )

// // 管理员重置用户密码（生成临时密码）
// func (s *UserService) ResetUserPasswordWithTemp(operatorID uint, targetID uint) (string, error) {
// 	ctx := context.Background()

// 	// 1. 查询目标用户
// 	targetUser, err := s.repo.FindByID(targetID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return "", apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", targetID)
// 		}
// 		return "", fmt.Errorf("查询目标用户失败: %w", err)
// 	}

// 	// 2. 查询操作者
// 	operator, err := s.repo.FindByID(operatorID)
// 	if err != nil {
// 		return "", fmt.Errorf("查询操作者信息失败: %w", err)
// 	}

// 	// 3. 安全检查：不能重置超级管理员的密码（除非自己是超管）
// 	if targetUser.Role == model.RoleSuperAdmin && operator.Role != model.RoleSuperAdmin {
// 		return "", apperrors.Wrap(apperrors.ErrInsufficientPermission, "只有超级管理员才能重置其他超级管理员的密码")
// 	}

// 	// 4. 普通管理员不能重置其他管理员的密码
// 	if operator.Role != model.RoleSuperAdmin && targetUser.Role == model.RoleAdmin {
// 		return "", apperrors.Wrap(apperrors.ErrInsufficientPermission, "只有超级管理员才能重置其他管理员的密码")
// 	}

// 	// 5. 生成随机临时密码
// 	tempPassword, err := generateSecurePassword(tempPasswordLength)
// 	if err != nil {
// 		return "", fmt.Errorf("生成临时密码失败: %w", err)
// 	}

// 	// 6. 加密密码
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(tempPassword), bcrypt.DefaultCost)
// 	if err != nil {
// 		return "", fmt.Errorf("密码加密失败: %w", err)
// 	}

// 	// 7. 更新密码（UpdatePassword 内部会自动清理 Token）
// 	if err := s.repo.UpdatePassword(ctx, targetID, string(hashedPassword)); err != nil {
// 		return "", fmt.Errorf("更新密码失败: %w", err)
// 	}

// 	// 8. 记录密码重置标记（可选，用于判断是否需要强制修改密码）
// 	expiresAt := time.Now().Add(tempPasswordTTL)
// 	if err := s.repo.SetTempPasswordFlag(ctx, targetID, true, expiresAt); err != nil {
// 		// 记录日志但不中断流程
// 		// log.Printf("设置临时密码标记失败: %v", err)
// 	}

// 	// 9. 记录审计日志
// 	s.logAudit(ctx, operatorID, targetID, "reset_password_temp",
// 		fmt.Sprintf("管理员 %d 为用户 %d (%s) 生成了临时密码", operatorID, targetID, targetUser.Username))

// 	return tempPassword, nil
// }

// // generateSecurePassword 生成安全随机密码
// func generateSecurePassword(length int) (string, error) {
// 	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

// 	result := make([]byte, length)
// 	for i := 0; i < length; i++ {
// 		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
// 		if err != nil {
// 			return "", err
// 		}
// 		result[i] = chars[n.Int64()]
// 	}

// 	return string(result), nil
// }

// // ResetUserPassword 管理员重置用户密码（指定密码）
// func (s *UserService) ResetUserPassword(operatorID uint, targetID uint, newPassword string) error {
// 	ctx := context.Background()

// 	// 1. 查询目标用户
// 	targetUser, err := s.repo.FindByID(targetID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", targetID)
// 		}
// 		return fmt.Errorf("查询目标用户失败: %w", err)
// 	}

// 	// 2. 查询操作者
// 	operator, err := s.repo.FindByID(operatorID)
// 	if err != nil {
// 		return fmt.Errorf("查询操作者信息失败: %w", err)
// 	}

// 	// 3. 安全检查：不能重置超级管理员的密码（除非自己是超管）
// 	if targetUser.Role == model.RoleSuperAdmin && operator.Role != model.RoleSuperAdmin {
// 		return apperrors.Wrap(apperrors.ErrInsufficientPermission, "只有超级管理员才能重置其他超级管理员的密码")
// 	}

// 	// 4. 普通管理员不能重置其他管理员的密码
// 	if operator.Role != model.RoleSuperAdmin && targetUser.Role == model.RoleAdmin {
// 		return apperrors.Wrap(apperrors.ErrInsufficientPermission, "只有超级管理员才能重置其他管理员的密码")
// 	}

// 	// 5. 验证新密码强度（可选但推荐）
// 	if err := s.validatePasswordStrength(newPassword); err != nil {
// 		return err
// 	}

// 	// 6. 加密新密码
// 	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
// 	if err != nil {
// 		return fmt.Errorf("密码加密失败: %w", err)
// 	}

// 	// 7. 更新密码（UpdatePassword 内部会自动清理 Token）
// 	if err := s.repo.UpdatePassword(ctx, targetID, string(hashedPassword)); err != nil {
// 		return fmt.Errorf("更新密码失败: %w", err)
// 	}

// 	// 8. 记录审计日志
// 	s.logAudit(ctx, operatorID, targetID, "reset_password",
// 		fmt.Sprintf("管理员 %d 重置了用户 %d (%s) 的密码", operatorID, targetID, targetUser.Username))

// 	return nil
// }

// // SetBlocked 管理员封禁/解封用户
// func (s *UserService) SetBlocked(targetID uint, operatorID uint, isBlocked bool) error {
// 	ctx := context.Background()

// 	// 1. 查询目标用户
// 	targetUser, err := s.repo.FindByID(targetID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", targetID)
// 		}
// 		return fmt.Errorf("查询目标用户失败: %w", err)
// 	}

// 	// 2. 查询操作者
// 	operator, err := s.repo.FindByID(operatorID)
// 	if err != nil {
// 		return fmt.Errorf("查询操作者信息失败: %w", err)
// 	}

// 	// 3. 安全检查：不能封禁自己
// 	if targetID == operatorID {
// 		return apperrors.Wrap(apperrors.ErrCannotModifySelf, "不能封禁自己的账号")
// 	}

// 	// 4. 安全检查：不能封禁超级管理员
// 	if targetUser.Role == model.RoleSuperAdmin {
// 		return apperrors.Wrap(apperrors.ErrCannotChangeOwnerRole, "不能封禁超级管理员")
// 	}

// 	// 5. 安全检查：普通管理员不能封禁其他管理员
// 	if operator.Role != model.RoleSuperAdmin && targetUser.Role == model.RoleAdmin {
// 		return apperrors.Wrap(apperrors.ErrInsufficientPermission, "只有超级管理员才能封禁其他管理员")
// 	}

// 	// 6. 幂等性检查
// 	if targetUser.IsBlocked == isBlocked {
// 		return nil
// 	}

// 	// 7. 更新封禁状态（UpdateBlocked 内部会在封禁时自动清理 Token）
// 	if err := s.repo.UpdateBlocked(ctx, targetID, isBlocked); err != nil {
// 		return fmt.Errorf("更新用户封禁状态失败: %w", err)
// 	}

// 	// 8. 记录审计日志
// 	action := "unblock_user"
// 	if isBlocked {
// 		action = "block_user"
// 	}
// 	s.logAudit(ctx, operatorID, targetID, action,
// 		fmt.Sprintf("管理员 %d %s用户 %d (%s)", operatorID, map[bool]string{true: "封禁", false: "解封"}[isBlocked], targetID, targetUser.Username))

// 	return nil
// }

// // SetActive 管理员设置用户激活状态
// func (s *UserService) SetActive(targetID uint, operatorID uint, isActive bool) error {
// 	ctx := context.Background()

// 	// 1. 查询目标用户
// 	targetUser, err := s.repo.FindByID(targetID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return apperrors.Wrapf(apperrors.ErrUserNotFound, "ID: %d", targetID)
// 		}
// 		return fmt.Errorf("查询目标用户失败: %w", err)
// 	}

// 	// 2. 安全检查：不能停用自己的账号
// 	if targetID == operatorID && !isActive {
// 		return apperrors.Wrap(apperrors.ErrCannotModifySelf, "不能停用自己的账号")
// 	}

// 	// 3. 安全检查：不能停用超级管理员
// 	if targetUser.Role == model.RoleSuperAdmin && !isActive {
// 		return apperrors.Wrap(apperrors.ErrCannotChangeOwnerRole, "不能停用超级管理员")
// 	}

// 	// 4. 幂等性检查
// 	if targetUser.IsActive == isActive {
// 		return nil
// 	}

// 	// 5. 更新激活状态
// 	if err := s.repo.UpdateActive(ctx, targetID, isActive); err != nil {
// 		return fmt.Errorf("更新用户激活状态失败: %w", err)
// 	}

// 	// 6. 如果停用账号，清理 Token
// 	if !isActive {
// 		_ = s.repo.InvalidateUserTokens(ctx, targetID)
// 	}

// 	// 7. 记录审计日志
// 	action := "activate_user"
// 	if !isActive {
// 		action = "deactivate_user"
// 	}
// 	s.logAudit(ctx, operatorID, targetID, action,
// 		fmt.Sprintf("管理员 %d %s用户 %d (%s)", operatorID, map[bool]string{true: "激活", false: "停用"}[isActive], targetID, targetUser.Username))

// 	return nil
// }

// // ========== 辅助方法 ==========

// // validatePasswordStrength 验证密码强度
// func (s *UserService) validatePasswordStrength(password string) error {
// 	if len(password) < 6 {
// 		return apperrors.Wrap(apperrors.ErrInvalidPassword, "密码长度至少6位")
// 	}
// 	if len(password) > 32 {
// 		return apperrors.Wrap(apperrors.ErrInvalidPassword, "密码长度不能超过32位")
// 	}
// 	// 可选：检查是否包含数字、字母等
// 	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
// 	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
// 	if !hasDigit || !hasLetter {
// 		return apperrors.Wrap(apperrors.ErrInvalidPassword, "密码必须包含数字和字母")
// 	}
// 	return nil
// }

// // logAudit 记录审计日志
// func (s *UserService) logAudit(ctx context.Context, operatorID, targetID uint, action, detail string) {
// 	// 如果有日志库，使用日志库记录
// 	// s.logger.Info("audit_log",
// 	//     "operator_id", operatorID,
// 	//     "target_id", targetID,
// 	//     "action", action,
// 	//     "detail", detail,
// 	// )

// 	// 或者存入数据库审计表
// 	// s.auditRepo.Create(ctx, &model.AuditLog{
// 	//     OperatorID: operatorID,
// 	//     TargetID:   targetID,
// 	//     Action:     action,
// 	//     Detail:     detail,
// 	// })
// }

// // ── 内部工具 ─────────────────────────────────────────────────────────────────

// func avatarURL(username string) string {
// 	return "https://api.dicebear.com/8.x/lorelei/svg?seed=" + username
// }

// type LoginResult struct {
// 	Token string    `json:"-"` // json:"-" 防止意外序列化到响应
// 	User  *UserInfo `json:"user"`
// }

//	type UserInfo struct {
//		ID    uint   `json:"id"`
//		Name  string `json:"name"`
//		Email string `json:"email"`
//		Role  string `json:"role"`
//	}
package service
