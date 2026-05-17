package upload

import (
	"crypto/sha256"
	"io"
	"path/filepath"
	"strings"
	"tiny-forum/internal/model/do"

	"github.com/google/uuid"
)

func computeHash(r io.Reader) ([]byte, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func generateStoredName(original string) string {
	ext := strings.ToLower(filepath.Ext(original))
	return uuid.New().String() + ext
}

func extractMimeMajor(mime string) do.MimeTypeMajor {
	switch {
	case strings.HasPrefix(mime, "image/"):
		return do.MimeImage
	case strings.HasPrefix(mime, "video/"):
		return do.MimeVideo
	case strings.HasPrefix(mime, "audio/"):
		return do.MimeAudio
	case strings.Contains(mime, "pdf") || strings.Contains(mime, "document") || strings.Contains(mime, "text"):
		return do.MimeDocument
	default:
		return do.MimeOther
	}
}
