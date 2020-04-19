package backup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/tunnel"
)

const (
	UploadTimeout = 30 * time.Minute
)

type FSService struct {
	tnl    tunnel.Executioner
	s3     *s3.S3
	config internal.BackupConfig
	files  []string
}

// Prepare populates files that needs to be backedup
// to check if dir exists - fmt.Sprintf(`[ -d "%s" ] && echo "1"`, path)
func (svc *FSService) Prepare() error {
	for _, path := range svc.config.Fspath {
		out, err := svc.tnl.RunWithOutput(fmt.Sprintf(`find %s -type f`, path))
		if err != nil {
			fmt.Print(out)
			return err
		}
		svc.files = append(svc.files, strings.Split(out, "\n")...)
	}
	return nil
}

// Migrate push fs items to s3 bucket
func (svc *FSService) Migrate() error {
	uploader := s3manager.NewUploaderWithClient(svc.s3)

	for _, filepath := range svc.files {
		f, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer f.Close()

		object := svc.s3object(f, filepath)
		ctx, cancel := context.WithTimeout(context.Background(), UploadTimeout)
		defer cancel()

		_, err = uploader.UploadWithContext(ctx, object)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return nil
}

func (svc *FSService) s3object(f *os.File, path string) *s3manager.UploadInput {
	return &s3manager.UploadInput{
		Bucket: aws.String(svc.config.Bucket),
		Key:    aws.String(filepath.Join(svc.config.Prefix, path)),
		ACL:    aws.String("private"),
		Body:   f,
	}
}

func NewFSService(tnl tunnel.Executioner, s3Client *s3.S3, config internal.BackupConfig) *FSService {
	fs := new(FSService)
	fs.tnl = tnl
	fs.s3 = s3Client
	fs.config = config
	fs.files = []string{}

	return fs
}
