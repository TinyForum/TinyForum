package board

import (
	"strconv"

	"tiny-forum/pkg/response"

	"github.com/gin-gonic/gin"
)

// DeletePost 版主删除帖子
// @Summary 删除帖子（版主）
// @Description 版主或管理员删除指定板块下的帖子
// @Tags 帖子管理
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param post_id path int true "帖子ID"
// @Success 200 {object} vo.BasicResponse{data=object} "删除成功"
// @Failure 400 {object} vo.BasicResponse "无效的板块ID或帖子ID"
// @Failure 401 {object} vo.BasicResponse "未授权"
// @Failure 403 {object} vo.BasicResponse "无权限（需要版主或管理员权限）"
// @Failure 404 {object} vo.BasicResponse "帖子不存在"
// @Router /boards/{id}/posts/{post_id} [delete]
func (h *BoardHandler) DeletePost(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	userID := c.GetUint("user_id")
	role, _ := c.Get("user_role")
	isAdmin := role == "admin" || role == "super_admin"

	if err := h.boardSvc.DeletePost(uint(boardID), uint(postID), userID, isAdmin); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}

// PinPost 版主置顶/取消置顶帖子
// @Summary 置顶/取消置顶帖子
// @Description 版主或管理员置顶或取消置顶指定板块下的帖子
// @Tags 帖子管理
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "板块ID"
// @Param post_id path int true "帖子ID"
// @Param body body object true "置顶选项"
// @Param body.pin_in_board body bool true "是否置顶" example(true)
// @Success 200 {object} vo.BasicResponse{data=object} "操作成功"
// @Failure 400 {object} vo.BasicResponse "无效的板块ID或帖子ID"
// @Failure 401 {object} vo.BasicResponse "未授权"
// @Failure 403 {object} vo.BasicResponse "无权限（需要版主或管理员权限）"
// @Failure 404 {object} vo.BasicResponse "帖子不存在"
// @Router /boards/{id}/posts/{post_id}/pin [put]
func (h *BoardHandler) PinPost(c *gin.Context) {
	boardID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的板块ID")
		return
	}
	postID, err := strconv.ParseUint(c.Param("post_id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的帖子ID")
		return
	}

	var body struct {
		PinInBoard bool `json:"pin_in_board"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.boardSvc.PinPost(uint(boardID), uint(postID), body.PinInBoard); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "操作成功"})
}
