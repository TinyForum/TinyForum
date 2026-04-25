package middleware

import (
	"bytes"
	"encoding/json"
	"io"

	"tiny-forum/internal/service/check"
	"tiny-forum/pkg/response"
	"tiny-forum/pkg/sensitive"

	"github.com/gin-gonic/gin"
)

// contextKey 私有类型，避免与其他包的 context key 冲突
type contextKey string

const (
	keyReviewRequired contextKey = "content_review_required" // 是否需要人工审核
	keyHitWords       contextKey = "content_hit_words"       // 命中关键词
	keyShadowed       contextKey = "content_shadowed"        // 内容被屏蔽
	keyBlocked        contextKey = "content_blocked"         // 内容被拦截
	keyReplace        contextKey = "content_replaced"        // 内容被替换
)

// ContentCheckMiddleware 内容安全前置检测中间件（同步 Pre-check）
// 从请求 body 中提取指定 JSON 字段，合并检测后统一决策：
//   - 一级情况：【安全】（文章级） → 放行，不注入审核标记
//   - 二级情况：【攻击】（文章级） → 任意字段命中 LevelReplace （且无 Block/Review）→ 放行，修改内容，标记用户违规，不注入审核标记
//   - 三级情况：【风险】（文章级）→ 任意字段命中 LevelReview （且无 Block）→ 放行，向 context 注入审核标记
//   - 四级情况：【屏蔽】（用户级）→ 任意字段命中 LevelShadowed （且无 Block/Review）→ 隐藏，向 context 注入隐藏标记
//   - 五级情况：【拦截】（用户级）→ 任意字段命中 LevelBlock → 返回 400，请求不进入 handler，标记用户风控行为
//
// fields: 需要检测的 JSON 字段名，如 []string{"title", "content"}
func ContentCheckMiddleware(checkSvc check.ContentCheckService, fields []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, restore, err := peekBody(c)
		if err != nil {
			// 读 body 失败不阻断业务，交由 handler 自行处理
			c.Next()
			return
		}

		restore() // 确保 handler 仍能读到完整 body

		var parsed map[string]any
		if err := json.Unmarshal(body, &parsed); err != nil {
			// 非 JSON body（如 form-data）直接放行
			c.Next()
			return
		}

		// 聚合所有字段的检测结果，统一决策
		var (
			maxLevel sensitive.Level
			allHits  []string
			seen     = make(map[string]bool)
		)
		for _, field := range fields {
			text := extractString(parsed, field)
			if text == "" {
				continue
			}
			res := checkSvc.CheckText(text)
			if res.Level > maxLevel {
				maxLevel = res.Level
			}
			for _, w := range res.HitWords {
				if !seen[w] {
					seen[w] = true
					allHits = append(allHits, w)
				}
			}
		}

		switch maxLevel {
		//  禁止发布
		case sensitive.LevelBlock:
			c.Set(string(keyBlocked), true)
			c.Set(string(keyHitWords), allHits)
			response.BadRequest(c, "内容存在风险，请修改后重新提交")
			c.Abort()
			return
			// 需要人工审核
		case sensitive.LevelReview:
			c.Set(string(keyReviewRequired), true)
			c.Set(string(keyHitWords), allHits)
			c.Next()
			return
			// 内容被隐藏
		case sensitive.LevelShadowed:
			c.Set(string(keyShadowed), true)
			c.Set(string(keyHitWords), allHits)
			return
			// 内容被替换
		case sensitive.LevelReplace:
			c.Set(string(keyReviewRequired), true)
			c.Set(string(keyHitWords), allHits)
			return
		}
		c.Next()
	}
}

// MARK: 内容核查
// IsReviewRequired 从 *gin.Context 中读取内容审核标记。
// 返回：(是否需要人工审核, 命中的敏感词列表)
func IsReviewRequired(c *gin.Context) (required bool, hitWords []string) {
	val, exists := c.Get(string(keyReviewRequired))
	if !exists {
		return false, nil
	}
	required, _ = val.(bool)
	if !required {
		return false, nil
	}
	if raw, ok := c.Get(string(keyHitWords)); ok {
		hitWords, _ = raw.([]string)
	}
	if hitWords == nil {
		hitWords = []string{}
	}
	return true, hitWords
}

// IsShadowed 从 *gin.Context 中读取内容隐藏标记。
func IsShadowed(c *gin.Context) (shadowed bool, hitWords []string) {
	val, exists := c.Get(string(keyShadowed))
	if !exists {
		return false, nil
	}
	shadowed, _ = val.(bool)
	if !shadowed {
		return false, nil
	}
	if raw, ok := c.Get(string(keyHitWords)); ok {
		hitWords, _ = raw.([]string)
	}
	return true, hitWords
}

// IsBlocked 从 *gin.Context 中读取内容屏蔽标记。
func IsBlocked(c *gin.Context) (blocked bool, hitWords []string) {
	val, exists := c.Get(string(keyBlocked))
	if !exists {
		return false, nil
	}
	blocked, _ = val.(bool)
	if !blocked {
		return false, nil
	}
	if raw, ok := c.Get(string(keyHitWords)); ok {
		hitWords, _ = raw.([]string)
	}
	return true, hitWords
}

// IsReplaced 从 *gin.Context 中读取内容替换标记。
func IsReplaced(c *gin.Context) (replaced bool, hitWords []string) {
	val, exists := c.Get(string(keyReplace))
	if !exists {
		return false, nil
	}
	replaced, _ = val.(bool)
	if !replaced {
		return false, nil
	}
	if raw, ok := c.Get(string(keyHitWords)); ok {
		hitWords, _ = raw.([]string)
	}
	return true, hitWords
}

// MARK: helpers

// peekBody 读取 body 并返回原始字节与还原函数。
// 调用者必须在使用完字节后调用 restore()，否则 handler 将读不到 body。
func peekBody(c *gin.Context) ([]byte, func(), error) {
	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, func() {}, err
	}
	restore := func() {
		c.Request.Body = io.NopCloser(bytes.NewReader(data))
	}
	restore() // 立即还原，让调用者在任意时刻都可再次调用
	return data, restore, nil
}

// extractString 从解析后的 JSON map 中安全提取字符串字段。
func extractString(m map[string]any, key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	s, _ := v.(string)
	return s
}
