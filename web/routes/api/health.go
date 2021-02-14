package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
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

func (h *HealthHandler) Register(router *httprouter.Router, _ *alice.Chain) {
	router.GET("/api/health", h.Get)
}

func (h *HealthHandler) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("content-type", "text-plain")
	w.Write([]byte(fmt.Sprintf("app=1\nversion=%s", h.config.Version)))
}
