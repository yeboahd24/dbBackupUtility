package cmd

import (
	"fmt"

	"github.com/urfave/cli/v2"
	"github.com/yeboahd24/dbBackupUitility/pkg/config"
)

func ConfigCommand() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "Manage configuration settings",
		Subcommands: []*cli.Command{
			{
				Name:  "validate",
				Usage: "Validate configuration file",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Aliases:  []string{"c"},
						Usage:    "Path to config file (optional, will auto-detect if not provided)",
						Required: false,
					},
				},
				Action: func(c *cli.Context) error {
					cfg, err := config.LoadConfig(c.String("config"))
					if err != nil {
						return fmt.Errorf("configuration validation failed: %w", err)
					}

					// Validate database configuration
					if cfg.Database.Type == "" {
						return fmt.Errorf("database type is required")
					}
					if cfg.Database.Host == "" {
						return fmt.Errorf("database host is required")
					}
					if cfg.Database.Port == 0 {
						return fmt.Errorf("database port is required")
					}
					if cfg.Database.Username == "" {
						return fmt.Errorf("database username is required")
					}
					if cfg.Database.Database == "" {
						return fmt.Errorf("database name is required")
					}

					// Validate storage configuration
					if cfg.Storage.Type == "" {
						return fmt.Errorf("storage type is required")
					}
					switch cfg.Storage.Type {
					case "local":
						if cfg.Storage.Path == "" {
							return fmt.Errorf("storage path is required for local storage")
						}
					case "s3":
						if cfg.Storage.Bucket == "" {
							return fmt.Errorf("storage bucket is required for S3 storage")
						}
						if cfg.Storage.Region == "" {
							return fmt.Errorf("storage region is required for S3 storage")
						}
					default:
						return fmt.Errorf("unsupported storage type: %s", cfg.Storage.Type)
					}

					// If notification is enabled, validate webhook URL
					if cfg.Notification.Enabled && cfg.Notification.SlackWebhook == "" {
						return fmt.Errorf("slack webhook URL is required when notifications are enabled")
					}

					fmt.Println("Configuration file is valid")
					return nil
				},
			},
		},
	}
}
