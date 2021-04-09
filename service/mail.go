package service

import (
	"bytes"
	"fmt"
	"github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/types"
	"io/ioutil"
	"os"
	"text/template"
)

// Service with methods for sending system mails, not to be confused with SendService

const (
	tplPath             = "templates"
	tplNameVerifyUser   = "user_verification"
	tplNameVerifySender = "sender_verification"
)

type MailService struct {
	config      *config.Config
	sendService *SendService
}

func NewMailService() *MailService {
	return &MailService{
		config:      config.Get(),
		sendService: NewSendService(),
	}
}

func (s *MailService) SendUserVerification(user *types.User, token string) error {
	tpl, err := s.loadTemplate(tplNameVerifyUser)
	if err != nil {
		return err
	}

	type data struct {
		VerifyLink string
	}

	payload := &data{
		VerifyLink: fmt.Sprintf(
			"%s/api/user/verify?token=%s",
			s.config.Web.GetPublicUrl(),
			token,
		),
	}

	var rendered bytes.Buffer
	if err := tpl.Execute(&rendered, payload); err != nil {
		return err
	}

	mail := &types.Mail{
		From:    types.MailAddress(fmt.Sprintf("MailWhale System <system@%s>", s.config.Mail.Domain)),
		To:      []types.MailAddress{types.MailAddress(user.ID)},
		Subject: "Verify your MailWhale account",
	}
	mail.WithHTML(rendered.String())

	return s.sendService.Send(mail)
}

func (s *MailService) SendSenderVerification(user *types.User, sender types.SenderAddress, token string) error {
	tpl, err := s.loadTemplate(tplNameVerifySender)
	if err != nil {
		return err
	}

	type data struct {
		UserId        string
		SenderAddress string
		VerifyLink    string
	}

	verifyLink := fmt.Sprintf(
		"%s/api/user/verify?token=%s",
		s.config.Web.GetPublicUrl(),
		token,
	)
	payload := &data{
		UserId:        user.ID,
		SenderAddress: sender.Raw(),
		VerifyLink:    verifyLink,
	}

	var rendered bytes.Buffer
	if err := tpl.Execute(&rendered, payload); err != nil {
		return err
	}

	mail := &types.Mail{
		From:    types.MailAddress(fmt.Sprintf("MailWhale System <system@%s>", s.config.Mail.Domain)),
		To:      []types.MailAddress{sender.MailAddress},
		Subject: "Verify your e-mail address for MailWhale",
	}
	mail.WithHTML(rendered.String())

	return s.sendService.Send(mail)
}

func (s *MailService) loadTemplate(tplName string) (*template.Template, error) {
	tplFile, err := os.Open(fmt.Sprintf("%s/%s.tpl.html", tplPath, tplName))
	if err != nil {
		return nil, err
	}
	defer tplFile.Close()

	tplData, err := ioutil.ReadAll(tplFile)
	if err != nil {
		return nil, err
	}

	return template.New(tplName).Parse(string(tplData))
}
