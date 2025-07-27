package jobs

import (
	"github.com/google/uuid"
)

type Service interface {
	CreateJob(j Job) error
	GetByUID(uid string) (Job, error)
	Update(uid string, updated Job) (Job, error)
	UpdateStatus(uid, status string) (Job, error)
	GetActiveJobsByCompany(companyID string) ([]Job, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) CreateJob(j Job) error {
	return s.repo.Create(j)
}

func (s *service) GetByUID(uid string) (Job, error) {
	return s.repo.FindByUID(uid)
}

func (s *service) Update(uid string, updated Job) (Job, error) {
	return s.repo.Update(uid, updated)
}

func (s *service) UpdateStatus(uid, status string) (Job, error) {
	return s.repo.UpdateStatus(uid, status)
}

func (s *service) GetActiveJobsByCompany(companyID string) ([]Job, error) {
	id := uuid.MustParse(companyID)
	return s.repo.FindLatestByCompany(id)
}
