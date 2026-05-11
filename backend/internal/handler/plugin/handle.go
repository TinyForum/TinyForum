package plugin

import (
	"strconv"
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/service/plugin"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/logger"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc plugin.PluginService
}

func NewHandler(svc plugin.PluginService) *Handler {
	return &Handler{svc: svc}
}

// ListPlugins 获取系统插件列表
// @Summary 获取所有插件列表
// @Description 获取系统中所有已安装的插件（通常需要管理员权限，可根据业务调整）
// @Tags 插件管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页条数" default(20)
// @Param category query string false "分类过滤 (commerce, marketing, content...)"
// @Param type query string false "类型过滤 (ui, backend, lib, app, miniapp)"
// @Param enabled query bool false "是否启用"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /plugins [get]
func (h *Handler) ListPlugins(c *gin.Context) {
	var req request.ListPluginsRequest
	if err := req.Bind(c); err != nil {
		response.HandleError(c, err)
		return
	}
	logger.Infof("ListPlugins request: %+v", req)
	// userID := c.GetUint("user_id")

	// 获取当前登录用户ID（中间件注入，类型为 uint）

	// 构建分页查询对象 PageQuery[PluginQueryBO]
	pageQuery := &bo.PageQuery[bo.PluginQueryBO]{
		Page:     req.Page,
		PageSize: req.PageSize,
		SortBy:   req.SortBy,
		Order:    req.Order,
		Options: bo.PluginQueryBO{
			Name:     req.Keyword,
			Slug:     req.Keyword,
			AuthorID: req.AuthorID,
			Category: req.Category,
			Tags:     req.Tags,
			Type:     req.Type,
			Keyword:  req.Keyword,
			Status:   do.PluginStatus(req.Status),
			Version:  req.Version,
		},
	}

	// 调用 Service（参数类型匹配）
	pageResult, err := h.svc.ListPlugins(c.Request.Context(), pageQuery)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, pageResult)
}

// UploadPlugin 上传并安装插件
// @Summary 安装插件
// @Description 上传一个 ZIP 压缩包，包含 manifest.json 和插件资源，进行插件安装
// @Tags 插件管理
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "插件 ZIP 文件"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /plugins/upload [post]
func (h *Handler) UploadPlugin(c *gin.Context) {
	userID := c.GetUint("user_id")
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file is required")
		return
	}
	meta, err := h.svc.Create(c.Request.Context(), file, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, meta)
}

// ListMyPlugins 获取当前用户已安装的插件
// @Summary 当前用户安装的插件
// @Description 获取当前登录用户已安装（上传）的插件列表，通常是通过 author_id 查询
// @Tags 插件管理
// @Accept json
// @Produce json
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /plugins/user/me [get]
func (h *Handler) ListMyPlugins(c *gin.Context) {
	userID := c.GetUint("user_id")
	plugins, err := h.svc.ListUserPlugins(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, plugins)
}

// DeletePlugin 删除插件
// @Summary 删除插件
// @Description 删除插件
// @Tags 插件管理
// @Accept json
// @Produce json
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /pligins/delete [get]
func (h *Handler) DeletePlugin(c *gin.Context) {
	// 1. 获取当前用户ID
	userID := c.GetUint("user_id")

	// 2. 获取路径参数中的插件ID
	pluginIDStr := c.Param("id")
	pluginID, err := strconv.ParseUint(pluginIDStr, 10, 32)
	if err != nil {
		response.HandleError(c, apperrors.ErrValidation)
	}

	// 3. 调用服务层删除
	if err := h.svc.DeletePlugin(c.Request.Context(), uint(pluginID), userID); err != nil {
		response.HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "plugin deleted successfully"})
}
