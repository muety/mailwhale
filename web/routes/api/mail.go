package api

import (
	"encoding/json"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/types/dto"
	"github.com/muety/mailwhale/util"
	"net/http"
)

const routeMail = "/api/mail"

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
	mux.Handle(routeMail, h)
}

func (h *MailHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.get(w, r)
	} else if r.Method == http.MethodPost {
		h.post(w, r)
	}
}

func (h *MailHandler) get(w http.ResponseWriter, r *http.Request) {
	// TODO: implement
	util.RespondError(w, r, http.StatusMethodNotAllowed)
}

func (h *MailHandler) post(w http.ResponseWriter, r *http.Request) {
	var payload dto.MailSendRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest)
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
		util.RespondError(w, r, http.StatusInternalServerError)
		return
	}

	util.RespondEmpty(w, r, 0)
}
