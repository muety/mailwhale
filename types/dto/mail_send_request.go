package dto

import "github.com/muety/mailwhale/types"

type MailSendRequest struct {
	From         types.MailAddress   `json:"from"`
	To           types.MailAddresses `json:"to"`
	Subject      string              `json:"subject"`
	Text         string              `json:"text"`
	Html         string              `json:"html"`
	TemplateId   string              `json:"template_id"`
	TemplateVars map[string]string   `json:"template_vars"`
}
