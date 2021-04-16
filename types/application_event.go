package types

import "time"

type EventType int

const (
	MailSent EventType = iota
)

type ApplicationEvent struct {
	Type      EventType   `json:"event_type"`
	UserId    string      `json:"user_id" boltholdIndex:"UserId"`
	ClientId  string      `json:"client_id" boltholdIndex:"ClientId"`
	CreatedAt time.Time   `json:"created_at"`
	Payload   interface{} `json:"payload"`
}

type MailSentPayload struct {
	From MailAddress   `json:"from"`
	To   MailAddresses `json:"to"`
	Size int           `json:"size"`
}

func (p *MailSentPayload) FromMail(mail *Mail) *MailSentPayload {
	p.From = mail.From
	p.To = mail.To
	p.Size = len([]byte(mail.Body))
	return p
}

// more event types and payloads can be added here
