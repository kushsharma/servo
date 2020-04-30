package backup

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/kushsharma/servo/internal"
	log "github.com/sirupsen/logrus"

	rcmd "github.com/rclone/rclone/cmd"
	rops "github.com/rclone/rclone/fs/operations"
	rsync "github.com/rclone/rclone/fs/sync"
)

type FSService struct {
	config internal.FSBackupConfig
	files  []string
}

// Prepare do any prerequisit for files that needs to be backedup
// TODO: we can gzip files before uploading
func (svc *FSService) Prepare() error {
	return nil
}

// Migrate push fs items to s3 bucket
func (svc *FSService) Migrate() error {
	if len(svc.config.Path) == 0 {
		return nil
	}

	errs := []error{}
	for _, sourcePath := range svc.config.Path {
		ctx, cancel := context.WithTimeout(context.Background(), uploadTimeout)
		defer cancel()

		destinationPath := filepath.Join(svc.config.Bucket, svc.config.Prefix, sourcePath)
		copyCommand := fmt.Sprintf("%s:%s %s:%s", svc.config.SourceConnection, sourcePath, svc.config.TargetConnection, destinationPath)

		fsrc, srcFileName, fdst := rcmd.NewFsSrcFileDst(strings.Split(copyCommand, " "))
		if srcFileName == "" {
			if err := rsync.CopyDir(ctx, fdst, fsrc, false); err != nil {
				errs = append(errs, err)
			}
		} else {
			if err := rops.CopyFile(ctx, fdst, fsrc, srcFileName, srcFileName); err != nil {
				errs = append(errs, err)
			}
		}
		log.Debug(".")
	}
	return internal.ErrMerge(errs)
}

// Close nil
func (svc *FSService) Close() error {
	return nil
}

func NewFSService(config internal.FSBackupConfig) *FSService {
	fs := new(FSService)
	fs.config = config
	fs.files = []string{}

	return fs
}

/*
// to check if dir exists - fmt.Sprintf(`[ -d "%s" ] && echo "1"`, path)
func (svc *FSService) findFiles() error {
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

func (svc *FSService) migrateManual() error {
	if len(svc.files) == 0 {
		return nil
	}

	uploader := s3manager.NewUploaderWithClient(svc.s3)
	for _, filepath := range svc.files {
		f, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer f.Close()

		object := svc.s3object(f, filepath)
		ctx, cancel := context.WithTimeout(context.Background(), uploadTimeout)
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
*/
