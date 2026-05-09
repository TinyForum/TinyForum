// package attachment

// import (
// 	"crypto/rand"
// 	"encoding/hex"
// 	"fmt"
// 	"io"
// 	"mime/multipart"
// 	"net/http"
// 	"path/filepath"
// 	"strings"
// 	"time"
// 	"tiny-forum/internal/model/do"

// 	"github.com/disintegration/imaging"
// )

// // validatePluginFile 验证插件文件（仅允许 ZIP）
// func (s *service) validatePluginFile(fileHeader *multipart.FileHeader) (string, error) {
// 	// 获取扩展名
// 	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
// 	if ext != ".zip" {
// 		return "", fmt.Errorf("只允许上传 .zip 格式的插件包")
// 	}
// 	// // 可选：限制大小
// 	// if fileHeader.Size > s.maxPluginSize {
// 	// 	return "", fmt.Errorf("插件包大小超过限制（最大 %d MB）", s.maxPluginSize>>20)
// 	// }
// 	return ext, nil
// }

// // 处理图片（缩放、格式转换）
// func processImage(srcPath, dstPath string, maxWidth, maxHeight int) error {
// 	src, err := imaging.Open(srcPath)
// 	if err != nil {
// 		return err
// 	}

// 	bounds := src.Bounds()
// 	width := bounds.Dx()
// 	height := bounds.Dy()

// 	// 如果图片尺寸超过限制，进行缩放
// 	if width > maxWidth || height > maxHeight {
// 		src = imaging.Fit(src, maxWidth, maxHeight, imaging.Lanczos)
// 	}

// 	return imaging.Save(src, dstPath)
// }

// // 检查文件类型
// func (s *service) validateFile(header *multipart.FileHeader, fileType do.FileType) (string, error) {
// 	// 检查大小
// 	if header.Size > s.maxSize {
// 		return "", fmt.Errorf("文件过大，最大允许 %d MB", s.maxSize/(1024*1024))
// 	}

// 	// 检查扩展名
// 	ext := strings.ToLower(filepath.Ext(header.Filename))
// 	if !s.allowedExt[ext] {
// 		return "", fmt.Errorf("不支持的文件类型: %s", ext)
// 	}

// 	// 嗅探真实MIME类型
// 	file, err := header.Open()
// 	if err != nil {
// 		return "", fmt.Errorf("无法打开文件: %w", err)
// 	}
// 	defer file.Close()

// 	sniffBuf := make([]byte, 512)
// 	_, err = file.Read(sniffBuf)
// 	if err != nil && err != io.EOF {
// 		return "", fmt.Errorf("读取文件头失败: %w", err)
// 	}

// 	mime := http.DetectContentType(sniffBuf)

// 	// 根据类型额外验证
// 	switch fileType {
// 	case "avatar":
// 		if !strings.HasPrefix(mime, "image/") {
// 			return "", fmt.Errorf("头像必须是图片格式")
// 		}
// 	case "post_image":
// 		if !strings.HasPrefix(mime, "image/") {
// 			return "", fmt.Errorf("帖子图片必须是图片格式")
// 		}
// 	case "comment_attachment":
// 		// 附件允许更多类型
// 		allowedMimes := []string{"image/", "application/pdf", "application/zip"}
// 		allowed := false
// 		for _, allowedMime := range allowedMimes {
// 			if strings.HasPrefix(mime, allowedMime) {
// 				allowed = true
// 				break
// 			}
// 		}
// 		if !allowed {
// 			return "", fmt.Errorf("不支持的附件类型")
// 		}
// 	}

// 	return ext, nil
// }

// // 生成唯一文件ID
// func generateFileID() string {
// 	b := make([]byte, 16)
// 	rand.Read(b)
// 	return hex.EncodeToString(b)
// }

// // 生成随机文件名
// func randomFileName(ext string) string {
// 	b := make([]byte, 8)
// 	rand.Read(b)
// 	return fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), hex.EncodeToString(b), ext)
// }

package attachment