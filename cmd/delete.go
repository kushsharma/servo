package cmd

import (
	"errors"
	"fmt"

	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/logtool"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	errInvalidSSHMachineConfig = errors.New("invalid machine configs provided")
)

func initDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use: "delete [log/s3]",
	}
	cmd.AddCommand(subDeleteS3())
	cmd.AddCommand(subDeleteLog())
	return cmd
}

func subDeleteS3() *cobra.Command {
	return &cobra.Command{
		Use:     "s3",
		Short:   "delete all the files provided in bucket that match with key prefix",
		Example: "servo delete [s3,log] temp/path/key/prefix bucketname",
		Args:    cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			prefix := args[0]
			bucket := args[1]
			log.Info("starting s3 deletion tool...")

			appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
			if !ok {
				return errors.New("unable to find application config")
			}

			s3Config := &aws.Config{
				Credentials: credentials.NewStaticCredentials(appConfig.Remotes.S3.Key, appConfig.Remotes.S3.Secret, ""),
				Endpoint:    aws.String(appConfig.Remotes.S3.Endpoint),
				Region:      aws.String("us-east-1"),
			}
			// Create S3 service client
			awsSession := session.New(s3Config)
			s3Client := s3.New(awsSession)
			if err := deleteS3(s3Client, bucket, prefix); err != nil {
				return err
			}
			return nil
		},
	}
}

// delete all the files in provided s3 bucket prefix
func deleteS3(client *s3.S3, bucket, prefix string) error {
	listInput := &s3.ListObjectsInput{
		Bucket:  aws.String(bucket),
		Prefix:  aws.String(prefix),
		MaxKeys: aws.Int64(500),
	}

	// list all files available in folder
	fileKeys := []string{}
	err := client.ListObjectsPages(listInput, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		for _, item := range page.Contents {
			fileKeys = append(fileKeys, *item.Key)
			log.Debugf("file found: %s", *item.Key)
		}
		return !lastPage
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				log.Error(s3.ErrCodeNoSuchBucket, aerr.Error())
			default:
				log.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Error(err.Error())
		}
		return err
	}

	if len(fileKeys) == 0 {
		return errors.New("no file found to delete in provided prefix")
	}

	// take user confirmation
	fmt.Printf("deleting %d items, type \"yes\" for confirmation: ", len(fileKeys))
	var confirmation string
	fmt.Scanln(&confirmation)
	if confirmation != "yes" {
		fmt.Println("operation aborted")
		return nil
	}

	deleteObjects := []*s3.ObjectIdentifier{}
	for _, key := range fileKeys {
		deleteObjects = append(deleteObjects, &s3.ObjectIdentifier{
			Key: aws.String(key),
		})
	}

	// batch delete files
	deleteInput := &s3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &s3.Delete{
			Objects: deleteObjects,
			Quiet:   aws.Bool(false),
		},
	}
	_, err = client.DeleteObjects(deleteInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Error(err.Error())
		}
		return err
	}

	log.Infof("all objects in %s are deleted from bucket %s\n", prefix, bucket)
	return nil
}

func subDeleteLog() *cobra.Command {
	return &cobra.Command{
		Use:   "log",
		Short: "delete logs older than provided days",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Info("starting log cleaner...")

			appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
			if !ok {
				return errInvalidSSHMachineConfig
			}

			for _, machine := range appConfig.Machines {
				// tnl, err := createTunnel(machine)
				// if err != nil {
				// 	return err
				// }
				// defer tnl.Close()

				// logToolService := logtool.NewShellService(tnl)
				// if err := logClean(logToolService, machine.Clean); err != nil {
				// 	return err
				// }

				logToolService := logtool.NewRcloneService()
				if err := logClean(logToolService, machine.Clean); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

//logClean remove files that are unnecessary and older than x days
func logClean(svc logtool.LogManager, config internal.CleanConfig) error {

	//files to be cleaned older than x days
	daysOld := config.OlderThan

	errs := []error{}
	for _, path := range config.Path {
		if files, err := svc.Cleanable(path, daysOld); err == nil {
			log.Debugf("%d file found for deletion: %v", len(files), files)
		} else {
			errs = append(errs, err)
		}

		if !DryRun {
			if err := svc.Clean(path, daysOld); err != nil {
				errs = append(errs, err)
			} else {
				log.Info("logs cleaned successfully for %s", path)
			}
		}
	}

	return internal.ErrMerge(errs)
}
