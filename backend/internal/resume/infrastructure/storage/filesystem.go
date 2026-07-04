// Package storage persists raw resume bytes. The filesystem implementation is
// for local/MVP use; swap for an S3 adapter in production (implement
// resume/domain.Storage).
package storage

import (
	"context"
	"os"
	"path/filepath"
)

type FileSystem struct {
	root string
}

// NewFileSystem stores files under RESUME_UPLOAD_DIR (default ./uploads/resumes).
func NewFileSystem() *FileSystem {
	root := os.Getenv("RESUME_UPLOAD_DIR")
	if root == "" {
		root = filepath.Join("uploads", "resumes")
	}
	return &FileSystem{root: root}
}

func (fs *FileSystem) Save(_ context.Context, key string, data []byte) error {
	path := filepath.Join(fs.root, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
