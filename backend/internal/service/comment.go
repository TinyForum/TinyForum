package service

import (
	"errors"

	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
)

type CommentService struct {
	commentRepo *repository.CommentRepository
	postRepo    repository.PostRepository
	userRepo    *repository.UserRepository
	notifSvc    *NotificationService
	questionSvc *QuestionService // 添加问答服务依赖（可选）
	voteRepo    *repository.VoteRepository
}

func NewCommentService(
	commentRepo *repository.CommentRepository,
	postRepo repository.PostRepository,
	userRepo *repository.UserRepository,
	notifSvc *NotificationService,
	voteRepo *repository.VoteRepository,
) *CommentService {
	return &CommentService{
		commentRepo: commentRepo,
		postRepo:    postRepo,
		userRepo:    userRepo,
		notifSvc:    notifSvc,
		voteRepo:    voteRepo,
	}
}

type CreateCommentInput struct {
	PostID   uint   `json:"post_id" binding:"required"`
	Content  string `json:"content" binding:"required,min=1,max=2000"`
	ParentID *uint  `json:"parent_id"`
}

func (s *CommentService) Create(authorID uint, input CreateCommentInput) (*model.Comment, error) {
	// Validate post exists
	post, err := s.postRepo.FindByID(input.PostID)
	if err != nil {
		return nil, errors.New("帖子不存在")
	}

	// 验证父评论
	if input.ParentID != nil && *input.ParentID != 0 {
		if err := s.commentRepo.ValidateParentComment(*input.ParentID, input.PostID); err != nil {
			return nil, err
		}
	}
	comment := &model.Comment{
		Content:  input.Content,
		PostID:   input.PostID,
		AuthorID: authorID,
		ParentID: input.ParentID,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	_ = s.userRepo.AddScore(authorID, 3)

	// Notify post author
	if post.AuthorID != authorID {
		s.notifSvc.Create(post.AuthorID, &authorID, model.NotifyComment,
			"有人评论了你的帖子《"+post.Title+"》", &input.PostID, "post")
	}

	// Notify parent comment author
	if input.ParentID != nil {
		parent, err := s.commentRepo.FindByID(*input.ParentID)
		if err == nil && parent.AuthorID != authorID {
			s.notifSvc.Create(parent.AuthorID, &authorID, model.NotifyReply,
				"有人回复了你的评论", input.ParentID, "comment")
		}
	}

	return s.commentRepo.FindByID(comment.ID)
}

func (s *CommentService) Delete(commentID, userID uint, isAdmin bool) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}
	if comment.AuthorID != userID && !isAdmin {
		return errors.New("无权限删除此评论")
	}
	return s.commentRepo.Delete(commentID)
}

func (s *CommentService) List(postID uint, page, pageSize int) ([]model.Comment, int64, error) {
	return s.commentRepo.ListByPost(postID, page, pageSize)
}

func (s *CommentService) GetCommentCount(postID uint) (int64, error) {
	return s.commentRepo.CountByPost(postID)
}

// ========== 新增的问答相关方法 ==========

// CreateAnswer 创建回答（仅限问答帖）
func (s *CommentService) CreateAnswer(authorID uint, input CreateCommentInput) (*model.Comment, error) {
	// Validate post exists and is question
	post, err := s.postRepo.FindByID(input.PostID)
	if err != nil {
		return nil, errors.New("帖子不存在")
	}
	if post.Type != "question" {
		return nil, errors.New("该帖子不是问答类型，请使用普通评论")
	}

	comment := &model.Comment{
		Content:  input.Content,
		PostID:   input.PostID,
		AuthorID: authorID,
		ParentID: input.ParentID,
		IsAnswer: true, // 标记为答案
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}

	// 回答奖励积分
	_ = s.userRepo.AddScore(authorID, 2)

	// Notify post author
	if post.AuthorID != authorID {
		s.notifSvc.Create(post.AuthorID, &authorID, model.NotifyComment,
			"有人回答了你的问题《"+post.Title+"》", &input.PostID, "post")
	}

	return s.commentRepo.FindByID(comment.ID)
}

// MarkAsAnswer 标记/取消标记为答案
func (s *CommentService) MarkAsAnswer(commentID, userID uint, isAdmin bool, isAnswer bool) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("评论不存在")
	}

	// Get post to check permissions
	post, err := s.postRepo.FindByID(comment.PostID)
	if err != nil {
		return errors.New("帖子不存在")
	}

	// Only post author or admin can mark as answer
	if post.AuthorID != userID && !isAdmin {
		return errors.New("无权限操作")
	}

	return s.commentRepo.MarkAsAnswer(commentID, isAnswer)
}

// GetAnswersByPostID 获取帖子的所有答案
func (s *CommentService) GetAnswersByPostID(postID uint, page, pageSize int, sortBy string) ([]model.Comment, int64, error) {
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
	default: // vote
		return s.commentRepo.GetAnswersByPostID(postID, pageSize, offset)
	}
}

// GetAnswerVoteCount 获取答案的投票数
func (s *CommentService) GetAnswerVoteCount(commentID uint) (int, error) {
	// 这个会通过 repository 的 questionRepo 来实现
	// 或者直接从 comment 的 vote_count 字段获取
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return 0, err
	}
	return comment.VoteCount, nil
}

// 获取回答
func (s *CommentService) GetAnswerByID(commentID uint) (*model.Comment, error) {
	return s.commentRepo.FindByID(commentID)
}

// DeleteAnswer 删除回答
func (s *CommentService) DeleteAnswer(commentID, userID uint, isAdmin bool) error {
	// 1. 查找回答
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("回答不存在")
	}

	// 2. 验证是否为回答类型
	if !comment.IsAnswer {
		return errors.New("该评论不是回答")
	}

	// 3. 权限检查
	// 3.1 管理员可以删除任何回答
	if isAdmin {
		return s.commentRepo.Delete(commentID)
	}

	// 3.2 回答作者可以删除自己的回答
	if comment.AuthorID == userID {
		return s.commentRepo.Delete(commentID)
	}

	// 3.3 问题作者可以删除自己问题下的回答
	post, err := s.postRepo.FindByID(comment.PostID)
	if err == nil && post.AuthorID == userID {
		return s.commentRepo.Delete(commentID)
	}

	// 4. 无权限
	return errors.New("无权限删除此回答")
}

// RemoveVote 取消投票
// func (s *CommentService) RemoveVote(answerID uint, userID uint) (*model.Comment, error) {
// 	// 1. 查找回答是否存在
// 	comment, err := s.commentRepo.FindByID(answerID)
// 	if err != nil {
// 		return nil, errors.New("回答不存在")
// 	}

// 	// 2. 验证是否为回答类型
// 	if !comment.IsAnswer {
// 		return nil, errors.New("该评论不是回答")
// 	}

// 	// 3. 查找用户的投票记录
// 	vote, err := s.commentRepo.FindVote(answerID, userID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("尚未投票，无法取消")
// 		}
// 		return nil, err
// 	}

// 	// 4. 取消投票：删除投票记录并更新回答的投票数
// 	err = s.commentRepo.RemoveVote(answerID, userID, vote.Value)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// 5. 返回更新后的回答
// 	return s.commentRepo.FindByID(answerID)
// }

// VoteAnswer 投票回答
func (s *CommentService) VoteAnswer(answerID uint, userID uint, voteType string) (*model.Comment, error) {
	// 1. 查找回答
	comment, err := s.commentRepo.FindByID(answerID)
	if err != nil {
		return nil, errors.New("回答不存在")
	}

	// 2. 验证是否为回答
	if !comment.IsAnswer {
		return nil, errors.New("该评论不是回答")
	}

	// 3. 不能给自己投票
	if comment.AuthorID == userID {
		return nil, errors.New("不能给自己的回答投票")
	}

	// 4. 转换投票类型
	var value int
	switch voteType {
	case "up":
		value = 1
	case "down":
		value = -1
	default:
		return nil, errors.New("无效的投票类型")
	}

	// 5. 获取用户当前投票状态
	currentVote, err := s.voteRepo.GetUserVote(answerID, userID)
	if err != nil {
		return nil, err
	}

	// 6. 处理投票逻辑
	if currentVote == value {
		// 相同投票：取消投票（toggle）
		if err := s.voteRepo.RemoveVote(answerID, userID); err != nil {
			return nil, err
		}
	} else if currentVote == 0 {
		// 未投票：创建新投票
		if err := s.voteRepo.CreateOrUpdateVote(answerID, userID, value); err != nil {
			return nil, err
		}
	} else {
		// 改变投票：更新现有投票
		if err := s.voteRepo.CreateOrUpdateVote(answerID, userID, value); err != nil {
			return nil, err
		}
	}

	// 7. 返回更新后的回答
	return s.commentRepo.FindByID(answerID)
}

// RemoveVote 取消投票
func (s *CommentService) RemoveVote(answerID uint, userID uint) (*model.Comment, error) {
	// 1. 查找回答
	comment, err := s.commentRepo.FindByID(answerID)
	if err != nil {
		return nil, errors.New("回答不存在")
	}

	// 2. 验证是否为回答
	if !comment.IsAnswer {
		return nil, errors.New("该评论不是回答")
	}

	// 3. 检查是否投过票
	currentVote, err := s.voteRepo.GetUserVote(answerID, userID)
	if err != nil {
		return nil, err
	}
	if currentVote == 0 {
		return nil, errors.New("尚未投票，无法取消")
	}

	// 4. 删除投票
	if err := s.voteRepo.RemoveVote(answerID, userID); err != nil {
		return nil, err
	}

	// 5. 返回更新后的回答
	return s.commentRepo.FindByID(answerID)
}

// GetUserVoteStatus 获取用户投票状态
func (s *CommentService) GetUserVoteStatus(answerID uint, userID uint) (int, error) {
	return s.voteRepo.GetUserVote(answerID, userID)
}

// GetVoteStatistics 获取投票统计
func (s *CommentService) GetVoteStatistics(answerID uint) (upCount, downCount int, err error) {
	// 获取所有投票用户
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

// UnacceptAnswer 取消接受答案
func (s *CommentService) UnacceptAnswer(answerID, userID uint, isAdmin bool) error {
	// 1. 查找回答
	answer, err := s.commentRepo.FindByID(answerID)
	if err != nil {
		return errors.New("回答不存在")
	}

	// 2. 验证是否为回答类型
	if !answer.IsAnswer {
		return errors.New("该评论不是回答")
	}

	// 3. 查找关联的问题
	post, err := s.postRepo.FindByID(answer.PostID)
	if err != nil {
		return errors.New("问题不存在")
	}

	// 4. 验证问题是否为问答类型
	if post.Type != "question" {
		return errors.New("该帖子不是问答类型")
	}

	// 5. 权限检查：只有问题作者或管理员可以取消接受答案
	if post.AuthorID != userID && !isAdmin {
		return errors.New("没有权限操作，只有问题作者可以取消接受答案")
	}

	// 6. 检查回答是否已被接受
	if !answer.IsAccepted {
		return errors.New("该回答未被接受为答案")
	}

	// 7. 取消接受答案
	if err := s.commentRepo.UnacceptAnswer(answerID); err != nil {
		return err
	}

	// 8. 可选：扣除相关积分奖励
	// _ = s.userRepo.AddScore(answer.AuthorID, -10) // 扣除之前奖励的积分

	// 9. 发送通知（可选）
	if post.AuthorID != userID {
		s.notifSvc.Create(answer.AuthorID, &userID, model.NotifyAcceptCancel,
			"你的答案在问题《"+post.Title+"》中被取消接受", &answer.PostID, "post")
	}

	return nil
}
