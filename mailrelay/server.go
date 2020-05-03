package mailrelay

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/emersion/go-smtp"
	"github.com/kushsharma/servo/internal"
)

const (
	// AwsRegion where mails will be relayed to
	AwsRegion = "ap-south-1"
)

//Start relay server
func Start(config internal.RemoteConfig) error {

	// Create a new session and specify an AWS Region.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AwsRegion),
		Credentials: credentials.NewStaticCredentials(config.SES.Key, config.SES.Secret, ""),
	})
	if err != nil {
		return err
	}

	// Create an SES client in the session.
	smtpBackend := new(Backend)
	smtpBackend.awsSES = ses.New(sess)
	smtpBackend.authUser = config.SMTP.Username
	smtpBackend.authPassword = config.SMTP.Password
	smtpBackend.sessionFactory = func(svc *ses.SES) *BackendSession {
		return &BackendSession{
			awsSES: svc,
		}
	}

	if err := StartSMTP(config.SMTP, smtpBackend); err != nil {
		return err
	}

	return nil
}

// The Backend implements SMTP server methods.
type Backend struct {
	authUser       string
	authPassword   string
	awsSES         *ses.SES
	sessionFactory SessionFactory
}

// Login handles a login command with username and password.
func (b *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if username == fmt.Sprintf("%s %s", b.authUser, b.authPassword) {
		return b.sessionFactory(b.awsSES), nil
	}
	if username == b.authUser && password == b.authPassword {
		return b.sessionFactory(b.awsSES), nil
	}
	return nil, errors.New("Invalid username or password")
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (b *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, smtp.ErrAuthRequired
}
