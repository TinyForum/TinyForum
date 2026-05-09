package storage

import (
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	baseDir string
}

func NewLocalStorage(baseDir string) *LocalStorage {
	os.MkdirAll(baseDir, 0755)
	return &LocalStorage{baseDir: baseDir}
}

func (s *LocalStorage) fullPath(relative string) string {
	return filepath.Join(s.baseDir, relative)
}

func (s *LocalStorage) Save(src io.Reader, destPath string) (int64, error) {
	full := s.fullPath(destPath)
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return 0, err
	}
	dst, err := os.Create(full)
	if err != nil {
		return 0, err
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

func (s *LocalStorage) Get(destPath string) (io.ReadCloser, error) {
	return os.Open(s.fullPath(destPath))
}

func (s *LocalStorage) Delete(destPath string) error {
	return os.Remove(s.fullPath(destPath))
}

func (s *LocalStorage) Exists(destPath string) (bool, error) {
	_, err := os.Stat(s.fullPath(destPath))
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (s *LocalStorage) DeleteDir(dirPath string) error {
	return os.RemoveAll(s.fullPath(dirPath))
}