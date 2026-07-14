// internal/handler/config/config_handler.go
package config

import (
	"strconv"
	"strings"
	"tiny-forum/internal/service/config"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"
	"tiny-forum/pkg/utils"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

type ConfigHandler struct {
	configSvc *config.ConfigService
}

func NewConfigHandler(configSvc *config.ConfigService) *ConfigHandler {
	return &ConfigHandler{
		configSvc: configSvc,
	}
}

// ListConfigs 获取配置列表
// @Summary 获取所有配置文件列表
// @Tags 配置管理
// @Produce json
// @Success 200 {object} response.Response{data=[]config.ConfigFileInfo}
// @Router /api/admin/config/list [get]
func (h *ConfigHandler) ListConfigs(c *gin.Context) {
	files, err := h.configSvc.ListConfigFiles()
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, files)
}

// GetConfig 获取单个配置内容
// @Summary 获取配置文件内容
// @Tags 配置管理
// @Produce json
// @Param file path string true "配置文件名 (如 basic.yml)"
// @Param format query string false "返回格式: yaml(默认) 或 kv"
// @Success 200 {object} response.Response{data=object{content=string}}
// @Router /api/admin/config/{file} [get]
func (h *ConfigHandler) GetConfig(c *gin.Context) {
	fileName := c.Param("file")
	if fileName == "" {
		response.HandleError(c, apperrors.ErrInternalError.WithMessagef("文件名为空"))
		return
	}

	// 规范化文件名
	fileName = h.normalizeFileName(fileName)

	content, err := h.configSvc.GetConfigContent(fileName)
	if err != nil {
		response.HandleError(c, apperrors.ErrInternalError.WithMessagef("获取配置内容失败: %s", err.Error()))
		return
	}

	// 检查是否请求键值对格式
	format := c.DefaultQuery("format", "yaml")

	if format == "kv" {
		// 解析 YAML 为键值对
		kvData, err := h.parseYAMLToKV(content)
		if err != nil {
			response.HandleError(c, apperrors.ErrInternalError.WithMessagef("解析配置为键值对失败: %s", err.Error()))
			return
		}

		response.Success(c, gin.H{
			"file":   fileName,
			"format": "kv",
			"config": kvData,
		})
		return
	}

	// 默认返回 YAML 格式
	response.Success(c, gin.H{
		"file":    fileName,
		"format":  "yaml",
		"content": content,
	})
}

// GetConfigKV 获取配置的键值对格式
// @Summary 获取配置文件的键值对格式
// @Tags 配置管理
// @Produce json
// @Param file path string true "配置文件名 (如 basic.yml)"
// @Success 200 {object} response.Response{data=object{config=map[string]string}}
// @Router /api/admin/config/{file}/kv [get]
func (h *ConfigHandler) GetConfigKV(c *gin.Context) {
	fileName := c.Param("file")
	if fileName == "" {
		response.HandleError(c, apperrors.ErrInternalError.WithMessagef("文件名为空"))
		return
	}

	// 规范化文件名
	fileName = h.normalizeFileName(fileName)

	content, err := h.configSvc.GetConfigContent(fileName)
	if err != nil {
		response.HandleError(c, apperrors.ErrInternalError.WithMessagef("获取配置内容失败: %s", err.Error()))
		return
	}

	// 解析 YAML 为键值对
	kvData, err := h.parseYAMLToKV(content)
	if err != nil {
		response.HandleError(c, apperrors.ErrInternalError.WithMessagef("解析配置为键值对失败: %s", err.Error()))
		return
	}

	response.Success(c, gin.H{
		"file":   fileName,
		"format": "kv",
		"config": kvData,
	})
}

// UpdateConfig 更新配置
// @Summary 更新配置文件
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param file path string true "配置文件名 (如 basic.yml)"
// @Param request body object{content=string} true "配置内容"
// @Success 200 {object} response.Response
// @Router /api/admin/config/{file} [put]
func (h *ConfigHandler) UpdateConfig(c *gin.Context) {
	fileName := c.Param("file")
	if fileName == "" {
		response.HandleError(c, apperrors.ErrInternalError.WithMessagef("文件名为空"))
		return
	}

	// 规范化文件名
	fileName = h.normalizeFileName(fileName)

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, apperrors.ErrInvalidRequest.WithMessagef("请求参数错误: %s", err.Error()))
		return
	}

	// 获取操作人（从 JWT 中获取）
	operator := c.GetString("username")
	if operator == "" {
		operator = "admin"
	}

	if err := h.configSvc.UpdateConfigContent(fileName, req.Content, operator); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "配置更新成功，已自动重载",
		"file":    fileName,
	})
}

// UpdateConfigKV 通过键值对更新配置
// @Summary 通过键值对更新配置文件
// @Tags 配置管理
// @Accept json
// @Produce json
// @Param file path string true "配置文件名 (如 basic.yml)"
// @Param request body object{config=map[string]string} true "键值对配置"
// @Success 200 {object} response.Response
// @Router /api/admin/config/{file}/kv [put]
func (h *ConfigHandler) UpdateConfigKV(c *gin.Context) {
	fileName := c.Param("file")
	if fileName == "" {
		response.HandleError(c, apperrors.ErrInternalError.WithMessagef("文件名为空"))
		return
	}

	// 规范化文件名
	fileName = h.normalizeFileName(fileName)

	// 先检查请求体
	var req struct {
		Config map[string]string `json:"config"`
	}

	// 使用 ShouldBindJSON 并捕获错误
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, apperrors.ErrInvalidRequest.WithMessagef("JSON解析失败: %s", err.Error()))
		return
	}

	// 验证 config 是否为空
	if req.Config == nil || len(req.Config) == 0 {
		response.HandleError(c, apperrors.ErrInvalidRequest.WithMessagef("config 不能为空"))
		return
	}

	// 获取当前配置内容
	currentContent, err := h.configSvc.GetConfigContent(fileName)
	if err != nil {
		response.HandleError(c, apperrors.ErrInternalError.WithMessagef("获取当前配置失败: %s", err.Error()))
		return
	}

	// 将键值对转换为 YAML
	newContent, err := h.mergeKVToYAML(currentContent, req.Config)
	if err != nil {
		response.HandleError(c, apperrors.ErrInternalError.WithMessagef("转换配置失败: %s", err.Error()))
		return
	}

	// 获取操作人
	operator := c.GetString("username")
	if operator == "" {
		operator = "admin"
	}

	// 更新配置
	if err := h.configSvc.UpdateConfigContent(fileName, newContent, operator); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "配置更新成功，已自动重载",
		"file":    fileName,
		"updated": len(req.Config),
	})
}

// ReloadConfig 手动重载配置
// @Summary 手动重载所有配置
// @Tags 配置管理
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/admin/config/reload [post]
func (h *ConfigHandler) ReloadConfig(c *gin.Context) {
	operator := c.GetString("username")
	if operator == "" {
		operator = "admin"
	}

	if err := h.configSvc.ReloadConfig(operator); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "配置重载成功",
	})
}

// GetHistory 获取配置变更历史
// @Summary 获取配置变更历史
// @Tags 配置管理
// @Produce json
// @Param file query string false "配置文件名 (不传则查所有)"
// @Param limit query int false "数量限制" default(50)
// @Success 200 {object} response.Response{data=[]config.ConfigHistory}
// @Router /api/admin/config/history [get]
func (h *ConfigHandler) GetHistory(c *gin.Context) {
	fileName := c.Query("file")
	if fileName != "" && fileName != "all" {
		fileName = h.normalizeFileName(fileName)
	}

	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit <= 0 {
		limit = 50
	}

	histories, err := h.configSvc.GetConfigHistory(fileName, limit)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, histories)
}

// normalizeFileName 规范化文件名
func (h *ConfigHandler) normalizeFileName(fileName string) string {
	// 移除路径分隔符
	fileName = strings.TrimPrefix(fileName, "/")
	fileName = strings.TrimPrefix(fileName, "\\")

	// 如果已经有扩展名，直接返回
	if strings.HasSuffix(fileName, ".yml") || strings.HasSuffix(fileName, ".yaml") {
		return fileName
	}

	// 添加 .yml 扩展名
	return fileName + ".yml"
}

// parseYAMLToKV 解析 YAML 为键值对
func (h *ConfigHandler) parseYAMLToKV(yamlContent string) (map[string]interface{}, error) {
	var data interface{}
	if err := yaml.Unmarshal([]byte(yamlContent), &data); err != nil {
		return nil, err
	}

	return utils.FlattenConfig("", data), nil
}

// mergeKVToYAML 将键值对合并到 YAML
func (h *ConfigHandler) mergeKVToYAML(currentYAML string, kvData map[string]string) (string, error) {
	// 解析当前 YAML
	var currentData map[string]interface{}
	if err := yaml.Unmarshal([]byte(currentYAML), &currentData); err != nil {
		return "", err
	}

	if currentData == nil {
		currentData = make(map[string]interface{})
	}

	// 应用键值对更新
	for key, value := range kvData {
		if err := utils.SetNestedValue(currentData, key, value); err != nil {
			return "", err
		}
	}

	// 序列化为 YAML
	result, err := yaml.Marshal(currentData)
	if err != nil {
		return "", err
	}

	return string(result), nil
}
