package service

import (
	"errors"
	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/types"
	"io"
)

type SendService struct {
	config *conf.Config
	auth   sasl.Client
}

func NewSendService() *SendService {
	config := conf.Get()
	return &SendService{
		config: config,
		auth: sasl.NewPlainClient(
			"",
			config.Smtp.Username,
			config.Smtp.Password,
		),
	}
}

func (s *SendService) Send(mail *types.Mail) error {
	return sendMail(
		s.config.Smtp.ConnStr(),
		s.config.Smtp.TLS,
		s.auth,
		mail.From.Raw(),
		mail.To.RawStrings(),
		mail.Reader(),
	)
}

func sendMail(addr string, tls bool, a sasl.Client, from string, to []string, r io.Reader) error {
	dial := smtp.Dial
	if tls {
		dial = func(addr string) (*smtp.Client, error) {
			return smtp.DialTLS(addr, nil)
		}
	}

	c, err := dial(addr)
	if err != nil {
		return err
	}

	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(nil); err != nil {
			return err
		}
	}
	if a != nil {
		if ok, _ := c.Extension("AUTH"); !ok {
			return errors.New("smtp: server doesn't support AUTH")
		}
		if err = c.Auth(a); err != nil {
			return err
		}
	}
	if err = c.Mail(from, nil); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
