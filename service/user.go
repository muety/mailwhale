package service

import (
	"errors"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/util"
	"github.com/timshannon/bolthold"
)

type UserService struct {
	config *conf.Config
	store  *bolthold.Store
}

func NewUserService() *UserService {
	return &UserService{
		config: conf.Get(),
		store:  conf.GetStore(),
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

func (s *UserService) Create(signup *types.Signup) (*types.User, error) {
	user := &types.User{
		ID:       signup.Email,
		Password: util.HashBcrypt(signup.Password, s.config.Security.Pepper),
	}
	if !user.IsValid() {
		return nil, errors.New("can't create user (empty password or invalid e-mail address)")
	}
	if err := s.store.Insert(user.ID, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(user *types.User) (*types.User, error) {
	user.Password = util.HashBcrypt(user.Password, s.config.Security.Pepper)
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
