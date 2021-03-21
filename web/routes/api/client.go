package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/util"
	"github.com/muety/mailwhale/web/handlers"
	"net/http"
)

const routeClient = "/api/client"

type ClientHandler struct {
	config        *conf.Config
	clientService *service.ClientService
	userService   *service.UserService
}

func NewClientHandler() *ClientHandler {
	return &ClientHandler{
		config:        conf.Get(),
		clientService: service.NewClientService(),
		userService:   service.NewUserService(),
	}
}

func (h *ClientHandler) Register(router *mux.Router) {
	r := router.PathPrefix(routeClient).Subrouter()
	r.Use(
		handlers.NewAuthMiddleware(h.clientService, h.userService, []string{types.PermissionManageClient}),
	)
	r.Path("/{id}").Methods(http.MethodGet).HandlerFunc(h.getById)
	r.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(h.delete)
	r.Path("").Methods(http.MethodGet).HandlerFunc(h.get)
	r.Path("").Methods(http.MethodPost).HandlerFunc(h.post)
}

func (h *ClientHandler) get(w http.ResponseWriter, r *http.Request) {
	reqClient := r.Context().Value(conf.KeyClient).(*types.Client)

	clients, err := h.clientService.GetByUser(reqClient.UserId)
	if err != nil {
		util.RespondError(w, r, http.StatusInternalServerError, err)
		return
	}
	for _, c := range clients {
		c.Sanitize()
	}
	util.RespondJson(w, http.StatusOK, clients)
}

func (h *ClientHandler) getById(w http.ResponseWriter, r *http.Request) {
	reqClient := r.Context().Value(conf.KeyClient).(*types.Client)

	client, err := h.clientService.GetById(mux.Vars(r)["id"])
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}
	if client.UserId != reqClient.UserId {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}
	client = client.Sanitize()
	util.RespondJson(w, http.StatusOK, client)
}

func (h *ClientHandler) post(w http.ResponseWriter, r *http.Request) {
	reqClient := r.Context().Value(conf.KeyClient).(*types.Client)

	var payload types.Client
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest, err)
		return
	}

	payload.UserId = reqClient.UserId

	if err := payload.Validate(); err != nil {
		util.RespondErrorMessage(w, r, http.StatusBadRequest, err)
		return
	}

	client, err := h.clientService.Create(&payload)
	if err != nil {
		util.RespondError(w, r, http.StatusConflict, err)
		return
	}
	util.RespondJson(w, http.StatusCreated, client)
}

func (h *ClientHandler) delete(w http.ResponseWriter, r *http.Request) {
	reqClient := r.Context().Value(conf.KeyClient).(*types.Client)

	client, err := h.clientService.GetById(mux.Vars(r)["id"])
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}

	if reqClient.UserId != client.UserId {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}

	if err := h.clientService.Delete(mux.Vars(r)["id"]); err != nil {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}

	util.RespondEmpty(w, r, http.StatusNoContent)
}
