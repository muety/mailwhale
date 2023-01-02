package types

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
)

const (
	typeHtml  = "text/html; charset=UTF-8"
	typePlain = "text/plain; charset=UTF-8"
)

type Mail struct {
	From      MailAddress
	To        MailAddresses
	Subject   string
	Body      string
	Type      string
	Date      time.Time
	MessageID string
}

func (m *Mail) WithText(text string) *Mail {
	m.Body = text
	m.Type = typePlain
	return m
}

func (m *Mail) WithHTML(html string) *Mail {
	m.Body = html
	m.Type = typeHtml
	return m
}

func (m *Mail) Sanitized() *Mail {
	if m.Type == "" {
		m.Type = typePlain
	}
	if m.Date.IsZero() {
		m.Date = time.Now()
	}
	if m.MessageID == "" {
		m.MessageID = fmt.Sprintf("<%s@%s>", uuid.New().String(), m.From.Domain())
	}
	return m
}

func (m *Mail) String() string {
	return fmt.Sprintf("To: %s\r\n"+
		"From: %s\r\n"+
		"Subject: %s\r\n"+
		"Message-ID: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: %s\r\n"+
		"Content-Transfer-Encoding: 8bit\r\n"+
		"Date: %s\r\n"+
		"\r\n"+
		"%s\r\n",
		strings.Join(m.To.Strings(), ", "),
		m.From.String(),
		m.Subject,
		m.MessageID,
		m.Type,
		m.Date.Format(time.RFC1123Z),
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
