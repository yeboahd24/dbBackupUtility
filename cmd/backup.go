package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/yeboahd24/dbBackupUitility/pkg/backup"
	"github.com/yeboahd24/dbBackupUitility/pkg/config"
)

func BackupCommand() *cli.Command {
	return &cli.Command{
		Name:  "backup",
		Usage: "Perform database backup",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Path to config file (optional, will auto-detect if not provided)",
				Required: false,
			},
			&cli.StringFlag{
				Name:    "type",
				Aliases: []string{"t"},
				Usage:   "Backup type (full, incremental, differential)",
				Value:   "full",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "Output file path (required when storage is disabled)",
			},
		},
		Action: func(c *cli.Context) error {
			ctx := context.Background()
			cfg, err := config.LoadConfig(c.String("config"))
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Initialize database backuper
			backuper := backup.NewPostgresBackup(cfg.Database)
			if err := backuper.Connect(ctx); err != nil {
				return err
			}
			defer backuper.Close()

			// Perform backup
			backupType := backup.BackupType(c.String("type"))
			reader, err := backuper.Backup(ctx, backupType)
			if err != nil {
				return err
			}

			// Handle backup storage
			if cfg.Storage.Enabled {
				// Initialize remote storage
				storage, err := initializeStorage(cfg.Storage)
				if err != nil {
					return err
				}

				// Store in remote storage
				filename := fmt.Sprintf("backup_%s_%s.dump",
					cfg.Database.Database,
					time.Now().Format("20060102150405"))

				return storage.Store(ctx, filename, reader)
			} else {
				// Store locally if output path is provided
				outputPath := c.String("output")
				if outputPath == "" {
					return fmt.Errorf("output path is required when storage is disabled")
				}

				// Create output file
				file, err := os.Create(outputPath)
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}
				defer file.Close()

				// Copy backup data to file
				_, err = io.Copy(file, reader)
				if err != nil {
					return fmt.Errorf("failed to write backup to file: %w", err)
				}

				fmt.Printf("Backup saved to: %s\n", outputPath)
				return nil
			}
		},
	}
}
