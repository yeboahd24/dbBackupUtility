package backup

import (
	"context"
	"io"
)

// BackupType represents the type of backup to perform
type BackupType string

const (
	Full         BackupType = "full"
	Incremental  BackupType = "incremental"
	Differential BackupType = "differential"
)

// DatabaseBackuper defines the interface for database backup operations
type DatabaseBackuper interface {
	Connect(ctx context.Context) error
	Backup(ctx context.Context, backupType BackupType) (io.Reader, error)
	Restore(ctx context.Context, backupFile io.Reader) error
	Close() error
}

// StorageProvider defines the interface for backup storage operations
type StorageProvider interface {
	Store(ctx context.Context, name string, data io.Reader) error
	Retrieve(ctx context.Context, name string) (io.ReadCloser, error)
	List(ctx context.Context) ([]string, error)
}

// Compressor defines the interface for backup compression
type Compressor interface {
	Compress(data io.Reader) (io.Reader, error)
	Decompress(data io.Reader) (io.Reader, error)
}

// NotificationService defines the interface for backup notifications
type NotificationService interface {
	Notify(message string) error
}
