package cmd

import (
	"errors"
	"fmt"

	"github.com/kushsharma/servo/internal"
	"github.com/rclone/rclone/fs/operations"
	clonecmd "github.com/rclone/rclone/fs/operations"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ()

func initTest() *cobra.Command {
	return &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("starting test tool...")
			fsrc, srcFileName, fdst := clonecmd.NewFsSrcFileDst(args)
			// appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
			// if !ok {
			// 	return errors.New("unable to find application config")
			// }

			err := operation.CopyFile(context.Background(), fdst, fsrc, srcFileName, srcFileName)

			return err
		},
	}
}
