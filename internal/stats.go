package internal

var (
	AppStats *ApplicationStats
)

// ApplicationStats contains basic service stats
type ApplicationStats struct {
	TimesLogCleaned    int `json:"times_log_cleaned"`
	TimesBackedUp      int `json:"times_backed_up"`
	TimesLogError      int `json:"times_log_error"`
	TimesBackupFSError int `json:"times_backup_fs_error"`
	TimesBackupDBError int `json:"times_backup_db_error"`

	Version string `json:"version"`
}

// InitStat initializes basic application stats
func InitStat(version string) *ApplicationStats {
	AppStats = new(ApplicationStats)
	AppStats.Version = version
	return AppStats
}

func (s *ApplicationStats) LogCleaned() {
	s.TimesLogCleaned++
}

func (s *ApplicationStats) Backedup() {
	s.TimesBackedUp++
}

func (s *ApplicationStats) LogCleanError() {
	s.TimesLogError++
}

func (s *ApplicationStats) BackupFSError() {
	s.TimesBackupFSError++
}

func (s *ApplicationStats) BackupDBError() {
	s.TimesBackupDBError++
}
