package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// InitCommands initializes application cli interface
func InitCommands() {
	rootCmd := &cobra.Command{
		Use:     viper.GetString("appname"),
		Version: viper.GetString("version"),
	}

	rootCmd.AddCommand(initLog())
	rootCmd.AddCommand(initVersion())
	rootCmd.Execute()
}
