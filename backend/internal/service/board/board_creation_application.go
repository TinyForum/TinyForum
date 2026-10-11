package board

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ApplyBoardCreation 用户申请创建新板块
func (s *boardService) ApplyBoardCreation(_ context.Context, req request.ApplyBoardCreationRequest) error {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return errors.New("板块名称不能为空")
	}
	if err := validateSlug(req.Slug); err != nil {
		return err
	}

	// 板块已存在（slug 冲突）时不允许申请
	if existing, _ := s.boardRepo.FindBySlug(req.Slug); existing != nil && existing.ID != 0 {
		return errors.New("该板块标识已存在")
	}

	// 同一 slug 同时只允许一条待审核申请
	pending, err := s.boardRepo.FindPendingBoardApplicationBySlug(req.Slug)
	if err != nil {
		return fmt.Errorf("查询待审核申请失败: %w", err)
	}
	if pending != nil {
		return errors.New("该板块标识已有待审核的申请，请等待管理员处理")
	}

	app := &do.BoardCreationApplication{
		UserID:      req.UserID,
		Name:        name,
		Slug:        req.Slug,
		Description: strings.TrimSpace(req.Description),
		Reason:      req.Reason,
		Status:      do.ApplicationPending,
	}
	if err := s.boardRepo.CreateBoardApplication(app); err != nil {
		return fmt.Errorf("提交申请失败: %w", err)
	}
	return nil
}

// GetUserBoardApplications 获取当前用户的板块创建申请记录
func (s *boardService) GetUserBoardApplications(userID uint, page, pageSize int) ([]do.BoardCreationApplication, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return s.boardRepo.GetBoardApplicationsByUserID(userID, page, pageSize)
}

// ListBoardApplications 管理员分页查询板块创建申请列表
func (s *boardService) ListBoardApplications(status do.ApplicationStatus, page, pageSize int) ([]do.BoardCreationApplication, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	if status != "" && !status.IsValid() {
		return nil, 0, fmt.Errorf("无效的申请状态: %s", status)
	}
	return s.boardRepo.ListBoardApplications(status, page, pageSize)
}

// ReviewBoardCreationApplication 管理员审批板块创建申请
// 批准时在同一事务内创建板块并落库审批状态，并锁定申请行避免并发重复审批。
func (s *boardService) ReviewBoardCreationApplication(ctx context.Context, req request.ReviewBoardCreationRequest, reviewerID uint) error {
	var (
		applicantID  uint
		createdBoard *do.Board
	)

	err := s.txManager.ExecuteInTransaction(ctx, func(tx *gorm.DB) error {
		// 锁定申请行，串行化并发审批，避免同一申请被重复批准而创建多个板块
		var app do.BoardCreationApplication
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&app, req.ApplicationID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("申请不存在")
			}
			return fmt.Errorf("获取申请失败: %w", err)
		}
		if !app.IsPending() {
			return errors.New("该申请已被处理")
		}
		applicantID = app.UserID

		// 拒绝分支：直接落库状态
		if !req.Approve {
			if err := app.Reject(reviewerID, req.ReviewNote); err != nil {
				return err
			}
			if err := tx.Save(&app).Error; err != nil {
				return fmt.Errorf("更新申请状态失败: %w", err)
			}
			return nil
		}

		// 批准分支：校验 slug 并创建板块，与审批状态落库保持原子
		if err := validateSlug(app.Slug); err != nil {
			return err
		}
		var count int64
		if err := tx.Model(&do.Board{}).Where("slug = ?", app.Slug).Count(&count).Error; err != nil {
			return fmt.Errorf("校验板块标识失败: %w", err)
		}
		if count > 0 {
			return errors.New("该板块标识已存在")
		}

		board := do.Board{
			Name:        app.Name,
			Slug:        app.Slug,
			Description: app.Description,
			ViewRole:    do.RoleUser,
			PostRole:    do.RoleUser,
			ReplyRole:   do.RoleUser,
		}
		if err := tx.Create(&board).Error; err != nil {
			return fmt.Errorf("创建板块失败: %w", err)
		}

		if err := app.Approve(reviewerID, req.ReviewNote); err != nil {
			return err
		}
		app.CreatedBoardID = &board.ID
		if err := tx.Save(&app).Error; err != nil {
			return fmt.Errorf("更新申请状态失败: %w", err)
		}
		createdBoard = &board
		return nil
	})
	if err != nil {
		return err
	}

	// 事务提交成功后再发送通知，避免通知发出而数据未落库
	if createdBoard != nil {
		s.sendBoardCreationNotification(applicantID, reviewerID, "恭喜！你申请的板块已创建成功", &createdBoard.ID)
	} else {
		s.sendBoardCreationNotification(applicantID, reviewerID, "你申请的板块创建申请已被拒绝："+req.ReviewNote, nil)
	}
	return nil
}

// sendBoardCreationNotification 发送板块创建申请结果通知
func (s *boardService) sendBoardCreationNotification(userID, reviewerID uint, message string, boardID *uint) {
	s.notifSvc.Create(
		userID,
		&reviewerID,
		do.NotifySystem,
		message,
		boardID,
		"board",
	)
}
