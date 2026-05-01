package board

import (
	"errors"
	"strconv"

	"tiny-forum/internal/model/dto"
	boardService "tiny-forum/internal/service/board"
	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Create 创建板块
// @Summary 创建新板块
// @Description 管理员创建一个新的论坛板块，需要管理员权限
// @Tags 板块管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param body body board.CreateBoardInput true "板块信息"
// @Success 200 {object} response.Response{data=po.Board} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Router /boards [post]
func (h *BoardHandler) Create(c *gin.Context) {
	var input boardService.CreateBoardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	board, err := h.boardSvc.Create(input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, board)
}

// Update 更新板块（仅管理员）
// @Summary 更新板块信息
// @Description 管理员更新指定板块的信息，需要管理员权限
// @Tags 板块管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID" minimum(1) example(1)
// @Param body body board.CreateBoardInput true "板块信息"
// @Success 200 {object} response.Response{data=po.Board} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误或板块ID无效"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "板块不存在"
// @Router /boards/{id} [put]
func (h *BoardHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	var input boardService.CreateBoardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	board, err := h.boardSvc.Update(uint(id), input)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, board)
}

// Delete 删除板块（仅管理员）
// @Summary 删除板块（仅管理员）
// @Description 删除指定板块，需要管理员权限
// @Tags 板块管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Success 200 {object} response.Response{data=object} "删除成功"
// @Failure 400 {object} response.Response "无效的板块ID"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权限"
// @Failure 404 {object} response.Response "板块不存在"
// @Router /boards/{id} [delete]
func (h *BoardHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	if err := h.boardSvc.Delete(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}

// GetByID 获取板块详情
// @Summary 获取板块详情
// @Description 根据ID获取板块详细信息
// @Tags 板块管理
// @Produce json
// @Param id path int true "板块ID"
// @Success 200 {object} response.Response{data=po.Board} "获取成功"
// @Failure 400 {object} response.Response "无效的板块ID"
// @Failure 404 {object} response.Response "板块不存在"
// @Router /boards/{id} [get]
func (h *BoardHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	board, err := h.boardSvc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, board)
}

// GetBoardBySlug 根据 slug 获取板块
// @Summary 根据板块标识符获取板块
// @Description 根据板块标识符（slug）获取板块信息
// @Tags 板块管理
// @Produce json
// @Param slug path string true "板块标识符"
// @Success 200 {object} response.Response{data=po.Board} "获取成功"
// @Failure 404 {object} response.Response "板块不存在"
// @Router /boards/slug/{slug} [get]
func (h *BoardHandler) GetBoardBySlug(c *gin.Context) {
	slug := c.Param("slug")
	board, err := h.boardSvc.GetBoardBySlug(slug)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.Success(c, board)
}

// List 获取板块列表
// @Summary 获取板块列表
// @Description 分页获取板块列表
// @Tags 板块管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]po.Board}} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /boards [get]
func (h *BoardHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	boards, total, err := h.boardSvc.List(page, pageSize)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.SuccessPage(c, boards, total, page, pageSize)
}

// GetTree 获取板块树形结构
// @Summary 获取板块树形结构
// @Description 获取所有板块的树形层级结构
// @Tags 板块管理
// @Produce json
// @Success 200 {object} response.Response{data=[]po.BoardTree} "获取成功"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /boards/tree [get]
func (h *BoardHandler) GetTree(c *gin.Context) {
	tree, err := h.boardSvc.GetTree()
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, tree)
}

// GetPostsBySlug 获取板块下的帖子列表
// @Summary 获取板块下的帖子列表
// @Description 分页获取指定板块下的所有帖子
// @Tags 板块管理
// @Produce json
// @Param slug path string true "板块标识符"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData{list=[]model.Post}} "获取成功"
// @Failure 400 {object} response.Response "板块 slug 不能为空"
// @Failure 404 {object} response.Response "板块不存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /boards/slug/{slug}/posts [get]
func (h *BoardHandler) GetPostsBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "板块 slug 不能为空")
		return
	}

	var req dto.GetBoardPostsRequest
	// 绑定查询参数 (page, page_size)
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "分页参数错误")
		return
	}
	// 设置默认值（若未传或为0）
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 20
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	posts, total, err := h.boardSvc.GetPostsBySlug(slug, req.Page, req.PageSize)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.NotFound(c, "板块不存在")
		} else {
			response.InternalError(c, err.Error())
		}
		return
	}

	response.SuccessPage(c, posts, total, req.Page, req.PageSize)
}
