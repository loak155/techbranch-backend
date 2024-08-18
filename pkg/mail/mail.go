package mail

import (
	"bytes"
	"fmt"
	"net/smtp"
	"strings"
	"text/template"
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

func (m *MailManager) SendPreSignUpMail(to []string, token string) error {
	subject := "ユーザー仮登録の確認"
	url := "http://localhost:8080/v1/signup?token=" + token

	tmpl, err := template.ParseFiles("./pkg/mail/pre-signup.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	writer := new(bytes.Buffer)
	tmpl.Execute(writer, url)

	return m.SendMailWithHTML(to, subject, writer.String())
}
