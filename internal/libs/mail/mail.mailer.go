package mail

import (
	"homedy/config"

	"github.com/jaytaylor/html2text"
	"gopkg.in/gomail.v2"
)

type Mailer struct {
	dialer *gomail.Dialer
}

func NewMailer(host, user, pass string, port int) *Mailer {
	return &Mailer{
		dialer: gomail.NewDialer(host, port, user, pass),
	}
}

func NewAppMailer() *Mailer {
	return NewMailer(config.MAIL_HOST, config.MAIL_USER, config.MAIL_PASS, 587 /* 587 = TLS */)
}

func (m *Mailer) Dialer() *gomail.Dialer {
	return m.dialer
}

// send with custom plain
func (m *Mailer) SendWithPlain(to, subject, plain, html string) error {
	msg := m.Start()
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plain)
	msg.AddAlternative("text/html", html)

	return m.dialer.DialAndSend(msg)
}

// auto convert html to plain
func (m *Mailer) Send(to, subject, html string) error {
	plain, err := html2text.FromString(html, html2text.Options{
		PrettyTables: true,
	})
	if err != nil {
		plain = "Please view this email in an HTML-supported client."
	}

	return m.SendWithPlain(to, subject, plain, html)
}

// start a message and set base headers
func (m *Mailer) Start() *gomail.Message {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.dialer.Username)
	return msg
}
