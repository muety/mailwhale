package service

import (
	"errors"
	"github.com/mileusna/spf"
	conf "github.com/muety/mailwhale/config"
	"net"
	"strings"
)

type SpfService struct {
	config *conf.Config
}

func NewSpfService() *SpfService {
	return &SpfService{
		config: conf.Get(),
	}
}

func (s *SpfService) Validate(senderAddress string) error {
	if err := s.checkDelegate(senderAddress); err != nil {
		if err := s.checkIp(senderAddress); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (s *SpfService) checkIp(senderAddress string) error {
	domain := strings.Split(senderAddress, "@")[1]
	for _, ip := range s.config.Mail.SPF.AuthorizedIPs {
		if result := spf.CheckHost(net.ParseIP(ip), domain, senderAddress, ""); result == spf.Pass {
			return nil
		}
	}
	return errors.New("spf ip check did not pass")
}

func (s *SpfService) checkDelegate(senderAddress string) error {
	domain := strings.Split(senderAddress, "@")[1]
	spfRecord, result := spf.LookupSPF(domain)
	if result == spf.None || result == spf.TempError || result == spf.PermError {
		return errors.New("spf lookup failed")
	}
	for _, d := range s.config.Mail.SPF.AuthorizedDelegates {
		if strings.Contains(spfRecord, "include:"+d) {
			return nil
		}
	}
	return errors.New("spf delegate check did not pass")
}
