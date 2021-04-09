package api

import (
	"encoding/json"
	"errors"
	"github.com/emvi/logbuch"
	"github.com/gorilla/mux"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/types/dto"
	"github.com/muety/mailwhale/util"
	"github.com/muety/mailwhale/web/handlers"
	"net/http"
)

const routeUser = "/api/user"

type UserHandler struct {
	config              *conf.Config
	clientService       *service.ClientService
	userService         *service.UserService
	verificationService *service.VerificationService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		config:              conf.Get(),
		clientService:       service.NewClientService(),
		userService:         service.NewUserService(),
		verificationService: service.NewVerificationService(),
	}
}

func (h *UserHandler) Register(router *mux.Router) {
	r := router.PathPrefix(routeUser).Subrouter()
	r.Path("").Methods(http.MethodPost).HandlerFunc(h.post)
	r.Path("/verify").Methods(http.MethodGet).HandlerFunc(h.verify)

	auth := handlers.NewAuthMiddleware(h.clientService, h.userService, []string{types.PermissionManageUser})
	r2 := r.PathPrefix("").Subrouter()
	r2.Use(auth)

	r2.Path("/me").Methods(http.MethodGet).HandlerFunc(h.getMe)
	r2.Path("/me").Methods(http.MethodPut).HandlerFunc(h.updateMe)
}

func (h *UserHandler) getMe(w http.ResponseWriter, r *http.Request) {
	var user *types.User
	if u := r.Context().Value(conf.KeyUser); u != nil {
		user = u.(*types.User)
	}
	if user == nil {
		util.RespondError(w, r, http.StatusNotFound, errors.New("user not found"))
		return
	}
	util.RespondJson(w, http.StatusOK, user)
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

func (h *UserHandler) updateMe(w http.ResponseWriter, r *http.Request) {
	var payload dto.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest, err)
		return
	}

	var user *types.User
	if u := r.Context().Value(conf.KeyUser); u != nil {
		user = u.(*types.User)
	}
	if user == nil {
		util.RespondError(w, r, http.StatusNotFound, errors.New("user not found"))
		return
	}

	update := *user
	update.Password = payload.Password
	update.Senders = payload.GetSenders(user)

	user, err := h.userService.Update(user, &update)
	if err != nil {
		util.RespondErrorMessage(w, r, http.StatusBadRequest, err)
		return
	}

	logbuch.Info("updated user '%s'", user.ID)
	util.RespondJson(w, http.StatusOK, user)
}

func (h *UserHandler) verify(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		util.RespondErrorMessage(w, r, http.StatusUnauthorized, errors.New("invalid verification token"))
		return
	}

	verification, err := h.verificationService.GetByToken(token)
	if err != nil {
		util.RespondErrorMessage(w, r, http.StatusUnauthorized, errors.New("invalid verification token"))
		return
	}

	user, err := h.userService.GetById(verification.UserId)
	if err != nil {
		util.RespondErrorMessage(w, r, http.StatusNotFound, errors.New("user not found"))
		return
	}

	if verification.Scope == types.VerificationScopeSender {
		update := *user // copy
		update.Password = ""

		for i, s := range update.Senders {
			if s.MailAddress.String() == verification.Subject {
				update.Senders[i].Verified = true
				break
			}
		}
		if _, err := h.userService.Update(user, &update); err != nil {
			util.RespondErrorMessage(w, r, http.StatusNotFound, err)
			return
		}
		go h.verificationService.Delete(verification.Token)
		logbuch.Info("verified sender address '%s' for user '%s'", verification.Subject, user.ID)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
