package sensitive

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ---- 词库加载 ----

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

func (f *filter) LoadDictFile(path string, level Level) error {
	words, err := readLines(path)
	if err != nil {
		return fmt.Errorf("sensitive: load %s: %w", path, err)
	}
	f.AddLevelWords(level, words)
	return nil
}

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

func fileLevel(name string) Level {
	lower := strings.ToLower(name)
	if strings.HasPrefix(lower, "review_") {
		return LevelReview
	}
	return LevelBlock
}

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
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", name, err))
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
