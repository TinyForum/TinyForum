package article

import (
	"errors"
	"testing"

	"tiny-forum/internal/model/do"
	postRepo "tiny-forum/internal/repository/article"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

// fakeArticleRepo 仅实现测试用到的 FindByArticleID / Update。
type fakeArticleRepo struct {
	postRepo.ArticleRepository
	post    *do.Article
	findErr error
	updated *do.Article
}

func (f *fakeArticleRepo) FindByArticleID(id uint) (*do.Article, error) {
	if f.findErr != nil {
		return nil, f.findErr
	}
	return f.post, nil
}

func (f *fakeArticleRepo) Update(p *do.Article) error {
	f.updated = p
	return nil
}

func TestSetStatus_InvalidStatus(t *testing.T) {
	svc := &articleService{postRepo: &fakeArticleRepo{}}
	if err := svc.SetStatus(1, do.CreationStatus("bogus")); !errors.Is(err, apperrors.ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest, got %v", err)
	}
}

func TestSetStatus_NotFound(t *testing.T) {
	svc := &articleService{postRepo: &fakeArticleRepo{findErr: gorm.ErrRecordNotFound}}
	if err := svc.SetStatus(1, do.CreationStatusHidden); !errors.Is(err, apperrors.ErrPostNotFound) {
		t.Fatalf("expected ErrPostNotFound, got %v", err)
	}
}

func TestSetStatus_PersistsStatus(t *testing.T) {
	repo := &fakeArticleRepo{post: &do.Article{Creation: do.Creation{CreationStatus: do.CreationStatusDraft}}}
	svc := &articleService{postRepo: repo}

	if err := svc.SetStatus(1, do.CreationStatusHidden); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.post.Creation.CreationStatus != do.CreationStatusHidden {
		t.Fatalf("expected status hidden, got %s", repo.post.Creation.CreationStatus)
	}
	if repo.updated != repo.post {
		t.Fatalf("expected Update to be called with mutated post")
	}
}
