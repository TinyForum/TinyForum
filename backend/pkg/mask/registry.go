package mask

import (
	"strconv"
	"strings"
	"sync"
)

var (
	strategies = make(map[string]MaskFunc)
	mu         sync.RWMutex
)

func init() {
	// 注册内置策略
	for name, fn := range builtinStrategies {
		if name == "regex" {
			continue // regex 特殊处理
		}
		RegisterStrategy(name, fn)
	}
}

// RegisterStrategy 注册全局脱敏策略（线程安全）
func RegisterStrategy(name string, fn MaskFunc) {
	mu.Lock()
	defer mu.Unlock()
	strategies[name] = fn
}

// GetStrategy 获取策略函数
func GetStrategy(name string) (MaskFunc, bool) {
	mu.RLock()
	defer mu.RUnlock()
	fn, ok := strategies[name]
	return fn, ok
}

// parseTag 解析 mask 标签，返回策略名和参数 map
// 格式示例: "name" 或 "address,keep=6" 或 "regex,pattern=(\\d{3})\\d{4}(\\d{4}),replace=$1****$2"
func parseTag(tag string) (string, map[string]string) {
	if tag == "" {
		return "", nil
	}
	parts := strings.Split(tag, ",")
	name := parts[0]
	params := make(map[string]string)
	for i := 1; i < len(parts); i++ {
		kv := strings.SplitN(parts[i], "=", 2)
		if len(kv) == 2 {
			params[kv[0]] = kv[1]
		} else {
			params[kv[0]] = ""
		}
	}
	return name, params
}

// applyStrategy 根据策略名和参数对字符串进行脱敏
func applyStrategy(name string, s string, params map[string]string) string {
	if s == "" {
		return s
	}
	switch name {
	case "regex":
		pattern := params["pattern"]
		replace := params["replace"]
		if pattern == "" || replace == "" {
			return s
		}
		return applyRegex(s, pattern, replace)
	case "address":
		keepStr := params["keep"]
		keep := 6 // 默认保留6个字符
		if keepStr != "" {
			if k, err := strconv.Atoi(keepStr); err == nil && k > 0 {
				keep = k
			}
		}
		return maskAddressKeepN(s, keep)
	default:
		fn, ok := GetStrategy(name)
		if !ok {
			// 未注册的策略降级为 full
			return maskFull(s)
		}
		return fn(s)
	}
}
