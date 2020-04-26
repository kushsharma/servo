package routes

import (
	"encoding/json"
	"net/http"
)

var (
	AppStats = new(ApplicationStats)
)

// ApplicationStats contains basic service stats
type ApplicationStats struct {
	TimesLogCleaned    int `json:"times_log_cleaned"`
	TimesBackedUp      int `json:"times_backed_up"`
	TimesLogError      int `json:"times_log_error"`
	TimesBackupFSError int `json:"times_backup_fs_error"`
	TimesBackupDBError int `json:"times_backup_db_error"`
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

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AppStats)
}
