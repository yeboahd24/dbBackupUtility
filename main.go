package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yeboahd24/dbBackupUitility/cmd"
)

func main() {
	app := &cli.App{
		Name:  "dbbackup",
		Usage: "A versatile database backup utility",
		Commands: []*cli.Command{
			cmd.BackupCommand(),
			cmd.RestoreCommand(),
			cmd.ConfigCommand(),
			cmd.HelpCommand(),
		},
		// Add global flags if needed
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Usage:   "Enable verbose output",
			},
		},
		// Add default help text
		Description: `Database Backup Utility provides a robust solution for managing database backups.

COMMANDS:
   backup   Perform database backup
   restore  Restore database from backup
   config   Manage configuration settings
   help     Shows detailed help information for commands

Run 'dbbackup help <command>' for more information about a command.`,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
