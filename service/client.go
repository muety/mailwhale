package service

import (
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/types"
	"github.com/timshannon/bolthold"
	"time"
)

type ClientService struct {
	config       *conf.Config
	store        *bolthold.Store
	eventService *ApplicationEventService
}

func NewClientService() *ClientService {
	return &ClientService{
		config:       conf.Get(),
		store:        conf.GetStore(),
		eventService: NewApplicationEventService(),
	}
}

func (s *ClientService) GetByUser(userId string) (clients []*types.Client, err error) {
	err = s.store.Find(&clients, bolthold.Where("UserId").Eq(userId).Index("UserId"))
	if clients == nil {
		clients = make([]*types.Client, 0)
	}

	for _, c := range clients {
		events, _ := s.eventService.GetByClientAndType(c.ID, types.MailSent) // TODO: make more efficient
		c.CountMails = len(events)
	}

	return clients, err
}

func (s *ClientService) GetById(id string) (*types.Client, error) {
	var client types.Client
	err := s.store.Get(id, &client)
	if err == nil {
		events, _ := s.eventService.GetByClientAndType(client.ID, types.MailSent)
		client.CountMails = len(events)
	}
	return &client, err
}

func (s *ClientService) Create(client *types.Client) (*types.Client, error) {
	client, clientDto := s.preprocess(client)
	if err := s.store.Insert(client.ID, client); err != nil {
		return nil, err
	}
	return clientDto, nil
}

func (s *ClientService) Delete(id string) error {
	return s.store.Delete(id, &types.Client{})
}

func (s *ClientService) preprocess(client *types.Client) (*types.Client, *types.Client) {
	client.ID = types.NewClientId()
	client.CreatedAt = time.Now()
	client.CountMails = 0
	apiKey, hash := types.NewClientApiKey(s.config.Security.Pepper)
	client.ApiKey = &hash

	if client.Permissions == nil {
		client.Permissions = []string{}
	}

	clientDto := *client
	clientDto.ApiKey = &apiKey

	return client, &clientDto
}
