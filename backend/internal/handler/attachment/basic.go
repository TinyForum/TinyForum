package attachment

import "tiny-forum/internal/service/attachment"

type AttachmentHandler struct {
	svc attachment.AttachmentService
}

func NewAttachmentHandler(svc attachment.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{svc: svc}
}
