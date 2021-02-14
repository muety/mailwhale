package api

import (
	"encoding/json"
	"github.com/emvi/logbuch"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/types/dto"
	"net/http"
)

type MailHandler struct {
	config      *conf.Config
	sendService *service.SendService
}

func NewMailHandler(sendService *service.SendService) *MailHandler {
	return &MailHandler{
		config:      conf.Get(),
		sendService: sendService,
	}
}

func (h *MailHandler) Register(mux *http.ServeMux) {
	mux.Handle("/api/mail", h)
}

func (h *MailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.handleGet(w, r)
	} else if r.Method == http.MethodPost {
		h.handlePost(w, r)
	}
}

func (h *MailHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	w.WriteHeader(http.StatusMethodNotAllowed)
	_, _ = w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
}

func (h *MailHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	var payload dto.MailSendRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logbuch.Error("failed to decode mail request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}

	mail := &types.Mail{
		From:    payload.From,
		To:      payload.To,
		Subject: payload.Subject,
	}

	if payload.Text != "" {
		mail.WithText(payload.Text)
	}
	if payload.Html != "" {
		mail.WithHTML(payload.Html)
	}

	if err := h.sendService.Send(mail); err != nil {
		logbuch.Error("failed to send mail: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(http.StatusText(http.StatusOK)))
}
