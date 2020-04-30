package mailrelay

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/kushsharma/servo/internal"
	"github.com/spf13/viper"
)

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "noreply@softnuke.com"

	// Replace recipient@example.com with a "To" address. If your account
	// is still in the sandbox, this address must be verified.
	Recipient = "kush.darknight@gmail.com"

	// Specify a configuration set. If you do not want to use a configuration
	// set, comment out the following constant and the
	// ConfigurationSetName: aws.String(ConfigurationSet) argument below
	ConfigurationSet = "ConfigSet"

	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	AwsRegion = "ap-south-1"

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
)

func sendMail(e *mail.Envelope) error {
	appConfig, ok := viper.Get("app").(internal.ApplicationConfig)
	if !ok {
		return errors.New("unable to find application config")
	}

	fmt.Print(e.MailFrom)

	// Create a new session and specify an AWS Region.
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(AwsRegion),
		Credentials: credentials.NewStaticCredentials(appConfig.Remotes.SES.Key, appConfig.Remotes.SES.Secret, ""),
	})
	if err != nil {
		return err
	}

	// Create an SES client in the session.
	svc := ses.New(sess)
	fmt.Print(svc.APIVersion)

	var msg bytes.Buffer
	msg.Write(e.Data.Bytes())
	msg.WriteString("\r\n")
	fmt.Print(msg.String())

	// mailOb := email.NewEmail()
	// mailOb.From = Sender
	// mailOb.To = []string{Recipient}
	// mailOb.Subject = e.Subject
	// mailOb.Headers = e.Header
	// mailOb.HTML = e.Data.Bytes()
	// mailBytes, err := mailOb.Bytes()
	// if err != nil {
	// 	return err
	// }
	//fmt.Print(string(mailBytes))

	input := &ses.SendRawEmailInput{
		Destinations: []*string{
			aws.String("To:" + Recipient),
		},
		RawMessage: &ses.RawMessage{
			Data: msg.Bytes(),
		},
	}

	//output, err := svc.SendRawEmail(input)
	fmt.Print(input)

	return nil
}

func checkUserIdentity(svc *ses.SES) {

	//get list of emails which are allowed to be sender on emails from ses
	result, err := svc.ListIdentities(&ses.ListIdentitiesInput{IdentityType: aws.String("EmailAddress")})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, email := range result.Identities {
		var e = []*string{email}
		verified, err := svc.GetIdentityVerificationAttributes(&ses.GetIdentityVerificationAttributesInput{Identities: e})
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		for _, va := range verified.VerificationAttributes {
			if *va.VerificationStatus == "Success" {
				fmt.Println(*email)
			}
		}
	}
}

func createFormatedMail(svc *ses.SES) {

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(Recipient),
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
		Source: aws.String(Sender),
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

	fmt.Println("Email Sent!")
	fmt.Print(sendMailResult)
}
