package sensitive

import (
	"log"
	"tiny-forum/pkg/utils"
)

// ---- Check（DFA + 可选 LLM 复判）----

// MARK: Check
// Check 检测文本，返回分级结果。
//
// 流程：
//  1. DFA 扫描，无命中 → LevelClean，直接返回
//  2. 有 LevelBlock 命中 → 直接返回 LevelBlock，不调 LLM（节省资源）
//  3. 仅有 LevelReview 命中 → 调 Ollama 复判：
//     - LLM 判定"安全" → 降级为 LevelClean
//     - LLM 判定"风险" → 保持 LevelReview
//     - LLM 判定"违规" → 升级为 LevelBlock
//     - LLM 调用失败  → 保守保持 LevelReview（不降级）
func (f *filter) Check(html string) CheckResult {
	text, err := utils.HTMLToText(html)
	if err != nil {
		return CheckResult{
			Level:    LevelBlock,
			HitWords: []string{},
			Text:     "HTML 解析失败",
		}
	}

	result := CheckResult{
		Level:    LevelClean,
		HitWords: []string{},
		Text:     text,
	}

	// --- Step 1: DFA 扫描 ---
	hitMap := f.dfa.FindAllCount(text)
	if len(hitMap) == 0 {
		return result
	}

	f.mu.RLock()
	var hits []string
	hasBlock := false
	hasReview := false
	for word := range hitMap {
		hits = append(hits, word)
		switch {
		case f.blockWords[word]:
			hasBlock = true
		case f.reviewWords[word]:
			hasReview = true
		default:
			hasBlock = true // 未显式分级 → 默认 block
		}
	}
	f.mu.RUnlock()

	result.HitWords = hits

	// --- Step 2: LevelBlock → 直接拦截，跳过 LLM ---
	if hasBlock {
		result.Level = LevelBlock
		result.Text = f.dfa.Replace(text, '*')
		return result
	}

	// --- Step 3: LevelReview → LLM 复判 ---
	if hasReview {
		result.Level = LevelReview
		result.Text = f.dfa.Replace(text, '*')

		if f.ollama != nil {
			llmLevel, err := f.llmReview(text, hits)
			if err != nil {
				// LLM 失败：保守保持 Review，记录日志
				log.Printf("[sensitive] LLM 复判失败，保守保持 LevelReview: %v", err)
			} else {
				result.Level = llmLevel
				if llmLevel == LevelClean {
					result.Text = text // 安全内容恢复原文
					result.HitWords = []string{}
				}
			}
		}
		return result
	}

	return result
}
