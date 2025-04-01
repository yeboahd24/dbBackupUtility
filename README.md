# Database Backup Utility

A robust command-line utility for managing database backups with support for PostgreSQL and MySQL databases. Features include local and S3 storage options, backup types (full, incremental, differential), and Slack notifications.

Please note that this utility is intended for use in controlled environments and should be adapted to fit specific security and operational requirements.

Roadmap Project Link: https://roadmap.sh/projects/database-backup-utility

## Features

- Support for PostgreSQL and MySQL databases
- Local and Amazon S3 storage options
- Multiple backup types (full, incremental, differential)
- Slack notifications for backup status
- Configurable through YAML files
- Auto-detection of config file location
- Comprehensive error handling and logging

## Prerequisites

- Go 1.24 or higher
- PostgreSQL client tools (`pg_dump`, `pg_restore`) for PostgreSQL backups
- MySQL client tools (`mysqldump`) for MySQL backups
- AWS credentials (if using S3 storage)

## Installation

```bash
# Clone the repository
git clone https://github.com/yeboahd24/dbBackupUtility.git

cd dbBackupUitility

# Build the binary
go build -o dbbackup main.go
```

## Configuration

Create a `config.yml` file in one of the following locations:
- Current directory
- `$HOME/.dbbackup/`
- `/etc/dbbackup/`
- `$XDG_CONFIG_HOME/dbbackup/`

Example configuration:

```yaml
database:
  type: postgres  # or mysql
  host: localhost
  port: 5432
  username: postgres
  password: "your_password"
  database: your_database

storage:
  enabled: false  # Set to true for S3 storage
  type: s3
  bucket: your-backup-bucket
  region: us-west-2

notification:
  slack_webhook: https://hooks.slack.com/services/xxx/yyy/zzz
  enabled: false
```

## Usage

### Backup Database

```bash
# Basic backup with local storage
./dbbackup backup --type full --output backup.dump

# Backup with custom config file
./dbbackup backup -c /path/to/config.yml -t full -o backup.dump

# Backup to S3 (when storage.enabled is true)
./dbbackup backup --type full
```

### Restore Database

```bash
# Restore from local backup file
./dbbackup restore --file backup.dump

# Restore with custom config
./dbbackup restore -c /path/to/config.yml -f backup.dump

# Restore from S3 (when storage.enabled is true)
./dbbackup restore --file backup_name.dump
```

### Validate Configuration

```bash
# Validate default config file
./dbbackup config validate

# Validate specific config file
./dbbackup config validate -c /path/to/config.yml
```

## AWS Configuration

When using S3 storage, configure AWS credentials using one of these methods:

1. Environment variables:
```bash
export AWS_ACCESS_KEY_ID="your_access_key"
export AWS_SECRET_ACCESS_KEY="your_secret_key"
export AWS_REGION="your_region"
```

2. AWS credentials file (`~/.aws/credentials`):
```ini
[default]
aws_access_key_id = your_access_key
aws_secret_access_key = your_secret_key
```

## Backup Types

- `full`: Complete backup of the database
- `incremental`: Backup of changes since the last backup (if supported by the database)
- `differential`: Backup of changes since the last full backup (if supported by the database)

## Security Considerations

1. Use environment variables or secure secret management for sensitive credentials
2. Ensure backup files are stored securely
3. Use strong database passwords
4. Implement proper AWS IAM policies when using S3
5. Regularly rotate credentials
6. Enable encryption for backups in transit and at rest

## Error Handling

The utility provides detailed error messages for common issues:

- Database connection failures
- Storage access issues
- Configuration problems
- Backup/restore operation failures

## Best Practices

1. Regularly test restore procedures
2. Monitor backup success/failure through Slack notifications
3. Implement backup rotation/retention policies
4. Keep the utility and dependencies updated
5. Maintain proper documentation of backup procedures
6. Regular validation of backup integrity

## Development

```bash
# Run tests
go test ./...

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o dbbackup-linux-amd64 main.go
GOOS=darwin GOARCH=amd64 go build -o dbbackup-darwin-amd64 main.go
GOOS=windows GOARCH=amd64 go build -o dbbackup-windows-amd64.exe main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request



## Support

For issues and feature requests, please create an issue in the GitHub repository.
