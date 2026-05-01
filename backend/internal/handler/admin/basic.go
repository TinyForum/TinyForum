package admin

import "tiny-forum/internal/service/admin"

type AdminHandler struct {
	service admin.AdminService
}

func NewAnnouncementHandler(svc admin.AdminService) *AdminHandler {
	return &AdminHandler{service: svc}
}
