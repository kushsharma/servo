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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/tunnel"
)

const (
	databaseDumpCommand = `#!/bin/sh
mysqldump -u{{.User}} {{.Password}} --all-databases --single-transaction --quick --lock-tables=false --triggers | gzip > {{.Name}}
`
	tempShellFileName = "/tmp/servo_db_dump.sh"
	backupS3Directory = "db"
)

var (
	errRemovingTemporaryFiles = errors.New("error in removing temporary files")
)

type DBService struct {
	tnl    tunnel.Executioner
	s3     *s3.S3
	config internal.BackupConfig
	file   string
}

type dumpTemplateInput struct {
	User     string
	Password string
	Name     string
}

// Prepare dumps database to a file
func (svc *DBService) Prepare() error {
	dbDumpTemplate, err := template.New("dbdump").Parse(databaseDumpCommand)
	if err != nil {
		return err
	}

	input := new(dumpTemplateInput)
	input.User = svc.config.DB.User
	if svc.config.DB.Password != "" {
		input.Password = fmt.Sprintf("-p%s", svc.config.DB.Password)
	}
	input.Name = fmt.Sprintf("/tmp/db_dump_%s.sql.gz", time.Now().Format("2006-01-02_15-04-05"))
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

	uploader := s3manager.NewUploaderWithClient(svc.s3)
	f, err := os.Open(svc.file)
	if err != nil {
		return err
	}
	defer f.Close()

	object := svc.s3object(f, svc.file)
	ctx, cancel := context.WithTimeout(context.Background(), UploadTimeout)
	defer cancel()

	_, err = uploader.UploadWithContext(ctx, object)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		//clean db file once upload is complete
		err = os.Remove(svc.file)
		if err != nil {
			fmt.Printf("%v: %v\n", errRemovingTemporaryFiles, err)
		}
	}()

	return nil
}

func (svc *DBService) s3object(f *os.File, path string) *s3manager.UploadInput {
	filename := filepath.Base(path)
	return &s3manager.UploadInput{
		Bucket: aws.String(svc.config.Bucket),
		Key:    aws.String(filepath.Join(svc.config.Prefix, backupS3Directory, filename)),
		ACL:    aws.String("private"),
		Body:   f,
	}
}

func NewDBService(tnl tunnel.Executioner, s3Client *s3.S3, config internal.BackupConfig) *DBService {
	db := new(DBService)
	db.tnl = tnl
	db.s3 = s3Client
	db.config = config
	db.file = ""

	return db
}
