package mail

import (
	"gopkg.in/gomail.v2"
	"manga-explorer/internal/app/common"
)

func NewSMTPMailService(config *common.Config) IService {
	return &sMTPService{config: config}
}

type sMTPService struct {
	config *common.Config
}

func (s sMTPService) SendEmail(mail *Mail) error {
	// Construct gomail message
	m := gomail.NewMessage()
	m.SetHeader("From", mail.From)
	m.SetHeader("To", mail.To)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody(mail.BodyType, mail.Body)

	// Open connection and send mail
	dialer := gomail.NewDialer(s.config.SMTPHost, int(s.config.SMTPPort), s.config.SMTPUser, s.config.SMTPPass)
	err := dialer.DialAndSend(m)

	return err
}
