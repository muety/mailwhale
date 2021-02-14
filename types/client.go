package types

type Client struct {
	Name           string        `json:"name"`
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
