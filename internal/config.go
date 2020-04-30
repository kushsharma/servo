package internal

// ApplicationConfig is populated from yaml config file
type ApplicationConfig struct {
	Machines []MachineConfig `yaml:"machines"`
	Remotes  RemoteConfig    `yaml:"remotes"`
}

type MachineConfig struct {
	Name     string       `yaml:"name"`
	Schedule string       `yaml:"schedule"`
	Clean    CleanConfig  `yaml:"clean"`
	Backup   BackupConfig `yaml:"backup"`
}

type RemoteConfig struct {
	SSH  SSHAuthConfig `yaml:"ssh"`
	S3   S3Config      `yaml:"s3"`
	SES  SESConfig     `yaml:"ses"`
	SMTP SMTPConfig    `yaml:"smtp"`
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

type SESConfig struct {
	Key    string
	Secret string
}

type CleanConfig struct {
	SourceConnection string   `yaml:"source"`
	OlderThan        int      `yaml:"olderthan"`
	Path             []string `yaml:"path"`
}

type BackupConfig struct {
	FS FSBackupConfig `yaml:"fs"`
	DB DBBackupConfig `yaml:"db"`
}

type FSBackupConfig struct {
	SourceConnection string   `yaml:"source"`
	TargetConnection string   `yaml:"target"`
	Bucket           string   `yaml:"bucket"`
	Prefix           string   `yaml:"prefix"`
	Path             []string `yaml:"path"`
}

type DBBackupConfig struct {
	TargetConnection string `yaml:"target"`
	Bucket           string `yaml:"bucket"`
	Prefix           string `yaml:"prefix"`
	Auth             DBAuth `yaml:"auth"`
}

type DBAuth struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
}

type SMTPConfig struct {
	LocalListenIP   string `yaml:"listen_ip"`
	LocalListenPort int    `yaml:"listen_port"`
	SMTPUsername    string `yaml:"user"`
	SMTPPassword    string `yaml:"password"`
	SMTPServer      string `yaml:"server"`
	SMTPPort        string `yaml:"port"`
}
