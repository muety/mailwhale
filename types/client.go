package types

type Client struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	ApiKey      string   `json:"-"` // hashed api key
}

type ClientWithApiKey struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	ApiKey      string   `json:"api_key"`
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
