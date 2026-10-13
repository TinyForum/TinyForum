package comment

import (
	"errors"
	"testing"

	"tiny-forum/internal/model/bo"
	"tiny-forum/internal/model/do"
	postRepo "tiny-forum/internal/repository/article"
	apperrors "tiny-forum/pkg/errors"
)

// fakeArticleRepo 仅实现测试用到的 FindByArticleID，其余方法由内嵌接口兜底。
type fakeArticleRepo struct {
	postRepo.ArticleRepository
	post *do.Article
}

func (f *fakeArticleRepo) FindByArticleID(id uint) (*do.Article, error) {
	return f.post, nil
}

// TestCreateComment_RejectsLockedPost 锁定帖子必须拒绝新增评论。
func TestCreateComment_RejectsLockedPost(t *testing.T) {
	svc := &commentService{
		postRepo: &fakeArticleRepo{post: &do.Article{Creation: do.Creation{Locked: true}}},
	}

	_, err := svc.CreateComment(1, bo.CreateCommentInput{PostID: 1, Content: "hello"})
	if !errors.Is(err, apperrors.ErrPostLocked) {
		t.Fatalf("expected ErrPostLocked, got %v", err)
	}
}
