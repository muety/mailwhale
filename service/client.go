package service

import (
	"github.com/google/uuid"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/types"
	"github.com/timshannon/bolthold"
)

type ClientService struct {
	config *conf.Config
	store  *bolthold.Store
}

func NewClientService() *ClientService {
	return &ClientService{
		config: conf.Get(),
		store:  conf.GetStore(),
	}
}

func (s *ClientService) GetAll() (clients []*types.Client, err error) {
	err = s.store.Find(&clients, &bolthold.Query{})
	if clients == nil {
		clients = make([]*types.Client, 0)
	}
	return clients, err
}

func (s *ClientService) GetByName(name string) (client *types.Client, err error) {
	err = s.store.Get(name, client)
	return client, err
}

func (s *ClientService) GetByApiKey(key string) (client *types.Client, err error) {
	err = s.store.FindOne(client, bolthold.Where("ApiKey").Eq(key))
	return client, err
}

func (s *ClientService) Create(client *types.Client) (*types.Client, error) {
	client.ApiKey = uuid.New().String()
	if err := s.store.Insert(client.Name, client); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *ClientService) Update(client *types.Client) (*types.Client, error) {
	client.ApiKey = uuid.New().String()
	if err := s.store.Update(client.Name, client); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *ClientService) Delete(name string) error {
	return s.store.Delete(name, &types.Client{})
}
