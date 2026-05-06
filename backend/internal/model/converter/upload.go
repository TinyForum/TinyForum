package converter

import (
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/request"
)

// converter/upload_converter.go
func UploadPluginRequestToUploadPluginBo(req request.UploadPluginRequest, userID uint) bo.PluginUpdateBO {
	return bo.PluginUpdateBO{
		UserID:     userID,
		FileHeader: req.File, // req.File 是 *multipart.FileHeader
	}
}
