package comment

import (
	"tiny-forum/internal/model"
)

// List 获取评论列表（分页）
func (s *commentService) List(postID uint, page, pageSize int) ([]model.Comment, int64, error) {
	return s.commentRepo.ListByPost(postID, page, pageSize)
}

// GetCommentCount 获取评论总数
func (s *commentService) GetCommentCount(postID uint) (int64, error) {
	return s.commentRepo.CountByPost(postID)
}

// GetAnswerByID 获取回答详情
func (s *commentService) GetAnswerByID(commentID uint) (*model.Comment, error) {
	return s.commentRepo.FindByID(commentID)
}

// GetAnswersByPostID 获取帖子的所有答案（支持排序）
func (s *commentService) GetAnswersByPostID(postID uint, page, pageSize int, sortBy string) ([]model.Comment, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	switch sortBy {
	case "newest":
		return s.commentRepo.GetAnswersByPostIDOrderByNewest(postID, pageSize, offset)
	case "oldest":
		return s.commentRepo.GetAnswersByPostIDOrderByOldest(postID, pageSize, offset)
	default:
		return s.commentRepo.GetAnswersByPostID(postID, pageSize, offset)
	}
}

// GetAnswerVoteCount 获取答案的投票数（直接从 Comment 字段读取）
func (s *commentService) GetAnswerVoteCount(commentID uint) (int, error) {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return 0, err
	}
	return comment.VoteCount, nil
}

// GetVoteStatistics 获取投票统计（赞成/反对数）
func (s *commentService) GetVoteStatistics(answerID uint) (upCount, downCount int, err error) {
	upUsers, err := s.voteRepo.GetVoteUsers(answerID, 1)
	if err != nil {
		return 0, 0, err
	}
	downUsers, err := s.voteRepo.GetVoteUsers(answerID, -1)
	if err != nil {
		return 0, 0, err
	}
	return len(upUsers), len(downUsers), nil
}
