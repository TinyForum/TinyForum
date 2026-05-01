package admin

import "tiny-forum/internal/model/do"

func (s *adminService) ListUsers(page, pageSize int, keyword string) ([]do.User, int64, error) {
	return s.userSvc.List(page, pageSize, keyword)
}
func (s *adminService) SetActiveUser(targetID uint, operatorID uint, isActive bool) error {
	return s.userSvc.SetActive(targetID, operatorID, isActive)
}

func (s *adminService) SetBlockedUser(targetID uint, operatorID uint, isBlocked bool) error {
	return s.userSvc.SetBlocked(targetID, operatorID, isBlocked)

}

func (s *adminService) DeleteUser(operatorID uint, targetID uint) error {
	return s.userSvc.DeleteUser(operatorID, targetID)
}

func (s *adminService) SetRoleUser(operatorID, targetID uint, newRole string) error {
	return s.userSvc.SetRole(operatorID, targetID, newRole)
}
