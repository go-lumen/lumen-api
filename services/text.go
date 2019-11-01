package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/gin-gonic/gin"
	"go-lumen/lumen-api/models"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
)

const (
	textSenderKey = "textSender"
)

// GetTextSender retrieves text sender
func GetTextSender(c context.Context) TextSender {
	return c.Value(textSenderKey).(TextSender)
}

// TextSender creates a text sender interface
type TextSender interface {
	SendAlertText(c *gin.Context, user *models.User, message string, templateLink string) error
	SendText(ctx *gin.Context, data models.TextData) error
}

// FakeTextSender structure
type FakeTextSender struct{}

// TextSenderParams with various text sender params
type TextSenderParams struct {
	senderEmail string
	senderName  string
	apiID       string
	apiKey      string
	apiUrl      string
}

// NewTextSender instantiates of the sender
func NewTextSender(config *viper.Viper) TextSender {
	return &TextSenderParams{
		config.GetString("mail_sender_address"),
		config.GetString("mail_sender_name"),
		config.GetString("aws_api_id"),
		config.GetString("aws_api_key"),
		config.GetString("api_url"),
	}
}

// SendAlertText sends a simple alert
func (s *TextSenderParams) SendAlertText(c *gin.Context, user *models.User, message string, templateLink string) error {
	data := models.TextData{PhoneNumber: user.Phone, Message: message}
	if s.SendText(c, data) != nil {
		logrus.Warnln(`Send text error`)
	}

	return nil
}

// SendText sends any type of text
func (s *TextSenderParams) SendText(ctx *gin.Context, data models.TextData) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-west-1")},
	)

	// logrus.Infoln("Amazon Creds: " + s.apiID + s.apiKey)
	creds := credentials.NewStaticCredentials(s.apiID, s.apiKey, "")

	// Creates an SES session.
	svc := sns.New(sess, &aws.Config{Credentials: creds})

	// Assembling the text and attempting to send the email.
	params := &sns.PublishInput{
		Subject:     aws.String(data.Subject),
		Message:     aws.String(data.Message),
		PhoneNumber: aws.String(data.PhoneNumber),
	}
	_, err = svc.Publish(params)

	if err != nil {
		logrus.Warnln(err.Error())
		return err
	}

	// logrus.Infoln("SNS Text Sent to " + data.PhoneNumber)
	// logrus.Infoln(resp)

	return nil
}
