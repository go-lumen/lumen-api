package services

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"html/template"
	"io/ioutil"

	"github.com/adrien3d/base-api/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

const (
	emailSenderKey = "emailSender"
)

// GetEmailSender retrieves text sender
func GetEmailSender(c context.Context) EmailSender {
	return c.Value(emailSenderKey).(EmailSender)
}

// EmailSender creates a text sender interface
type EmailSender interface {
	//SendUserValidationEmail(user *models.User, subject string, templateLink string) error
	//SendAlertEmail(user *models.User, device *models.Device, observation *models.Observation, subject string, templateLink string) error
	SendEmailFromTemplate(ctx *gin.Context, data *models.EmailData, templateLink string) error
	SendEmail(data *models.EmailData) error
}

// FakeEmailSender structure
type FakeEmailSender struct{}

// EmailSenderParams with various text sender params
type EmailSenderParams struct {
	senderEmail string
	senderName  string
	apiID       string
	apiKey      string
	apiUrl      string
}

// NewEmailSender instantiates of the sender
func NewEmailSender(config *viper.Viper) EmailSender {
	return &EmailSenderParams{
		config.GetString("mail_sender_address"),
		config.GetString("mail_sender_name"),
		config.GetString("aws_api_id"),
		config.GetString("aws_api_key"),
		config.GetString("api_url"),
	}
}

// SendEmail sends a mail
func (s *EmailSenderParams) SendEmail(data *models.EmailData) error {
	file, err := ioutil.ReadFile("./templates/html/mail_skeleton.html")
	if err != nil {
		return err
	}

	htmlTemplate := template.Must(template.New("emailTemplate").Parse(string(file)))

	buffer := new(bytes.Buffer)
	err = htmlTemplate.Execute(buffer, data)
	if err != nil {
		logrus.Warnln(err)
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	creds := credentials.NewStaticCredentials(s.apiID, s.apiKey, "")

	// Create an SES session.
	svc := ses.New(sess, &aws.Config{Credentials: creds})

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(data.ReceiverMail),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(buffer.String()),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(data.Subject),
			},
		},
		Source: aws.String(s.senderEmail),
	}

	// Attempt to send the email.
	_, err = svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				logrus.Warnln(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				logrus.Warnln(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				logrus.Warnln(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				logrus.Warnln(aerr.Error())
			}
		} else {
			logrus.Warnln(err.Error())
		}

		logrus.Warnln(err)
		return err
	}

	logrus.Infoln("SES Email Sent to " + data.ReceiverName + " at address: " + data.ReceiverMail)

	return nil
}

// SendEmailFromTemplate sends an email from template
func (s *EmailSenderParams) SendEmailFromTemplate(ctx *gin.Context, data *models.EmailData, templateLink string) error {
	file, err := ioutil.ReadFile(templateLink)
	if err != nil {
		return err
	}

	htmlTemplate := template.Must(template.New("emailTemplate").Parse(string(file)))

	buffer := new(bytes.Buffer)
	err = htmlTemplate.Execute(buffer, data)
	if err != nil {
		logrus.Warnln(err)
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	creds := credentials.NewStaticCredentials(s.apiID, s.apiKey, "")

	// Create an SES session.
	svc := ses.New(sess, &aws.Config{Credentials: creds})

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(data.ReceiverMail),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(buffer.String()),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(data.Subject),
			},
		},
		Source: aws.String(s.senderEmail),
	}

	// Attempt to send the email.
	_, err = svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				logrus.Warnln(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				logrus.Warnln(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				logrus.Warnln(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				logrus.Warnln(aerr.Error())
			}
		} else {
			logrus.Warnln(err.Error())
		}

		logrus.Warnln(err)
		return err
	}

	logrus.Infoln("SES Email Sent to " + data.ReceiverName + " at address: " + data.ReceiverMail)

	return nil
}
