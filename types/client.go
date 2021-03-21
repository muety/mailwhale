package types

import (
	"errors"
	"fmt"
	"github.com/muety/mailwhale/util"
)

const (
	PermissionSendMail       = "send_mail"
	PermissionManageClient   = "manage_client"
	PermissionManageTemplate = "manage_template"
)

func AllPermissions() []string {
	return []string{
		PermissionSendMail,
		PermissionManageClient,
		PermissionManageTemplate,
	}
}

type Client struct {
	ID             string        `json:"id"`
	Description    string        `json:"description"`
	UserId         string        `json:"-" boltholdIndex:"UserId"`
	Permissions    []string      `json:"permissions"`
	DefaultSender  MailAddress   `json:"default_sender"`
	AllowedSenders MailAddresses `json:"allowed_senders"` // none (except default) means all
	ApiKey         *string       `json:"api_key"`         // caution: usually you want to hide this!
}

func (c *Client) HasPermission(permission string) bool {
	if c.Permissions == nil {
		return false
	}
	for _, p := range c.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

func (c *Client) HasPermissionAnyOf(permissions []string) bool {
	for _, p := range permissions {
		if c.HasPermission(p) {
			return true
		}
	}
	return false
}

func (c *Client) AllowsSender(sender MailAddress) bool {
	if sender == "" {
		return false
	}
	if c.AllowedSenders == nil || len(c.AllowedSenders) == 0 {
		return true
	}
	for _, m := range c.AllowedSenders {
		if m.Raw() == sender.Raw() {
			return true
		}
	}
	return false
}

func (c *Client) Sanitize() *Client {
	c.ApiKey = nil
	return c
}

func (c *Client) Validate() error {
	allPerms := AllPermissions()
	if c.Permissions == nil || len(c.Permissions) == 0 {
		return errors.New(fmt.Sprintf("client needs to be given at least one type of privileges, available are: %v", allPerms))
	}
	for _, p := range c.Permissions {
		if !util.ContainsString(p, allPerms) {
			return errors.New(fmt.Sprintf("permission '%s' is invalid", p))
		}
	}
	if c.DefaultSender != "" && !c.DefaultSender.Valid() {
		return errors.New("invalid default sender address")
	}
	if c.AllowedSenders != nil && len(c.AllowedSenders) > 0 {
		for _, e := range c.AllowedSenders {
			if !e.Valid() {
				return errors.New("invalid allowed sender address")
			}
		}
	}
	return nil
}
