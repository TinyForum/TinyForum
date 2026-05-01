package user

import (
	"errors"
	"fmt"
	"tiny-forum/internal/model/po"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

// GetProfile 获取自己的资料（简化）
func (s *userService) GetProfile(userID uint) (*po.User, error) {
	return s.repo.FindByID(userID)
}

// UpdateProfile 更新个人资料
func (s *userService) UpdateProfile(userID uint, input po.UpdateProfileInput) error {
	fields := map[string]interface{}{}
	if input.Bio != "" {
		fields["bio"] = input.Bio
	}
	if input.Avatar != "" {
		fields["avatar"] = input.Avatar
	}
	if input.Email != "" {
		fields["email"] = input.Email
	}
	if len(fields) == 0 {
		return nil
	}
	return s.repo.UpdateFields(userID, fields)
}

// GetUserProfile 获取他人资料（含关注统计）
func (s *userService) GetUserProfile(targetID, viewerID uint) (*UserProfileResponse, error) {
	user, err := s.repo.FindByID(targetID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound.WithMessagef("ID: %d", targetID)
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	resp := &UserProfileResponse{
		User:           user,
		FollowerCount:  s.repo.GetFollowerCount(targetID),
		FollowingCount: s.repo.GetFollowingCount(targetID),
	}

	if viewerID > 0 {
		resp.IsFollowing = s.repo.IsFollowing(viewerID, targetID)
	}

	return resp, nil
}

// GetUserBasicInfo 获取用户基本信息
func (s *userService) GetUserBasicInfo(userID uint) (*po.User, error) {
	return s.repo.GetUserBasicInfoById(userID)
}

// GetUserRoleById 获取用户角色
func (s *userService) GetUserRoleById(userID uint) (string, error) {
	return s.repo.GetUserRoleById(userID)
}
