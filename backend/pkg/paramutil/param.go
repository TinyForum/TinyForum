// pkg/utils/param.go
package paramutils

import (
	"strconv"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// GetUintParam 从 URL 路径获取 uint 参数，失败返回 0 和 false
func GetUintParam(c *gin.Context, key string) (uint, bool) {
	val := c.Param(key)
	if val == "" {
		return 0, false
	}
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, false
	}
	return uint(id), true
}

// MustGetUintParam 从 URL 路径获取 uint 参数，失败时自动返回错误响应
func MustGetUintParam(c *gin.Context, key string) (uint, bool) {
	id, ok := GetUintParam(c, key)
	if !ok {
		response.ValidationFailed(c, []response.ValidationError{
			{Field: key, Message: "无效的ID格式"},
		})
		return 0, false
	}
	return id, true
}
