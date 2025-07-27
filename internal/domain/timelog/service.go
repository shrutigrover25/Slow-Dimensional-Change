package timelog

import (
	"github.com/google/uuid"
)

type Service interface {
	Create(t Timelog) (Timelog, error)
	GetByUID(uid string) (Timelog, error)
	Update(uid string, updated Timelog) (Timelog, error)
	Delete(uid string) error
	GetByContractor(id string) ([]Timelog, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(t Timelog) (Timelog, error) {
	t.UID = uuid.New()
	t.Version = 1
	return s.repo.Insert(t)
}

func (s *service) GetByUID(uid string) (Timelog, error) {
	return s.repo.FindByUID(uid)
}

func (s *service) Update(uid string, updated Timelog) (Timelog, error) {
	return s.repo.Update(uid, updated)
}

func (s *service) Delete(uid string) error {
	return s.repo.SoftDelete(uid)
}

func (s *service) GetByContractor(id string) ([]Timelog, error) {
	return s.repo.FindLatestByContractor(uuid.MustParse(id))
}
