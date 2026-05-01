package board

import (
	"errors"
	"fmt"

	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/po"
	apperrors "tiny-forum/pkg/errors"
)

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

func (s *boardService) Create(input CreateBoardInput) (*po.Board, error) {
	if input.ParentID != nil && *input.ParentID != 0 {
		parent, err := s.boardRepo.FindByID(*input.ParentID)
		if err != nil || parent == nil || parent.ID == 0 {
			return nil, fmt.Errorf("父板块不存在: id=%d", *input.ParentID)
		}
	} else {
		input.ParentID = nil
	}

	if err := validateRoles(input.ViewRole, input.PostRole, input.ReplyRole); err != nil {
		return nil, err
	}
	if existing, _ := s.boardRepo.FindBySlug(input.Slug); existing != nil && existing.ID != 0 {
		return nil, errors.New("板块标识已存在")
	}
	if err := validateSlug(input.Slug); err != nil {
		return nil, err
	}

	board := &po.Board{
		Name:        input.Name,
		Slug:        input.Slug,
		Description: input.Description,
		Icon:        input.Icon,
		Cover:       input.Cover,
		ParentID:    input.ParentID,
		SortOrder:   input.SortOrder,
		ViewRole:    po.UserRole(input.ViewRole),
		PostRole:    po.UserRole(input.PostRole),
		ReplyRole:   po.UserRole(input.ReplyRole),
	}
	if board.ViewRole == "" {
		board.ViewRole = po.RoleUser
	}
	if board.PostRole == "" {
		board.PostRole = po.RoleUser
	}
	if board.ReplyRole == "" {
		board.ReplyRole = po.RoleUser
	}

	if err := s.boardRepo.Create(board); err != nil {
		return nil, fmt.Errorf("创建板块失败: %w", err)
	}
	return s.boardRepo.FindByID(board.ID)
}

func (s *boardService) Update(id uint, input CreateBoardInput) (*po.Board, error) {
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
	board.ViewRole = po.UserRole(input.ViewRole)
	board.PostRole = po.UserRole(input.PostRole)
	board.ReplyRole = po.UserRole(input.ReplyRole)

	if err := s.boardRepo.Update(board); err != nil {
		return nil, err
	}
	return board, nil
}

func (s *boardService) Delete(id uint) error {
	result, err := s.boardRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("删除板块失败: %w", err)
	}
	if result == 0 {
		return apperrors.ErrBoardNotFound
	}
	return nil
}

func (s *boardService) GetByID(id uint) (*po.Board, error) {
	return s.boardRepo.FindByID(id)
}

func (s *boardService) GetBoardBySlug(slug string) (*po.Board, error) {
	return s.boardRepo.FindBySlug(slug)
}

// service/board_service.go
func (s *boardService) GetPostsBySlug(slug string, page, pageSize int) ([]*dto.GetBoardPostsResponse, int64, error) {
	return s.boardRepo.GetPostsBySlug(slug, page, pageSize)
}

func (s *boardService) List(page, pageSize int) ([]po.Board, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.boardRepo.List(pageSize, offset)
}

func (s *boardService) GetTree() ([]po.Board, error) {
	return s.boardRepo.GetTree()
}

func (s *boardService) GetPosts(boardID uint, page, pageSize int) ([]po.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.postRepo.GetByBoardID(boardID, pageSize, offset)
}
