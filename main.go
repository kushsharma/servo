package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kushsharma/servo/cmd"
	"github.com/kushsharma/servo/internal"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	// Version is app version
	Version = ""
	// Build is build date of this executable
	Build = ""
	// AppName of this executable
	AppName = "servo"

	// TODO: provide option to read this from cli flag
	cfgFile = ""

	// AppConfig stores all the application specific configuration
	// required for various auth and actions
	AppConfig internal.ApplicationConfig
)

func main() {
	initConfig()
	cmd.InitCommands()
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

		viper.SetConfigName("." + AppName) //.servo

		// search in current dir
		viper.AddConfigPath(".")
		// search config in home directory
		viper.AddConfigPath(home)
	}

	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(strings.ToUpper(AppName))
	viper.AutomaticEnv()

	viper.Set("author", "Kush Kumar Sharma <thekushsharma@gmail.com>")
	viper.Set("appname", AppName)
	viper.Set("version", Version)

	if err := viper.ReadInConfig(); err == nil {
		configFilePath := viper.ConfigFileUsed()
		fmt.Println("Using config file:", configFilePath)

		// ready yaml seperately because viper parsing sucks
		configByte, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			panic(err)
		}
		if err := yaml.Unmarshal(configByte, &AppConfig); err != nil {
			panic(err)
		}
		viper.Set("app", AppConfig)
	} else {
		panic(err)
	}
}
