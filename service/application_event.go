package service

import (
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/types"
	"github.com/timshannon/bolthold"
	"time"
)

type ApplicationEventService struct {
	config *conf.Config
	store  *bolthold.Store
}

func NewApplicationEventService() *ApplicationEventService {
	return &ApplicationEventService{
		config: conf.Get(),
		store:  conf.GetStore(),
	}
}

func (s *ApplicationEventService) Create(event *types.ApplicationEvent) (*types.ApplicationEvent, error) {
	event.CreatedAt = time.Now()
	if err := s.store.Insert(bolthold.NextSequence(), event); err != nil {
		return nil, err
	}
	return event, nil
}

func (s *ApplicationEventService) GetByUserAndType(userId string, ttype types.EventType) (events []*types.ApplicationEvent, err error) {
	err = s.store.Find(&events, bolthold.Where("UserId").Eq(userId).Index("UserId"))
	if events == nil {
		events = make([]*types.ApplicationEvent, 0)
	}
	// TODO: use indexed query for filtering
	return s.FilterByType(events, ttype), err
}

func (s *ApplicationEventService) GetByClientAndType(clientId string, ttype types.EventType) (events []*types.ApplicationEvent, err error) {
	err = s.store.Find(&events, bolthold.Where("ClientId").Eq(clientId).Index("ClientId"))
	if events == nil {
		events = make([]*types.ApplicationEvent, 0)
	}
	// TODO: use indexed query for filtering
	return s.FilterByType(events, ttype), err
}

func (s *ApplicationEventService) FilterByType(events []*types.ApplicationEvent, ttype types.EventType) []*types.ApplicationEvent {
	filtered := make([]*types.ApplicationEvent, 0, len(events))
	for _, e := range events {
		if e.Type == ttype {
			filtered = append(filtered, e)
		}
	}
	return filtered
}
