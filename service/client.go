package service

import (
	"github.com/google/uuid"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/types"
	"github.com/muety/mailwhale/util"
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

func (s *ClientService) GetByName(name string) (*types.Client, error) {
	var client types.Client
	err := s.store.Get(name, &client)
	return &client, err
}

func (s *ClientService) Create(client *types.Client) (*types.Client, error) {
	client, clientDto := s.preprocess(client)
	if err := s.store.Insert(client.Name, client); err != nil {
		return nil, err
	}
	return clientDto, nil
}

func (s *ClientService) Update(client *types.Client) (*types.Client, error) {
	client, clientDto := s.preprocess(client)
	if err := s.store.Update(client.Name, client); err != nil {
		return nil, err
	}
	return clientDto, nil
}

func (s *ClientService) Delete(name string) error {
	return s.store.Delete(name, &types.Client{})
}

func (s *ClientService) createApiKey() (key, hash string) {
	key = uuid.New().String()
	hash = util.HashBcrypt(key, s.config.Security.Pepper)
	return key, hash
}

func (s *ClientService) preprocess(client *types.Client) (*types.Client, *types.Client) {
	apiKey, hash := s.createApiKey()
	client.ApiKey = &hash

	if client.Permissions == nil {
		client.Permissions = []string{}
	}

	if client.AllowedSenders == nil {
		client.AllowedSenders = []types.MailAddress{}
	}

	if client.DefaultSender != "" && (!client.AllowsSender(client.DefaultSender) || len(client.AllowedSenders) == 0) {
		client.AllowedSenders = append(client.AllowedSenders, types.MailAddress(client.DefaultSender.Raw()))
	}

	clientDto := *client
	clientDto.ApiKey = &apiKey

	return client, &clientDto
}
