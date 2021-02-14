package middleware

import (
	"context"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/util"
	"net/http"
)

type AuthMiddleware struct {
	handler       http.Handler
	config        *conf.Config
	clientService *service.ClientService
	permissions   []string
}

func NewAuthMiddleware(clientService *service.ClientService, permissions []string) func(handler http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &AuthMiddleware{
			handler:       h,
			config:        conf.Get(),
			clientService: clientService,
			permissions:   permissions,
		}
	}
}

func (m AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name, apiKey, ok := r.BasicAuth()
	if !ok {
		util.RespondEmpty(w, r, http.StatusUnauthorized)
		return
	}

	client, err := m.clientService.GetByName(name)
	if err != nil || client == nil {
		util.RespondEmpty(w, r, http.StatusUnauthorized)
		return
	}

	if util.CompareBcrypt(*client.ApiKey, apiKey, m.config.Security.Pepper) {
		if m.permissions == nil || len(m.permissions) == 0 || client.HasPermissionAnyOf(m.permissions) {
			ctx := context.WithValue(r.Context(), conf.KeyClient, client)
			m.handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}
	}

	util.RespondEmpty(w, r, http.StatusUnauthorized)
}
