package api

import (
	"encoding/json"
	"errors"
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/util"
	"github.com/muety/mailwhale/web/handlers"
	"net/http"
)

const routeUser = "/api/user"

type UserHandler struct {
	config        *conf.Config
	clientService *service.ClientService
	userService   *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		config:        conf.Get(),
		clientService: service.NewClientService(),
		userService:   service.NewUserService(),
	}
}

func (h *UserHandler) Register(router *mux.Router) {
	r := router.PathPrefix(routeUser).Subrouter()
	r.Path("").Methods(http.MethodPost).HandlerFunc(h.post)

	auth := handlers.NewAuthMiddleware(h.clientService, h.userService, []string{types.PermissionManageUser})
	r2 := r.PathPrefix("").Subrouter()
	r2.Use(auth)

	r2.Path("/{id}").Methods(http.MethodPut).HandlerFunc(h.update)
}

func (h *UserHandler) post(w http.ResponseWriter, r *http.Request) {
	if !h.config.Security.AllowSignup {
		util.RespondErrorMessage(w, r, http.StatusMethodNotAllowed, errors.New("user registration is disabled on this server"))
		return
	}

	var payload types.Signup
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest, err)
		return
	}

	user, err := h.userService.Create(&payload)
	if err != nil {
		util.RespondError(w, r, http.StatusBadRequest, err)
		return
	}

	logbuch.Info("created user '%s'", user.ID)
	util.RespondJson(w, http.StatusCreated, user)
}

func (h *UserHandler) update(w http.ResponseWriter, r *http.Request) {
	reqClient := r.Context().Value(conf.KeyClient).(*types.Client)

	var payload types.Signup
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest, err)
		return
	}

	if payload.Email != reqClient.UserId {
		util.RespondEmpty(w, r, http.StatusForbidden)
		return
	}

	user, err := h.userService.GetById(reqClient.UserId)
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}

	user.Password = payload.Password

	user, err = h.userService.Update(user)
	if err != nil {
		util.RespondError(w, r, http.StatusInternalServerError, err)
		return
	}

	logbuch.Info("updated user '%s'", user.ID)
	util.RespondJson(w, http.StatusOK, user)
}
