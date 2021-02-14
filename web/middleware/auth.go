package middleware

import (
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/util"
	"net/http"
)

type AuthMiddleware struct {
	handle        http.Handler
	config        *conf.Config
	clientService *service.ClientService
	permissions   []string
}

func NewAuthMiddleware(clientService *service.ClientService, permissions []string) func(handlerFunc http.HandlerFunc) http.Handler {
	return func(h http.HandlerFunc) http.Handler {
		return &AuthMiddleware{
			handle:        h,
			config:        conf.Get(),
			clientService: clientService,
			permissions:   permissions,
		}
	}
}

func (m AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name, apiKey, ok := r.BasicAuth()
	if !ok {
		util.RespondError(w, r, http.StatusUnauthorized)
		return
	}

	client, err := m.clientService.GetByName(name)
	if err != nil || client == nil {
		util.RespondError(w, r, http.StatusUnauthorized)
		return
	}

	if util.CompareBcrypt(client.ApiKey, apiKey, m.config.Security.Pepper) {
		if m.permissions == nil || len(m.permissions) == 0 || client.HasPermissionAnyOf(m.permissions) {
			m.handle.ServeHTTP(w, r)
			return
		}
	}

	util.RespondError(w, r, http.StatusUnauthorized)
}
