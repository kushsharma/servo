package cmd

import (
	"errors"

	"github.com/kushsharma/servo/backup"
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/tunnel"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ()

func initBackup() *cobra.Command {
	bcmd := &cobra.Command{
		Use:  "backup",
		Long: "backup file system and database provided in config",
	}
	bcmd.AddCommand(subBackupFS())
	bcmd.AddCommand(subBackupDB())
	return bcmd
}

func subBackupDB() *cobra.Command {
	return &cobra.Command{
		Use:     "db",
		Example: "servo backup db",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("starting db backup tool...")

			appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
			if !ok {
				return errors.New("unable to find application config")
			}

			for _, machine := range appConfig.Machines {
				localTnl := tunnel.NewLocalTunnel()
				defer localTnl.Close()
				dbService := backup.NewDBService(localTnl, machine.Backup.DB)
				if err := backupDB(dbService); err != nil {
					return err
				}
				log.Infof("db backup completed successfully for %s", machine.Name)
			}

			return nil
		},
	}
}

func subBackupFS() *cobra.Command {
	return &cobra.Command{
		Use:     "fs",
		Example: "servo backup fs --dry-run --verbose",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("starting fs backup tool...")

			appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
			if !ok {
				return errors.New("unable to find application config")
			}

			for _, machine := range appConfig.Machines {
				fsService := backup.NewFSService(machine.Backup.FS)
				if err := backupFS(fsService); err != nil {
					return err
				}
				log.Infof("fs backup completed successfully for %s", machine.Name)
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
