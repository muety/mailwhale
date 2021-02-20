package api

import (
	"fmt"
	"github.com/gorilla/mux"
	conf "github.com/muety/mailwhale/config"
	"net/http"
)

const routeHealth = "/api/health"

type HealthHandler struct {
	config *conf.Config
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		config: conf.Get(),
	}
}

func (h *HealthHandler) Register(router *mux.Router) {
	router.Path(routeHealth).Methods(http.MethodGet).HandlerFunc(h.Get)
}

func (h *HealthHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text-plain")
	w.Write([]byte(fmt.Sprintf("app=1\nversion=%s", h.config.Version)))
}
