package auth

import (
	"context"
	"tiny-forum/internal/model/do"
)

func (s *authService) FinduUserEmailByID(userID uint) (string, error) {
	user := &do.User{}
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return "", err
	}
	return user.Email, nil

}

func (s *authService) IsUserExist(ctx context.Context, email string) (bool, error) {
	return s.userRepo.IsUserExistsByEmail(email)
}
