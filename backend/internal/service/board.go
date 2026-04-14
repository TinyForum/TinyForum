package service

import (
	"errors"
	"fmt"
	"regexp"
	"time"
	apperrors "tiny-forum/internal/errors"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
)

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

type CreateBoardInput struct {
	Name        string `json:"name" binding:"required,min=2,max=50"`
	Slug        string `json:"slug" binding:"required,min=2,max=50"`
	Description string `json:"description" binding:"max=500"`
	Icon        string `json:"icon" binding:"max=100"`
	Cover       string `json:"cover" binding:"max=500"`
	ParentID    *uint  `json:"parent_id"`
	SortOrder   int    `json:"sort_order"`
	ViewRole    string `json:"view_role"`
	PostRole    string `json:"post_role"`
	ReplyRole   string `json:"reply_role"`
}

func (s *BoardService) Create(input CreateBoardInput) (*model.Board, error) {
	// 1. 验证 ParentID（关键修复）
	if input.ParentID != nil && *input.ParentID != 0 {
		// 检查父板块是否存在
		parentBoard, err := s.boardRepo.FindByID(*input.ParentID)
		if err != nil {
			return nil, fmt.Errorf("父板块不存在: id=%d", *input.ParentID)
		}
		if parentBoard == nil || parentBoard.ID == 0 {
			return nil, fmt.Errorf("父板块不存在: id=%d", *input.ParentID)
		}
	} else {
		// 确保根板块的 ParentID 为 nil
		input.ParentID = nil
	}

	// 2. 验证角色权限值
	validRoles := map[model.UserRole]bool{
		model.RoleGuest:     true,
		model.RoleUser:      true,
		model.RoleMember:    true,
		model.RoleModerator: true,
		model.RoleAdmin:     true,
	}

	if input.ViewRole != "" {
		if !validRoles[model.UserRole(input.ViewRole)] {
			return nil, errors.New("无效的查看角色")
		}
	}
	if input.PostRole != "" {
		if !validRoles[model.UserRole(input.PostRole)] {
			return nil, errors.New("无效的发帖角色")
		}
	}
	if input.ReplyRole != "" {
		if !validRoles[model.UserRole(input.ReplyRole)] {
			return nil, errors.New("无效的回复角色")
		}
	}

	// 3. 检查 Slug 唯一性
	existing, _ := s.boardRepo.FindBySlug(input.Slug)
	if existing != nil && existing.ID != 0 {
		return nil, errors.New("板块标识已存在")
	}

	// 4. 验证 Slug 格式
	if err := validateSlug(input.Slug); err != nil {
		return nil, err
	}

	// 5. 验证板块名称长度
	if len(input.Name) == 0 || len(input.Name) > 50 {
		return nil, errors.New("板块名称长度必须在1-50字符之间")
	}

	// 6. 创建板块对象
	board := &model.Board{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Icon:        input.Icon,
		Cover:       input.Cover,
		ParentID:    input.ParentID, // 现在已经是处理过的值
		SortOrder:   input.SortOrder,
		ViewRole:    model.UserRole(input.ViewRole),
		PostRole:    model.UserRole(input.PostRole),
		ReplyRole:   model.UserRole(input.ReplyRole),
	}

	// 7. 设置默认角色
	if board.ViewRole == "" {
		board.ViewRole = model.RoleUser
	}
	if board.PostRole == "" {
		board.PostRole = model.RoleUser
	}
	if board.ReplyRole == "" {
		board.ReplyRole = model.RoleUser
	}

	// 8. 创建板块
	if err := s.boardRepo.Create(board); err != nil {
		return nil, fmt.Errorf("创建板块失败: %w", err)
	}

	// 9. 返回完整的板块信息
	return s.boardRepo.FindByID(board.ID)
}

// 辅助函数：验证 Slug 格式
func validateSlug(slug string) error {
	if len(slug) == 0 || len(slug) > 50 {
		return errors.New("板块标识长度必须在1-50字符之间")
	}

	// 只允许小写字母、数字、横线和下划线
	matched, _ := regexp.MatchString(`^[a-z0-9\-_]+$`, slug)
	if !matched {
		return errors.New("板块标识只能包含小写字母、数字、横线和下划线")
	}
	return nil
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

// 获取所有帖子
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

// Moderator methods
type AddModeratorInput struct {
	UserID             uint `json:"user_id" binding:"required"`
	BoardID            uint `json:"board_id" binding:"required"`
	CanDeletePost      bool `json:"can_delete_post"`
	CanPinPost         bool `json:"can_pin_post"`
	CanEditAnyPost     bool `json:"can_edit_any_post"`
	CanManageModerator bool `json:"can_manage_moderator"`
	CanBanUser         bool `json:"can_ban_user"`
}

func (s *BoardService) AddModerator(input AddModeratorInput, operatorID uint) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(input.UserID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// Check if already moderator
	isMod, _ := s.boardRepo.IsModerator(input.UserID, input.BoardID)
	if isMod {
		return errors.New("用户已经是版主")
	}

	mod := &model.Moderator{
		UserID:             input.UserID,
		BoardID:            input.BoardID,
		CanDeletePost:      input.CanDeletePost,
		CanPinPost:         input.CanPinPost,
		CanEditAnyPost:     input.CanEditAnyPost,
		CanManageModerator: input.CanManageModerator,
		CanBanUser:         input.CanBanUser,
	}

	if err := s.boardRepo.AddModerator(mod); err != nil {
		return err
	}

	// Send notification
	s.notifSvc.Create(user.ID, &operatorID, model.NotifySystem,
		"你被任命为版主", &input.BoardID, "board")

	return nil
}

func (s *BoardService) RemoveModerator(userID, boardID uint) error {
	return s.boardRepo.RemoveModerator(userID, boardID)
}

func (s *BoardService) GetModerators(boardID uint) ([]model.Moderator, error) {
	return s.boardRepo.GetModerators(boardID)
}

func (s *BoardService) IsModerator(userID, boardID uint) (bool, error) {
	return s.boardRepo.IsModerator(userID, boardID)
}

// Ban methods
type BanUserInput struct {
	UserID    uint       `json:"user_id" binding:"required"`
	BoardID   uint       `json:"board_id" binding:"required"`
	Reason    string     `json:"reason" binding:"required,max=500"`
	ExpiresAt *time.Time `json:"expires_at"`
}

func (s *BoardService) BanUser(input BanUserInput, bannerID uint) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(input.UserID)
	if err != nil {
		return errors.New("用户不存在")
	}

	// Check if already banned
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

	// Send notification
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

// DeletePost 删除帖子（版主）
func (s *BoardService) DeletePost(boardID, postID, userID uint, isAdmin bool) error {
	// 检查帖子是否属于该板块
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	if post.BoardID != boardID {
		return errors.New("帖子不属于该板块")
	}

	// 检查权限（版主或管理员）
	isMod, _ := s.boardRepo.IsModerator(userID, boardID)
	if !isMod && !isAdmin {
		return errors.New("无权限删除此帖子")
	}

	// 记录版主操作日志
	log := &model.ModeratorLog{
		ModeratorID: userID,
		BoardID:     boardID,
		Action:      "delete_post",
		TargetType:  "post",
		TargetID:    postID,
		Reason:      "版主删除",
	}
	s.boardRepo.CreateModeratorLog(log)

	return s.postRepo.Delete(postID)
}

// PinPost 置顶帖子（版主）
func (s *BoardService) PinPost(boardID, postID uint, pin bool) error {
	// 检查帖子是否属于该板块
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	if post.BoardID != boardID {
		return errors.New("帖子不属于该板块")
	}

	return s.postRepo.TogglePinInBoard(postID, pin)
}
