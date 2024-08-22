package mail

import (
	"fmt"
	"net/smtp"
	"strings"
)

const (
	GmailHost = "smtp.gmail.com"
	GmailPort = 587
)

type Manager struct {
	host string
	port int
	from string
	auth smtp.Auth
}

func NewManager(host string, port int, from, password string) *Manager {
	var auth smtp.Auth
	if from != "" && password != "" {
		auth = smtp.PlainAuth("", from, password, host)
	}
	return &Manager{
		host: host,
		port: port,
		from: from,
		auth: auth,
	}
}

func (g *Manager) SendMail(to []string, subject, body string) error {
	msg := []byte(strings.ReplaceAll(fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", g.from, strings.Join(to, ","), subject, body), "\n", "\r\n"))
	return smtp.SendMail(fmt.Sprintf("%s:%d", g.host, g.port), g.auth, g.from, to, msg)
}

func (g *Manager) SendMailWithHTML(to []string, subject, body string) error {
	msg := []byte(strings.ReplaceAll(fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\nContent-Type: text/html; charset=UTF-8\n\n%s", g.from, strings.Join(to, ","), subject, body), "\n", "\r\n"))
	return smtp.SendMail(fmt.Sprintf("%s:%d", g.host, g.port), g.auth, g.from, to, msg)
}
