package cmd

import (
	"bytes"
	"context"
	"fmt"

	rcmd "github.com/rclone/rclone/cmd"
	rops "github.com/rclone/rclone/fs/operations"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//localfs:./version s3DO:moonwaretech/temp/ --ignore-existing

var ()

func initTest() *cobra.Command {
	return &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("starting test tool...")

			fsrc := rcmd.NewFsSrc([]string{"localfs:./"})
			var stdout bytes.Buffer
			err := rops.ListLong(context.Background(), fsrc, &stdout)
			log.Info(stdout.String())

			return err
		},
	}
}
