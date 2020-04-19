package cmd

import (
	"errors"
	"fmt"

	"github.com/kushsharma/servo/backup"
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/sshtunnel"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ()

func initBackup() *cobra.Command {
	return &cobra.Command{
		Use: "backup",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("starting backup tool...")

			appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
			if !ok {
				return errors.New("unable to find application config")
			}

			for _, machine := range appConfig.Machines {
				sshclient, err := sshtunnel.ConnectWithKeyPassphrase(machine.Auth)
				if err != nil {
					return err
				}
				defer sshclient.Close()

				fsService := backup.NewFSService(sshclient, machine.Backup)
				if err := backupFS(fsService); err != nil {
					return err
				}
			}

			return nil
		},
	}
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

	return err
}
