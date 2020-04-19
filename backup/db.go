package backup

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/tunnel"
)

type DBService struct {
	tnl    tunnel.Executioner
	s3     *s3.S3
	config internal.BackupConfig
	files  []string
}

func (svc *DBService) Prepare() error {

	return nil
}

// Migrate push fs db items to s3 bucket
func (svc *DBService) Migrate() error {

	return nil
}

func NewDBService(tnl tunnel.Executioner, s3Client *s3.S3, config internal.BackupConfig) *DBService {
	db := new(DBService)
	db.tnl = tnl
	db.s3 = s3Client
	db.config = config
	db.files = []string{}

	return db
}
