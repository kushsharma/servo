package cmd

import (
	"errors"
	"fmt"

	"github.com/kushsharma/servo/logtool"
	"github.com/kushsharma/servo/sshtunnel"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrInvalidSSHMachineConfig = errors.New("invalid ssh machine configs provided")
)

func initLog(ltService logtool.LogManager) *cobra.Command {
	return &cobra.Command{
		Use: "log",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("logtool initialized for " + viper.GetString("config.schedule"))

			//fetch logs
			// logs, err := ltService.List("/Users/rick/Dev/web/hire/site/laravel/storage/logs")
			// if err != nil {
			// 	return err
			// }

			machinesRaw, ok := viper.GetStringMap("ssh")["machines"]
			if !ok {
				return errors.New("no ssh configs provided")
			}
			machines, ok := machinesRaw.([]interface{})
			if !ok {
				return ErrInvalidSSHMachineConfig
			}

			sshConfig := &sshtunnel.SSHAuthConfig{}

			for _, machineRaw := range machines {
				machine, ok := machineRaw.(map[interface{}]interface{})
				if !ok {
					return ErrInvalidSSHMachineConfig
				}

				user, ok := machine["user"]
				if ok {
					sshConfig.User = user.(string)
				}
				addr, ok := machine["address"]
				if ok {
					sshConfig.Address = addr.(string)
				}
				key, ok := machine["key"]
				if ok {
					sshConfig.KeyFile = key.(string)
				}
				pass, ok := machine["password"]
				if ok {
					sshConfig.KeyPassword = pass.(string)
				}

				//DEBUG
				break
			}

			fmt.Println(sshConfig)
			sshclient, err := sshtunnel.ConnectWithKeyPassphrase(sshConfig)
			if err != nil {
				return err
			}
			rscript := sshclient.Cmd("whoami")
			if str, err := rscript.Output(); err == nil {
				fmt.Println("output from ssh: " + string(str))
				//}
				//if err := rscript.Run(); err == nil {
			} else {
				return err
			}

			sshclient.Close()
			return nil
		},
	}
}
