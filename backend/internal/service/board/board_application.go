package board

import (
	"context"
	"errors"
	"fmt"
	"time"

	"tiny-forum/internal/model"
)

type ApplicationStatusDetail struct {
	HasApplication bool            `json:"has_application"`
	ApplicationID  uint            `json:"application_id,omitempty"`
	Status         string          `json:"status,omitempty"`
	Reason         string          `json:"reason,omitempty"`
	CreatedAt      time.Time       `json:"created_at,omitempty"`
	ReviewNote     string          `json:"review_note,omitempty"`
	ReviewerID     uint            `json:"reviewer_id,omitempty"`
	ReviewedAt     *time.Time      `json:"reviewed_at,omitempty"`
	CanCancel      bool            `json:"can_cancel"`
	CanResubmit    bool            `json:"can_resubmit"`
	RequestedPerms map[string]bool `json:"requested_perms,omitempty"`
	CanApply       bool            `json:"can_apply"`
}

type ReviewApplicationInput struct {
	ApplicationID      uint   `json:"application_id" binding:"required"`
	Approve            bool   `json:"approve"`
	ReviewNote         string `json:"review_note" binding:"max=500"`
	CanDeletePost      *bool  `json:"can_delete_post"`
	CanPinPost         *bool  `json:"can_pin_post"`
	CanEditAnyPost     *bool  `json:"can_edit_any_post"`
	CanManageModerator *bool  `json:"can_manage_moderator"`
	CanBanUser         *bool  `json:"can_ban_user"`
}

func (s *BoardService) ApplyModerator(input model.ApplyModeratorInput) error {
	isMod, _ := s.boardRepo.IsModerator(input.UserID, input.BoardID)
	if isMod {
		return errors.New("你已经是该板块的版主")
	}
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

func (s *BoardService) GetUserApplications(userID uint, page, pageSize int) ([]model.ModeratorApplication, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return s.boardRepo.GetApplicationsByUserID(userID, page, pageSize)
}

func (s *BoardService) ReviewApplication(_ context.Context, input ReviewApplicationInput, reviewerID uint) error {
	app, err := s.boardRepo.GetApplicationByID(input.ApplicationID)
	if err != nil || app == nil {
		return errors.New("申请不存在")
	}
	if app.Status != model.ApplicationPending {
		return errors.New("该申请已被处理")
	}
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
		s.notifSvc.Create(app.UserID, &reviewerID, model.NotifySystem,
			fmt.Sprintf("你的版主申请已被拒绝：%s", input.ReviewNote), &app.BoardID, "board")
		return nil
	}
	isMod, _ := s.boardRepo.IsModerator(app.UserID, app.BoardID)
	if !isMod {
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
	s.notifSvc.Create(app.UserID, &reviewerID, model.NotifySystem,
		"恭喜！你的版主申请已通过", &app.BoardID, "board")
	return nil
}

func (s *BoardService) ListApplications(boardID *uint, status model.ApplicationStatus, page, pageSize int) ([]model.ModeratorApplication, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.boardRepo.ListApplications(boardID, status, page, pageSize)
}
