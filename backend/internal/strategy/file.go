package strategy

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"tiny-forum/internal/model/do"
)

type FileTypeHandler interface {
	Validate(file *multipart.FileHeader) error
	Process(src io.Reader, meta *do.Attachment) (io.Reader, error)
	GetStoragePath(meta *do.Attachment) string
}

type BaseHandler struct {
	MaxSize     int64
	AllowedMime []string
}

func (h *BaseHandler) Validate(file *multipart.FileHeader) error {
	if file.Size > h.MaxSize {
		return fmt.Errorf("file too large: max %d bytes", h.MaxSize)
	}
	ct := file.Header.Get("Content-Type")
	for _, allowed := range h.AllowedMime {
		if strings.HasPrefix(ct, allowed) {
			return nil
		}
	}
	return fmt.Errorf("unsupported mime type: %s", ct)
}

func (h *BaseHandler) Process(src io.Reader, meta *do.Attachment) (io.Reader, error) {
	return src, nil // 默认不处理
}

// PostImageHandler 帖子图片
type PostImageHandler struct {
	BaseHandler
	MaxWidth int
}

func (h *PostImageHandler) GetStoragePath(meta *do.Attachment) string {
	return filepath.Join("post_images", fmt.Sprintf("%d", meta.UserID), meta.StoredName)
}

// PluginAssetHandler 插件静态资源（JS/CSS/图片等）
type PluginAssetHandler struct {
	BaseHandler
}

func (h *PluginAssetHandler) GetStoragePath(meta *do.Attachment) string {
	return filepath.Join("plugins", meta.PluginID, meta.StoredName)
}

// AvatarHandler 用户头像
type AvatarHandler struct {
	BaseHandler
	MaxWidth  int
	MaxHeight int
}

func (h *AvatarHandler) GetStoragePath(meta *do.Attachment) string {
	return filepath.Join("avatars", fmt.Sprintf("%d", meta.UserID), meta.StoredName)
}

// HandlerRegistry 注册表
type HandlerRegistry struct {
	handlers map[do.FileType]FileTypeHandler
}

func NewHandlerRegistry() *HandlerRegistry {
	reg := &HandlerRegistry{
		handlers: make(map[do.FileType]FileTypeHandler),
	}
	// 注册内置策略
	reg.Register(do.FileTypePostImage, &PostImageHandler{
		BaseHandler: BaseHandler{MaxSize: 10 << 20, AllowedMime: []string{"image/"}},
		MaxWidth:    4096,
	})
	reg.Register(do.FileTypePluginAsset, &PluginAssetHandler{
		BaseHandler: BaseHandler{MaxSize: 50 << 20, AllowedMime: []string{"application/javascript", "text/css", "image/", "application/json"}},
	})
	reg.Register(do.FileTypeAvatar, &AvatarHandler{
		BaseHandler: BaseHandler{MaxSize: 2 << 20, AllowedMime: []string{"image/"}},
		MaxWidth:    512,
		MaxHeight:   512,
	})
	// 可继续注册其他类型
	return reg
}

func (r *HandlerRegistry) Register(fileType do.FileType, handler FileTypeHandler) {
	r.handlers[fileType] = handler
}

func (r *HandlerRegistry) Get(fileType do.FileType) (FileTypeHandler, error) {
	handler, ok := r.handlers[fileType]
	if !ok {
		return nil, fmt.Errorf("unsupported file type: %s", fileType)
	}
	return handler, nil
}