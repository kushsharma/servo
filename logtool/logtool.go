package logtool

// I call this logtool instead of log because log term is normally
// used for storing application logs, creating a package folder with
// that name could cause confusion

//LogManager manages logs in a machine instance
type LogManager interface {
	Delete(path string) error
	Clean(path string, daysold int) error
	DryClean(path string, daysold int) ([]string, error)
}
