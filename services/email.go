package services

import (
	"github.com/adrien3d/base-api/models"
	"github.com/matcornic/hermes/v2"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
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
	SendEmail( /*buffer *bytes.Buffer*/ content string, data *models.EmailData) error
	SendActivationEmail(user *models.User, apiUrl string, appName string, frontUrl string) error
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

// SendEmail sends an mail
func (s *EmailSenderParams) SendEmail( /*buffer *bytes.Buffer*/ content string, data *models.EmailData) error {
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
					Data:    aws.String(content),
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

// SendActivationEmail allows to send an email to user to activate his account
func (s *EmailSenderParams) SendActivationEmail(user *models.User, apiUrl string, appName string, frontUrl string) error {
	currentYear := string(time.Now().Year())

	h := hermes.Hermes{
		Theme: new(hermes.Flat),
		Product: hermes.Product{ // Appears in header & footer of e-mails
			Name: appName,
			Link: frontUrl,
			//Logo: ``,
			Copyright: `Copyright © ` + currentYear + ` ` + appName + `. All rights reserved.`,
		},
	}

	email := hermes.Email{
		Body: hermes.Body{
			Name: user.FirstName + ` ` + user.LastName,
			Intros: []string{
				`Welcome to ` + appName + `! We're very excited to have you on board.`,
			},
			Dictionary: []hermes.Entry{
				{Key: "Firstname", Value: user.FirstName},
				{Key: "Lastname", Value: user.LastName},
			},
			Actions: []hermes.Action{
				{
					Instructions: `To get started with ` + appName + `, please click here:`,
					Button: hermes.Button{
						/*Color: ``,
						TextColor: ``,*/
						Text: "Confirm your account",
						Link: apiUrl,
					},
				},
			},
			Outros: []string{
				`If you received this mail and it was not intended to you, please ignore it.`,
			},
		},
	}

	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		logrus.Warnln(err)
		panic(err)
	}

	data := models.EmailData{ReceiverMail: user.Email, ReceiverName: user.FirstName + " " + user.LastName, User: user, Subject: `Welcome to ` + appName + `! We're very excited to have you on board.`, ApiUrl: apiUrl, AppName: appName}

	return s.SendEmail(emailBody, &data)
}
