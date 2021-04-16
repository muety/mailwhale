package service

import (
	"errors"
	"fmt"
	"github.com/emvi/logbuch"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/types/dto"
	"github.com/muety/mailwhale/util"
	"github.com/timshannon/bolthold"
	"strings"
	"time"
)

type UserService struct {
	config              *conf.Config
	store               *bolthold.Store
	spfService          *SpfService
	mailService         *MailService
	verificationService *VerificationService
}

func NewUserService() *UserService {
	return &UserService{
		config:              conf.Get(),
		store:               conf.GetStore(),
		spfService:          NewSpfService(),
		mailService:         NewMailService(),
		verificationService: NewVerificationService(),
	}
}

func (s *UserService) GetAll() (users []*types.User, err error) {
	err = s.store.Find(&users, &bolthold.Query{})
	if users == nil {
		users = make([]*types.User, 0)
	}
	return users, err
}

func (s *UserService) GetById(id string) (*types.User, error) {
	var user types.User
	err := s.store.Get(id, &user)
	return &user, err
}

func (s *UserService) Create(signup *dto.Signup) (*types.User, error) {
	user := &types.User{
		ID:        signup.Email,
		Password:  util.HashBcrypt(signup.Password, s.config.Security.Pepper),
		Senders:   []types.SenderAddress{},
		CreatedAt: time.Now(),
	}
	if !user.IsValid() {
		return nil, errors.New("can't create user (empty password or invalid e-mail address)")
	}
	if err := s.store.Insert(user.ID, user); err != nil {
		return nil, err
	}
	if s.config.Security.VerifyUsers {
		go s.verifyUser(user)
	}
	return user, nil
}

func (s *UserService) Update(user *types.User, update *types.User) (*types.User, error) {
	if update.Password != "" {
		user.Password = util.HashBcrypt(update.Password, s.config.Security.Pepper)
	}

	newSenders := s.extractNewSenders(user, update)
	if s.config.Security.VerifySenders {
		if err := s.spfCheckSenders(newSenders); err != nil {
			return nil, err
		}
		go s.verifySenders(user, newSenders)
	}
	user.Senders = update.Senders
	user.Verified = update.Verified

	if !user.IsValid() {
		return nil, errors.New("user data invalid")
	}
	if err := s.store.Update(user.ID, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Delete(id string) error {
	return s.store.Delete(id, &types.User{})
}

func (s *UserService) extractNewSenders(user *types.User, update *types.User) []types.SenderAddress {
	newSenders := make([]types.SenderAddress, 0)
	for _, sender := range update.Senders {
		if user.HasSender(sender.MailAddress) {
			continue
		}
		newSenders = append(newSenders, sender)
	}
	return newSenders
}

func (s *UserService) spfCheckSenders(senders []types.SenderAddress) error {
	// TODO: parallelize
	for _, sender := range senders {
		senderDomain := strings.Split(sender.MailAddress.Raw(), "@")[1]
		if err := s.spfService.Validate(senderDomain); err != nil {
			return errors.New(fmt.Sprintf("failed to verify spf entry for domain '%s'", senderDomain))
		}
	}
	return nil
}

func (s *UserService) verifyUser(user *types.User) error {
	verification, err := s.verificationService.Create(types.NewVerification(
		user,
		types.VerificationScopeUser,
		user.ID,
	))
	if err != nil {
		logbuch.Error("failed to create user verification token for '%s'", user.ID)
		return err
	}
	if err := s.mailService.SendUserVerification(user, verification.Token); err != nil {
		logbuch.Error("failed to send user verification to '%s'", user.ID)
		return err
	} else {
		logbuch.Info("sent user verification mail for '%s'", user.ID)
	}
	return nil
}

// generates verification tokens for senders addresses and sends them via mail
func (s *UserService) verifySenders(user *types.User, senders []types.SenderAddress) error {
	for _, sender := range senders {
		verification, err := s.verificationService.Create(types.NewVerification(
			user,
			types.VerificationScopeSender,
			sender.String(),
		))
		if err != nil {
			logbuch.Error("failed to create sender verification token for '%s'", sender.MailAddress.String())
			return err
		}
		if err := s.mailService.SendSenderVerification(user, sender, verification.Token); err != nil {
			logbuch.Error("failed to send sender verification to '%s'", sender.MailAddress.String())
			return err
		} else {
			logbuch.Info("sent sender verification mail for user '%s' to '%s'", user.ID, sender.String())
		}
	}
	return nil
}
