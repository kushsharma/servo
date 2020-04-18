package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kushsharma/servo/cmd"
	"github.com/kushsharma/servo/logtool"
	"github.com/spf13/viper"
)

var (
	// Version is app version
	Version = ""
	// Build is build date of this executable
	Build = ""
	// AppName of this executable
	AppName = "servo"

	cfgFile = ""
)

func main() {
	initConfig()

	logToolService := logtool.NewService()

	cmd.InitCommands(logToolService)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		// search config in home directory
		viper.AddConfigPath(home)
		// search in current dir
		viper.AddConfigPath(".")
		viper.SetConfigName("." + AppName) //.servo
	}

	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(strings.ToUpper(AppName))
	viper.AutomaticEnv()

	viper.Set("author", "Kush Kumar Sharma <thekushsharma@gmail.com>")
	viper.Set("appname", AppName)
	viper.Set("version", Version)

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		panic(err)
	}
}
