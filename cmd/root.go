package cmd

import (
	"errors"

	// import fs backends
	_ "github.com/rclone/rclone/backend/alias"
	_ "github.com/rclone/rclone/backend/cache"
	_ "github.com/rclone/rclone/backend/local"
	_ "github.com/rclone/rclone/backend/s3"

	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/tunnel"
	rfs "github.com/rclone/rclone/fs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Debug  = false
	DryRun = false
)

// InitCommands initializes application cli interface
func InitCommands() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     viper.GetString("appname"),
		Version: viper.GetString("version"),
	}
	rootCmd.AddCommand(initBackup())
	rootCmd.AddCommand(initDelete())
	rootCmd.AddCommand(initVersion())
	rootCmd.AddCommand(initTest())

	rootCmd.PersistentFlags().BoolVarP(&DryRun, "dry-run", "d", false, "does not actually perform the action")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "verbose", "v", false, "debug level logs")
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if Debug {
			log.SetLevel(log.DebugLevel)
		}
		return nil
	}

	if err := overrideConfigCallback(); err != nil {
		panic(err)
	}
	return rootCmd
}

func createTunnel(machine internal.MachineConfig) (tunnel.Executioner, error) {
	if machine.ConnectionType == "remote" {
		return tunnel.NewSSHTunnel(machine.Auth)
	}

	return tunnel.NewLocalTunnel(), nil
}

// overrideConfigCallback will provide custom config getter for rclone
func overrideConfigCallback() error {
	appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
	if !ok {
		return errors.New("unable to find application config")
	}

	rfs.ConfigFileGet = func(section, key string) (string, bool) {
		if section == "local" {
			switch key {
			case "type":
				return "local", true
			}
		} else if section == "s3" {
			switch key {
			case "type":
				return "s3", true
			case "provider":
				return "Digital Ocean", true
			case "access_key_id":
				return appConfig.S3.Key, true
			case "secret_access_key":
				return appConfig.S3.Secret, true
			case "endpoint":
				return appConfig.S3.Endpoint, true
			case "acl":
				return "private", true
			}
		}

		return "", false
	}

	return nil
}
