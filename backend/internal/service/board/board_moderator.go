package board

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"tiny-forum/internal/model"
)

type AddModeratorInput struct {
	UserID             uint `json:"user_id"              binding:"required"`
	BoardID            uint `json:"board_id"             binding:"required"`
	CanDeletePost      bool `json:"can_delete_post"`
	CanPinPost         bool `json:"can_pin_post"`
	CanEditAnyPost     bool `json:"can_edit_any_post"`
	CanManageModerator bool `json:"can_manage_moderator"`
	CanBanUser         bool `json:"can_ban_user"`
}

type UpdateModeratorPermissionsInput struct {
	UserID             uint `json:"user_id"              binding:"required"`
	BoardID            uint `json:"board_id"             binding:"required"`
	CanDeletePost      bool `json:"can_delete_post"`
	CanPinPost         bool `json:"can_pin_post"`
	CanEditAnyPost     bool `json:"can_edit_any_post"`
	CanManageModerator bool `json:"can_manage_moderator"`
	CanBanUser         bool `json:"can_ban_user"`
}

type ModeratorBoardWithPerms struct {
	model.Board
	Permissions model.ModeratorPermissions `json:"permissions"`
}

func (s *BoardService) AddModerator(_ context.Context, input AddModeratorInput, operatorID uint) error {
	user, err := s.userRepo.FindByID(input.UserID)
	if err != nil {
		return errors.New("用户不存在")
	}
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
	_ = s.boardRepo.CancelUserApplications(input.UserID, input.BoardID)
	s.writeLog(operatorID, input.BoardID, "add_moderator", "user", input.UserID, "直接任命版主")
	boardID := input.BoardID
	s.notifSvc.Create(user.ID, &operatorID, model.NotifySystem,
		"你已被任命为版主", &boardID, "board")
	return nil
}

func (s *BoardService) RemoveModerator(_ context.Context, userID, boardID uint, operatorID uint) error {
	isMod, _ := s.boardRepo.IsModerator(userID, boardID)
	if !isMod {
		return errors.New("该用户不是此板块的版主")
	}
	if err := s.boardRepo.RemoveModerator(userID, boardID); err != nil {
		return fmt.Errorf("移除版主失败: %w", err)
	}
	s.writeLog(operatorID, boardID, "remove_moderator", "user", userID, "移除版主")
	s.notifSvc.Create(userID, &operatorID, model.NotifySystem,
		"你已被移除版主职务", &boardID, "board")
	return nil
}

func (s *BoardService) GetModerators(boardID uint) ([]model.Moderator, error) {
	return s.boardRepo.GetModerators(boardID)
}

func (s *BoardService) IsModerator(userID, boardID uint) (bool, error) {
	return s.boardRepo.IsModerator(userID, boardID)
}

func (s *BoardService) UpdateModeratorPermissions(_ context.Context, input UpdateModeratorPermissionsInput, operatorID uint) error {
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
	s.writeLogWithValues(operatorID, input.BoardID,
		"update_moderator_perms", "moderator", mod.ID,
		"更新版主权限",
		fmt.Sprintf("%+v", oldPerms),
		fmt.Sprintf("%+v", newPerms),
	)
	s.notifSvc.Create(input.UserID, &operatorID, model.NotifySystem,
		"你的版主权限已被更新", &input.BoardID, "board")
	return nil
}

func (s *BoardService) CheckModeratorPermission(_ context.Context, userID, boardID uint, permission string) (bool, error) {
	mod, err := s.boardRepo.FindModeratorByUserAndBoard(userID, boardID)
	if err != nil {
		return false, nil
	}
	return mod.HasPermission(permission), nil
}

func (s *BoardService) GetModeratorBoardsWithPermissions(userID uint) ([]ModeratorBoardWithPerms, error) {
	repoResults, err := s.boardRepo.GetModeratorBoardsWithPermissions(userID)
	if err != nil {
		return nil, err
	}
	results := make([]ModeratorBoardWithPerms, len(repoResults))
	for i, repo := range repoResults {
		var perms model.ModeratorPermissions
		if repo.Permissions != "" {
			if err := json.Unmarshal([]byte(repo.Permissions), &perms); err != nil {
				perms = model.ModeratorPermissions{
					CanDeletePost:      false,
					CanPinPost:         false,
					CanEditAnyPost:     false,
					CanManageModerator: false,
					CanBanUser:         false,
				}
			}
		}
		results[i] = ModeratorBoardWithPerms{
			Board:       repo.Board,
			Permissions: perms,
		}
	}
	return results, nil
}
