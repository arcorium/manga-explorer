package service

import (
	"github.com/stretchr/testify/require"
	"manga-explorer/internal/common/constant"
	"manga-explorer/internal/common/status"
	"manga-explorer/internal/infrastructure/mail"
	"manga-explorer/internal/util"
	"os"
	"reflect"
	"testing"
)

var mailer = &mailerSMTPService{config: SMTPMailerConfig{
	Host: "sandbox.smtp.mailtrap.io",
	Port: 2525,
	User: "b4b05a7f41be5d",
	Pass: "8612e5c9faf549",
},
	sender: constant.SenderEmail,
}

func Test_mailerSMTPService_SendEmail(t *testing.T) {
	html, err := mail.NewHTML("forgot-password.gohtml", "localhost:9999/asdasd", util.DropError(os.Getwd()))
	require.NoError(t, err)
	html.Recipients = []string{util.GenerateRandomString(30) + "@gmail.com"}
	html.Subject = "Test Complex HTML"

	type args struct {
		mail *mail.Mail
	}
	tests := []struct {
		name string
		args args
		want status.Object
	}{
		{
			name: "Normal Plain",
			args: args{
				mail: &mail.Mail{
					Recipients: []string{util.GenerateRandomString(30) + "@gmail.com"},
					Subject:    "Test",
					BodyType:   mail.BodyTypePlain,
					Body:       "Hello, this is test",
				},
			},
			want: status.InternalSuccess(),
		},
		{
			name: "Normal Complex HTML",
			args: args{
				mail: html,
			},
			want: status.InternalSuccess(),
		},
		{
			name: "Normal Template HTML",
			args: args{
				mail: &mail.Mail{
					Recipients: []string{util.GenerateRandomString(30) + "@gmail.com"},
					Subject:    "Test HTML",
					BodyType:   mail.BodyTypeHTML,
					Body:       "<b>Hello, this is HTML test</b>",
				},
			},
			want: status.InternalSuccess(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := mailer
			if got := s.SendEmail(tt.args.mail); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SendEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
