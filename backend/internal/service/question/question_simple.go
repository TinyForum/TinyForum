package question

import (
	"time"
)

// SimpleAuthor 精简的作者信息
type SimpleAuthor struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// SimpleTag 精简的标签信息
type SimpleTag struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// QuestionSimpleResponse 问题精简列表响应
type QuestionSimpleResponse struct {
	ID               uint          `json:"id"`
	Title            string        `json:"title"`
	Summary          string        `json:"summary"`
	ViewCount        int           `json:"view_count"`
	RewardScore      int           `json:"reward_score"`
	AnswerCount      int           `json:"answer_count"`
	AcceptedAnswerID *uint         `json:"accepted_answer_id"`
	Author           *SimpleAuthor `json:"author"`
	Tags             []SimpleTag   `json:"tags"`
	CreatedAt        time.Time     `json:"created_at"`
}

// GetQuestionSimpleList 获取问题精简列表
func (s *QuestionService) GetQuestionSimpleList(pageSize, offset int, boardID *uint, filter, sort, keyword string) ([]QuestionSimpleResponse, int64, error) {
	questions, total, err := s.questionRepo.FindSimpleQuestions(pageSize, offset, boardID, filter, sort, keyword)
	if err != nil {
		return nil, 0, err
	}
	if len(questions) == 0 {
		return []QuestionSimpleResponse{}, total, nil
	}
	authorIDSet := make(map[uint]bool)
	postIDSet := make(map[uint]bool)
	for _, q := range questions {
		if q.AuthorID > 0 {
			authorIDSet[q.AuthorID] = true
		}
		if q.PostID > 0 {
			postIDSet[q.PostID] = true
		}
	}
	authorIDs := make([]uint, 0, len(authorIDSet))
	for id := range authorIDSet {
		authorIDs = append(authorIDs, id)
	}
	postIDs := make([]uint, 0, len(postIDSet))
	for id := range postIDSet {
		postIDs = append(postIDs, id)
	}
	authorMap := make(map[uint]*SimpleAuthor)
	if len(authorIDs) > 0 {
		authors, err := s.userRepo.FindByIDs(authorIDs)
		if err == nil {
			for i := range authors {
				authorMap[authors[i].ID] = &SimpleAuthor{
					ID:     authors[i].ID,
					Name:   authors[i].Username,
					Avatar: authors[i].Avatar,
				}
			}
		}
	}
	tagMap := make(map[uint][]SimpleTag)
	if len(postIDs) > 0 {
		tagsMap, err := s.tagRepo.FindTagsByPostIDs(postIDs)
		if err == nil {
			for postID, tags := range tagsMap {
				simpleTags := make([]SimpleTag, len(tags))
				for i, tag := range tags {
					simpleTags[i] = SimpleTag{
						ID:   tag.ID,
						Name: tag.Name,
					}
				}
				tagMap[postID] = simpleTags
			}
		} else {
			tagMap = make(map[uint][]SimpleTag)
		}
	}
	result := make([]QuestionSimpleResponse, len(questions))
	for i, q := range questions {
		result[i] = QuestionSimpleResponse{
			ID:               q.ID,
			Title:            q.Title,
			Summary:          q.Summary,
			ViewCount:        q.ViewCount,
			RewardScore:      q.RewardScore,
			AnswerCount:      q.AnswerCount,
			AcceptedAnswerID: q.AcceptedAnswerID,
			CreatedAt:        q.CreatedAt,
			Tags:             []SimpleTag{},
		}
		if author, ok := authorMap[q.AuthorID]; ok {
			result[i].Author = author
		}
		if tags, ok := tagMap[q.PostID]; ok {
			result[i].Tags = tags
		}
	}
	return result, total, nil
}

// GetQuestionSimpleByID 获取单个问题详情（精简版）
func (s *QuestionService) GetQuestionSimpleByID(questionID uint) (*QuestionSimpleResponse, error) {
	question, err := s.questionRepo.FindQuestionSimpleByID(questionID)
	if err != nil {
		return nil, err
	}
	var simpleAuthor *SimpleAuthor
	author, err := s.userRepo.FindByID(question.AuthorID)
	if err == nil && author != nil {
		simpleAuthor = &SimpleAuthor{
			ID:     author.ID,
			Name:   author.Username,
			Avatar: author.Avatar,
		}
	}
	simpleTags := []SimpleTag{}
	tags, err := s.tagRepo.FindTagsByPostID(question.PostID)
	if err == nil {
		simpleTags = make([]SimpleTag, len(tags))
		for i, tag := range tags {
			simpleTags[i] = SimpleTag{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}
	}
	result := &QuestionSimpleResponse{
		ID:               question.ID,
		Title:            question.Title,
		Summary:          question.Summary,
		ViewCount:        question.ViewCount,
		RewardScore:      question.RewardScore,
		AnswerCount:      question.AnswerCount,
		AcceptedAnswerID: question.AcceptedAnswerID,
		Author:           simpleAuthor,
		Tags:             simpleTags,
		CreatedAt:        question.CreatedAt,
	}
	return result, nil
}
