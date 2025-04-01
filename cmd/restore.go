package cmd

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yeboahd24/dbBackupUitility/pkg/backup"
	"github.com/yeboahd24/dbBackupUitility/pkg/config"
	"github.com/yeboahd24/dbBackupUitility/pkg/storage"
)

func initializeStorage(cfg config.StorageConfig) (backup.StorageProvider, error) {
	switch cfg.Type {
	case "local":
		return storage.NewLocalStorage(cfg.Path)
	case "s3":
		ctx := context.Background()
		return storage.NewS3Storage(ctx, cfg.Bucket, cfg.Region)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}

func RestoreCommand() *cli.Command {
	return &cli.Command{
		Name:  "restore",
		Usage: "Restore database from backup",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Path to config file (optional, will auto-detect if not provided)",
				Required: false,
			},
			&cli.StringFlag{
				Name:     "file",
				Aliases:  []string{"f"},
				Usage:    "Backup file to restore",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			ctx := context.Background()
			cfg, err := config.LoadConfig(c.String("config"))
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Initialize database backuper based on type
			var backuper backup.DatabaseBackuper
			switch cfg.Database.Type {
			case "postgres":
				backuper = backup.NewPostgresBackup(cfg.Database)
			case "mysql":
				backuper = backup.NewMySQLBackup(cfg.Database)
			default:
				return fmt.Errorf("unsupported database type: %s", cfg.Database.Type)
			}

			var reader io.ReadCloser
			backupFile := c.String("file")

			if cfg.Storage.Enabled {
				// Get from remote storage
				storage, err := initializeStorage(cfg.Storage)
				if err != nil {
					return fmt.Errorf("failed to initialize storage: %w", err)
				}

				reader, err = storage.Retrieve(ctx, backupFile)
				if err != nil {
					return fmt.Errorf("failed to retrieve backup file: %w", err)
				}
			} else {
				// Open local file
				file, err := os.Open(backupFile)
				if err != nil {
					return fmt.Errorf("failed to open backup file: %w", err)
				}
				reader = file
			}
			defer reader.Close()

			// Perform restore
			if err := backuper.Connect(ctx); err != nil {
				return fmt.Errorf("failed to connect to database: %w", err)
			}
			defer backuper.Close()

			if err := backuper.Restore(ctx, reader); err != nil {
				return fmt.Errorf("failed to restore backup: %w", err)
			}

			fmt.Println("Database restored successfully")
			return nil
		},
	}
}
