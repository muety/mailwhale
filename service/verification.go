package service

import (
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/types"
	"github.com/timshannon/bolthold"
)

type VerificationService struct {
	config *conf.Config
	store  *bolthold.Store
}

func NewVerificationService() *VerificationService {
	return &VerificationService{
		config: conf.Get(),
		store:  conf.GetStore(),
	}
}

func (s *VerificationService) GetByToken(token string) (*types.Verification, error) {
	var verification types.Verification
	err := s.store.Get(token, &verification)
	return &verification, err
}

func (s *VerificationService) Create(verification *types.Verification) (*types.Verification, error) {
	if err := s.store.Insert(verification.Token, verification); err != nil {
		return nil, err
	}
	return verification, nil
}

func (s *VerificationService) Delete(token string) error {
	return s.store.Delete(token, &types.Verification{})
}
