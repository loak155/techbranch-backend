package mail

import (
	"bytes"
	"fmt"
	"text/template"
)

type PresignupMailManager struct {
	mailManager *MailManager
	subject     string
	tmpl        *template.Template
	signupURL   string
}

type TemplateData struct {
	Username string
	URL      string
}

func NewPresignupMailManager(address, from, username, password, subject, templateFilePath, signupURL string) (*PresignupMailManager, error) {

	mailManager := NewMailManager(address, from, username, password)

	tmpl, err := template.ParseFiles(templateFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &PresignupMailManager{
		mailManager: mailManager,
		subject:     subject,
		tmpl:        tmpl,
		signupURL:   signupURL,
	}, nil
}

func (m *PresignupMailManager) SendPreSignUpMail(to []string, username, token string) error {
	url := m.signupURL + token

	tmplData := TemplateData{
		Username: username,
		URL:      url,
	}

	writer := new(bytes.Buffer)
	m.tmpl.Execute(writer, tmplData)

	return m.mailManager.SendMailWithHTML(to, m.subject, writer.String())
}
