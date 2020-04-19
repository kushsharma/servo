package cmd

import (
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/tunnel"
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
	rootCmd.AddCommand(initBackup())
	rootCmd.AddCommand(initVersion())
	rootCmd.Execute()
}

func createTunnel(machine internal.MachineConfig) tunnel.Executioner {

	return nil
}
