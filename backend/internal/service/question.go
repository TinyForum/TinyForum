package service

import (
	"errors"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"
)

type QuestionService struct {
	questionRepo *repository.QuestionRepository
	postRepo     *repository.PostRepository
	commentRepo  *repository.CommentRepository
	userRepo     *repository.UserRepository
	notifSvc     *NotificationService
}

func NewQuestionService(
	questionRepo *repository.QuestionRepository,
	postRepo *repository.PostRepository,
	commentRepo *repository.CommentRepository,
	userRepo *repository.UserRepository,
	notifSvc *NotificationService,
) *QuestionService {
	return &QuestionService{
		questionRepo: questionRepo,
		postRepo:     postRepo,
		commentRepo:  commentRepo,
		userRepo:     userRepo,
		notifSvc:     notifSvc,
	}
}

type CreateQuestionInput struct {
	PostID      uint `json:"post_id" binding:"required"`
	RewardScore int  `json:"reward_score" binding:"min=0,max=100"`
}

func (s *QuestionService) CreateQuestion(input CreateQuestionInput) error {
	// Check if post exists and is question
	post, err := s.postRepo.FindByID(input.PostID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	if !post.IsQuestion {
		return errors.New("帖子不是问答类型")
	}

	// Deduct reward score from author
	if input.RewardScore > 0 {
		if err := s.userRepo.AddScore(post.AuthorID, -input.RewardScore); err != nil {
			return errors.New("积分不足")
		}
	}

	question := &model.Question{
		PostID:      input.PostID,
		RewardScore: input.RewardScore,
	}

	return s.questionRepo.Create(question)
}

func (s *QuestionService) AcceptAnswer(postID, commentID uint, userID uint) error {
	// Check if user is post author
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return errors.New("帖子不存在")
	}
	if post.AuthorID != userID {
		return errors.New("只有发帖人可以采纳答案")
	}

	// Check if comment exists and is answer
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return errors.New("回答不存在")
	}
	if !comment.IsAnswer {
		return errors.New("该评论不是回答")
	}

	// Get question
	question, err := s.questionRepo.FindByPostID(postID)
	if err != nil {
		return errors.New("问答信息不存在")
	}

	// If already has accepted answer, return error
	if question.AcceptedAnswerID != nil {
		return errors.New("已经采纳过答案了")
	}

	// Set accepted answer
	if err := s.questionRepo.SetAcceptedAnswer(postID, commentID); err != nil {
		return err
	}

	// Mark comment as accepted
	if err := s.commentRepo.MarkAsAccepted(commentID); err != nil {
		return err
	}

	// Reward the answer author
	if question.RewardScore > 0 {
		s.userRepo.AddScore(comment.AuthorID, question.RewardScore)
	}

	// Notify answer author
	s.notifSvc.Create(comment.AuthorID, &userID, model.NotifySystem,
		"你的回答被采纳为最佳答案", &postID, "post")

	return nil
}

type VoteAnswerInput struct {
	CommentID uint   `json:"comment_id" binding:"required"`
	VoteType  string `json:"vote_type" binding:"required,oneof=up down"`
}

type VoteAnswerResult struct {
	VoteType  string `json:"vote_type"`  // up, down, or "" (取消)
	VoteCount int    `json:"vote_count"` // 当前总票数
	Action    string `json:"action"`     // added, removed, updated
}

func (s *QuestionService) VoteAnswer(userID uint, input VoteAnswerInput) (*VoteAnswerResult, error) {
	// Check if comment exists and is answer
	comment, err := s.commentRepo.FindByID(input.CommentID)
	if err != nil {
		return nil, errors.New("回答不存在")
	}
	if !comment.IsAnswer {
		return nil, errors.New("只能对回答进行投票")
	}

	// 不能给自己的答案投票
	if comment.AuthorID == userID {
		return nil, errors.New("不能给自己的答案投票")
	}

	// Check if user already voted
	existingVote, _ := s.questionRepo.FindAnswerVote(userID, input.CommentID)

	var result VoteAnswerResult
	var action string

	if existingVote != nil && existingVote.ID != 0 {
		// If same vote type, remove vote (toggle)
		if existingVote.VoteType == input.VoteType {
			if err := s.questionRepo.DeleteAnswerVote(userID, input.CommentID); err != nil {
				return nil, err
			}
			action = "removed"
			result.VoteType = ""
		} else {
			// Update vote type
			existingVote.VoteType = input.VoteType
			if err := s.questionRepo.CreateAnswerVote(existingVote); err != nil {
				return nil, err
			}
			action = "updated"
			result.VoteType = input.VoteType
		}
	} else {
		// Create new vote
		vote := &model.AnswerVote{
			UserID:    userID,
			CommentID: input.CommentID,
			VoteType:  input.VoteType,
		}
		if err := s.questionRepo.CreateAnswerVote(vote); err != nil {
			return nil, err
		}
		action = "added"
		result.VoteType = input.VoteType
	}

	// Update vote count
	voteCount, _ := s.questionRepo.GetAnswerVoteCount(input.CommentID)
	s.commentRepo.UpdateVoteCount(input.CommentID, voteCount)

	result.VoteCount = voteCount
	result.Action = action

	// 发送通知（仅当添加或更新投票时）
	if action != "removed" {
		s.notifSvc.Create(comment.AuthorID, &userID, model.NotifyLike,
			"有人给你的答案投票了", &input.CommentID, "comment")
	}

	return &result, nil
}

// GetAnswerVoteStatus 获取用户对答案的投票状态
func (s *QuestionService) GetAnswerVoteStatus(userID, commentID uint) (map[string]interface{}, error) {
	vote, err := s.questionRepo.FindAnswerVote(userID, commentID)
	if err != nil {
		// 没有投票记录
		return map[string]interface{}{
			"has_voted":  false,
			"vote_type":  "",
			"vote_count": 0,
		}, nil
	}

	// 获取当前票数
	voteCount, _ := s.questionRepo.GetAnswerVoteCount(commentID)

	return map[string]interface{}{
		"has_voted":  true,
		"vote_type":  vote.VoteType,
		"vote_count": voteCount,
	}, nil
}

func (s *QuestionService) GetQuestionWithAnswers(postID uint, page, pageSize int) (*model.Question, []model.Comment, int64, error) {
	question, err := s.questionRepo.FindByPostID(postID)
	if err != nil {
		return nil, nil, 0, err
	}

	answers, total, err := s.commentRepo.GetAnswersByPostID(postID, pageSize, (page-1)*pageSize)
	return question, answers, total, err
}
