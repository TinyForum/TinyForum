package question

import (
	"tiny-forum/internal/model/dto"
	"tiny-forum/internal/model/vo"
)

// GetQuestionSimpleList 获取问题精简列表
func (s *questionService) GetQuestionSimpleList(pageSize, offset int, boardID *uint, filter, sort, keyword string) ([]vo.QuestionSimpleVO, int64, error) {
	questions, total, err := s.questionRepo.FindSimpleQuestions(pageSize, offset, boardID, filter, sort, keyword)
	if err != nil {
		return nil, 0, err
	}
	if len(questions) == 0 {
		return []vo.QuestionSimpleVO{}, total, nil
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
	authorMap := make(map[uint]*dto.SimpleAuthor)
	if len(authorIDs) > 0 {
		authors, err := s.userRepo.FindByIDs(authorIDs)
		if err == nil {
			for i := range authors {
				authorMap[authors[i].ID] = &dto.SimpleAuthor{
					ID:        authors[i].ID,
					Name:      authors[i].Username,
					AvatarUrl: authors[i].AvatarUrl,
				}
			}
		}
	}
	tagMap := make(map[uint][]dto.SimpleTag)
	if len(postIDs) > 0 {
		tagsMap, err := s.tagRepo.FindTagsByPostIDs(postIDs)
		if err == nil {
			for postID, tags := range tagsMap {
				simpleTags := make([]dto.SimpleTag, len(tags))
				for i, tag := range tags {
					simpleTags[i] = dto.SimpleTag{
						ID:   tag.ID,
						Name: tag.Name,
					}
				}
				tagMap[postID] = simpleTags
			}
		} else {
			tagMap = make(map[uint][]dto.SimpleTag)
		}
	}
	result := make([]vo.QuestionSimpleVO, len(questions))
	for i, q := range questions {
		result[i] = vo.QuestionSimpleVO{
			ID:               q.ID,
			Title:            q.Title,
			Summary:          q.Summary,
			ViewCount:        q.ViewCount,
			RewardScore:      q.RewardScore,
			AnswerCount:      q.AnswerCount,
			AcceptedAnswerID: q.AcceptedAnswerID,
			CreatedAt:        q.CreatedAt,
			Tags:             []dto.SimpleTag{},
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
func (s *questionService) GetQuestionSimpleByID(questionID uint) (*vo.QuestionSimpleVO, error) {
	question, err := s.questionRepo.FindQuestionSimpleByID(questionID)
	if err != nil {
		return nil, err
	}
	var simpleAuthor *dto.SimpleAuthor
	author, err := s.userRepo.FindByID(question.AuthorID)
	if err == nil && author != nil {
		simpleAuthor = &dto.SimpleAuthor{
			ID:        author.ID,
			Name:      author.Username,
			AvatarUrl: author.AvatarUrl,
		}
	}
	simpleTags := []dto.SimpleTag{}
	tags, err := s.tagRepo.FindTagsByPostID(question.PostID)
	if err == nil {
		simpleTags = make([]dto.SimpleTag, len(tags))
		for i, tag := range tags {
			simpleTags[i] = dto.SimpleTag{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}
	}
	result := &vo.QuestionSimpleVO{
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
