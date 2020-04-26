package cmd

import (
	"context"
	"time"

	robscure "github.com/rclone/rclone/fs/config/obscure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//localfs:./version s3DO:moonwaretech/temp/ --ignore-existing

var ()

func initTest() *cobra.Command {
	return &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Info("starting test tool...")
			_, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			defer cancel()

			str, err := robscure.Reveal("gWAnMklpEUt407FzwqgOX9tzeR2b")
			log.Info(str)

			return err
		},
	}
}
