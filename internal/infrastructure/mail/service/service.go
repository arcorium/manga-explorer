package service

import (
  "manga-explorer/internal/common/status"
  "manga-explorer/internal/infrastructure/mail"
)

type IMail interface {
  // SendEmail sending mail asynchronously
  SendEmail(mail *mail.Mail) status.Object

  Sender() string
}
