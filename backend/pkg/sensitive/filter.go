package sensitive

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Level 敏感词命中等级
type Level int

const (
	LevelClean    Level = 0 // [安全] 无敏感词：完全放行，不触发任何处理。例："你好"、"今天天气不错"
	LevelReplace  Level = 1 // 【攻击】替换敏感词：将命中词替换为 * 或自定义字符。例："傻逼" → "**"，"他妈" → "**"（轻度脏话或情绪词，替换后不影响阅读，但消除攻击性。）
	LevelReview   Level = 2 // 【风险】人工审核：帖子进入审核队列，通过后才公开。例："代考"、"出售卡号"、"兼职刷单"（涉及交易、服务或中等风险内容，需人工判断是否合规。）
	LevelShadowed Level = 3 // 【防对抗】遮蔽词：发布者本人可见，其他用户不可见。例："加微信"、"点击链接"、"私聊我"（高频广告词或试探性违规，让用户以为发布成功但实际不扩散，减少对抗。）
	LevelBlock    Level = 4 // 【禁止】直接拦截：拒绝发布并提示违规。例："法轮功"、"六四"、"贩卖毒品"、"儿童色情"（法律法规明确禁止的内容，必须彻底拦截。）
)

// CheckResult 检测结果
type CheckResult struct {
	Level    Level    // 最高命中等级
	HitWords []string // 命中的词列表
	Text     string   // 替换后的文本（* 号替换）
}

// DictLoadResult 目录加载结果摘要
type DictLoadResult struct {
	BlockFiles   []string // 识别为 block 级的文件
	ReviewFiles  []string // 识别为 review 级的文件
	ShadowFiles  []string // 识别为 shadow 级的文件
	ReplaceFiles []string // 识别为 replace 级的文件
	Errors       []string // 加载失败的文件及原因
	TotalWords   int      // 成功加载的词条总数
}

// Filter 敏感词过滤器接口
type Filter interface {
	// DFA 基础能力
	IsSensitive(text string) bool
	FindOne(text string) string
	FindAll(text string) []string
	FindAllCount(text string) map[string]int
	Replace(text string, replaceChar rune) string
	Remove(text string) string
	AddWord(words ...string) error
	DelWord(words ...string) error

	// 分级能力
	Check(text string) CheckResult
	AddLevelWords(level Level, words []string)

	// 加载词库
	// LoadDictFile 加载单个词库文件（每行一词，# 开头为注释）
	// level 指定该文件中所有词的拦截等级
	LoadDictFile(path string, level Level) error

	// LoadDictContent 从字符串内容加载词库
	// level 指定该批词的拦截等级
	LoadDictContent(content string, level Level) error

	// LoadDictDir 扫描目录下所有 .txt 文件，按文件名前缀自动识别等级：
	//   block_*.txt  → LevelBlock
	//   review_*.txt → LevelReview
	//   其余 *.txt   → LevelBlock（默认兜底）
	LoadDictDir(dir string) (DictLoadResult, error)
}

// filter 分级过滤器，包装 DFAFilter 并维护 block/review 词集
type filter struct {
	dfa         *DFAFilter
	blockWords  map[string]bool
	reviewWords map[string]bool
	mu          sync.RWMutex
}

// NewFilter 创建新的过滤器
func NewFilter() Filter {
	return &filter{
		dfa:         NewDFAFilter(),
		blockWords:  make(map[string]bool),
		reviewWords: make(map[string]bool),
	}
}

var (
	globalFilter Filter
	globalOnce   sync.Once
)

// GetGlobalFilter 获取全局单例过滤器
func GetGlobalFilter() Filter {
	globalOnce.Do(func() {
		globalFilter = NewFilter()
	})
	return globalFilter
}

// ---- DFA 基础能力转发 ----

func (f *filter) IsSensitive(text string) bool            { return f.dfa.IsSensitive(text) }
func (f *filter) FindOne(text string) string              { return f.dfa.FindOne(text) }
func (f *filter) FindAll(text string) []string            { return f.dfa.FindAll(text) }
func (f *filter) FindAllCount(text string) map[string]int { return f.dfa.FindAllCount(text) }
func (f *filter) Replace(text string, r rune) string      { return f.dfa.Replace(text, r) }
func (f *filter) Remove(text string) string               { return f.dfa.Remove(text) }

// AddWord 添加词，默认 block 级（除非已在 review map 中）
func (f *filter) AddWord(words ...string) error {
	if err := f.dfa.AddWord(words...); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, w := range words {
		if !f.reviewWords[w] {
			f.blockWords[w] = true
		}
	}
	return nil
}

// DelWord 删除词
func (f *filter) DelWord(words ...string) error {
	if err := f.dfa.DelWord(words...); err != nil {
		return err
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, w := range words {
		delete(f.blockWords, w)
		delete(f.reviewWords, w)
	}
	return nil
}

// AddLevelWords 添加指定等级的词汇
func (f *filter) AddLevelWords(level Level, words []string) {
	_ = f.dfa.AddWord(words...)
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, w := range words {
		switch level {
		case LevelBlock:
			f.blockWords[w] = true
			delete(f.reviewWords, w)
		case LevelReview:
			f.reviewWords[w] = true
			delete(f.blockWords, w)
		}
	}
}

// ---- 词库加载 ----

// readLines 读取文件，返回有效词条（过滤空行和 # 注释行）
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words = append(words, line)
	}
	return words, scanner.Err()
}

// LoadDictFile 加载单个词库文件，level 指定该文件的拦截等级
func (f *filter) LoadDictFile(path string, level Level) error {
	words, err := readLines(path)
	if err != nil {
		return fmt.Errorf("sensitive: load %s: %w", path, err)
	}
	f.AddLevelWords(level, words)
	return nil
}

// LoadDictContent 从字符串内容加载词库，level 指定等级
func (f *filter) LoadDictContent(content string, level Level) error {
	var words []string
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		words = append(words, line)
	}
	f.AddLevelWords(level, words)
	return nil
}

// fileLevel 根据文件名推断 Level
//
//	review_*.txt → LevelReview
//	其余 *.txt   → LevelBlock（含 block_*.txt）
func fileLevel(name string) Level {
	if strings.HasPrefix(strings.ToLower(name), "review_") {
		return LevelReview
	}
	return LevelBlock
}

// MARK: LoadDict
// LoadDictDir 扫描目录，按文件名前缀自动识别等级并加载所有 .txt 词库
func (f *filter) LoadDictDir(dir string) (DictLoadResult, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return DictLoadResult{}, fmt.Errorf("sensitive: open dir %s: %w", dir, err)
	}

	var result DictLoadResult

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".txt") {
			continue
		}

		path := filepath.Join(dir, name)
		level := fileLevel(name)

		words, err := readLines(path)
		if err != nil {
			result.Errors = append(result.Errors,
				fmt.Sprintf("%s: %v", name, err))
			continue
		}

		f.AddLevelWords(level, words)
		result.TotalWords += len(words)

		if level == LevelReview {
			result.ReviewFiles = append(result.ReviewFiles, name)
		} else {
			result.BlockFiles = append(result.BlockFiles, name)
		}
	}

	return result, nil
}

// ---- Check ----

// MARK: Check
// Check 检测文本，返回分级结果
func (f *filter) Check(text string) CheckResult {
	result := CheckResult{
		Level:    LevelClean,
		HitWords: []string{},
		Text:     text,
	}

	hitMap := f.dfa.FindAllCount(text)
	if len(hitMap) == 0 {
		return result
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

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
			// 直接通过 DFA AddWord 加入但未显式分级 → 默认 block
			hasBlock = true
		}
	}

	result.HitWords = hits
	if hasBlock {
		result.Level = LevelBlock
	} else if hasReview {
		result.Level = LevelReview
	}
	result.Text = f.dfa.Replace(text, '*')
	return result
}
