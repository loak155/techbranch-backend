package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

type MailManager struct {
	address string
	from    string
	auth    smtp.Auth
}

func NewMailManager(address, from, username, password string) *MailManager {
	auth := smtp.CRAMMD5Auth(username, password)
	return &MailManager{
		address: address,
		from:    from,
		auth:    auth,
	}
}

func (m *MailManager) SendMail(to []string, subject, body string) error {
	msg := []byte(strings.ReplaceAll(fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", m.from, strings.Join(to, ","), subject, body), "\n", "\r\n"))
	return smtp.SendMail(m.address, m.auth, m.from, to, msg)
}

func (m *MailManager) SendMailWithHTML(to []string, subject, body string) error {
	msg := []byte(strings.ReplaceAll(fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nContent-Type: text/html; charset=UTF-8\n\n%s", m.from, strings.Join(to, ","), subject, body), "\n", "\r\n"))
	return smtp.SendMail(m.address, m.auth, m.from, to, msg)
}
