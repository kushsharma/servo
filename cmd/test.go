package cmd

import (
	"errors"

	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/mailrelay"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//localfs:./version s3DO:moonwaretech/temp/ --ignore-existing

func initTest() *cobra.Command {
	return &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Info("starting test tool...")

			appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
			if !ok {
				return errors.New("unable to find application config")
			}

			err = mailrelay.Start(appConfig.Remotes)
			if err != nil {
				return err
			}

			return err
		},
	}
}
