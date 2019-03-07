package models

// EmailData holds informations required to send an email
type EmailData struct {
	ReceiverMail string
	ReceiverName string
	User         *User
	Subject      string
	Body         string
	ApiUrl       string
	AppName      string
}

// TextData holds informations required to send a text
type TextData struct {
	PhoneNumber string
	Subject     string
	Message     string
}
