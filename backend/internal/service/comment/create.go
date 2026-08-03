package comment

import (
	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	apperrors "tiny-forum/pkg/errors"
)

// Create 创建普通评论
func (s *commentService) CreateComment(authorID uint, input bo.CreateCommentInput) (*do.Comment, error) {
	// 获取请求参数
	post, err := s.postRepo.FindByArticleID(input.PostID)
	if err != nil {
		return nil, apperrors.ErrPostNotFound
	}

	if input.ParentID != nil && *input.ParentID != 0 {
		if err := s.commentRepo.ValidateParentComment(*input.ParentID, input.PostID); err != nil {
			return nil, err
		}
	}
	comment := &do.Comment{
		WorksID: input.PostID,
		Reply: &do.Reply{
			AuthorID: authorID,
			ParentID: input.ParentID,
			Content:  input.Content,
			TargetID: input.PostID,
		},
	}

	if err := s.commentRepo.CreateComment(comment); err != nil {
		return nil, err
	}

	// 发布事件（触发零代码机器人）
	s.botSvc.PublishEvent("comment.created", map[string]any{
		"comment_id":   comment.ID,
		"post_id":      input.PostID,
		"post_content": input.Content,
		"author_id":    authorID,
		"parent_id":    input.ParentID,
	})

	_ = s.userRepo.AddScore(authorID, 3)

	if post.Creation.AuthorID != authorID {
		s.notifSvc.Create(post.Creation.AuthorID, &authorID, do.NotifyComment,
			"有人评论了你的帖子《"+post.Creation.Title+"》", &input.PostID, "post")
	}

	if input.ParentID != nil {
		parent, err := s.commentRepo.FindByCommentID(*input.ParentID)
		if err == nil && parent.Reply.AuthorID != authorID {
			s.notifSvc.Create(parent.Reply.AuthorID, &authorID, do.NotifyReply,
				"有人回复了你的评论", input.ParentID, "comment")
		}
	}

	// 异步记录评论行为，供推荐系统使用
	go s.recSvc.RecordBehavior(authorID, request.RecordBehaviorRequest{
		TargetID:     input.PostID,
		TargetType:   "creation",
		BehaviorType: string(do.BehaviorComment),
		Value:        3.0,
	})

	return s.commentRepo.FindByCommentID(comment.ID)
}

// CreateAnswer 创建回答（仅限问答帖）
func (s *commentService) CreateAnswer(authorID uint, input bo.CreateAnswerInput) (*do.Answer, error) {
	// 查找 question id
	question, err := s.postRepo.FindQuestionByQuestionID(input.QuestionID)
	if err != nil {
		return nil, apperrors.ErrPostNotFound
	}
	// if post.Creation.Type != "question" {
	// 	return nil, errors.New("该帖子不是问答类型，请使用普通评论")
	// }

	comment := &do.Answer{
		QuestionID: input.QuestionID,
		Reply: &do.Reply{
			Content:  input.Content,
			AuthorID: authorID,
			ParentID: input.ParentID,
		},
	}

	if err := s.commentRepo.CreateAnswer(comment); err != nil {
		return nil, err
	}

	// 发布事件（触发零代码机器人）
	s.botSvc.PublishEvent("comment.created", map[string]any{
		"comment_id":  comment.ID,
		"post_id":     input.QuestionID,
		"content":     input.Content,
		"author_id":   authorID,
	})

	_ = s.userRepo.AddScore(authorID, 2)

	if question.Creation.AuthorID != authorID {
		s.notifSvc.Create(question.Creation.AuthorID, &authorID, do.NotifyComment,
			"有人回答了你的问题《"+question.Creation.Title+"》", &input.QuestionID, "post")
	}

	return s.commentRepo.FindByAnswerID(comment.ID)
}
