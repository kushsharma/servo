package backup

import (
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/sshtunnel"
)

type FSService struct {
	ssh          *sshtunnel.Client
	backupConfig internal.BackupConfig
}

func (svc *FSService) Prepare() error {

	return nil
}

func (svc *FSService) Migrate() error {

	return nil
}

func NewFSService(client *sshtunnel.Client, config internal.BackupConfig) *FSService {
	fs := new(FSService)
	fs.ssh = client
	fs.backupConfig = config

	return fs
}
