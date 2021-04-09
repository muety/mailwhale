package handlers

import (
	"context"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/util"
	"net/http"
)

type AuthMiddleware struct {
	handler       http.Handler
	config        *conf.Config
	clientService *service.ClientService
	userService   *service.UserService
	permissions   []string
}

func NewAuthMiddleware(clientService *service.ClientService, userService *service.UserService, permissions []string) func(handler http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &AuthMiddleware{
			handler:       h,
			config:        conf.Get(),
			clientService: clientService,
			userService:   userService,
			permissions:   permissions,
		}
	}
}

func (m *AuthMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	clientOrUser, credentials, ok := r.BasicAuth()
	if !ok {
		util.RespondEmpty(w, r, http.StatusUnauthorized)
		return
	}

	var user *types.User
	var client *types.Client

	if util.IsEmail(clientOrUser) {
		// Case 1: Principal is a user (users have all privileges)
		u, err := m.userService.GetById(clientOrUser)
		if err != nil {
			util.RespondEmpty(w, r, http.StatusUnauthorized)
			return
		}

		if !util.CompareBcrypt(u.Password, credentials, m.config.Security.Pepper) {
			util.RespondEmpty(w, r, http.StatusUnauthorized)
			return
		}

		// Case 1.1: User requested to be authenticated as a specific client
		if useClientId := r.Header.Get("X-Client-Id"); useClientId != "" {
			c, err := m.clientService.GetById(useClientId)
			if err != nil || c.UserId != u.ID {
				util.RespondEmpty(w, r, http.StatusUnauthorized)
				return
			}
			client = c
		} else {
			// Case 1.2: Create dummy client with all permissions
			client = &types.Client{
				ID:          types.NewClientIdFrom(u.ID),
				UserId:      u.ID,
				Permissions: types.AllPermissions(),
			}
		}
		user = u
	} else {
		// Case 2: Principal is an API client
		c, err := m.clientService.GetById(clientOrUser)
		if err != nil {
			util.RespondEmpty(w, r, http.StatusUnauthorized)
			return
		}

		if !util.CompareBcrypt(*c.ApiKey, credentials, m.config.Security.Pepper) {
			util.RespondEmpty(w, r, http.StatusUnauthorized)
			return
		}

		client = c
		user, _ = m.userService.GetById(client.UserId)
	}

	if m.permissions == nil || len(m.permissions) == 0 || client.HasPermissionAnyOf(m.permissions) {
		ctx := context.WithValue(r.Context(), conf.KeyClient, client)
		if user != nil {
			ctx = context.WithValue(ctx, conf.KeyUser, user)
		}
		m.handler.ServeHTTP(w, r.WithContext(ctx))
		return
	}

	util.RespondEmpty(w, r, http.StatusUnauthorized)
}
