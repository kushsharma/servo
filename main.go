package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/kushsharma/servo/cmd"
	"github.com/kushsharma/servo/internal"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var (
	// Version is app version
	Version = "0"
	// Build is build date of this executable
	Build = "0"
	// AppName of this executable
	AppName = "servo"

	// TODO: provide option to read this from cli flag
	cfgFile = ""

	// AppConfig stores all the application specific configuration
	// required for various auth and actions
	AppConfig internal.ApplicationConfig

	//LogFilePath where output of this application is written into
	LogFilePath = AppName + ".log"
)

func main() {
	// Output to stdout instead of the default stderr
	// Can be any io.Writer

	logFile, err := os.OpenFile(LogFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file for logging: %v", err)
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	initConfig()
	internal.InitStat(Version)
	rootCmd := cmd.InitCommands()
	rootCmd.Execute()
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			log.Panic("unable to find home directory")
			panic(err)
		}

		viper.SetConfigName("." + AppName) //.servo

		// search in current dir
		viper.AddConfigPath(".")
		// search config in home directory
		viper.AddConfigPath(home)
		// search config in ~/.config/servo/.servo
		viper.AddConfigPath(filepath.Join(home, ".config", "servo"))
	}

	viper.SetConfigType("yaml")
	viper.SetEnvPrefix(strings.ToUpper(AppName))
	viper.AutomaticEnv()

	viper.Set("author", "Kush Kumar Sharma <thekushsharma@gmail.com>")
	viper.Set("appname", AppName)
	viper.Set("version", Version)

	if err := viper.ReadInConfig(); err == nil {
		configFilePath := viper.ConfigFileUsed()
		log.Debugf("using config file: %s", configFilePath)

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
