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

// @Summary 获取待审核帖子列表
// @Description 管理员获取需要审核的帖子列表（状态为待审核）
// @Tags 帖子管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "搜索关键词"
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Post}} "获取成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/posts/pending [get]
func (h *PostHandler) AdminGetModerationRequire(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	opts := dto.PostListOptions{
		// Status:  model.PostStatusPending,
		ModerationStatus: model.ModerationStatusPending,
		Keyword:          keyword,
	}

	posts, total, err := h.postSvc.AdminList(page, pageSize, opts)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}

// @Summary 审核通过帖子
// @Description 管理员审核通过指定帖子
// @Tags 帖子管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Success 200 {object} response.Response{data=object{message=string,post_id=int}} "审核通过成功"
// @Failure 400 {object} response.Response "无效的帖子ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "帖子不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/audit/tasks/{id}/approve [put]
func (h *PostHandler) AdminApprovePost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	if err := h.postSvc.AdminSetReviewPost(uint(postID), model.ModerationStatusApproved); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "帖子不存在")
		} else {
			response.InternalError(c, "审核通过失败："+err.Error())
		}
		return
	}

	response.Success(c, gin.H{
		"message": "审核通过成功",
		"post_id": postID,
	})
}

// @Summary 审核拒绝帖子
// @Description 管理员审核拒绝指定帖子
// @Tags 帖子管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "帖子ID"
// @Param body body object false "拒绝原因" example({"reason": "内容不合规"})
// @Success 200 {object} response.Response{data=object{message=string,post_id=int}} "审核拒绝成功"
// @Failure 400 {object} response.Response "无效的帖子ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "帖子不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /admin/audit/tasks/{id}/reject [put]
func (h *PostHandler) AdminRejectPost(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	// 可选：从 body 获取拒绝原因
	var req struct {
		Reason string `json:"reason" binding:"max=500"`
	}
	_ = c.ShouldBindJSON(&req)

	if err := h.postSvc.AdminSetReviewPost(uint(postID), model.ModerationStatusRejected); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "帖子不存在")
		} else {
			response.InternalError(c, "审核拒绝失败："+err.Error())
		}
		return
	}

	// 这里可以保存拒绝原因到日志或扩展字段
	response.Success(c, gin.H{
		"message": "审核拒绝成功",
		"post_id": postID,
	})
}
