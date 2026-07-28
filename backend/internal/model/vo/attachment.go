package vo

import (
	"tiny-forum/internal/model/dto"
)

type ListFilesVO struct {
	List  []*dto.FileInfo `json:"list"`
	Total int64           `json:"total"`
}
