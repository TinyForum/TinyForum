// 文件：routes/plugin_files.go
package routes

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"tiny-forum/pkg/logger"

	"github.com/gin-gonic/gin"
)

func ListPluginFiles(c *gin.Context) {

	// 允许访问的根目录
	rootDir := "./store/plugins" // 请使用绝对路径或配置

	// wd, _ := os.Getwd()
	// logger.Infof("工作目录: %s", wd)
	absRoot, _ := filepath.Abs(rootDir)
	// logger.Infof("期望的插件目录: %s", absRoot)
	// 可选：限制只读用户或游客
	// if !isAllowed(c) { c.JSON(403, gin.H{"error": "forbidden"}) }

	// 获取子路径参数（安全限制）
	subPath := c.Query("path")
	if subPath == "" {
		subPath = "."
	}
	// 防止路径遍历
	if strings.Contains(subPath, "..") || strings.Contains(subPath, "\\") ||
		strings.HasPrefix(subPath, "/") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path"})
		return
	}

	fullPath := filepath.Join(rootDir, subPath)
	logger.Infof("路径：", fullPath)
	// 确保最终路径仍在 rootDir 内
	// absRoot, _ := filepath.Abs(rootDir)
	absFull, _ := filepath.Abs(fullPath)
	if !strings.HasPrefix(absFull, absRoot+string(os.PathSeparator)) && absFull != absRoot {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	entries, err := os.ReadDir(fullPath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "directory not found"})
		return
	}

	var fileList []FileInfo
	for _, entry := range entries {
		name := entry.Name()
		// 跳过隐藏文件（Unix 以点开头）
		if strings.HasPrefix(name, ".") {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			continue
		}
		// 可选：只允许特定扩展名
		ext := strings.ToLower(filepath.Ext(name))
		allowedExts := map[string]bool{
			".js": true, ".css": true, ".html": true, ".json": true,
			".png": true, ".jpg": true, ".jpeg": true, ".gif": true,
			".txt": true, ".md": true, ".zip": false, // zip 不显示
		}
		if !entry.IsDir() {
			if allowed, ok := allowedExts[ext]; !ok || !allowed {
				continue // 跳过不允许的文件类型
			}
		}
		fileList = append(fileList, FileInfo{
			Name:  name,
			IsDir: entry.IsDir(),
			Size:  info.Size(),
			Ext:   ext,
		})
	}
	parent := ""
	if subPath != "." {
		parent = filepath.Dir(subPath)
	}

	c.JSON(http.StatusOK, gin.H{
		"path":   subPath,
		"parent": parent,
		"files":  fileList,
	})

}

// 注册路由（放在需要鉴权的 Group 内或公开路由，取决于需求）
// api.GET("/plugins/files", ListPluginFiles)
