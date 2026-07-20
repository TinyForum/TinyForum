// internal/infra/config/dynamic.go
package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
)

// DynamicConfig 动态配置管理器
type DynamicConfig struct {
	// 当前配置（原子操作，保证并发安全）
	current atomic.Value // 存储 *Config

	// 配置目录
	configDir string

	// 配置文件列表
	configFiles []string

	// 文件监听器
	watcher *fsnotify.Watcher

	// 回调函数列表
	callbacks []ConfigChangeCallback
	mu        sync.RWMutex

	// 上下文控制
	ctx    context.Context
	cancel context.CancelFunc

	// 防抖
	debounceTimer *time.Timer
	debounceMu    sync.Mutex

	// 初始化状态
	initialized bool
}

// ConfigChangeCallback 配置变更回调函数
// 参数：变更的文件名，旧配置，新配置
type ConfigChangeCallback func(fileName string, oldConfig, newConfig *Config)

// NewDynamicConfig 创建动态配置管理器
func NewDynamicConfig(configDir string) (*DynamicConfig, error) {
	// 检查配置目录是否存在
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("config directory not found: %s", configDir)
	}

	// 创建文件监听器
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	dc := &DynamicConfig{
		configDir:   configDir,
		configFiles: []string{"basic", "private", "risk_control", "postgres", "redis", "ai"},
		watcher:     watcher,
		ctx:         ctx,
		cancel:      cancel,
		callbacks:   make([]ConfigChangeCallback, 0),
	}

	// 1. 首次加载配置
	if err := dc.loadAllConfigs(); err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to load initial config: %w", err)
	}

	// 2. 启动文件监听
	if err := dc.startWatching(); err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to start watching: %w", err)
	}

	dc.initialized = true
	log.Printf("[DynamicConfig] Initialized successfully, watching directory: %s", configDir)

	return dc, nil
}

// loadAllConfigs 加载所有配置文件
func (dc *DynamicConfig) loadAllConfigs() error {
	// 使用你现有的 Load 函数加载配置
	cfg, err := Load(dc.configDir)
	if err != nil {
		return err
	}

	// 存储到 atomic.Value
	dc.current.Store(cfg)
	return nil
}

// startWatching 启动文件监听
func (dc *DynamicConfig) startWatching() error {
	// 监听配置目录
	if err := dc.watcher.Add(dc.configDir); err != nil {
		return err
	}

	// 启动监听协程
	go dc.watchLoop()

	return nil
}

// watchLoop 监听循环
func (dc *DynamicConfig) watchLoop() {
	for {
		select {
		case <-dc.ctx.Done():
			return

		case event, ok := <-dc.watcher.Events:
			if !ok {
				return
			}

			// 只关注写入、创建和重命名事件
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Rename) == 0 {
				continue
			}

			// 检查是否是配置文件
			fileName := filepath.Base(event.Name)
			if !dc.isConfigFile(fileName) {
				continue
			}

			log.Printf("[DynamicConfig] Detected change: %s", fileName)

			// 防抖处理
			dc.handleDebouncedReload(event.Name)

		case err, ok := <-dc.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("[DynamicConfig] Watcher error: %v", err)
		}
	}
}

// isConfigFile 检查是否是配置文件
func (dc *DynamicConfig) isConfigFile(fileName string) bool {
	ext := filepath.Ext(fileName)
	if ext != ".yml" && ext != ".yaml" {
		return false
	}

	// 去掉扩展名检查是否在配置列表中
	name := strings.TrimSuffix(fileName, ext)
	for _, cfgFile := range dc.configFiles {
		if cfgFile == name {
			return true
		}
	}
	return false
}

// handleDebouncedReload 防抖处理配置重载
func (dc *DynamicConfig) handleDebouncedReload(filePath string) {
	dc.debounceMu.Lock()
	defer dc.debounceMu.Unlock()

	// 重置定时器
	if dc.debounceTimer != nil {
		dc.debounceTimer.Stop()
	}

	// 延迟 100ms 执行，避免频繁触发
	dc.debounceTimer = time.AfterFunc(100*time.Millisecond, func() {
		dc.reloadConfig(filePath)
	})
}

// reloadConfig 重新加载配置
func (dc *DynamicConfig) reloadConfig(filePath string) {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	log.Printf("[DynamicConfig] Reloading config due to change: %s", filepath.Base(filePath))

	// 获取旧配置
	oldConfig := dc.Get()

	// 重新加载所有配置
	if err := dc.loadAllConfigs(); err != nil {
		log.Printf("[DynamicConfig] Failed to reload config: %v", err)
		return
	}

	// 获取新配置
	newConfig := dc.Get()

	// 触发回调
	fileName := filepath.Base(filePath)
	for _, callback := range dc.callbacks {
		// 使用 recover 防止回调 panic 导致程序崩溃
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[DynamicConfig] Callback panic recovered: %v", r)
				}
			}()
			callback(fileName, oldConfig, newConfig)
		}()
	}

	log.Printf("[DynamicConfig] Config reloaded successfully")
}

// Get 获取当前配置
func (dc *DynamicConfig) Get() *Config {
	if val := dc.current.Load(); val != nil {
		return val.(*Config)
	}
	return nil
}

// GetBasic 获取基础配置
func (dc *DynamicConfig) GetBasic() *ConfigBasic {
	cfg := dc.Get()
	if cfg == nil {
		return nil
	}
	return &cfg.Basic
}

// GetPrivate 获取私有配置
func (dc *DynamicConfig) GetPrivate() *ConfigPrivate {
	cfg := dc.Get()
	if cfg == nil {
		return nil
	}
	return &cfg.Private
}

// GetPostgres 获取PostgreSQL配置
func (dc *DynamicConfig) GetPostgres() *ConfigPostgres {
	cfg := dc.Get()
	if cfg == nil {
		return nil
	}
	return &cfg.Postgres
}

// GetRedis 获取Redis配置
func (dc *DynamicConfig) GetRedis() *ConfigRedis {
	cfg := dc.Get()
	if cfg == nil {
		return nil
	}
	return &cfg.Redis
}

// GetRiskControl 获取风控配置
func (dc *DynamicConfig) GetRiskControl() *ConfigRiskControl {
	cfg := dc.Get()
	if cfg == nil {
		return nil
	}
	return &cfg.RiskControl
}

func (dc *DynamicConfig) GetAIConfig() *ConfigAI {
	cfg := dc.Get()
	if cfg == nil {
		return nil
	}
	return &cfg.AI
}

// OnChange 注册配置变更回调
func (dc *DynamicConfig) OnChange(callback ConfigChangeCallback) {
	dc.mu.Lock()
	defer dc.mu.Unlock()
	dc.callbacks = append(dc.callbacks, callback)
}

// Reload 手动重新加载配置（可用于信号触发）
func (dc *DynamicConfig) Reload() error {
	dc.mu.Lock()
	defer dc.mu.Unlock()

	log.Printf("[DynamicConfig] Manual reload triggered")

	oldConfig := dc.Get()
	if err := dc.loadAllConfigs(); err != nil {
		return err
	}
	newConfig := dc.Get()

	// 触发回调
	for _, callback := range dc.callbacks {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("[DynamicConfig] Callback panic recovered: %v", r)
				}
			}()
			callback("manual_reload", oldConfig, newConfig)
		}()
	}

	log.Printf("[DynamicConfig] Manual reload completed")
	return nil
}

// Close 关闭配置管理器
func (dc *DynamicConfig) Close() error {
	dc.cancel()
	if dc.watcher != nil {
		return dc.watcher.Close()
	}
	return nil
}

// IsInitialized 检查是否已初始化
func (dc *DynamicConfig) IsInitialized() bool {
	return dc.initialized
}
