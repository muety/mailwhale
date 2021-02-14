package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/util"
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

func (h *ClientHandler) Register(router *httprouter.Router) {
	router.GET(routeClient+"/:name", h.getByName)
	router.PUT(routeClient+"/:name", h.put)
	router.DELETE(routeClient+"/:name", h.delete)
	router.GET(routeClient, h.getAll)
	router.POST(routeClient, h.post)
}

func (h *ClientHandler) getAll(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	clients, err := h.clientService.GetAll()
	if err != nil {
		util.RespondError(w, r, http.StatusInternalServerError)
		return
	}
	util.RespondJson(w, http.StatusOK, clients)
}

func (h *ClientHandler) getByName(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, err := h.clientService.GetByName(ps.ByName("name"))
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound)
		return
	}
	util.RespondJson(w, http.StatusOK, client)
}

func (h *ClientHandler) post(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var payload types.Client
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest)
		return
	}

	client, err := h.clientService.Create(&payload)
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound)
		return
	}
	util.RespondJson(w, http.StatusOK, client)
}

func (h *ClientHandler) put(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var payload types.Client
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest)
		return
	}
	payload.Name = ps.ByName("name")

	client, err := h.clientService.Update(&payload)
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound)
		return
	}
	util.RespondJson(w, http.StatusOK, client)
}

func (h *ClientHandler) delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := h.clientService.Delete(ps.ByName("name"))
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound)
		return
	}
	util.RespondEmpty(w, r, http.StatusNoContent)
}
