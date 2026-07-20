package sensitive

import "unicode"

func newTrieNode() *trieNode {
	return &trieNode{children: make(map[rune]*trieNode)}
}

// normalize 规范化字符：去除不可见字符、统一大小写、全角转半角
func normalize(r rune) (rune, bool) {
	// 跳过不可见 / 零宽 / 控制字符
	if r == '\u200b' || r == '\u200c' || r == '\u200d' ||
		r == '\ufeff' || r == '\u00ad' || unicode.Is(unicode.Cf, r) {
		return 0, false
	}
	// 全角字母数字 → 半角
	if r >= 0xFF01 && r <= 0xFF5E {
		r = r - 0xFEE0
	}
	// 统一小写
	r = unicode.ToLower(r)
	return r, true
}

// normalizeRunes 将文本规范化为 rune 切片，同时保留原始位置映射
// 返回：规范化后的 rune 列表，以及每个规范化 rune 对应原始 rune 的索引
func normalizeRunes(text string) ([]rune, []int) {
	runes := []rune(text)
	var norm []rune
	var idx []int
	for i, r := range runes {
		nr, ok := normalize(r)
		if !ok {
			continue
		}
		norm = append(norm, nr)
		idx = append(idx, i)
	}
	return norm, idx
}

func (f *dfaFilter) delWord(node *trieNode, runes []rune, depth int) bool {
	if depth == len(runes) {
		if !node.isEnd {
			return false
		}
		node.isEnd = false
		return len(node.children) == 0
	}
	r := runes[depth]
	nr, ok := normalize(r)
	if !ok {
		return false
	}
	child, exists := node.children[nr]
	if !exists {
		return false
	}
	shouldDelete := f.delWord(child, runes, depth+1)
	if shouldDelete {
		delete(node.children, nr)
		return !node.isEnd && len(node.children) == 0
	}
	return false
}

// match 从规范化 rune 列表的 start 位置开始尝试最长匹配
// 返回：匹配到的原始文本（通过原始 rune 列表重建），以及匹配结束的规范化索引（不含），-1 表示无匹配
func (f *dfaFilter) match(norm []rune, origRunes []rune, origIdx []int, start int) (string, int) {
	node := f.root
	lastEnd := -1
	lastEndOrig := -1

	for i := start; i < len(norm); i++ {
		child, ok := node.children[norm[i]]
		if !ok {
			break
		}
		node = child
		if node.isEnd {
			lastEnd = i + 1
			if i < len(origIdx) {
				lastEndOrig = origIdx[i] + 1
			}
		}
	}

	if lastEnd == -1 {
		return "", -1
	}

	// 通过原始 rune 重建命中词
	origStart := origIdx[start]
	word := string(origRunes[origStart:lastEndOrig])
	return word, lastEnd
}
