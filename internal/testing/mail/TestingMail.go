package mail

import (
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"gopkg.in/mail.v2"
)

type TestingMail struct {
	Name  string
	Email string
}

func (m TestingMail) Message() *mail.Message {
	msg := mail.NewMessage()
	msg.SetHeader("To", m.Email)
	msg.SetHeader("Subject", "Hello First!")
	msg.SetBody("text/html", xtremepkg.MailHTMLTemplate("testing_email.html", m))

	return msg
}
