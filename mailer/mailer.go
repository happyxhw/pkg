package mailer

import (
	"gopkg.in/gomail.v2"
)

var ma *Mailer

func InitDefaultMailer(cfg *Config) {
	ma = NewMailer(cfg)
}

func DefaultMailer() *Mailer {
	return ma
}

// Config for mailer
type Config struct {
	Host     string
	Port     int
	From     string
	User     string
	Password string
}

// Mailer send email
type Mailer struct {
	from   string
	dialer *gomail.Dialer
}

// NewMailer return a mailer
func NewMailer(c *Config) *Mailer {
	dialer := gomail.NewDialer(
		c.Host, c.Port, c.User, c.Password,
	)
	return &Mailer{
		from:   c.From,
		dialer: dialer,
	}
}

// Send dial and send email
func (m *Mailer) Send(to, subj, body string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", m.from)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subj)
	mail.SetBody("text/html", body)
	return m.dialer.DialAndSend(mail)
}
