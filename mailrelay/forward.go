package mailrelay

import (
	"fmt"
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/emersion/go-smtp"
	log "github.com/sirupsen/logrus"
)

const (
	// Specify a configuration set. If you do not want to use a configuration
	// set, comment out the following constant and the
	// ConfigurationSetName: aws.String(ConfigurationSet) argument below
	ConfigurationSet = "ConfigSet"

	// The subject line for the email.
	Subject = "Amazon SES Test (AWS SDK for Go) 2"

	// The HTML body for the email.
	HtmlBody = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
		"<a href='https://aws.amazon.com/ses/'>Amazon SES</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	//The email body for recipients with non-HTML email clients.
	TextBody = "This email was sent with Amazon SES using the AWS SDK for Go."

	// The character encoding for the email.
	CharSet = "UTF-8"

	FromRegex = "From: [a-zA-Z0-9@+>< _.]+\\n"
	MailRegex = "([a-zA-Z0-9+._-]+@[a-zA-Z0-9._-]+\\.[a-zA-Z0-9_-]+)"
)

type SessionFactory func(*ses.SES) *BackendSession

// BackendSession is returned after successful login.
type BackendSession struct {
	awsSES *ses.SES
	from   string
	to     string
	data   []byte
}

func (s *BackendSession) Mail(from string, opts smtp.MailOptions) error {
	s.from = from
	return nil
}

func (s *BackendSession) Rcpt(to string) error {
	s.to = to
	return nil
}

func (s *BackendSession) Data(r io.Reader) (err error) {
	if b, err := ioutil.ReadAll(r); err == nil {
		s.data = b
	}
	return err
}

func (s *BackendSession) Reset() {
	if len(s.data) > 0 {
		s.SendMail()
	}
	s.from = ""
	s.to = ""
	s.data = []byte{}
}

func (s *BackendSession) Logout() error {
	return s.SendMail()
}

func (s *BackendSession) SendMail() error {
	if len(s.data) == 0 {
		return nil
	}
	mailData := string(s.data)

	fromPatt, _ := regexp.Compile(FromRegex)
	fromMatch := fromPatt.FindString(mailData)

	mailPatt, _ := regexp.Compile(MailRegex)
	emailFrom := mailPatt.FindString(fromMatch)

	log.Debugf("sending mail from %s : %s", emailFrom, mailData)
	if valid, err := checkUserIdentity(s.awsSES, emailFrom); err != nil || !valid {
		berr := fmt.Errorf("unable to verify valid email sender: %v", err)
		log.Error(berr)
		return berr
	}

	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: s.data,
		},
	}

	_, err := s.awsSES.SendRawEmail(input)
	if err != nil {
		log.Error(err)
		return err
	}

	s.data = []byte{}
	return err
}

//checkUserIdentity check if we are allowed to send email using this from field
func checkUserIdentity(svc *ses.SES, from string) (bool, error) {
	fromSplit := strings.Split(from, "@")
	if len(fromSplit) == 1 {
		return false, nil
	}
	domain := fromSplit[1]

	var e = []*string{aws.String(domain)}
	verified, err := svc.GetIdentityVerificationAttributes(&ses.GetIdentityVerificationAttributesInput{Identities: e})
	if err != nil {
		return false, err
	}
	for _, va := range verified.VerificationAttributes {
		if *va.VerificationStatus == "Success" {
			return true, nil
		}
	}

	return false, nil
}

func createFormatedMail(svc *ses.SES, sender, recipient string) {

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(Subject),
			},
		},
		Source: aws.String(sender),
		// Comment or remove the following line if you are not using a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	sendMailResult, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	log.Info("Email Sent!")
	log.Info(sendMailResult)
}
