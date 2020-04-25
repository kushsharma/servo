package cmd

import (
	"errors"
	"fmt"

	"github.com/kushsharma/servo/backup"
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/tunnel"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ()

func initBackup() *cobra.Command {
	bcmd := &cobra.Command{
		Use: "backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("starting backup tool...")

			appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
			if !ok {
				return errors.New("unable to find application config")
			}

			for _, machine := range appConfig.Machines {
				fsService := backup.NewFSService(machine.Backup.FS)
				if err := backupFS(fsService); err != nil {
					return err
				}
				fmt.Printf("fs backup completed successfully for %s\n", machine.Name)

				localTnl := tunnel.NewLocalTunnel()
				defer localTnl.Close()
				dbService := backup.NewDBService(localTnl, machine.Backup.DB)
				if err := backupDB(dbService); err != nil {
					return err
				}
				fmt.Printf("db backup completed successfully for %s\n", machine.Name)
			}

			return nil
		},
	}
	return bcmd
}

//backupFS
func backupFS(svc backup.BackupService) error {

	err := svc.Prepare()
	if err != nil {
		return err
	}

	err = svc.Migrate()
	if err != nil {
		return err
	}

	return svc.Close()
}

//backupDB
func backupDB(svc backup.BackupService) error {

	err := svc.Prepare()
	if err != nil {
		return err
	}

	if !DryRun {
		err = svc.Migrate()
		if err != nil {
			return err
		}
	}

	return svc.Close()
}
