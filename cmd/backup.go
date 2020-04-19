package cmd

import (
	"errors"
	"fmt"

	"github.com/kushsharma/servo/backup"
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/sshtunnel"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

			s3Config := &aws.Config{
				Credentials: credentials.NewStaticCredentials(appConfig.S3.Key, appConfig.S3.Secret, ""),
				Endpoint:    aws.String("https://moonware.nyc3.digitaloceanspaces.com"),
				Region:      aws.String("us-east-1"),
			}
			// Create S3 service client
			newSession := session.New(s3Config)
			s3Client := s3.New(newSession)

			for _, machine := range appConfig.Machines {
				sshclient, err := sshtunnel.ConnectWithKeyPassphrase(machine.Auth)
				if err != nil {
					return err
				}
				defer sshclient.Close()

				fsService := backup.NewFSService(sshclient, s3Client, machine.Backup)
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
