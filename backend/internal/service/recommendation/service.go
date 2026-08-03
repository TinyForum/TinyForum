package recommendation

import (
	"context"
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"tiny-forum/internal/model/do"
	"tiny-forum/internal/model/request"
	"tiny-forum/internal/model/vo"
	recRepo "tiny-forum/internal/repository/recommendation"
	"tiny-forum/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RecommendationService 推荐服务接口
type RecommendationService interface {
	// 行为追踪
	RecordBehavior(userID uint, req request.RecordBehaviorRequest) error
	BatchRecordBehaviors(userID uint, events []request.RecordBehaviorRequest) error
	RecordView(ctx context.Context, userID, creationID uint, sessionID string) error

	// 推荐
	GetRecommendations(ctx context.Context, userID uint, query request.RecommendationQuery) (*vo.RecommendationResponse, error)
	GetPersonalizedFeed(ctx context.Context, userID uint, page, pageSize int) (*vo.RecommendationResponse, error)

	// 反馈
	SubmitFeedback(userID uint, req request.RecommendationFeedbackRequest) error
	SubmitBatchFeedback(userID uint, req request.BatchFeedbackRequest) error

	// 内容分析
	AnalyzeContent(ctx context.Context, creationID uint) error

	// 画像
	GetInterestProfile(userID uint) (*vo.UserInterestProfileVO, error)
	RefreshUserProfile(userID uint) error

	// 管理员分析
	GetOverviewStats(ctx context.Context) (*vo.RecOverviewStats, error)
	GetBehaviorStats(ctx context.Context, days int) (*vo.BehaviorStats, error)
	GetUserAnalysis(ctx context.Context, days int) (*vo.UserAnalysis, error)
	GetContentPerformance(ctx context.Context) (*vo.ContentPerformance, error)
	GetRiskAnalysis(ctx context.Context) (*vo.RiskAnalysis, error)
	GetComprehensiveUserAnalysis(ctx context.Context) (*vo.ComprehensiveUserAnalysis, error)

	// 管理员 — 用户分析（6 大模块）
	GetUAOverview(ctx context.Context) (*vo.UAOverview, error)
	GetUASegments(ctx context.Context) (*vo.UASegments, error)
	GetUAProfiles(ctx context.Context, page, pageSize int, keyword, tier, sortBy string) (*vo.UAProfileList, error)
	GetUABehavior(ctx context.Context) (*vo.UABehavior, error)
	GetUACohorts(ctx context.Context) (*vo.UACohorts, error)
	GetUARisk(ctx context.Context) (*vo.UARisk, error)
}

type recommendationService struct {
	repo recRepo.RecommendationRepository
	db   *gorm.DB
}

func NewRecommendationService(repo recRepo.RecommendationRepository, db *gorm.DB) RecommendationService {
	return &recommendationService{repo: repo, db: db}
}

// --- 行为追踪 ---

func (s *recommendationService) RecordBehavior(userID uint, req request.RecordBehaviorRequest) error {
	behaviorType := do.BehaviorType(req.BehaviorType)
	if !behaviorType.IsValid() {
		behaviorType = do.BehaviorView
	}

	event := &do.UserBehaviorEvent{
		UserID:       userID,
		TargetID:     req.TargetID,
		TargetType:   req.TargetType,
		BehaviorType: behaviorType,
		Value:        req.Value,
		SessionID:    req.SessionID,
		ContextJSON:  req.ContextJSON,
		CreatedTS:    time.Now().UnixMilli(),
	}
	if event.Value == 0 {
		event.Value = 1.0
	}
	if event.ContextJSON == "" {
		event.ContextJSON = "{}"
	}

	return s.repo.RecordBehavior(event)
}

func (s *recommendationService) BatchRecordBehaviors(userID uint, events []request.RecordBehaviorRequest) error {
	doEvents := make([]do.UserBehaviorEvent, 0, len(events))
	now := time.Now()
	for _, e := range events {
		bType := do.BehaviorType(e.BehaviorType)
		if !bType.IsValid() {
			bType = do.BehaviorView
		}
		val := e.Value
		if val == 0 {
			val = 1.0
		}
		ctxJSON := e.ContextJSON
		if ctxJSON == "" {
			ctxJSON = "{}"
		}
		doEvents = append(doEvents, do.UserBehaviorEvent{
			UserID:       userID,
			TargetID:     e.TargetID,
			TargetType:   e.TargetType,
			BehaviorType: bType,
			Value:        val,
			SessionID:    e.SessionID,
			ContextJSON:  ctxJSON,
			CreatedTS:    now.UnixMilli(),
		})
	}
	return s.repo.BatchRecordBehaviors(doEvents)
}

func (s *recommendationService) RecordView(ctx context.Context, userID, creationID uint, sessionID string) error {
	return s.RecordBehavior(userID, request.RecordBehaviorRequest{
		TargetID:    creationID,
		TargetType:  "creation",
		BehaviorType: string(do.BehaviorView),
		Value:       0.3,
		SessionID:   sessionID,
	})
}

// --- 推荐引擎 ---

type recommendationStrategy struct {
	Name    string
	Weight  float64
	Enabled bool
}

var defaultStrategies = []recommendationStrategy{
	{Name: "content_similarity", Weight: 0.35, Enabled: true},
	{Name: "collaborative", Weight: 0.25, Enabled: true},
	{Name: "popular", Weight: 0.20, Enabled: true},
	{Name: "fresh", Weight: 0.10, Enabled: true},
	{Name: "explore", Weight: 0.10, Enabled: true},
}

func (s *recommendationService) GetRecommendations(ctx context.Context, userID uint, query request.RecommendationQuery) (*vo.RecommendationResponse, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 || query.PageSize > 50 {
		query.PageSize = 20
	}

	sessionID := generateSessionID(userID)

	scored := s.computeScores(ctx, userID)

	// 排序并分页
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].TotalScore > scored[j].TotalScore
	})

	start := (query.Page - 1) * query.PageSize
	end := start + query.PageSize
	if start > len(scored) {
		start = len(scored)
	}
	if end > len(scored) {
		end = len(scored)
	}

	pageItems := scored[start:end]
	items := s.enrichItems(ctx, pageItems)

	return &vo.RecommendationResponse{
		Items:       items,
		Total:       int64(len(scored)),
		Page:        query.Page,
		PageSize:    query.PageSize,
		HasMore:     end < len(scored),
		SessionID:   sessionID,
		Strategy:    "hybrid",
		GeneratedAt: time.Now(),
	}, nil
}

func (s *recommendationService) GetPersonalizedFeed(ctx context.Context, userID uint, page, pageSize int) (*vo.RecommendationResponse, error) {
	return s.GetRecommendations(ctx, userID, request.RecommendationQuery{
		Page:     page,
		PageSize: pageSize,
	})
}

type scoredItem struct {
	CreationID    uint
	ContentScore  float64
	CollabScore   float64
	PopScore      float64
	FreshScore    float64
	ExploreScore  float64
	TotalScore    float64
	Reason        string
}

func (s *recommendationService) computeScores(ctx context.Context, userID uint) []scoredItem {
	now := time.Now()

	// 权重配置
	const (
		wContent       = 0.35
		wCollaborative = 0.25
		wPopular       = 0.20
		wFresh         = 0.10
		wExplore       = 0.10
	)

	// 获取用户近期的互动过的内容（排除）
	recentBehaviors, _ := s.repo.GetUserRecentBehaviors(userID, 200)
	excludeIDs := make(map[uint]bool)
	var likedTagIDs []uint
	for _, b := range recentBehaviors {
		excludeIDs[b.TargetID] = true
		// 收集用户喜欢的内容的标签
		if b.BehaviorType == do.BehaviorLike || b.BehaviorType == do.BehaviorFavorite {
			cf, err := s.repo.GetContentFeature(b.TargetID)
			if err == nil && cf.TagIDsJSON != "" {
				var tagIDs []uint
				if err := json.Unmarshal([]byte(cf.TagIDsJSON), &tagIDs); err == nil {
					likedTagIDs = append(likedTagIDs, tagIDs...)
				}
			}
		}
	}

	// 获取用户兴趣画像
	var userTagWeights map[string]float64
	profile, err := s.repo.GetInterestProfile(userID)
	if err == nil && profile.TagWeightsJSON != "" {
		json.Unmarshal([]byte(profile.TagWeightsJSON), &userTagWeights)
	}

	// 候选内容：热门 + 标签匹配 + 协同过滤
	candidateSet := make(map[uint]*scoredItem)

	// 1. 内容相似度：基于用户喜欢的标签
	if likedTagIDs != nil {
		features, _ := s.repo.GetContentFeaturesByTags(likedTagIDs, sliceMapKeys(excludeIDs), 100)
		for _, f := range features {
			if _, exists := candidateSet[f.CreationID]; !exists {
				candidateSet[f.CreationID] = &scoredItem{CreationID: f.CreationID}
			}
			candidateSet[f.CreationID].ContentScore = f.QualityScore * 1.5
		}
	}

	// 2. 协同过滤：相似用户喜欢的内容
	similarUsers, _ := s.repo.FindSimilarUsers(userID, 10)
	for _, suID := range similarUsers {
		suBehaviors, _ := s.repo.GetUserBehaviorsByType(suID, do.BehaviorLike, now.AddDate(0, 0, -30))
		for _, b := range suBehaviors {
			if excludeIDs[b.TargetID] {
				continue
			}
			if _, exists := candidateSet[b.TargetID]; !exists {
				candidateSet[b.TargetID] = &scoredItem{CreationID: b.TargetID}
			}
			candidateSet[b.TargetID].CollabScore += b.Value * 0.3
		}
	}

	// 3. 热门内容
	hotFeatures, _ := s.repo.GetHotContentFeatures(80, sliceMapKeys(excludeIDs))
	for _, f := range hotFeatures {
		if _, exists := candidateSet[f.CreationID]; !exists {
			candidateSet[f.CreationID] = &scoredItem{CreationID: f.CreationID}
		}
		candidateSet[f.CreationID].PopScore = f.HotScore * 0.8
		candidateSet[f.CreationID].FreshScore = math.Exp(-float64(now.Hour()) / 24.0) * f.FreshnessScore
	}

	// 4. 探索：随机新鲜内容
	var exploreIDs []uint
	exploreFeatures, _ := s.repo.GetHotContentFeatures(30, nil)
	for i := 0; i < 10 && i < len(exploreFeatures); i++ {
		idx := (len(exploreFeatures) * i / 10) % len(exploreFeatures)
		if !excludeIDs[exploreFeatures[idx].CreationID] {
			exploreIDs = append(exploreIDs, exploreFeatures[idx].CreationID)
		}
	}
	for _, eid := range exploreIDs {
		if _, exists := candidateSet[eid]; !exists {
			candidateSet[eid] = &scoredItem{CreationID: eid}
		}
		candidateSet[eid].ExploreScore = 0.5
	}

	// 计算总分
	var result []scoredItem
	for _, item := range candidateSet {
		// 内容得分考虑用户标签权重
		contentBoost := 1.0
		if userTagWeights != nil {
			cf, err := s.repo.GetContentFeature(item.CreationID)
			if err == nil && cf.TagIDsJSON != "" {
				var tagIDs []uint
				json.Unmarshal([]byte(cf.TagIDsJSON), &tagIDs)
				for _, tid := range tagIDs {
					if w, ok := userTagWeights[strconv.Itoa(int(tid))]; ok {
						contentBoost += w * 2.0
					}
				}
			}
		}
		item.ContentScore *= contentBoost

		item.TotalScore = item.ContentScore*wContent +
			item.CollabScore*wCollaborative +
			item.PopScore*wPopular +
			item.FreshScore*wFresh +
			item.ExploreScore*wExplore

		// 确定推荐理由
		maxScore := math.Max(item.ContentScore,
			math.Max(item.CollabScore,
				math.Max(item.PopScore,
					math.Max(item.FreshScore, item.ExploreScore))))
		switch {
		case maxScore == item.ContentScore:
			item.Reason = "你可能喜欢"
		case maxScore == item.CollabScore:
			item.Reason = "与你兴趣相似的人也喜欢"
		case maxScore == item.PopScore:
			item.Reason = "热门内容"
		case maxScore == item.FreshScore:
			item.Reason = "新鲜推荐"
		default:
			item.Reason = "探索发现"
		}

		result = append(result, *item)
	}

	return result
}

func (s *recommendationService) enrichItems(ctx context.Context, items []scoredItem) []vo.RecommendationItem {
	if len(items) == 0 {
		return nil
	}

	creationIDs := make([]uint, len(items))
	for i, item := range items {
		creationIDs[i] = item.CreationID
	}

	// 批量查询 Creation 基本信息
	type CreationInfo struct {
		ID        uint
		Title     string
		Summary   string
		CoverURL  string
		AuthorID  uint
		ViewCount int
		LikeCount int
		BoardID   uint
		CreatedAt time.Time
	}
	var creations []CreationInfo
	s.db.WithContext(ctx).Raw(`
		SELECT c.id, c.title, c.summary, c.cover_url AS cover_url,
		       c.author_id, c.view_count, c.like_count, c.board_id, c.created_at
		FROM creations c
		WHERE c.id IN ?
		  AND c.deleted_at IS NULL
		  AND c.creation_status = 'published'
		  AND c.moderation_status = 'approved'
	`, creationIDs).Scan(&creations)

	creationMap := make(map[uint]CreationInfo)
	for _, c := range creations {
		creationMap[c.ID] = c
	}

	// 批量查询作者信息
	authorIDs := make([]uint, 0, len(creations))
	for _, c := range creations {
		authorIDs = append(authorIDs, c.AuthorID)
	}
	type AuthorInfo struct {
		ID       uint
		Username string
		Nickname string
		Avatar   string
	}
	var authors []AuthorInfo
	if len(authorIDs) > 0 {
		s.db.WithContext(ctx).Raw(`
			SELECT id, username, nickname, avatar
			FROM users WHERE id IN ? AND deleted_at IS NULL
		`, authorIDs).Scan(&authors)
	}
	authorMap := make(map[uint]AuthorInfo)
	for _, a := range authors {
		authorMap[a.ID] = a
	}

	// 批量查询板块
	boardIDs := make([]uint, 0, len(creations))
	for _, c := range creations {
		if c.BoardID > 0 {
			boardIDs = append(boardIDs, c.BoardID)
		}
	}
	type BoardInfo struct {
		ID   uint
		Name string
	}
	var boards []BoardInfo
	if len(boardIDs) > 0 {
		s.db.WithContext(ctx).Raw(`
			SELECT id, name FROM boards WHERE id IN ? AND deleted_at IS NULL
		`, boardIDs).Scan(&boards)
	}
	boardMap := make(map[uint]BoardInfo)
	for _, b := range boards {
		boardMap[b.ID] = b
	}

	// 批量查询评论数
	type CommentCount struct {
		WorksID uint
		Count   int
	}
	var commentCounts []CommentCount
	s.db.WithContext(ctx).Raw(`
		SELECT works_id, COUNT(*) AS count FROM comments
		WHERE works_id IN ? AND deleted_at IS NULL
		GROUP BY works_id
	`, creationIDs).Scan(&commentCounts)
	commentCountMap := make(map[uint]int)
	for _, cc := range commentCounts {
		commentCountMap[cc.WorksID] = cc.Count
	}

	result := make([]vo.RecommendationItem, 0, len(items))
	for _, item := range items {
		c, ok := creationMap[item.CreationID]
		if !ok {
			continue
		}
		recItem := vo.RecommendationItem{
			CreationID:   c.ID,
			Title:        c.Title,
			Summary:      c.Summary,
			CoverURL:     c.CoverURL,
			AuthorID:     c.AuthorID,
			ViewCount:    c.ViewCount,
			LikeCount:    c.LikeCount,
			CommentCount: commentCountMap[c.ID],
			BoardID:      c.BoardID,
			Score:        item.TotalScore,
			Reason:       item.Reason,
			CreatedAt:    c.CreatedAt.Format(time.RFC3339),
		}
		if author, ok := authorMap[c.AuthorID]; ok {
			recItem.AuthorName = author.Nickname
			if recItem.AuthorName == "" {
				recItem.AuthorName = author.Username
			}
			recItem.AuthorAvatar = author.Avatar
		}
		if board, ok := boardMap[c.BoardID]; ok {
			recItem.BoardName = board.Name
		}
		result = append(result, recItem)
	}

	return result
}

// --- 反馈 ---

func (s *recommendationService) SubmitFeedback(userID uint, req request.RecommendationFeedbackRequest) error {
	fbType := do.FeedbackType(req.FeedbackType)
	if !fbType.IsValid() {
		fbType = do.FeedbackImpression
	}
	source := req.SourceType
	if source == "" {
		source = "recommend"
	}

	fb := &do.RecommendationFeedback{
		UserID:       userID,
		CreationID:   req.CreationID,
		FeedbackType: fbType,
		SourceType:   source,
		Position:     req.Position,
		SessionID:    req.SessionID,
	}

	// 如果用户点了"不感兴趣"，同步记录行为事件
	if fbType == do.FeedbackNotInterested {
		_ = s.RecordBehavior(userID, request.RecordBehaviorRequest{
			TargetID:    req.CreationID,
			TargetType:  "creation",
			BehaviorType: string(do.BehaviorNotInterested),
			Value:       -1.0,
			SessionID:   req.SessionID,
		})
	}

	return s.repo.RecordFeedback(fb)
}

func (s *recommendationService) SubmitBatchFeedback(userID uint, req request.BatchFeedbackRequest) error {
	fbs := make([]do.RecommendationFeedback, 0, len(req.Feedbacks))
	for _, r := range req.Feedbacks {
		fbType := do.FeedbackType(r.FeedbackType)
		if !fbType.IsValid() {
			fbType = do.FeedbackImpression
		}
		source := r.SourceType
		if source == "" {
			source = "recommend"
		}
		fbs = append(fbs, do.RecommendationFeedback{
			UserID:       userID,
			CreationID:   r.CreationID,
			FeedbackType: fbType,
			SourceType:   source,
			Position:     r.Position,
			SessionID:    r.SessionID,
		})
	}
	return s.repo.BatchRecordFeedbacks(fbs)
}

// --- 内容分析 ---

func (s *recommendationService) AnalyzeContent(ctx context.Context, creationID uint) error {
	// 查询 Creation 的基本信息
	type CreationRow struct {
		ID        uint
		BoardID   uint
		AuthorID  uint
		Title     string
		LikeCount int
		ViewCount int
		CreatedAt time.Time
	}
	var row CreationRow
	err := s.db.WithContext(ctx).Raw(`
		SELECT id, board_id, author_id, title, like_count, view_count, created_at
		FROM creations WHERE id = ? AND deleted_at IS NULL
	`, creationID).Scan(&row).Error
	if err != nil {
		logger.Error("推荐系统：无法查询 Creation", zap.Uint("creation_id", creationID), zap.Error(err))
		return err
	}

	// 查询关联标签
	type TagRow struct {
		ID uint
	}
	var tagRows []TagRow
	s.db.WithContext(ctx).Raw(`
		SELECT t.id FROM tags t
		INNER JOIN creation_tags ct ON ct.tag_id = t.id
		WHERE ct.creation_id = ? AND t.deleted_at IS NULL
	`, creationID).Scan(&tagRows)

	tagIDs := make([]uint, len(tagRows))
	for i, t := range tagRows {
		tagIDs[i] = t.ID
	}
	tagIDsJSON, _ := json.Marshal(tagIDs)

	// 计算质量分: 基于点赞 + 浏览
	qualityScore := 0.0
	if row.ViewCount > 0 {
		qualityScore = float64(row.LikeCount) / float64(row.ViewCount) * 10.0
	}
	qualityScore += math.Log2(float64(row.ViewCount+1)) * 0.5

	// 计算热度分: 时间衰减
	hoursAgo := time.Since(row.CreatedAt).Hours()
	hotScore := (qualityScore + math.Log2(float64(row.LikeCount+1))) * math.Exp(-hoursAgo/72.0)

	// 新鲜度
	freshnessScore := math.Exp(-hoursAgo / 168.0)

	feature := &do.ContentFeature{
		CreationID:     creationID,
		TagIDsJSON:     string(tagIDsJSON),
		BoardID:        row.BoardID,
		AuthorID:       row.AuthorID,
		QualityScore:   qualityScore,
		HotScore:       hotScore,
		FreshnessScore: freshnessScore,
	}

	return s.repo.UpsertContentFeature(feature)
}

// --- 用户画像 ---

func (s *recommendationService) GetInterestProfile(userID uint) (*vo.UserInterestProfileVO, error) {
	profile, err := s.repo.GetInterestProfile(userID)
	if err != nil {
		return &vo.UserInterestProfileVO{
			ActiveTags: []string{},
			TagWeights: make(map[string]float64),
		}, nil
	}

	var tagWeights map[string]float64
	var activeTags []string
	json.Unmarshal([]byte(profile.TagWeightsJSON), &tagWeights)
	json.Unmarshal([]byte(profile.ActiveTagsJSON), &activeTags)

	return &vo.UserInterestProfileVO{
		ActiveTags: activeTags,
		TagWeights: tagWeights,
		UpdatedAt:  profile.LastUpdatedAt,
	}, nil
}

func (s *recommendationService) RefreshUserProfile(userID uint) error {
	now := time.Now()
	recentBehaviors, err := s.repo.GetUserRecentBehaviors(userID, 500)
	if err != nil {
		return err
	}

	// 统计标签权重
	tagWeights := make(map[string]float64)
	tagCount := make(map[string]int)
	for _, b := range recentBehaviors {
		cf, err := s.repo.GetContentFeature(b.TargetID)
		if err != nil || cf.TagIDsJSON == "" {
			continue
		}
		var tagIDs []uint
		if err := json.Unmarshal([]byte(cf.TagIDsJSON), &tagIDs); err != nil {
			continue
		}
		weight := b.Value
		if b.BehaviorType == do.BehaviorLike || b.BehaviorType == do.BehaviorFavorite {
			weight *= 3.0
		} else if b.BehaviorType == do.BehaviorComment || b.BehaviorType == do.BehaviorShare {
			weight *= 2.0
		} else if b.BehaviorType == do.BehaviorNotInterested {
			weight *= -2.0
		}

		for _, tid := range tagIDs {
			key := strconv.Itoa(int(tid))
			tagWeights[key] += weight
			tagCount[key]++
		}
	}

	// 归一化
	for tag := range tagWeights {
		tagWeights[tag] /= float64(tagCount[tag]+1)
	}

	// 获取最高的标签
	type tagEntry struct {
		ID     string
		Weight float64
	}
	var entries []tagEntry
	for k, w := range tagWeights {
		if w > 0 {
			entries = append(entries, tagEntry{ID: k, Weight: w})
		}
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Weight > entries[j].Weight
	})

	maxTags := 5
	if len(entries) < maxTags {
		maxTags = len(entries)
	}
	activeTags := make([]string, maxTags)
	for i := 0; i < maxTags; i++ {
		id, _ := strconv.Atoi(entries[i].ID)
		_ = id // 实际可查 tag name
		activeTags[i] = entries[i].ID
	}

	tagWeightsJSON, _ := json.Marshal(tagWeights)
	activeTagsJSON, _ := json.Marshal(activeTags)

	profile := &do.UserInterestProfile{
		UserID:          userID,
		TagWeightsJSON:  string(tagWeightsJSON),
		BoardWeightsJSON: "{}",
		ActiveTagsJSON:  string(activeTagsJSON),
		LastUpdatedAt:   now.Format(time.RFC3339),
	}

	return s.repo.UpsertInterestProfile(profile)
}

// --- 辅助函数 ---

func generateSessionID(userID uint) string {
	return "rec_" + strconv.Itoa(int(userID)) + "_" + strconv.FormatInt(time.Now().Unix(), 36)
}

func sliceMapKeys(m map[uint]bool) []uint {
	keys := make([]uint, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

var (
	_ RecommendationService = (*recommendationService)(nil)
	_                       = sync.Mutex{}
	_                       = strings.ReplaceAll
	_                       = zap.String
)
