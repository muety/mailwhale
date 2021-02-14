package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/util"
	"github.com/muety/mailwhale/web/middleware"
	"net/http"
)

const routeClient = "/api/client"

type ClientHandler struct {
	config        *conf.Config
	clientService *service.ClientService
}

func NewClientHandler(clientService *service.ClientService) *ClientHandler {
	return &ClientHandler{
		config:        conf.Get(),
		clientService: clientService,
	}
}

func (h *ClientHandler) Register(router *httprouter.Router, baseChain *alice.Chain) {
	chain := baseChain.Extend(alice.New(
		middleware.NewAuthMiddleware(h.clientService, []string{conf.PermissionManageClient}),
	))
	router.Handler(http.MethodGet, routeClient+"/:name", chain.ThenFunc(h.getByName))
	router.Handler(http.MethodPut, routeClient+"/:name", chain.ThenFunc(h.put))
	router.Handler(http.MethodDelete, routeClient+"/:name", chain.ThenFunc(h.delete))
	router.Handler(http.MethodGet, routeClient, chain.ThenFunc(h.getAll))
	router.Handler(http.MethodPost, routeClient, chain.ThenFunc(h.post))
}

func (h *ClientHandler) getAll(w http.ResponseWriter, r *http.Request) {
	clients, err := h.clientService.GetAll()
	if err != nil {
		util.RespondError(w, r, http.StatusInternalServerError, err)
		return
	}
	util.RespondJson(w, http.StatusOK, clients)
}

func (h *ClientHandler) getByName(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())
	client, err := h.clientService.GetByName(ps.ByName("name"))
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}
	util.RespondJson(w, http.StatusOK, client)
}

func (h *ClientHandler) post(w http.ResponseWriter, r *http.Request) {
	var payload types.Client
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest, err)
		return
	}

	client, err := h.clientService.Create(&payload)
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}
	util.RespondJson(w, http.StatusOK, client)
}

func (h *ClientHandler) put(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())
	var payload types.Client
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest, err)
		return
	}
	payload.Name = ps.ByName("name")

	client, err := h.clientService.Update(&payload)
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}
	util.RespondJson(w, http.StatusOK, client)
}

func (h *ClientHandler) delete(w http.ResponseWriter, r *http.Request) {
	ps := httprouter.ParamsFromContext(r.Context())
	err := h.clientService.Delete(ps.ByName("name"))
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}
	util.RespondEmpty(w, r, http.StatusNoContent)
}
