package moderator

import (
	"context"
	"encoding/json"
	"errors"
	"tiny-forum/internal/model/do"

	"gorm.io/gorm"
)

func (r *moderatorRepository) UpdatePermissions(ctx context.Context, moderatorID uint, permissions do.Permission) error {
	permsJSON, err := json.Marshal(permissions)
	if err != nil {
		return err
	}
	return r.db.WithContext(ctx).
		Model(&do.Moderator{}).
		Where("id = ?", moderatorID).
		Update("permissions", permsJSON).Error
}

func (r *moderatorRepository) HasPermission(ctx context.Context, userID, boardID uint, permission string) (bool, error) {
	var moderator do.Moderator
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND board_id = ?", userID, boardID).
		First(&moderator).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	perms, err := moderator.GetPermissions()
	if err != nil {
		return false, err
	}
	switch permission {
	case "delete_post":
		return perms.CanDeletePost, nil
	case "pin_post":
		return perms.CanPinPost, nil
	case "edit_any_post":
		return perms.CanEditAnyPost, nil
	case "manage_moderator":
		return perms.CanManageModerator, nil
	case "ban_user":
		return perms.CanBanUser, nil
	default:
		return false, nil
	}
}
