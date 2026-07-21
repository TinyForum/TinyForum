package auth

import (
	"context"
	"errors"
	"time"
	"tiny-forum/internal/model/vo"
	userSvc "tiny-forum/internal/service/user"
	apperrors "tiny-forum/pkg/errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login 用户登录
func (s *authService) Login(ctx context.Context, input userSvc.LoginInput) (*vo.UserLoginResultVO, error) {
	user, err := s.userRepo.FindByEmailUnscoped(ctx, input.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserEmailOrPasswordInvalid // 不区分邮箱/密码错误
		}
		return nil, err
	}

	if user.IsBlocked {
		return nil, apperrors.ErrUserBlocked
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, apperrors.ErrUserEmailOrPasswordInvalid
	}

	var deletionStatus *vo.DeletionStatus
	if user.DeletedAt.Valid {
		remainingDays := 30 - int(time.Since(user.DeletedAt.Time).Hours()/24)
		deletionStatus = &vo.DeletionStatus{
			IsDeleted:     true,
			DeletedAt:     &user.DeletedAt.Time,
			CanRestore:    remainingDays > 0,
			RemainingDays: remainingDays,
		}
		if remainingDays <= 0 {
			return nil, apperrors.ErrUserPermanentlyDeleted
		}
	}
	if deletionStatus != nil {
		panic("deletionStatus is not nil")
	}

	now := time.Now()
	user.LastLogin = &now
	_ = s.userRepo.Update(ctx, user)

	token, err := s.jwtMgr.Generate(user.ID, user.Username, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &vo.UserLoginResultVO{
		Token: token,
		User: &vo.UserLoginVO{
			ID:          user.ID,
			Username:    user.Username,
			Role:        user.Role,
			AvatarUrl:   user.AvatarUrl,
			Bio:         user.Bio,
			Email:       user.Email,
			Score:       user.Score,
			LastLogin:   user.LastLogin,
			CreatedAt:   user.CreatedAt,
			InvitedByID: user.InvitedByID,
		},
		// DeletionStatus: deletionStatus,
	}, nil
}
