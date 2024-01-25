package service

import (
	"gopkg.in/gomail.v2"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/infrastructure/mail"
)

func NewSMTPMailService(sender string, config SMTPMailerConfig) IMail {
	return &mailerSMTPService{config: config, sender: sender}
}

// SMTPMailerConfig Only support Plain Auth
type SMTPMailerConfig struct {
	Host string
	Port uint16
	User string
	Pass string
}

type mailerSMTPService struct {
	config SMTPMailerConfig
	sender string
}

func (s mailerSMTPService) SendEmail(mail *mail.Mail) status.Object {
	// Construct gomail message
	m := gomail.NewMessage()
	m.SetHeader("From", s.Sender())
	m.SetHeader("To", mail.Recipients...)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody(mail.BodyType.String(), mail.Body)

	for _, v := range mail.EmbedFiles {
		m.Embed(v)
	}

	for _, v := range mail.AttachFiles {
		m.Attach(v)
	}

	// Open connection and send mail
	dialer := gomail.NewDialer(s.config.Host,
		int(s.config.Port),
		s.config.User,
		s.config.Pass)

	err := dialer.DialAndSend(m)
	if err != nil {
		return status.Error(status.MAIL_SEND_FAILED)
	}
	return status.InternalSuccess()
}

func (s mailerSMTPService) Sender() string {
	return s.sender
}
