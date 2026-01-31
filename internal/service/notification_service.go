package service

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"
)

type IMessageSender interface {
	Send(to string, subject string, body string) error
}

type emailSender struct {
	host     string
	port     int
	user     string
	password string
}

// NewEmailSender sekarang menerima konfigurasi SMTP
func NewEmailSender(host string, port int, user string, password string) IMessageSender {
	return &emailSender{
		host:     host,
		port:     port,
		user:     user,
		password: password,
	}
}

func (s *emailSender) Send(to string, subject string, body string) error {
	m := gomail.NewMessage()
	// m.SetHeader("From", s.user)
	m.SetHeader("From", "admin@myapp.dev")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.host, s.port, s.user, s.password)

	// Opsi tambahan jika menggunakan self-signed certificate atau SMTP lokal (// Mailtrap biasanya bekerja dengan InsecureSkipVerify: true di tahap dev)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
