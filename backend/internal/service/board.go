package service

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	apperrors "tiny-forum/internal/errors"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
)

// ── BoardService ──────────────────────────────────────────────────────────────

// BoardService 负责板块、版主、禁言三大模块的业务逻辑。
// 所有数据访问通过 boardRepo（*repository.BoardRepository）和
// userRepo（*repository.UserRepository）完成，不依赖独立的 moderatorRepo，
// 避免职责分散和接口签名冲突。
type BoardService struct {
	boardRepo *repository.BoardRepository
	userRepo  *repository.UserRepository
	postRepo  repository.PostRepository
	notifSvc  *NotificationService
}

func NewBoardService(
	boardRepo *repository.BoardRepository,
	userRepo *repository.UserRepository,
	postRepo repository.PostRepository,
	notifSvc *NotificationService,
) *BoardService {
	return &BoardService{
		boardRepo: boardRepo,
		userRepo:  userRepo,
		postRepo:  postRepo,
		notifSvc:  notifSvc,
	}
}

// ── Board CRUD ────────────────────────────────────────────────────────────────

type CreateBoardInput struct {
	Name        string `json:"name"        binding:"required,min=2,max=50"`
	Slug        string `json:"slug"        binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"max=500"`
	Icon        string `json:"icon"        binding:"max=100"`
	Cover       string `json:"cover"       binding:"max=500"`
	ParentID    *uint  `json:"parent_id"`
	SortOrder   int    `json:"sort_order"`
	ViewRole    string `json:"view_role"`
	PostRole    string `json:"post_role"`
	ReplyRole   string `json:"reply_role"`
}

func (s *BoardService) Create(input CreateBoardInput) (*model.Board, error) {
	// 验证父板块
	if input.ParentID != nil && *input.ParentID != 0 {
		parent, err := s.boardRepo.FindByID(*input.ParentID)
		if err != nil || parent == nil || parent.ID == 0 {
			return nil, fmt.Errorf("父板块不存在: id=%d", *input.ParentID)
		}
	} else {
		input.ParentID = nil
	}

	// 验证角色值
	if err := validateRoles(input.ViewRole, input.PostRole, input.ReplyRole); err != nil {
		return nil, err
	}

	// Slug 唯一性 + 格式
	if existing, _ := s.boardRepo.FindBySlug(input.Slug); existing != nil && existing.ID != 0 {
		return nil, errors.New("板块标识已存在")
	}
	if err := validateSlug(input.Slug); err != nil {
		return nil, err
	}

	board := &model.Board{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Icon:        input.Icon,
		Cover:       input.Cover,
		ParentID:    input.ParentID,
		SortOrder:   input.SortOrder,
		ViewRole:    model.UserRole(input.ViewRole),
		PostRole:    model.UserRole(input.PostRole),
		ReplyRole:   model.UserRole(input.ReplyRole),
	}
	// 设置默认角色
	if board.ViewRole == "" {
		board.ViewRole = model.RoleUser
	}
	if board.PostRole == "" {
		board.PostRole = model.RoleUser
	}
	if board.ReplyRole == "" {
		board.ReplyRole = model.RoleUser
	}

	if err := s.boardRepo.Create(board); err != nil {
		return nil, fmt.Errorf("创建板块失败: %w", err)
	}
	return s.boardRepo.FindByID(board.ID)
}

func (s *BoardService) Update(id uint, input CreateBoardInput) (*model.Board, error) {
	board, err := s.boardRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("板块不存在")
	}

	if input.Slug != board.Slug {
		existing, _ := s.boardRepo.FindBySlug(input.Slug)
		if existing != nil && existing.ID != id {
			return nil, errors.New("板块标识已存在")
		}
		if err := validateSlug(input.Slug); err != nil {
			return nil, err
		}
		board.Slug = input.Slug
	}

	board.Name = input.Name
	board.Description = input.Description
	board.Icon = input.Icon
	board.Cover = input.Cover
	board.ParentID = input.ParentID
	board.SortOrder = input.SortOrder
	board.ViewRole = model.UserRole(input.ViewRole)
	board.PostRole = model.UserRole(input.PostRole)
	board.ReplyRole = model.UserRole(input.ReplyRole)

	if err := s.boardRepo.Update(board); err != nil {
		return nil, err
	}
	return board, nil
}

func (s *BoardService) Delete(id uint) error {
	result, err := s.boardRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("删除板块失败: %w", err)
	}
	if result == 0 {
		return apperrors.ErrBoardNotFound
	}
	return nil
}

func (s *BoardService) GetByID(id uint) (*model.Board, error) {
	return s.boardRepo.FindByID(id)
}

func (s *BoardService) GetBoardBySlug(slug string) (*model.Board, error) {
	return s.boardRepo.FindBySlug(slug)
}

func (s *BoardService) GetPostsBySlug(slug string, page, pageSize int) ([]*model.Post, int64, error) {
	return s.boardRepo.GetPostsBySlug(slug, page, pageSize)
}

func (s *BoardService) List(page, pageSize int) ([]model.Board, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.boardRepo.List(pageSize, offset)
}

func (s *BoardService) GetTree() ([]model.Board, error) {
	return s.boardRepo.GetTree()
}

func (s *BoardService) GetPosts(boardID uint, page, pageSize int) ([]model.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.postRepo.GetByBoardID(boardID, pageSize, offset)
}

// ── Moderator 核心操作 ────────────────────────────────────────────────────────

// AddModeratorInput 添加版主的参数（管理员直接任命 / 审批通过后使用）。
type AddModeratorInput struct {
	UserID             uint `json:"user_id"              binding:"required"`
	BoardID            uint `json:"board_id"             binding:"required"`
	CanDeletePost      bool `json:"can_delete_post"`
	CanPinPost         bool `json:"can_pin_post"`
	CanEditAnyPost     bool `json:"can_edit_any_post"`
	CanManageModerator bool `json:"can_manage_moderator"`
	CanBanUser         bool `json:"can_ban_user"`
}

// AddModerator 直接任命版主（管理员操作）。
// 若用户已有 pending 申请，同时将其标记为 canceled。
func (s *BoardService) AddModerator(_ context.Context, input AddModeratorInput, operatorID uint) error {
	// 确认用户存在
	user, err := s.userRepo.FindByID(input.UserID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// 幂等：已是版主则直接返回
	isMod, _ := s.boardRepo.IsModerator(input.UserID, input.BoardID)
	if isMod {
		return errors.New("用户已经是版主")
	}

	mod := &model.Moderator{
		UserID:  input.UserID,
		BoardID: input.BoardID,
	}
	if err := mod.SetPermissions(model.ModeratorPermissions{
		CanDeletePost:      input.CanDeletePost,
		CanPinPost:         input.CanPinPost,
		CanEditAnyPost:     input.CanEditAnyPost,
		CanManageModerator: input.CanManageModerator,
		CanBanUser:         input.CanBanUser,
	}); err != nil {
		return fmt.Errorf("权限序列化失败: %w", err)
	}

	if err := s.boardRepo.AddModerator(mod); err != nil {
		return fmt.Errorf("添加版主失败: %w", err)
	}

	// 取消该用户在此板块的 pending 申请
	_ = s.boardRepo.CancelUserApplications(input.UserID, input.BoardID)

	// 记录操作日志
	s.writeLog(operatorID, input.BoardID, "add_moderator", "user", input.UserID, "直接任命版主")

	// 发送通知
	boardID := input.BoardID
	s.notifSvc.Create(user.ID, &operatorID, model.NotifySystem,
		"你已被任命为版主", &boardID, "board")

	return nil
}

// RemoveModerator 移除版主（管理员操作）。
func (s *BoardService) RemoveModerator(_ context.Context, userID, boardID uint, operatorID uint) error {
	isMod, _ := s.boardRepo.IsModerator(userID, boardID)
	if !isMod {
		return errors.New("该用户不是此板块的版主")
	}

	if err := s.boardRepo.RemoveModerator(userID, boardID); err != nil {
		return fmt.Errorf("移除版主失败: %w", err)
	}

	s.writeLog(operatorID, boardID, "remove_moderator", "user", userID, "移除版主")

	// 通知被移除者
	s.notifSvc.Create(userID, &operatorID, model.NotifySystem,
		"你已被移除版主职务", &boardID, "board")

	return nil
}

// GetModerators 获取板块版主列表。
func (s *BoardService) GetModerators(boardID uint) ([]model.Moderator, error) {
	return s.boardRepo.GetModerators(boardID)
}

// IsModerator 判断用户是否是版主。
func (s *BoardService) IsModerator(userID, boardID uint) (bool, error) {
	return s.boardRepo.IsModerator(userID, boardID)
}

// ── 权限升降级 ────────────────────────────────────────────────────────────────

// UpdateModeratorPermissionsInput 升降级版主权限的参数。
type UpdateModeratorPermissionsInput struct {
	UserID             uint `json:"user_id"              binding:"required"`
	BoardID            uint `json:"board_id"             binding:"required"`
	CanDeletePost      bool `json:"can_delete_post"`
	CanPinPost         bool `json:"can_pin_post"`
	CanEditAnyPost     bool `json:"can_edit_any_post"`
	CanManageModerator bool `json:"can_manage_moderator"`
	CanBanUser         bool `json:"can_ban_user"`
}

// UpdateModeratorPermissions 修改版主权限（升级或降级）。
// 调用方（handler）已通过 AdminRequired 中间件确保只有管理员可操作。
func (s *BoardService) UpdateModeratorPermissions(
	_ context.Context,
	input UpdateModeratorPermissionsInput,
	operatorID uint,
) error {
	mod, err := s.boardRepo.FindModeratorByUserAndBoard(input.UserID, input.BoardID)
	if err != nil {
		return errors.New("版主记录不存在")
	}

	oldPerms, _ := mod.GetPermissions()

	newPerms := model.ModeratorPermissions{
		CanDeletePost:      input.CanDeletePost,
		CanPinPost:         input.CanPinPost,
		CanEditAnyPost:     input.CanEditAnyPost,
		CanManageModerator: input.CanManageModerator,
		CanBanUser:         input.CanBanUser,
	}
	if err := mod.SetPermissions(newPerms); err != nil {
		return fmt.Errorf("权限序列化失败: %w", err)
	}

	if err := s.boardRepo.UpdateModerator(mod); err != nil {
		return fmt.Errorf("更新版主权限失败: %w", err)
	}

	// 记录日志（old/new 均存入 JSON 字符串）
	s.writeLogWithValues(operatorID, input.BoardID,
		"update_moderator_perms", "moderator", mod.ID,
		"更新版主权限",
		fmt.Sprintf("%+v", oldPerms),
		fmt.Sprintf("%+v", newPerms),
	)

	// 通知版主本人
	s.notifSvc.Create(input.UserID, &operatorID, model.NotifySystem,
		"你的版主权限已被更新", &input.BoardID, "board")

	return nil
}

// CheckModeratorPermission 检查版主是否拥有某项具体权限。
func (s *BoardService) CheckModeratorPermission(_ context.Context, userID, boardID uint, permission string) (bool, error) {
	mod, err := s.boardRepo.FindModeratorByUserAndBoard(userID, boardID)
	if err != nil {
		return false, nil // 非版主视为无权限
	}
	return mod.HasPermission(permission), nil
}

// ── 版主申请流程 ──────────────────────────────────────────────────────────────

// ApplyModerator 用户提交版主申请。
// 规则：同一用户在同一板块只能有一条 pending 申请。
func (s *BoardService) ApplyModerator(input model.ApplyModeratorInput) error {
	// 已是版主则无需申请
	isMod, _ := s.boardRepo.IsModerator(input.UserID, input.BoardID)
	if isMod {
		return errors.New("你已经是该板块的版主")
	}

	// 已有待审核申请
	existing, err := s.boardRepo.FindPendingApplication(input.UserID, input.BoardID)
	if err != nil {
		return fmt.Errorf("查询申请失败: %w", err)
	}
	if existing != nil {
		return errors.New("你已有一条待审核的申请，请等待管理员处理")
	}

	app := &model.ModeratorApplication{
		UserID:             input.UserID,
		BoardID:            input.BoardID,
		Reason:             input.Reason,
		Status:             model.ApplicationPending,
		ReqDeletePost:      input.ReqDeletePost,
		ReqPinPost:         input.ReqPinPost,
		ReqEditAnyPost:     input.ReqEditAnyPost,
		ReqManageModerator: input.ReqManageModerator,
		ReqBanUser:         input.ReqBanUser,
	}

	if err := s.boardRepo.CreateApplication(app); err != nil {
		return fmt.Errorf("提交申请失败: %w", err)
	}
	return nil
}

// CancelApplication 用户撤销自己的待审核申请。
func (s *BoardService) CancelApplication(applicationID, userID uint) error {
	app, err := s.boardRepo.GetApplicationByID(applicationID)
	if err != nil || app == nil {
		return errors.New("申请不存在")
	}
	if app.UserID != userID {
		return errors.New("无权操作此申请")
	}
	if app.Status != model.ApplicationPending {
		return errors.New("只能撤销待审核的申请")
	}

	app.Status = model.ApplicationCanceled
	return s.boardRepo.UpdateApplication(app)
}

// GetUserApplications 获取用户的所有申请记录
func (s *BoardService) GetUserApplications(userID uint, page, pageSize int) ([]model.ModeratorApplication, int64, error) {
	// 参数默认值
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	applications, total, err := s.boardRepo.GetApplicationsByUserID(userID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("查询申请记录失败: %w", err)
	}
	return applications, total, nil
}

// ApplicationStatusDetail 申请状态详情
type ApplicationStatusDetail struct {
	HasApplication bool       `json:"has_application"`
	ApplicationID  uint       `json:"application_id,omitempty"`
	Status         string     `json:"status,omitempty"` // pending, approved, rejected, canceled
	Reason         string     `json:"reason,omitempty"`
	CreatedAt      time.Time  `json:"created_at,omitempty"`
	ReviewNote     string     `json:"review_note,omitempty"`
	ReviewerID     uint       `json:"reviewer_id,omitempty"`
	ReviewedAt     *time.Time `json:"reviewed_at,omitempty"`

	// 可执行的操作
	CanCancel   bool `json:"can_cancel"`   // 是否可以撤销
	CanResubmit bool `json:"can_resubmit"` // 是否可以重新申请

	// 申请的权限
	RequestedPerms map[string]bool `json:"requested_perms,omitempty"`

	// 是否可以申请（没有申请或申请已拒绝/已撤销时可为true）
	CanApply bool `json:"can_apply"`
}

// ReviewApplicationInput 管理员审批申请的参数。
type ReviewApplicationInput struct {
	ApplicationID uint `json:"application_id" binding:"required"`
	// Approve == true 表示通过，false 表示拒绝
	Approve    bool   `json:"approve"`
	ReviewNote string `json:"review_note" binding:"max=500"`
	// 通过时可覆盖申请者请求的权限（管理员按需调整）
	CanDeletePost      *bool `json:"can_delete_post"`
	CanPinPost         *bool `json:"can_pin_post"`
	CanEditAnyPost     *bool `json:"can_edit_any_post"`
	CanManageModerator *bool `json:"can_manage_moderator"`
	CanBanUser         *bool `json:"can_ban_user"`
}

// ReviewApplication 管理员审批版主申请。
// 通过时以审批参数（优先）或申请者请求的权限创建 Moderator 记录。
// MARK: 审批
func (s *BoardService) ReviewApplication(_ context.Context, input ReviewApplicationInput, reviewerID uint) error {
	app, err := s.boardRepo.GetApplicationByID(input.ApplicationID)
	if err != nil || app == nil {
		return errors.New("申请不存在")
	}
	if app.Status != model.ApplicationPending {
		return errors.New("该申请已被处理")
	}

	// 更新申请状态
	if input.Approve {
		app.Status = model.ApplicationApproved
	} else {
		app.Status = model.ApplicationRejected
	}
	app.ReviewerID = &reviewerID
	app.ReviewNote = input.ReviewNote

	if err := s.boardRepo.UpdateApplication(app); err != nil {
		return fmt.Errorf("更新申请状态失败: %w", err)
	}

	if !input.Approve {
		// 拒绝：通知申请者
		s.notifSvc.Create(app.UserID, &reviewerID, model.NotifySystem,
			fmt.Sprintf("你的版主申请已被拒绝：%s", input.ReviewNote), &app.BoardID, "board")
		return nil
	}

	// 通过：若用户已是版主（并发情况）则跳过创建
	isMod, _ := s.boardRepo.IsModerator(app.UserID, app.BoardID)
	if !isMod {
		// 权限以审批参数为准，未传则 fallback 到申请者的请求值
		perms := model.ModeratorPermissions{
			CanDeletePost:      boolVal(input.CanDeletePost, app.ReqDeletePost),
			CanPinPost:         boolVal(input.CanPinPost, app.ReqPinPost),
			CanEditAnyPost:     boolVal(input.CanEditAnyPost, app.ReqEditAnyPost),
			CanManageModerator: boolVal(input.CanManageModerator, app.ReqManageModerator),
			CanBanUser:         boolVal(input.CanBanUser, app.ReqBanUser),
		}

		mod := &model.Moderator{UserID: app.UserID, BoardID: app.BoardID}
		if err := mod.SetPermissions(perms); err != nil {
			return fmt.Errorf("权限序列化失败: %w", err)
		}
		if err := s.boardRepo.AddModerator(mod); err != nil {
			return fmt.Errorf("创建版主失败: %w", err)
		}

		s.writeLog(reviewerID, app.BoardID, "approve_application", "user", app.UserID, "审批申请通过")
	}

	// 通知申请者
	s.notifSvc.Create(app.UserID, &reviewerID, model.NotifySystem,
		"恭喜！你的版主申请已通过", &app.BoardID, "board")

	return nil
}

// ListApplications 分页查询申请列表（管理员用）。
func (s *BoardService) ListApplications(
	boardID *uint,
	status model.ApplicationStatus,
	page, pageSize int,
) ([]model.ModeratorApplication, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.boardRepo.ListApplications(boardID, status, page, pageSize)
}

// MARK: 禁言
type BanUserInput struct {
	UserID    uint       `json:"user_id"  binding:"required"`
	BoardID   uint       `json:"board_id" binding:"required"`
	Reason    string     `json:"reason"   binding:"required,max=500"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (s *BoardService) BanUser(input BanUserInput, bannerID uint) error {
	user, err := s.userRepo.FindByID(input.UserID)
	if err != nil {
		return errors.New("用户不存在")
	}

	isBanned, _ := s.boardRepo.IsBanned(input.UserID, input.BoardID)
	if isBanned {
		return errors.New("用户已被禁言")
	}

	ban := &model.BoardBan{
		UserID:    input.UserID,
		BoardID:   input.BoardID,
		BannedBy:  bannerID,
		Reason:    input.Reason,
		ExpiresAt: input.ExpiresAt,
	}
	if err := s.boardRepo.BanUser(ban); err != nil {
		return err
	}

	s.notifSvc.Create(user.ID, &bannerID, model.NotifySystem,
		"你在板块中被禁言", &input.BoardID, "board")
	return nil
}

func (s *BoardService) UnbanUser(userID, boardID uint) error {
	return s.boardRepo.UnbanUser(userID, boardID)
}

func (s *BoardService) IsBanned(userID, boardID uint) (bool, error) {
	return s.boardRepo.IsBanned(userID, boardID)
}

// ── 帖子管理（版主） ──────────────────────────────────────────────────────────

// DeletePost 版主删除帖子。
func (s *BoardService) DeletePost(boardID, postID, userID uint, isAdmin bool) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	if post.BoardID != boardID {
		return errors.New("帖子不属于该板块")
	}

	isMod, _ := s.boardRepo.IsModerator(userID, boardID)
	if !isMod && !isAdmin {
		return errors.New("无权限删除此帖子")
	}

	s.writeLog(userID, boardID, "delete_post", "post", postID, "版主删除")
	return s.postRepo.Delete(postID)
}

// PinPost 版主置顶/取消置顶帖子。
func (s *BoardService) PinPost(boardID, postID uint, pin bool) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	if post.BoardID != boardID {
		return errors.New("帖子不属于该板块")
	}
	return s.postRepo.TogglePinInBoard(postID, pin)
}

// ── 辅助函数 ──────────────────────────────────────────────────────────────────

// validateSlug 校验 slug 格式（小写字母、数字、横线、下划线）。
func validateSlug(slug string) error {
	if len(slug) == 0 || len(slug) > 50 {
		return errors.New("板块标识长度必须在1-50字符之间")
	}
	matched, _ := regexp.MatchString(`^[a-z0-9\-_]+$`, slug)
	if !matched {
		return errors.New("板块标识只能包含小写字母、数字、横线和下划线")
	}
	return nil
}

// validateRoles 校验 ViewRole / PostRole / ReplyRole 是否合法。
func validateRoles(roles ...string) error {
	valid := map[model.UserRole]bool{
		model.RoleGuest:     true,
		model.RoleUser:      true,
		model.RoleMember:    true,
		model.RoleModerator: true,
		model.RoleAdmin:     true,
	}
	for _, r := range roles {
		if r != "" && !valid[model.UserRole(r)] {
			return fmt.Errorf("无效的角色值: %s", r)
		}
	}
	return nil
}

// boolVal 返回指针值（若非 nil）否则返回 fallback。
func boolVal(ptr *bool, fallback bool) bool {
	if ptr != nil {
		return *ptr
	}
	return fallback
}

// writeLog 记录版主操作日志（fire-and-forget，忽略错误）。
func (s *BoardService) writeLog(moderatorID, boardID uint, action, targetType string, targetID uint, reason string) {
	log := &model.ModeratorLog{
		ModeratorID: moderatorID,
		BoardID:     boardID,
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		Reason:      reason,
	}
	_ = s.boardRepo.CreateModeratorLog(log)
}

// writeLogWithValues 同 writeLog，额外写入 OldValue / NewValue。
func (s *BoardService) writeLogWithValues(
	moderatorID, boardID uint,
	action, targetType string, targetID uint,
	reason, oldValue, newValue string,
) {
	log := &model.ModeratorLog{
		ModeratorID: moderatorID,
		BoardID:     boardID,
		Action:      action,
		TargetType:  targetType,
		TargetID:    targetID,
		Reason:      reason,
		OldValue:    oldValue,
		NewValue:    newValue,
	}
	_ = s.boardRepo.CreateModeratorLog(log)
}
