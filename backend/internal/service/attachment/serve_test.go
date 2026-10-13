package attachment

import (
	"context"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"tiny-forum/internal/model/do"
	attachmentRepo "tiny-forum/internal/repository/attachment"
	"tiny-forum/internal/storage"
	apperrors "tiny-forum/pkg/errors"

	"gorm.io/gorm"
)

// fakeAttRepo 仅实现测试用到的 GetByFileID。
type fakeAttRepo struct {
	attachmentRepo.AttachmentRepository
	att *do.Attachment
	err error
}

func (f *fakeAttRepo) GetByFileID(ctx context.Context, fileID string) (*do.Attachment, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.att, nil
}

func TestOpenFile_ReturnsContent(t *testing.T) {
	store := storage.NewLocalStorage(t.TempDir())
	rel := filepath.Join("2026", "10", "a.txt")
	if _, err := store.Save(strings.NewReader("hello"), rel); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	svc := &service{
		repo:      &fakeAttRepo{att: &do.Attachment{FileID: "f1", StoredPath: rel, OriginalName: "a.txt", MimeType: "text/plain", Size: 5}},
		storage:   store,
		urlPrefix: "/uploads",
	}

	rc, info, err := svc.OpenFile(context.Background(), "f1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer rc.Close()

	got, _ := io.ReadAll(rc)
	if string(got) != "hello" {
		t.Fatalf("expected content hello, got %q", string(got))
	}
	if info.MimeType != "text/plain" {
		t.Fatalf("expected mime text/plain, got %q", info.MimeType)
	}
}

func TestOpenFile_NotFound(t *testing.T) {
	svc := &service{repo: &fakeAttRepo{err: gorm.ErrRecordNotFound}, storage: storage.NewLocalStorage(t.TempDir())}
	if _, _, err := svc.OpenFile(context.Background(), "missing"); !errors.Is(err, apperrors.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
