package logtool

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/kushsharma/servo/internal"
	"github.com/pkg/errors"
	rcmd "github.com/rclone/rclone/cmd"
	rfs "github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/operations"
	rops "github.com/rclone/rclone/fs/operations"
	log "github.com/sirupsen/logrus"
)

type RcloneService struct {
	config internal.CleanConfig
}

// List provides all the files in directory recursively
func (svc *RcloneService) List(path string) ([]FileInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	lsCommand := fmt.Sprintf("local:%s", path)
	fsrc := rcmd.NewFsSrc(strings.Split(lsCommand, " "))

	files := []FileInfo{}
	err := rops.ListFn(ctx, fsrc, func(o rfs.Object) {
		file := FileInfo{
			Name:         o.Remote(),
			Size:         o.Size(),
			ModifiedTime: o.ModTime(ctx),
		}
		files = append(files, file)
		log.Debug(fmt.Sprintf("%9dB %s %s\n", file.Size, file.ModifiedTime.Local().Format("2006-01-02T15:04:05.000000000"), file.Name))
	})

	return files, err
}

// Clean removes unwanted files
func (svc *RcloneService) Clean(path string, minDayAge int) (err error) {
	files, err := svc.Cleanable(path, minDayAge)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	//now actually delete them
	for _, file := range files {
		fullPath := fmt.Sprintf("%s:%s", svc.config.SourceConnection, filepath.Join(path, file.Name))

		fs, fileName := rcmd.NewFsFile(fullPath)
		if fileName == "" {
			return errors.Errorf("%s is a directory or doesn't exist", fullPath)
		}
		fileObj, err := fs.NewObject(ctx, fileName)
		if err != nil {
			return err
		}
		if err := operations.DeleteFile(ctx, fileObj); err != nil {
			return err
		}
		log.Debug(".")
	}
	log.Infof("deleted %d files from %s", len(files), path)
	return nil
}

// Cleanable don't actually removes the files, only filter and list them
func (svc *RcloneService) Cleanable(path string, minDayAge int) ([]FileInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	fullPath := fmt.Sprintf("%s:%s", svc.config.SourceConnection, path)
	fsrc := rcmd.NewFsSrc(strings.Split(fullPath, " "))

	files := []FileInfo{}
	err := rops.ListFn(ctx, fsrc, func(o rfs.Object) {
		modTime := o.ModTime(ctx)

		//filter files based on there modification time
		prevTime := time.Now().Local().AddDate(0, 0, -minDayAge)
		if modTime.After(prevTime) {
			return
		}

		file := FileInfo{
			Name:         o.Remote(),
			Size:         o.Size(),
			ModifiedTime: modTime,
		}
		files = append(files, file)
	})

	return files, err
}

// NewRcloneService returns a instance of RcloneService that implements LogMangager using rclone in backend
func NewRcloneService(conf internal.CleanConfig) *RcloneService {
	svc := new(RcloneService)
	svc.config = conf
	return svc
}
