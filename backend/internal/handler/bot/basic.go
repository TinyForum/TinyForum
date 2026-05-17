package bot

import "tiny-forum/internal/service/bot"

// Handler bot HTTP 处理器
type Handler struct {
	svc bot.Service
}

func NewHandler(svc bot.Service) *Handler {
	return &Handler{svc: svc}
}
