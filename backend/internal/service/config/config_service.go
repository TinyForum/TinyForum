package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tiny-forum/internal/infra/config"

	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

// ConfigService 配置管理服务
type ConfigService struct {
	configDir   string
	dynCfg      *config.DynamicConfig
	db          *gorm.DB
	configFiles []string
}

// ConfigFileInfo 配置文件信息
type ConfigFileInfo struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Size     int64     `json:"size"`
	ModTime  time.Time `json:"mod_time"`
	Format   string    `json:"format"`
	Editable bool      `json:"editable"`
}

// ConfigHistory 配置变更历史
type ConfigHistory struct {
	ID        uint   `gorm:"primarykey"`
	FileName  string `gorm:"index;not null"`
	Content   string `gorm:"type:text"`
	Operator  string `gorm:"index"`
	Operation string `gorm:"not null"` // create, update, delete, reload
	CreatedAt time.Time
}

func NewConfigService(configDir string, dynCfg *config.DynamicConfig, db *gorm.DB) *ConfigService {
	// 自动迁移历史表
	db.AutoMigrate(&ConfigHistory{})

	return &ConfigService{
		configDir:   configDir,
		dynCfg:      dynCfg,
		db:          db,
		configFiles: []string{"basic", "private", "risk_control", "postgres", "redis"},
	}
}

// ListConfigFiles 列出所有配置文件
func (s *ConfigService) ListConfigFiles() ([]ConfigFileInfo, error) {
	var files []ConfigFileInfo

	for _, name := range s.configFiles {
		path := filepath.Join(s.configDir, name+".yml")
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, err
		}

		files = append(files, ConfigFileInfo{
			Name:     name + ".yml",
			Path:     path,
			Size:     info.Size(),
			ModTime:  info.ModTime(),
			Format:   "yaml",
			Editable: true,
		})
	}

	return files, nil
}

// GetConfigContent 获取配置文件内容
func (s *ConfigService) GetConfigContent(fileName string) (string, error) {
	// 确保文件名安全
	if !s.isValidFileName(fileName) {
		return "", fmt.Errorf("invalid file name: %s", fileName)
	}

	path := filepath.Join(s.configDir, fileName)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// UpdateConfigContent 更新配置文件内容
func (s *ConfigService) UpdateConfigContent(fileName, content, operator string) error {
	// 1. 验证文件名
	if !s.isValidFileName(fileName) {
		return fmt.Errorf("invalid file name: %s", fileName)
	}

	path := filepath.Join(s.configDir, fileName)

	// 2. 验证 YAML 格式
	if err := s.validateYAML(content); err != nil {
		return fmt.Errorf("invalid YAML format: %w", err)
	}

	// 3. 读取旧内容（用于历史记录）
	// oldContent, _ := ioutil.ReadFile(path)

	// 4. 写入新内容
	if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
		return err
	}

	// 5. 记录变更历史
	history := ConfigHistory{
		FileName:  fileName,
		Content:   content,
		Operator:  operator,
		Operation: "update",
	}
	if err := s.db.Create(&history).Error; err != nil {
		// 记录失败不影响主流程
	}

	// 6. 触发热重载
	if err := s.dynCfg.Reload(); err != nil {
		return fmt.Errorf("config reload failed: %w", err)
	}

	return nil
}

// ReloadConfig 手动重载配置
func (s *ConfigService) ReloadConfig(operator string) error {
	// 记录重载操作
	history := ConfigHistory{
		FileName:  "all",
		Operator:  operator,
		Operation: "reload",
	}
	s.db.Create(&history)

	return s.dynCfg.Reload()
}

// GetConfigHistory 获取配置变更历史
func (s *ConfigService) GetConfigHistory(fileName string, limit int) ([]ConfigHistory, error) {
	var histories []ConfigHistory
	query := s.db.Order("created_at DESC")

	if fileName != "" && fileName != "all" {
		query = query.Where("file_name = ?", fileName)
	}

	if limit <= 0 {
		limit = 50
	}
	query = query.Limit(limit)

	if err := query.Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}

// validateYAML 验证 YAML 格式
func (s *ConfigService) validateYAML(content string) error {
	var data interface{}
	decoder := yaml.NewDecoder(bytes.NewReader([]byte(content)))
	return decoder.Decode(&data)
}

// isValidFileName 验证文件名是否安全
func (s *ConfigService) isValidFileName(fileName string) bool {
	// 只允许 .yml 文件
	if !strings.HasSuffix(fileName, ".yml") {
		return false
	}

	// 检查是否在允许列表中
	name := strings.TrimSuffix(fileName, ".yml")
	for _, allowed := range s.configFiles {
		if name == allowed {
			return true
		}
	}
	return false
}
