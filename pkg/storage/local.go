package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	basePath string
}

func NewLocalStorage(path string) (*LocalStorage, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{basePath: path}, nil
}

func (l *LocalStorage) Store(ctx context.Context, name string, data io.Reader) error {
	path := filepath.Join(l.basePath, name)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}
	return nil
}

func (l *LocalStorage) Retrieve(ctx context.Context, name string) (io.ReadCloser, error) {
	path := filepath.Join(l.basePath, name)
	return os.Open(path)
}

func (l *LocalStorage) List(ctx context.Context) ([]string, error) {
	var files []string
	err := filepath.Walk(l.basePath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	return files, err
}
