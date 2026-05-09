package dto

// UploadResponse 上传响应
type UploadResponse struct {
	FileID       string `json:"file_id"` // 存储标识
	URL          string `json:"url"`     // 访问URL
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	MimeType     string `json:"mime_type"`
}

