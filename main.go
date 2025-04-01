package main

import (
    "log"
    "os"

    "github.com/yeboahd24/dbBackupUitility/cmd"
    "github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
        Name:  "dbbackup",
        Usage: "A versatile database backup utility",
        Commands: []*cli.Command{
            cmd.BackupCommand(),
            cmd.RestoreCommand(),
            cmd.ConfigCommand(),
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}