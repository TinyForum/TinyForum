package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"tiny-forum/pkg/response"
	"tiny-forum/pkg/sensitive"

	riskservice "tiny-forum/internal/service/risk"

	"github.com/gin-gonic/gin"
)

// ContentCheckMiddleware 内容安全前置检测中间件（同步，Pre-check）
//
// 从请求 body 中读取指定字段进行敏感词检测：
//   - LevelBlock  → 直接返回 400，请求不进入 handler
//   - LevelReview → 放行，但在 context 中注入标记，handler 可据此将内容设为 pending
//
// fields: 需要检测的 JSON 字段名列表，例如 []string{"title", "content"}
func ContentCheckMiddleware(checkSvc *riskservice.ContentCheckService, fields []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 读取 body（读后必须复位，否则 handler 无法再读）
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.Next()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 解析为 map 提取指定字段
		var body map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			c.Next()
			return
		}

		for _, field := range fields {
			val, ok := body[field]
			if !ok {
				continue
			}
			text, ok := val.(string)
			if !ok || text == "" {
				continue
			}

			result := checkSvc.CheckText(text)

			if result.Level == sensitive.LevelBlock {
				response.BadRequest(c, "内容包含违禁词，请修改后重新提交")
				c.Abort()
				return
			}

			if result.Level == sensitive.LevelReview {
				// 注入标记到 context，handler 读取后将内容状态设为 pending
				c.Set("content_review_required", true)
				c.Set("content_hit_words", result.HitWords)
				// 不 Abort，继续处理
			}
		}

		c.Next()
	}
}

// IsReviewRequired 从 gin.Context 中读取是否需要人工审核标记
// IsReviewRequired 检查内容是否需要审核
// 返回值：是否需要审核，敏感词列表
func IsReviewRequired(c context.Context) (required bool, hitWords []string) {
	// 获取审核标记
	requiredVal := c.Value("content_review_required")
	if requiredVal == nil {
		return false, nil
	}

	required, ok := requiredVal.(bool)
	if !ok || !required {
		return false, nil
	}

	// 获取敏感词列表
	wordsVal := c.Value("content_hit_words")
	if wordsVal == nil {
		return true, nil
	}

	hitWords, ok = wordsVal.([]string)
	if !ok {
		// 类型断言失败，返回空切片而非 nil
		return true, []string{}
	}

	return true, hitWords
}
