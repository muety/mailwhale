package handlers

import (
	"context"
	"fmt"
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
	// TODO: Securecookie auth for web ui

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

		// Create dummy client with all permissions
		user = u
		client = &types.Client{
			ID:            fmt.Sprintf("_%s", user.ID),
			UserId:        user.ID,
			Permissions:   types.AllPermissions(),
			DefaultSender: types.MailAddress(user.ID),
		}
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
