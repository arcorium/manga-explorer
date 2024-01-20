package service

import "manga-explorer/internal/infrastructure/mail"

type IMail interface {
	// SendEmail sending mail asynchronously
	SendEmail(mail *mail.Mail) error
}
