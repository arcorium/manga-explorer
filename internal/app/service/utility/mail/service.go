package mail

type IService interface {
	SendEmail(mail *Mail) error
}
