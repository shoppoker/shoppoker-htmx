package email

import (
	"github.com/w1png/go-htmx-ecommerce-template/config"
	"gopkg.in/gomail.v2"
)

var EmailServerInstance *EmailServer

type EmailServer struct {
	dialer *gomail.Dialer
	from   string
}

func InitEmailServer() (err error) {
	EmailServerInstance, err = NewEmailServer()
	return
}

func NewEmailServer() (*EmailServer, error) {
	dialer := gomail.NewDialer(
		config.ConfigInstance.SMTPHost,
		config.ConfigInstance.SMTPPort,
		config.ConfigInstance.SMTPUser,
		config.ConfigInstance.SMTPPassword,
	)

	closer, err := dialer.Dial()
	if err != nil {
		return nil, err
	}

	if err := closer.Close(); err != nil {
		return nil, err
	}

	return &EmailServer{
		dialer: dialer,
		from:   config.ConfigInstance.SMTPFrom,
	}, nil
}

func (s EmailServer) SendHTML(to []string, subject, body, attachement_path string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	if attachement_path != "" {
		m.Attach(attachement_path)
	}

	return s.dialer.DialAndSend(m)
}
