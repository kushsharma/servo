package cmd

import (
	"context"
	"fmt"

	_ "github.com/rclone/rclone/backend/all" // import all backends
	rcmd "github.com/rclone/rclone/cmd"
	rconfig "github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/operations"
	"github.com/spf13/cobra"
)

var ()

func initTest() *cobra.Command {
	return &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("starting test tool...")

			rconfig.LoadConfig()
			fsrc, srcFileName, fdst := rcmd.NewFsSrcFileDst(args)
			// appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
			// if !ok {
			// 	return errors.New("unable to find application config")
			// }

			err := operations.CopyFile(context.Background(), fdst, fsrc, srcFileName, srcFileName)

			return err
		},
	}
}
