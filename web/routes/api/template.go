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

const routeTemplate = "/api/template"

type TemplateHandler struct {
	config          *conf.Config
	templateService *service.TemplateService
	clientService   *service.ClientService
	userService     *service.UserService
}

func NewTemplateHandler() *TemplateHandler {
	return &TemplateHandler{
		config:          conf.Get(),
		templateService: service.NewTemplateService(),
		clientService:   service.NewClientService(),
		userService:     service.NewUserService(),
	}
}

func (h *TemplateHandler) Register(router *mux.Router) {
	r := router.PathPrefix(routeTemplate).Subrouter()
	r.Use(
		handlers.NewAuthMiddleware(h.clientService, h.userService, []string{types.PermissionManageTemplate}),
	)
	r.Path("/{id}").Methods(http.MethodGet).HandlerFunc(h.getById)
	r.Path("/{id}").Methods(http.MethodDelete).HandlerFunc(h.delete)
	r.Path("/{id}/rendered").Methods(http.MethodPost).HandlerFunc(h.render)
	r.Path("/").Methods(http.MethodGet).HandlerFunc(h.get)
	r.Path("/").Methods(http.MethodPost).HandlerFunc(h.post)
}

func (h *TemplateHandler) get(w http.ResponseWriter, r *http.Request) {
	reqClient := r.Context().Value(conf.KeyClient).(*types.Client)
	templates, err := h.templateService.GetByUser(reqClient.UserId)
	if err != nil {
		util.RespondError(w, r, http.StatusInternalServerError, err)
		return
	}
	util.RespondJson(w, http.StatusOK, templates)
}

func (h *TemplateHandler) getById(w http.ResponseWriter, r *http.Request) {
	reqClient := r.Context().Value(conf.KeyClient).(*types.Client)
	template, err := h.templateService.GetById(mux.Vars(r)["id"])
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}
	if template.UserId != reqClient.UserId {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}
	util.RespondJson(w, http.StatusOK, template)
}

func (h *TemplateHandler) post(w http.ResponseWriter, r *http.Request) {
	reqClient := r.Context().Value(conf.KeyClient).(*types.Client)

	var payload types.Template
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest, err)
		return
	}

	payload.UserId = reqClient.UserId

	template, err := h.templateService.Create(&payload)
	if err != nil {
		util.RespondError(w, r, http.StatusConflict, err)
		return
	}
	util.RespondJson(w, http.StatusCreated, template)
}

func (h *TemplateHandler) delete(w http.ResponseWriter, r *http.Request) {
	reqClient := r.Context().Value(conf.KeyClient).(*types.Client)

	template, err := h.templateService.GetById(mux.Vars(r)["id"])
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}

	if reqClient.UserId != template.UserId {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}

	if err := h.clientService.Delete(mux.Vars(r)["id"]); err != nil {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}

	util.RespondEmpty(w, r, http.StatusNoContent)
}

func (h *TemplateHandler) render(w http.ResponseWriter, r *http.Request) {
	reqClient := r.Context().Value(conf.KeyClient).(*types.Client)
	template, err := h.templateService.GetById(mux.Vars(r)["id"])
	if err != nil {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}
	if template.UserId != reqClient.UserId {
		util.RespondError(w, r, http.StatusNotFound, err)
		return
	}

	var payload map[string]string
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest, err)
		return
	}

	util.RespondHtml(w, http.StatusOK, template.FillContent(payload))
}
