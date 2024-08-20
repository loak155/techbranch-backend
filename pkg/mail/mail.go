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

const (
	GmailHost = "smtp.gmail.com"
	GmailPort = 587
)

type GmailManager struct {
	host string
	port int
	from string
	auth smtp.Auth
}

func NewGmailManager(fromAddress, password string) *GmailManager {
	return &GmailManager{
		host: GmailHost,
		port: GmailPort,
		from: fromAddress,
		auth: smtp.PlainAuth("", fromAddress, password, GmailHost),
	}
}

func (g *GmailManager) SendMail(to []string, subject, body string) error {
	msg := []byte(strings.ReplaceAll(fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", g.from, strings.Join(to, ","), subject, body), "\n", "\r\n"))
	return smtp.SendMail(fmt.Sprintf("%s:%d", g.host, g.port), g.auth, g.from, to, msg)
}

func (g *GmailManager) SendMailWithHTML(to []string, subject, body string) error {
	msg := []byte(strings.ReplaceAll(fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nContent-Type: text/html; charset=UTF-8\n\n%s", g.from, strings.Join(to, ","), subject, body), "\n", "\r\n"))
	return smtp.SendMail(fmt.Sprintf("%s:%d", g.host, g.port), g.auth, g.from, to, msg)
}
