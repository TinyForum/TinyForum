package storage

import (
	"io"
)

type StorageDriver interface {
	Save(src io.Reader, destPath string) (int64, error)
	Get(destPath string) (io.ReadCloser, error)
	Delete(destPath string) error
	Exists(destPath string) (bool, error)
	DeleteDir(dirPath string) error // 递归删除目录
}
