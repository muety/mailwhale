package service

import (
	"github.com/google/uuid"
	conf "github.com/muety/mailwhale/config"
	"github.com/muety/mailwhale/types"
	"github.com/timshannon/bolthold"
)

type TemplateService struct {
	config *conf.Config
	store  *bolthold.Store
}

func NewTemplateService() *TemplateService {
	return &TemplateService{
		config: conf.Get(),
		store:  conf.GetStore(),
	}
}

func (s *TemplateService) GetByUser(userId string) (templates []*types.Template, err error) {
	err = s.store.Find(&templates, bolthold.Where("UserId").Eq(userId).Index("UserId"))
	if templates == nil {
		templates = make([]*types.Template, 0)
	}
	return templates, err
}

func (s *TemplateService) GetById(id string) (*types.Template, error) {
	var template types.Template
	err := s.store.Get(id, &template)
	return &template, err
}

func (s *TemplateService) Create(template *types.Template) (*types.Template, error) {
	template.ID = s.createId()
	if err := s.store.Insert(template.ID, template); err != nil {
		return nil, err
	}
	return template, nil
}

func (s *TemplateService) Update(template *types.Template) (*types.Template, error) {
	if err := s.store.Update(template.ID, template); err != nil {
		return nil, err
	}
	return template, nil
}

func (s *TemplateService) Delete(id string) error {
	return s.store.Delete(id, &types.Template{})
}

func (s *TemplateService) createId() string {
	return uuid.New().String()
}
