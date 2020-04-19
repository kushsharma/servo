package internal

// ApplicationConfig is populated from yaml config file
type ApplicationConfig struct {
	Machines []MachineConfig `yaml:"machines"`
	S3       S3Config        `yaml:"s3"`
}

type MachineConfig struct {
	Name           string
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
	Bucket   string `yaml:"bucket"`
	Endpoint string `yaml:"endpoint"`
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
