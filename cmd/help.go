package cmd

import (
    "fmt"

    "github.com/urfave/cli/v2"
)

func HelpCommand() *cli.Command {
    return &cli.Command{
        Name:    "help",
        Aliases: []string{"h"},
        Usage:   "Shows detailed help information for commands",
        Subcommands: []*cli.Command{
            {
                Name:  "backup",
                Usage: "Show detailed help for backup command",
                Action: func(c *cli.Context) error {
                    fmt.Println(`
BACKUP COMMAND
-------------
Performs database backup operations with various options.

Usage:
  dbbackup backup [options]

Options:
  --type, -t     Backup type (full, incremental, differential) (default: "full")
  --output, -o   Output file path for local storage
  --config, -c   Path to config file (optional)

Examples:
  1. Local backup:
     dbbackup backup --type full --output backup.dump

  2. S3 backup (when storage.enabled is true):
     dbbackup backup --type full

  3. Custom config:
     dbbackup backup -c /path/to/config.yml -t full -o backup.dump

Notes:
  - For S3 storage, ensure AWS credentials are properly configured
  - Incremental and differential backups depend on database support
  - Output path is required when storage.enabled is false
`)
                    return nil
                },
            },
            {
                Name:  "restore",
                Usage: "Show detailed help for restore command",
                Action: func(c *cli.Context) error {
                    fmt.Println(`
RESTORE COMMAND
--------------
Restores database from a backup file.

Usage:
  dbbackup restore [options]

Options:
  --file, -f     Backup file to restore from
  --config, -c   Path to config file (optional)

Examples:
  1. Local restore:
     dbbackup restore --file backup.dump

  2. S3 restore (when storage.enabled is true):
     dbbackup restore --file backup_name.dump

  3. Custom config:
     dbbackup restore -c /path/to/config.yml -f backup.dump

Notes:
  - Ensure target database exists and is accessible
  - User must have sufficient privileges for restore operation
  - For S3 restores, ensure AWS credentials are properly configured
`)
                    return nil
                },
            },
            {
                Name:  "config",
                Usage: "Show detailed help for config command",
                Action: func(c *cli.Context) error {
                    fmt.Println(`
CONFIG COMMAND
-------------
Manages and validates configuration settings.

Usage:
  dbbackup config validate [options]

Options:
  --config, -c   Path to config file (optional)

Examples:
  1. Validate default config:
     dbbackup config validate

  2. Validate specific config:
     dbbackup config validate -c /path/to/config.yml

Config File Locations:
  The utility searches for config.yml in these locations:
  1. Current directory
  2. $HOME/.dbbackup/
  3. /etc/dbbackup/
  4. $XDG_CONFIG_HOME/dbbackup/

Required Configuration:
  database:
    type: postgres|mysql
    host: <hostname>
    port: <port>
    username: <username>
    password: <password>
    database: <dbname>

  storage:
    enabled: true|false
    type: local|s3
    bucket: <bucket-name>    # for S3
    region: <region>         # for S3
    path: <local-path>       # for local

  notification:
    slack_webhook: <webhook-url>
    enabled: true|false
`)
                    return nil
                },
            },
        },
    }
}