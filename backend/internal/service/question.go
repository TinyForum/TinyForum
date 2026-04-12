package service

import (
	"errors"
	"fmt"

	apperrors "tiny-forum/internal/errors"
	"tiny-forum/internal/model"
	"tiny-forum/internal/repository"

	"gorm.io/gorm"
)

type QuestionService struct {
	questionRepo *repository.QuestionRepository
	postRepo     *repository.PostRepository
	commentRepo  *repository.CommentRepository
	userRepo     *repository.UserRepository
	notifSvc     *NotificationService
	db           *gorm.DB
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

// CreateQuestion 创建问答帖
func (s *QuestionService) CreateQuestion(userID uint, input model.CreateQuestionInput) (*model.QuestionResponse, error) {
	// 验证输入
	if err := s.validateCreateQuestion(input); err != nil {
		return nil, err
	}

	// 调用 Repository 的事务方法
	question, err := s.questionRepo.CreateWithTransaction(userID, input)
	if err != nil {
		return nil, fmt.Errorf("创建问答失败: %w", err)
	}

	return question, nil
}

// validateCreateQuestion 验证创建问答的输入
func (s *QuestionService) validateCreateQuestion(input model.CreateQuestionInput) error {
	if input.Title == "" {
		return errors.New("标题不能为空")
	}
	if input.Content == "" {
		return errors.New("内容不能为空")
	}
	if len(input.Title) > 100 {
		return errors.New("标题长度不能超过100个字符")
	}
	if len(input.Summary) > 500 {
		return errors.New("摘要长度不能超过500个字符")
	}
	if input.RewardScore < 0 || input.RewardScore > 100 {
		return errors.New("悬赏积分必须在0-100之间")
	}
	return nil
}

func (s *QuestionService) GetQuestionDetail(questionID uint) (*model.QuestionResponse, error) {
	question, err := s.questionRepo.FindByID(questionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("问答不存在")
		}
		return nil, fmt.Errorf("查询问答失败: %w", err)
	}

	return &model.QuestionResponse{
		ID:               question.ID,
		PostID:           question.PostID,
		Title:            question.Post.Title,
		Content:          question.Post.Content,
		Summary:          question.Post.Summary,
		Cover:            question.Post.Cover,
		BoardID:          question.Post.BoardID,
		AuthorID:         question.Post.AuthorID,
		RewardScore:      question.RewardScore,
		AnswerCount:      question.AnswerCount,
		AcceptedAnswerID: question.AcceptedAnswerID,
		Status:           string(question.Post.Status),
		CreatedAt:        question.CreatedAt,
		UpdatedAt:        question.UpdatedAt,
	}, nil
}

// GetQuestions 获取问答帖列表，支持只看未回答
func (s *QuestionService) GetQuestions(page, pageSize int, unanswered bool) ([]model.Post, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	if unanswered {
		return s.postRepo.GetUnansweredQuestions(pageSize, offset)
	}
	return s.postRepo.GetQuestions(pageSize, offset)
}

func (s *QuestionService) AcceptAnswer(postID, commentID uint, userID uint) error {
	post, err := s.postRepo.FindByID(postID)
	if err != nil {
		return fmt.Errorf("%w: 帖子不存在", apperrors.ErrPostNotFound)
	}
	if post.AuthorID != userID {
		return fmt.Errorf("%w: 只有发帖人可以采纳答案", apperrors.ErrAcceptForbidden)
	}

	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return fmt.Errorf("%w: 回答不存在", apperrors.ErrAnswerNotFound)
	}
	if !comment.IsAnswer {
		return errors.New("该评论不是回答")
	}

	question, err := s.questionRepo.FindByPostID(postID)
	if err != nil {
		return fmt.Errorf("%w: 问答信息不存在", apperrors.ErrQuestionNotFound)
	}

	if question.AcceptedAnswerID != nil {
		return errors.New("已经采纳过答案了")
	}

	if err := s.questionRepo.SetAcceptedAnswer(postID, commentID); err != nil {
		return err
	}

	if err := s.commentRepo.MarkAsAccepted(commentID); err != nil {
		return err
	}

	if question.RewardScore > 0 {
		s.userRepo.AddScore(comment.AuthorID, question.RewardScore)
	}

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
	VoteCount int    `json:"vote_count"` // 当前净票数
	Action    string `json:"action"`     // added, removed, updated
}

func (s *QuestionService) VoteAnswer(userID uint, input VoteAnswerInput) (*VoteAnswerResult, error) {
	comment, err := s.commentRepo.FindByID(input.CommentID)
	if err != nil {
		return nil, errors.New("回答不存在")
	}
	if !comment.IsAnswer {
		return nil, errors.New("只能对回答进行投票")
	}
	if comment.AuthorID == userID {
		return nil, errors.New("不能给自己的答案投票")
	}

	existingVote, _ := s.questionRepo.FindAnswerVote(userID, input.CommentID)

	var result VoteAnswerResult
	var action string

	if existingVote != nil && existingVote.ID != 0 {
		if existingVote.VoteType == input.VoteType {
			if err := s.questionRepo.DeleteAnswerVote(userID, input.CommentID); err != nil {
				return nil, err
			}
			action = "removed"
			result.VoteType = ""
		} else {
			existingVote.VoteType = input.VoteType
			if err := s.questionRepo.UpdateAnswerVote(existingVote); err != nil {
				return nil, err
			}
			action = "updated"
			result.VoteType = input.VoteType
		}
	} else {
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

	voteCount, _ := s.questionRepo.GetAnswerVoteCount(input.CommentID)
	s.commentRepo.UpdateVoteCount(input.CommentID, voteCount)

	result.VoteCount = voteCount
	result.Action = action

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
		return map[string]interface{}{
			"has_voted":  false,
			"vote_type":  "",
			"vote_count": 0,
		}, nil
	}

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
