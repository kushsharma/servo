package internal

import "github.com/kushsharma/servo/sshtunnel"

// ApplicationConfig is populated from yaml config file
type ApplicationConfig struct {
	Machines []MachineConfig `yaml:"machines"`
	S3       S3Config        `yaml:"s3"`
}

type MachineConfig struct {
	Name   string
	Auth   sshtunnel.SSHAuthConfig `yaml:"auth"`
	Clean  CleanConfig             `yaml:"clean"`
	Backup BackupConfig            `yaml:"backup"`
}

type S3Config struct {
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
	Bucket string `yaml:"bucket"`
	Host   string `yaml:"host"`
}

type CleanConfig struct {
	OlderThan int `yaml:"olderthan"`
	Path      []string
}

type BackupConfig struct {
	Schedule string
	Fspath   []string
	DB       DBConfig `yaml:"db"`
}

type DBConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}
