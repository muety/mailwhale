package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/service"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/types/dto"
	"github.com/muety/mailwhale/util"
	"github.com/muety/mailwhale/web/middleware"
	"net/http"
)

const routeMail = "/api/mail"

type MailHandler struct {
	config        *conf.Config
	sendService   *service.SendService
	clientService *service.ClientService
}

func NewMailHandler(sendService *service.SendService, clientService *service.ClientService) *MailHandler {
	return &MailHandler{
		config:        conf.Get(),
		sendService:   sendService,
		clientService: clientService,
	}
}

func (h *MailHandler) Register(router *httprouter.Router) {
	auth := middleware.NewAuthMiddleware(h.clientService, []string{conf.PermissionSendMail})
	router.HandlerFunc(http.MethodPost, routeMail, auth(h.post).ServeHTTP)
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
