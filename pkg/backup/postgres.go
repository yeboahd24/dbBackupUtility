package backup

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"

	_ "github.com/lib/pq"
	"github.com/yeboahd24/dbBackupUitility/pkg/config"
)

type PostgresBackup struct {
	config config.DatabaseConfig
	db     *sql.DB
}

func NewPostgresBackup(config config.DatabaseConfig) *PostgresBackup {
	return &PostgresBackup{config: config}
}

func (p *PostgresBackup) Connect(ctx context.Context) error {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		p.config.Host,
		p.config.Port,
		p.config.Username,
		p.config.Password,
		p.config.Database,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	p.db = db
	return nil
}

func (p *PostgresBackup) Backup(ctx context.Context, backupType BackupType) (io.Reader, error) {
	cmd := exec.CommandContext(ctx, "pg_dump",
		"-h", p.config.Host,
		"-p", fmt.Sprintf("%d", p.config.Port),
		"-U", p.config.Username,
		"-d", p.config.Database,
		"-F", "c", // Use custom format
	)

	// Set PGPASSWORD environment variable
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", p.config.Password))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("backup failed: %w", err)
	}

	return io.NopCloser(bytes.NewReader(output)), nil
}

func (p *PostgresBackup) Restore(ctx context.Context, backupFile io.Reader) error {
	// Create a temporary file to store the backup
	tmpFile, err := os.CreateTemp("", "postgres-backup-*.dump")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	// Copy backup data to temp file
	if _, err := io.Copy(tmpFile, backupFile); err != nil {
		return fmt.Errorf("failed to write backup to temp file: %w", err)
	}

	cmd := exec.CommandContext(ctx, "pg_restore",
		"-h", p.config.Host,
		"-p", fmt.Sprintf("%d", p.config.Port),
		"-U", p.config.Username,
		"-d", p.config.Database,
		"-c",      // Clean (drop) database objects before recreating
		"-F", "c", // Custom format
		tmpFile.Name(),
	)

	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", p.config.Password))

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("restore failed: %s: %w", string(output), err)
	}

	return nil
}

func (p *PostgresBackup) Close() error {
	if p.db != nil {
		return p.db.Close()
	}
	return nil
}
