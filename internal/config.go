package internal

// ApplicationConfig is populated from yaml config file
type ApplicationConfig struct {
	Machines []MachineConfig `yaml:"machines"`
	S3       S3Config        `yaml:"s3"`
}

type MachineConfig struct {
	Name           string        `yaml:"name"`
	ConnectionType string        `yaml:"conn"`
	Auth           SSHAuthConfig `yaml:"auth"`
	Clean          CleanConfig   `yaml:"clean"`
	Backup         BackupConfig  `yaml:"backup"`
}

type SSHAuthConfig struct {
	Address      string `yaml:"address"`
	User         string `yaml:"user"`
	AuthPassword string `yaml:"authpassword"`
	KeyFile      string `yaml:"keyfile"`
	KeyPassword  string `yaml:"keypassword"`
}

type S3Config struct {
	Key      string `yaml:"key"`
	Secret   string `yaml:"secret"`
	Endpoint string `yaml:"endpoint"`
}

type CleanConfig struct {
	OlderThan int      `yaml:"olderthan"`
	Path      []string `yaml:"path"`
}

type BackupConfig struct {
	Schedule string   `yaml:"schedule"`
	Bucket   string   `yaml:"bucket"`
	Prefix   string   `yaml:"prefix"`
	Fspath   []string `yaml:"fspath"`
	DB       DBConfig `yaml:"db"`
}

type DBConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// RclonePrepare exports env vars required for running rclone commands
func RclonePrepare() {
	//TODO
}

// RcloneClean removes rclone env vars once the program is done executing
func RcloneClean() {
	//TODO
}
