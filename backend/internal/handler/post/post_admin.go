package post

import (
	"errors"
	"strconv"

	"tiny-forum/internal/dto"
	"tiny-forum/internal/model"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AdminList 管理员获取帖子列表
// @Summary 管理员获取帖子列表
// @Description 管理员分页获取所有帖子列表，支持关键词搜索
// @Tags 管理接口
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Post}} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/posts [get]
func (h *PostHandler) AdminList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	opts := dto.PostListOptions{
		// Status:  model.PostStatusPending, // 关键：只查待审核
		Keyword: keyword,
		// 可按需添加其他筛选，如作者、标签等
	}
	posts, total, err := h.postSvc.AdminList(page, pageSize, opts)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}

// AdminTogglePin 管理员切换帖子置顶状态
// @Summary 切换帖子置顶状态
// @Description 管理员切换指定帖子的置顶状态
// @Tags 管理接口
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Success 200 {object} response.Response{data=object} "操作成功"
// @Failure 400 {object} response.Response "无效的帖子ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/posts/{id}/pin [put]
func (h *PostHandler) AdminTogglePin(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	if err := h.postSvc.TogglePin(uint(postID)); err != nil {
		if errors.Is(err, apperrors.ErrPostNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "操作成功"})
}

// AdminGetPending 获取待审核帖子列表（管理员）
// @Summary      获取待审核帖子列表
// @Description  分页获取所有状态为 pending 的帖子，支持按标题关键词搜索。需要管理员权限。
// @Tags         管理接口 - 帖子管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        page       query     int     false  "页码"                default(1)      minimum(1)
// @Param        page_size  query     int     false  "每页数量"            default(20)     minimum(1)  maximum(100)
// @Param        keyword    query     string  false  "搜索关键词（匹配标题）"
// @Success      200  {object}  response.Response{data=response.PageData{list=[]model.Post}}  "获取成功"
// @Failure      401  {object}  response.Response  "未授权（缺少或无效的 Token）"
// @Failure      403  {object}  response.Response  "无权限（非管理员）"
// @Failure      500  {object}  response.Response  "服务器内部错误"
// @Router       /admin/posts/pending [get]
func (h *PostHandler) AdminGetPending(c *gin.Context) {
	// role := c.Get("user_role")
	// if (role != "admin" || role != "super_admin") {
	// 	response.Forbidden(c, "未授权")
	// 	return
	// }
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	opts := dto.PostListOptions{
		// Status:  model.PostStatusPending,
		Risk:    model.AuditStatusPending,
		Keyword: keyword,
	}

	posts, total, err := h.postSvc.AdminList(page, pageSize, opts)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}

// AdminReviewPost 管理员审核帖子
// @Summary      管理员审核帖子
// @Description  管理员对指定帖子进行审核，可设置为 approved（通过）、rejected（拒绝）或 pending（待审）
// @Tags         管理接口 - 帖子管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id     path      int     true  "帖子ID"
// @Param        request body      object  true  "审核状态"  SchemaExample({"status":"approved"})
// @Success      200    {object}  response.Response{data=object}  "审核成功"
// @Failure      400    {object}  response.Response  "请求参数错误"
// @Failure      401    {object}  response.Response  "未授权"
// @Failure      403    {object}  response.Response  "无权限（非管理员）"
// @Failure      404    {object}  response.Response  "帖子不存在"
// @Failure      500    {object}  response.Response  "服务器内部错误"
// @Router       /admin/posts/{id}/review [put]
func (h *PostHandler) AdminReviewPost(c *gin.Context) {
	// 1. 权限检查（已有）
	role, _ := c.Get("user_role")
	isAdmin := role == "admin"
	if !isAdmin {
		response.Forbidden(c, "无权限，需要管理员角色")
		return
	}

	// 2. 解析帖子ID
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	// 3. 解析请求体（审核状态）
	var req struct {
		Status string `json:"status" binding:"required,oneof=approved rejected pending"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "审核状态无效，必须是 approved、rejected 或 pending")
		return
	}

	// 4. 调用 Service 执行审核（需要你的 Service 层提供 UpdatePostStatus 方法）
	// 示例：h.postSvc.UpdatePostStatus(uint(postID), req.Status)
	if err := h.postSvc.AdminSetReviewPost(uint(postID), model.ModerationStatus(req.Status)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "帖子不存在")
		} else {
			response.InternalError(c, "审核失败："+err.Error())
		}
		return
	}

	// 5. 返回成功
	response.Success(c, gin.H{
		"message": "审核操作成功",
		"post_id": postID,
		"status":  req.Status,
	})
}
