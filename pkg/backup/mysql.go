package backup

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"

	_ "github.com/go-sql-driver/mysql"
	"github.com/yeboahd24/dbBackupUitility/pkg/config"
)

type MySQLBackup struct {
	config config.DatabaseConfig
	db     *sql.DB
}

func NewMySQLBackup(config config.DatabaseConfig) *MySQLBackup {
	return &MySQLBackup{config: config}
}

func (m *MySQLBackup) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		m.config.Username,
		m.config.Password,
		m.config.Host,
		m.config.Port,
		m.config.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping MySQL: %w", err)
	}

	m.db = db
	return nil
}

func (m *MySQLBackup) Backup(ctx context.Context, backupType BackupType) (io.Reader, error) {
	cmd := exec.CommandContext(ctx, "mysqldump",
		"-h", m.config.Host,
		"-P", fmt.Sprintf("%d", m.config.Port),
		"-u", m.config.Username,
		"-p"+m.config.Password,
		m.config.Database,
	)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("backup failed: %w", err)
	}

	return io.NopCloser(bytes.NewReader(output)), nil
}

func (m *MySQLBackup) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

func (m *MySQLBackup) Restore(ctx context.Context, backupFile io.Reader) error {
	// Create a temporary file to store the backup
	tmpFile, err := os.CreateTemp("", "mysql-backup-*.sql")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Copy backup data to temp file
	if _, err := io.Copy(tmpFile, backupFile); err != nil {
		return fmt.Errorf("failed to write backup to temp file: %w", err)
	}

	// Execute mysql command to restore
	cmd := exec.CommandContext(ctx, "mysql",
		"-h", m.config.Host,
		"-P", fmt.Sprintf("%d", m.config.Port),
		"-u", m.config.Username,
		"-p"+m.config.Password,
		m.config.Database,
	)

	cmd.Stdin = tmpFile
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("restore failed: %s: %w", string(output), err)
	}

	return nil
}
