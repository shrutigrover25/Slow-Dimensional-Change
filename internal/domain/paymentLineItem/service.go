package payment

import "github.com/google/uuid"

type Service interface {
	Create(p PaymentLineItem) (PaymentLineItem, error)
	GetByUID(uid string) (PaymentLineItem, error)
	Update(uid string, p PaymentLineItem) (PaymentLineItem, error)
	Delete(uid string) error
	GetByContractor(id string) ([]PaymentLineItem, error)
}

type service struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &service{repo: r}
}

func (s *service) Create(p PaymentLineItem) (PaymentLineItem, error) {
	p.UID = uuid.New()
	p.Version = 1
	return s.repo.Insert(p)
}

func (s *service) GetByUID(uid string) (PaymentLineItem, error) {
	return s.repo.FindByUID(uid)
}

func (s *service) Update(uid string, p PaymentLineItem) (PaymentLineItem, error) {
	return s.repo.Update(uid, p)
}

func (s *service) Delete(uid string) error {
	return s.repo.SoftDelete(uid)
}

func (s *service) GetByContractor(id string) ([]PaymentLineItem, error) {
	return s.repo.FindLatestByContractor(uuid.MustParse(id))
}
