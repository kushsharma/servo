package cmd

import (
	"errors"
	"fmt"

	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/logtool"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	ErrInvalidSSHMachineConfig = errors.New("invalid machine configs provided")
)

func initLog() *cobra.Command {
	return &cobra.Command{
		Use: "log",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("starting logtool...")

			appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
			if !ok {
				return ErrInvalidSSHMachineConfig
			}

			for _, machine := range appConfig.Machines {
				tnl, err := createTunnel(machine)
				if err != nil {
					return err
				}
				logToolService := logtool.NewService(tnl)
				if err := logClean(logToolService, machine.Clean); err != nil {
					return err
				}
			}

			return nil
		},
	}
}

//logClean remove files that are unnecessary and older than x days
func logClean(svc logtool.LogManager, config internal.CleanConfig) error {

	//files to be cleaned older than x days
	daysOld := config.OlderThan

	errs := []error{}
	for _, path := range config.Path {
		files, err := svc.DryClean(path, daysOld)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		fmt.Print(files)
	}

	return internal.ErrMerge(errs)
}
