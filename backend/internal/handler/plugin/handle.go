package plugin

import (
	"tiny-forum/internal/service/plugin"
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
// @Tags plugin
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
	// 此处你需要根据实际 bo.PluginQueryBO 构造分页请求
	// result, err := h.svc.ListPlugins(c.Request.Context(), queryBO)
	// 简单返回示例
	response.Success(c, gin.H{"message": "TODO"})
}

// UploadPlugin 上传并安装插件
// @Summary 安装插件
// @Description 上传一个 ZIP 压缩包，包含 manifest.json 和插件资源，进行插件安装
// @Tags plugin
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
	meta, err := h.svc.Install(c.Request.Context(), file, userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, meta)
}

// ListMyPlugins 获取当前用户已安装的插件
// @Summary 当前用户安装的插件
// @Description 获取当前登录用户已安装（上传）的插件列表，通常是通过 author_id 查询
// @Tags plugin
// @Accept json
// @Produce json
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /users/me/plugins [get]
func (h *Handler) ListMyPlugins(c *gin.Context) {
	userID := c.GetInt64("user_id")
	plugins, err := h.svc.ListUserPlugins(c.Request.Context(), userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, plugins)
}