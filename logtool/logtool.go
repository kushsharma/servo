package logtool

import "time"

// I call this logtool instead of log because log term is normally
// used for storing application logs, creating a package folder with
// that name could cause confusion

const (
	uploadTimeout = 30 * time.Minute
)

type FileInfo struct {
	Name         string
	Size         int64
	ModifiedTime time.Time
}

//LogManager manages logs in a machine instance
type LogManager interface {
	List(path string) ([]FileInfo, error)
	Clean(path string, minDayAge int) error
	Cleanable(path string, minDayAge int) ([]FileInfo, error)
}
