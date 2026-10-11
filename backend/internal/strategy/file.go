package strategy

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
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

// Validate 验证文件大小与真实内容类型
// 不能信任客户端提交的 Content-Type，这里通过读取文件头（magic bytes）嗅探真实 MIME，
// 并拒绝可在浏览器中执行的内容（HTML/SVG/XML/JS），防止存储型 XSS。
func (h *BaseHandler) Validate(file *multipart.FileHeader) error {
	if file.Size > h.MaxSize {
		return fmt.Errorf("file too large: max %d bytes", h.MaxSize)
	}
	f, err := file.Open()
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return fmt.Errorf("read file: %w", err)
	}
	detected := strings.ToLower(http.DetectContentType(buf[:n]))

	// 允许的媒介类型前缀
	for _, allowed := range h.AllowedMime {
		if strings.HasPrefix(detected, strings.ToLower(allowed)) {
			// image/* 仍需单独排除可执行内容的子类型
			if strings.HasPrefix(detected, "image/") && strings.Contains(detected, "svg") {
				return fmt.Errorf("file type not allowed: %s", detected)
			}
			return nil
		}
	}
	// 黑名单：即使不在允许列表中，也必须明确拒绝可执行内容
	for _, bad := range []string{"text/html", "text/xml", "application/xhtml+xml", "application/javascript", "application/x-httpd-php"} {
		if strings.HasPrefix(detected, bad) {
			return fmt.Errorf("file type not allowed: %s", detected)
		}
	}
	return fmt.Errorf("unsupported mime type: %s", detected)
}

func (h *BaseHandler) Process(src io.Reader, meta *do.Attachment) (io.Reader, error) {
	return src, nil // 默认不处理
}

type PluginHandler struct {
	BaseHandler
}

func (h *PluginHandler) GetStoragePath(meta *do.Attachment) string {
	// 存储到特定目录，例如 uploads/plugins/
	return filepath.Join("plugins", meta.StoredName)
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

// VideoHandler 视频文件
type VideoHandler struct {
	BaseHandler
}

func (h *VideoHandler) GetStoragePath(meta *do.Attachment) string {
	return filepath.Join("videos", fmt.Sprintf("%d", meta.UserID), meta.StoredName)
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
	// reg.Register(do.FileTypePlugin, &PluginAssetHandler{
	// 	BaseHandler: BaseHandler{MaxSize: 50 << 20, AllowedMime: []string{"application/javascript", "text/css", "image/", "application/json"}},
	// })
	reg.Register(do.FileTypeAvatar, &AvatarHandler{
		BaseHandler: BaseHandler{MaxSize: 2 << 20, AllowedMime: []string{"image/"}},
		MaxWidth:    512,
		MaxHeight:   512,
	})
	reg.Register(do.FileTypePlugin, &PluginHandler{
		BaseHandler: BaseHandler{
			MaxSize:     50 << 20, // 50MB，根据需求调整
			AllowedMime: []string{"application/zip", "application/x-zip-compressed"},
		},
	})
	reg.Register(do.FileTypePostCover, &PostImageHandler{
		BaseHandler: BaseHandler{MaxSize: 10 << 20, AllowedMime: []string{"image/"}},
		MaxWidth:    4096,
	})
	reg.Register(do.FileTypeTopicCover, &PostImageHandler{
		BaseHandler: BaseHandler{MaxSize: 10 << 20, AllowedMime: []string{"image/"}},
		MaxWidth:    4096,
	})
	reg.Register(do.FileTypeVideo, &VideoHandler{
		BaseHandler: BaseHandler{
			MaxSize:     500 << 20, // 500MB
			AllowedMime: []string{"video/"},
		},
	})
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
