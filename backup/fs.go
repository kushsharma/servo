package backup

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/sshtunnel"
)

type FSService struct {
	ssh          *sshtunnel.Client
	s3           *s3.S3
	backupConfig internal.BackupConfig
}

func (svc *FSService) Prepare() error {
	spaces, err := svc.s3.ListBuckets(nil)
	if err != nil {
		return err
	}

	for _, b := range spaces.Buckets {
		fmt.Println(aws.StringValue(b.Name))
	}

	return nil
}

func (svc *FSService) Migrate() error {

	return nil
}

func NewFSService(sshclient *sshtunnel.Client, s3Client *s3.S3, config internal.BackupConfig) *FSService {
	fs := new(FSService)
	fs.ssh = sshclient
	fs.s3 = s3Client
	fs.backupConfig = config

	return fs
}
