package internal

// ApplicationConfig is populated from yaml config file
type ApplicationConfig struct {
	Machines []MachineConfig `yaml:"machines"`
	Remotes  RemoteConfig    `yaml:"remotes"`
}

type MachineConfig struct {
	Name   string       `yaml:"name"`
	Clean  CleanConfig  `yaml:"clean"`
	Backup BackupConfig `yaml:"backup"`
}

type RemoteConfig struct {
	SSH SSHAuthConfig `yaml:"auth"`
	S3  S3Config      `yaml:"s3"`
}

type SSHAuthConfig struct {
	Host            string `yaml:"host"`
	User            string `yaml:"user"`
	AuthPassword    string `yaml:"password"`
	KeyFile         string `yaml:"key_file"`
	KeyFilePassword string `yaml:"key_file_pass"`
}

type S3Config struct {
	Key      string `yaml:"key"`
	Secret   string `yaml:"secret"`
	Endpoint string `yaml:"endpoint"`
}

type CleanConfig struct {
	SourceConnection string   `yaml:"source"`
	OlderThan        int      `yaml:"olderthan"`
	Path             []string `yaml:"path"`
}

type BackupConfig struct {
	SourceConnection string   `yaml:"source"`
	TargetConnection string   `yaml:"target"`
	Bucket           string   `yaml:"bucket"`
	Prefix           string   `yaml:"prefix"`
	Fspath           []string `yaml:"fspath"`
	DB               DBConfig `yaml:"db"`
	Schedule         string   `yaml:"schedule"`
}

type DBConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string
}

// RclonePrepare exports env vars required for running rclone commands
func RclonePrepare() {
	//TODO
}

// RcloneClean removes rclone env vars once the program is done executing
func RcloneClean() {
	//TODO
}
