package backup

import (
	"io"
	"time"
)

const (
	uploadTimeout = 30 * time.Minute
)

type BackupService interface {
	Prepare() error
	Migrate() error
	io.Closer
}
