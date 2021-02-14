package api

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
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

func (h *MailHandler) Register(router *httprouter.Router, baseChain *alice.Chain) {
	chain := baseChain.Extend(alice.New(
		middleware.NewAuthMiddleware(h.clientService, []string{conf.PermissionSendMail}),
	))
	router.Handler(http.MethodPost, routeMail, chain.ThenFunc(h.post))
}

func (h *MailHandler) post(w http.ResponseWriter, r *http.Request) {
	var payload dto.MailSendRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, r, http.StatusBadRequest, err)
		return
	}

	client := r.Context().Value(conf.KeyClient).(*types.Client)

	sender := payload.From
	if sender == "" {
		sender = client.DefaultSender
	}

	if !client.AllowsSender(sender) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte("sender not allowed"))
		return
	}

	mail := &types.Mail{
		From:    sender,
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
		util.RespondError(w, r, http.StatusInternalServerError, err)
		return
	}

	util.RespondEmpty(w, r, 0)
}
