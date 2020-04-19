package backup

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/tunnel"
)

type FSService struct {
	tnl          tunnel.Executioner
	s3           *s3.S3
	backupConfig internal.BackupConfig
}

func (svc *FSService) Prepare() error {

	return nil
}

func (svc *FSService) Migrate() error {

	return nil
}

func NewFSService(tnl tunnel.Executioner, s3Client *s3.S3, config internal.BackupConfig) *FSService {
	fs := new(FSService)
	fs.tnl = tnl
	fs.s3 = s3Client
	fs.backupConfig = config

	return fs
}
