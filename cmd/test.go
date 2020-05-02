package cmd

import (
	"errors"

	"github.com/kushsharma/servo/internal"
	"github.com/kushsharma/servo/mailrelay"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

//localfs:./version s3DO:moonwaretech/temp/ --ignore-existing

func initTest() *cobra.Command {
	return &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			log.Info("starting test tool...")

			appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
			if !ok {
				return errors.New("unable to find application config")
			}

			// Create a new session and specify an AWS Region.
			// sess, err := session.NewSession(&aws.Config{
			// 	Region:      aws.String(AwsRegion),
			// 	Credentials: credentials.NewStaticCredentials(appConfig.Remotes.SES.Key, appConfig.Remotes.SES.Secret, ""),
			// })

			err = mailrelay.StartServerGOSMTP(appConfig.Remotes.SMTP)
			//err = mailrelay.StartServer(appConfig.Remotes.SMTP)
			if err != nil {
				return err
			}
			// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
			// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
			//signal.Notify(termChan, os.Interrupt)
			//signal.Notify(termChan, os.Kill)

			// // Block until we receive our signal.
			//<-termChan

			return err
		},
	}
}
