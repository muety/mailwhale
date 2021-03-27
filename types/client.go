package types

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/util"
	"strings"
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
	ID          string      `json:"id"`
	Description string      `json:"description"`
	UserId      string      `json:"-" boltholdIndex:"UserId"`
	Permissions []string    `json:"permissions"`
	Sender      MailAddress `json:"sender"`
	ApiKey      *string     `json:"api_key"` // caution: usually you want to hide this!
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

func (c *Client) Sanitize() *Client {
	c.ApiKey = nil
	c.Sender = c.SenderOrDefault()
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
	if c.Sender != "" && !c.Sender.Valid() {
		return errors.New("invalid default sender address")
	}
	return nil
}

func (c *Client) SenderOrDefault() MailAddress {
	if c.Sender != "" {
		return c.Sender
	}
	return c.DefaultSender()
}

func (c *Client) DefaultSender() MailAddress {
	return MailAddress(
		fmt.Sprintf(
			"%s+user@%s",
			strings.ToLower(c.ID[0:conf.ClientIdPrefixLength]),
			conf.Get().Mail.Domain,
		),
	)
}

func NewClientId() string {
	return base64.StdEncoding.EncodeToString([]byte(util.RandomString(16)))
}
func NewClientIdFrom(base string) string {
	return base64.StdEncoding.EncodeToString([]byte(util.RandomStringSeeded(16, base)))

}

func NewClientApiKey() (key, hash string) {
	key = uuid.New().String()
	hash = util.HashBcrypt(key, conf.Get().Security.Pepper)
	return key, hash
}
