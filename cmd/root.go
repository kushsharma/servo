package cmd

import (
	"github.com/kushsharma/servo/logtool"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// InitCommands initializes application cli interface
func InitCommands(ltService logtool.LogManager) {
	rootCmd := &cobra.Command{
		Use:     viper.GetString("appname"),
		Version: viper.GetString("version"),
	}

	rootCmd.AddCommand(initLog(ltService))
	rootCmd.AddCommand(initVersion())
	rootCmd.Execute()
}
