package dto

import "github.com/muety/mailwhale/types"

type MailSendRequest struct {
	From         types.MailAddress
	To           types.MailAddresses
	Subject      string
	Text         string
	Html         string
	TemplateId   string
	TemplateVars map[string]string
}
