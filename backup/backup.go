package backup

type BackupService interface {
	Prepare() error
	Migrate() error
}
