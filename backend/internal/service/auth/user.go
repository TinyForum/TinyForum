package auth

import "context"

func (s *authService) IsUserExist(ctx context.Context, email string) (bool, error) {
	return s.userRepo.IsUserExistsByEmail(email)
}
