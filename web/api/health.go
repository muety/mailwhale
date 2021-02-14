package api

import (
	"fmt"
	conf "github.com/muety/mailwhale/config"
	"net/http"
)

type HealthHandler struct {
	config *conf.Config
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		config: conf.Get(),
	}
}

func (h *HealthHandler) Register(mux *http.ServeMux) {
	mux.Handle("/api/health", h)
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text-plain")
	w.Write([]byte(fmt.Sprintf("app=1\nversion=%s", h.config.Version)))
}
