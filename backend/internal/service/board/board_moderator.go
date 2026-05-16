package board

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/pkg/logger"
)

// type AddModeratorInput struct {
// 	UserID             uint `json:"user_id"              binding:"required"`
// 	BoardID            uint `json:"board_id"             binding:"required"`
// 	CanDeletePost      bool `json:"can_delete_post"`
// 	CanPinPost         bool `json:"can_pin_post"`
// 	CanEditAnyPost     bool `json:"can_edit_any_post"`
// 	CanManageModerator bool `json:"can_manage_moderator"`
// 	CanBanUser         bool `json:"can_ban_user"`
// }

// type UpdateModeratorPermissionsInput struct {
// 	UserID             uint `json:"user_id"              binding:"required"`
// 	BoardID            uint `json:"board_id"             binding:"required"`
// 	CanDeletePost      bool `json:"can_delete_post"`
// 	CanPinPost         bool `json:"can_pin_post"`
// 	CanEditAnyPost     bool `json:"can_edit_any_post"`
// 	CanManageModerator bool `json:"can_manage_moderator"`
// 	CanBanUser         bool `json:"can_ban_user"`
// }

// type ModeratorBoardWithPerms struct {
// 	do.Board
// 	Permissions do.ModeratorPermissions `json:"permissions"`
// }

// AddModerator 直接添加版主（管理员操作）
func (s *boardService) AddModerator(ctx context.Context, input request.AddModeratorRequest, operatorID uint) error {
	// 1. 参数校验
	if err := input.Validate(); err != nil {
		return err
	}

	// 2. 检查用户是否存在
	user, err := s.userRepo.FindByID(input.UserID)
	if err != nil {
		return fmt.Errorf("用户不存在: %w", err)
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 3. 检查是否已经是版主（幂等性）
	isMod, err := s.boardRepo.IsModerator(input.UserID, input.BoardID)
	if err != nil {
		return fmt.Errorf("检查版主状态失败: %w", err)
	}
	if isMod {
		return errors.New("用户已经是版主")
	}

	// 4. 创建版主记录（直接使用权限切片）
	mod := &do.Moderator{
		UserID:      input.UserID,
		BoardID:     input.BoardID,
		Permissions: input.Permissions, // 直接赋值，无需序列化（PermissionSet 实现了 Scanner/Valuer）
	}

	// 5. 持久化
	if err := s.boardRepo.AddModerator(mod); err != nil {
		return fmt.Errorf("添加版主失败: %w", err)
	}

	// 6. 取消该用户在该板块的所有待处理申请（避免重复申请）
	if err := s.boardRepo.CancelUserApplications(input.UserID, input.BoardID); err != nil {
		// 非关键操作，仅记录日志，不阻断主流程
		logger.Warnf("取消用户申请失败", "user_id", input.UserID, "board_id", input.BoardID, "error", err)
	}

	// 7. 记录操作日志
	s.writeLog(operatorID, input.BoardID, "add_moderator", "user", input.UserID, "直接任命版主")

	// 8. 发送系统通知
	boardID := input.BoardID
	s.notifSvc.Create(
		user.ID,
		&operatorID,
		do.NotifySystem,
		fmt.Sprintf("你已被任命为版主，授予权限：%v", formatPermissions(input.Permissions)),
		&boardID,
		"board",
	)

	return nil
}

// formatPermissions 辅助函数，将权限切片格式化为可读字符串（用于通知）
// func formatPermissions(perms do.ModeratorPermissionSet) string {
// 	if len(perms) == 0 {
// 		return "无"
// 	}
// 	strs := make([]string, len(perms))
// 	for i, p := range perms {
// 		strs[i] = string(p)
// 	}
// 	return strings.Join(strs, "、")
// }

func (s *boardService) RemoveModerator(_ context.Context, userID, boardID uint, operatorID uint) error {
	isMod, _ := s.boardRepo.IsModerator(userID, boardID)
	if !isMod {
		return errors.New("该用户不是此板块的版主")
	}
	if err := s.boardRepo.RemoveModerator(userID, boardID); err != nil {
		return fmt.Errorf("移除版主失败: %w", err)
	}
	s.writeLog(operatorID, boardID, "remove_moderator", "user", userID, "移除版主")
	s.notifSvc.Create(userID, &operatorID, do.NotifySystem,
		"你已被移除版主职务", &boardID, "board")
	return nil
}

func (s *boardService) GetModerators(boardID uint) ([]do.Moderator, error) {
	return s.boardRepo.GetModerators(boardID)
}

func (s *boardService) IsModerator(userID, boardID uint) (bool, error) {
	return s.boardRepo.IsModerator(userID, boardID)
}

// UpdateModeratorPermissions 更新版主权限
func (s *boardService) UpdateModeratorPermissions(ctx context.Context, input request.UpdateModeratorPermissionsRequest, operatorID uint) error {
	// 1. 参数校验
	if input.UserID == 0 || input.BoardID == 0 {
		return errors.New("用户ID和版块ID不能为空")
	}
	// 权限合法性校验
	for _, perm := range input.Permissions {
		if !perm.IsValid() {
			return fmt.Errorf("无效的权限: %s", perm)
		}
	}
	// 去重校验（可选）
	seen := make(map[do.ModeratorPermission]bool)
	for _, perm := range input.Permissions {
		if seen[perm] {
			return fmt.Errorf("权限重复: %s", perm)
		}
		seen[perm] = true
	}

	// 2. 查找版主记录
	mod, err := s.boardRepo.FindModeratorByUserAndBoard(input.UserID, input.BoardID)
	if err != nil {
		return fmt.Errorf("查找版主记录失败: %w", err)
	}
	if mod == nil {
		return errors.New("版主记录不存在")
	}

	// 3. 记录旧权限（用于日志）
	oldPerms := mod.Permissions

	// 4. 更新权限（直接赋值切片）
	mod.Permissions = do.ModeratorPermissionSet(input.Permissions)

	// 5. 持久化
	if err := s.boardRepo.UpdateModerator(mod); err != nil {
		return fmt.Errorf("更新版主权限失败: %w", err)
	}

	// 6. 记录操作日志
	s.writeLogWithValues(operatorID, input.BoardID,
		"update_moderator_perms", "moderator", mod.ID,
		"更新版主权限",
		fmt.Sprintf("%v", oldPerms),
		fmt.Sprintf("%v", mod.Permissions),
	)

	// 7. 发送通知（异步/同步均可）
	s.notifSvc.Create(input.UserID, &operatorID, do.NotifySystem,
		fmt.Sprintf("你的版主权限已更新为：%v", formatPermissions(mod.Permissions)),
		&input.BoardID, "board")

	return nil
}

// CheckModeratorPermission 检查用户是否在指定板块拥有某权限
func (s *boardService) CheckModeratorPermission(ctx context.Context, userID, boardID uint, permission do.ModeratorPermission) (bool, error) {
	// 合法性校验
	if !permission.IsValid() {
		return false, fmt.Errorf("无效的权限标识: %s", permission)
	}

	mod, err := s.boardRepo.FindModeratorByUserAndBoard(userID, boardID)
	if err != nil {
		// 记录不存在视为无权限，不报错
		return false, nil
	}
	if mod == nil {
		return false, nil
	}
	return mod.HasPermission(permission), nil
}

type ModeratorBoardWithPerms struct {
	Board       do.Board                  `json:"board"`       // 板块信息
	Permissions do.ModeratorPermissionSet `json:"permissions"` // 权限集（切片）
}

// GetModeratorBoardsWithPermissions 获取用户管理的所有板块及对应权限
func (s *boardService) GetModeratorBoardsWithPermissions(userID uint) ([]ModeratorBoardWithPerms, error) {
	repoResults, err := s.boardRepo.GetModeratorBoardsWithPermissions(userID)
	if err != nil {
		return nil, fmt.Errorf("查询版主板块失败: %w", err)
	}

	results := make([]ModeratorBoardWithPerms, 0, len(repoResults))
	for _, repo := range repoResults {
		// 直接使用 repo.Permissions（类型为 do.ModeratorPermissionSet），无需反序列化
		// 如果 repo 中 Permissions 是 JSON 字符串，则需扫描；假设 repo 已处理为结构体
		// 注意：这里假设 boardRepo 返回的 result 已包含正确的 Permissions 字段（类型为 do.ModeratorPermissionSet）
		// 如果仍是字符串，需添加扫描逻辑（参考下方备选）
		results = append(results, ModeratorBoardWithPerms{
			Board:       repo.Board,
			Permissions: repo.Permissions, // 直接使用切片
		})
	}
	return results, nil
}

// 辅助函数：格式化权限切片用于日志/通知
func formatPermissions(perms do.ModeratorPermissionSet) string {
	if len(perms) == 0 {
		return "无"
	}
	strs := make([]string, len(perms))
	for i, p := range perms {
		strs[i] = string(p)
	}
	return strings.Join(strs, "、")
}
