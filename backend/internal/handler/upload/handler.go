// package upload

// import (
// 	"fmt"
// 	"io"
// 	"mime/multipart"
// 	"path/filepath"
// 	"strings"
// 	"tiny-forum/internal/model/do"
// )

// type FileTypeHandler interface {
// 	// Validate 校验文件内容、大小、MIME等
// 	Validate(file *multipart.FileHeader) error
// 	// Process 可选的预处理（如压缩图片），返回新的reader或原样
// 	Process(src io.Reader, meta *do.Attachment) (io.Reader, error)
// 	// GetStoragePath 生成存储相对路径（不含baseDir）
// 	GetStoragePath(meta *do.Attachment) string
// }

// type BaseHandler struct {
// 	MaxSize   int64
// 	AllowMimes []string
// }

// func (h *BaseHandler) Validate(file *multipart.FileHeader) error {
// 	if file.Size > h.MaxSize {
// 		return fmt.Errorf("文件过大，最大 %d 字节", h.MaxSize)
// 	}
// 	ct := file.Header.Get("Content-Type")
// 	for _, allowed := range h.AllowMimes {
// 		if strings.HasPrefix(ct, allowed) {
// 			return nil
// 		}
// 	}
// 	return fmt.Errorf("不支持的文件类型: %s", ct)
// }

// func (h *BaseHandler) Process(src io.Reader, meta *do.Attachment) (io.Reader, error) {
// 	// 默认不做处理，直接返回原reader
// 	return src, nil
// }

// // PostImageHandler 帖子图片
// type PostImageHandler struct {
// 	BaseHandler
// 	WidthLimit int
// }

// func (h *PostImageHandler) GetStoragePath(meta *do.Attachment) string {
// 	return filepath.Join("post_images", fmt.Sprintf("%d", meta.UserID), meta.StoredName)
// }

// // PluginAssetHandler 插件资源文件
// type PluginAssetHandler struct {
// 	BaseHandler
// }

// func (h *PluginAssetHandler) GetStoragePath(meta *do.Attachment) string {
// 	return filepath.Join("plugins", meta.PluginID, meta.StoredName)
// }

// // 注册表
// var handlers = map[do.FileType]FileTypeHandler{
// 	do.FileTypePostImage: &PostImageHandler{
// 		BaseHandler: BaseHandler{MaxSize: 10 << 20, AllowMimes: []string{"image/"}},
// 		WidthLimit:  4096,
// 	},
// 	do.FileTypePluginAsset: &PluginAssetHandler{
// 		BaseHandler: BaseHandler{MaxSize: 50 << 20, AllowMimes: []string{"application/javascript", "text/css", "image/"}},
// 	},
// 	// 可根据需要添加 FileTypeAvatar, FileTypePostFile 等
// }

package upload
