package types

import (
	"fmt"
	"strings"
)

// TODO: support multipart

type Mail struct {
	From    MailAddress   `json:"from"`
	To      MailAddresses `json:"to"`
	Subject string        `json:"subject"`
	Body    string        `json:"body"`
	Type    string        `json:"type"`
}

func (m *Mail) WithText(text string) *Mail {
	m.Body = text
	m.Type = "text/plain; charset=UTF-8"
	return m
}

func (m *Mail) WithHTML(html string) *Mail {
	m.Body = html
	m.Type = "text/html; charset=UTF-8"
	return m
}

func (m *Mail) String() string {
	return fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: %s\r\n"+
		"\r\n"+
		"%s\r\n",
		strings.Join(m.To.Strings(), ", "),
		m.From.String(),
		m.Subject,
		m.Type,
		m.Body,
	)
}

func (m *Mail) Reader() *strings.Reader {
	return strings.NewReader(m.String())
}

type MailSendRequest struct {
	To           MailAddresses     `json:"to"`
	Subject      string            `json:"subject"`
	Text         string            `json:"text"`
	Html         string            `json:"html"`
	TemplateId   string            `json:"template_id"`
	TemplateVars map[string]string `json:"template_vars"`
}
