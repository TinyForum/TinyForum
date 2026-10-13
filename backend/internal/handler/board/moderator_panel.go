package board

import (
	"strconv"

	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// ListModeratorPosts 版主面板：板块帖子列表
// @Summary 版主面板帖子列表
// @Description 分页获取指定板块下的帖子（不做已发布过滤），支持关键词与状态筛选
// @Tags 版主面板
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param keyword query string false "关键词"
// @Param status query string false "帖子状态：pending/published/deleted"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /moderator/boards/{id}/posts [get]
func (h *BoardHandler) ListModeratorPosts(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, apperrors.ErrInvalidRequest)
		return
	}
	var req request.ListModeratorPostsRequest
	if err := req.Bind(c); err != nil {
		response.HandleError(c, err)
		return
	}
	list, total, err := h.boardSvc.ListModeratorPosts(c.Request.Context(), uint(boardID), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.SuccessPage(c, list, total, req.Page, req.PageSize)
}

// ListModeratorReports 版主面板：板块举报列表
// @Summary 版主面板举报列表
// @Description 分页获取指定板块的举报（帖子或评论属于该板块），可按状态筛选
// @Tags 版主面板
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query string false "举报状态：pending/resolved/rejected"
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /moderator/boards/{id}/reports [get]
func (h *BoardHandler) ListModeratorReports(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, apperrors.ErrInvalidRequest)
		return
	}
	var req request.ListModeratorReportsRequest
	if err := req.Bind(c); err != nil {
		response.HandleError(c, err)
		return
	}
	list, total, err := h.boardSvc.ListModeratorReports(c.Request.Context(), uint(boardID), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.SuccessPage(c, list, total, req.Page, req.PageSize)
}

// ListModeratorBans 版主面板：板块封禁列表
// @Summary 版主面板封禁列表
// @Description 分页获取指定板块的封禁用户记录
// @Tags 版主面板
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} common.BasicResponse
// @Failure 400 {object} common.BasicResponse
// @Router /moderator/boards/{id}/bans [get]
func (h *BoardHandler) ListModeratorBans(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.HandleError(c, apperrors.ErrInvalidRequest)
		return
	}
	var req request.ListModeratorBansRequest
	if err := req.Bind(c); err != nil {
		response.HandleError(c, err)
		return
	}
	list, total, err := h.boardSvc.ListModeratorBans(c.Request.Context(), uint(boardID), &req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.SuccessPage(c, list, total, req.Page, req.PageSize)
}