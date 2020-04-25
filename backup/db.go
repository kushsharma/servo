package backup

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/tunnel"
	rcmd "github.com/rclone/rclone/cmd"
	rops "github.com/rclone/rclone/fs/operations"
	log "github.com/sirupsen/logrus"
)

const (
	databaseDumpCommand = `#!/bin/sh
mysqldump -u{{.User}} {{.Password}} --all-databases --single-transaction --quick --lock-tables=false --triggers | gzip > {{.Name}}
`
	tempShellFileName   = "/tmp/servo_db_dump.sh"
	backupS3Directory   = "db"
	backupFileTimestamp = "2006-01-02_15-04-05"
	sourceConnection    = "local" //only local connection is supported for source for now
)

var (
	errRemovingTemporaryFiles = errors.New("error in removing temporary files")
)

type DBService struct {
	tnl    tunnel.Executioner
	config internal.BackupConfig
	file   string
}

type dumpTemplateInput struct {
	User     string
	Password string
	Name     string
}

// Prepare dumps database to a file
func (svc *DBService) Prepare() (err error) {
	if svc.config.DB.User == "" {
		log.Info("db backup skipped, no config found")
		return nil
	}

	dbDumpTemplate, err := template.New("dbdump").Parse(databaseDumpCommand)
	if err != nil {
		return err
	}

	input := new(dumpTemplateInput)
	input.User = svc.config.DB.User
	if svc.config.DB.Password != "" {
		input.Password = fmt.Sprintf("-p%s", svc.config.DB.Password)
	}
	input.Name = fmt.Sprintf("/tmp/db_dump_%s.sql.gz", time.Now().Format(backupFileTimestamp))
	svc.file = input.Name

	buf := &bytes.Buffer{}
	if err := dbDumpTemplate.Execute(buf, input); err != nil {
		return err
	}

	//write a temporary file to execute rendered command otherwise it won't work
	if err = ioutil.WriteFile(tempShellFileName, buf.Bytes(), 0744); err != nil {
		return err
	}
	defer func() {
		//delete temp file
		err := os.Remove(tempShellFileName)
		if err != nil {
			fmt.Printf("%v: %v\n", errRemovingTemporaryFiles, err)
		}
	}()

	if out, err := svc.tnl.RunWithOutput(fmt.Sprintf("sh -c %s", tempShellFileName)); err != nil {
		return fmt.Errorf("%s %v", string(out), err)
	}
	return nil
}

// Migrate push fs db items to s3 bucket
func (svc *DBService) Migrate() error {
	if svc.file == "" {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), uploadTimeout)
	defer cancel()

	destinationPath := filepath.Join(svc.config.Bucket, svc.config.Prefix, svc.file)
	copyCommand := fmt.Sprintf("%s:%s %s:%s --ignore-existing", sourceConnection, svc.file, svc.config.TargetConnection, destinationPath)

	fsrc, srcFileName, fdst := rcmd.NewFsSrcFileDst(strings.Split(copyCommand, " "))
	if err := rops.CopyFile(ctx, fdst, fsrc, srcFileName, srcFileName); err != nil {
		return err
	}

	return nil
}

// Close clean db file once upload is complete
func (svc *DBService) Close() error {
	if svc.file != "" {
		if err := os.Remove(svc.file); err != nil {
			log.Errorf("%v: %v\n", errRemovingTemporaryFiles, err)
			return err
		}
	}

	return nil
}

func NewDBService(tnl tunnel.Executioner, config internal.BackupConfig) *DBService {
	db := new(DBService)
	db.tnl = tnl
	db.config = config
	db.file = ""

	return db
}
