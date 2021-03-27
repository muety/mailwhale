package service

import (
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

func (s *ClientService) GetByUser(userId string) (clients []*types.Client, err error) {
	err = s.store.Find(&clients, bolthold.Where("UserId").Eq(userId).Index("UserId"))
	if clients == nil {
		clients = make([]*types.Client, 0)
	}
	return clients, err
}

func (s *ClientService) GetById(id string) (*types.Client, error) {
	var client types.Client
	err := s.store.Get(id, &client)
	return &client, err
}

func (s *ClientService) Create(client *types.Client) (*types.Client, error) {
	client, clientDto := s.preprocess(client)
	if err := s.store.Insert(client.ID, client); err != nil {
		return nil, err
	}
	return clientDto.Sanitize(), nil
}

func (s *ClientService) Delete(id string) error {
	return s.store.Delete(id, &types.Client{})
}

func (s *ClientService) preprocess(client *types.Client) (*types.Client, *types.Client) {
	client.ID = types.NewClientId()
	apiKey, hash := types.NewClientApiKey()
	client.ApiKey = &hash

	if client.Permissions == nil {
		client.Permissions = []string{}
	}

	clientDto := *client
	clientDto.ApiKey = &apiKey

	return client, &clientDto
}
