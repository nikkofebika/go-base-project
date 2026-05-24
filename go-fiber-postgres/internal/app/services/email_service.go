package services

import (
	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendEmail(to []string, subject, body string) error
}

type emailService struct {
	host string
	port int
	user string
	pass string
	from string
}

func NewEmailService(host string, port int, user, pass, from string) EmailService {
	return &emailService{host: host, port: port, user: user, pass: pass, from: from}
}

func (s *emailService) SendEmail(to []string, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(s.host, s.port, s.user, s.pass)

	return d.DialAndSend(m)
}
