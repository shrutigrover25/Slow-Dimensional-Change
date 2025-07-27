package payment

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mercor/internal/scd"
)

type Repository interface {
	Insert(p PaymentLineItem) (PaymentLineItem, error)
	FindByUID(uid string) (PaymentLineItem, error)
	Update(uid string, p PaymentLineItem) (PaymentLineItem, error)
	SoftDelete(uid string) error
	FindLatestByContractor(contractorID uuid.UUID) ([]PaymentLineItem, error)
}

type repo struct {
	scd *scd.SCDManager[PaymentLineItem]
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{scd: scd.NewManager[PaymentLineItem](db)}
}

func (r *repo) Insert(p PaymentLineItem) (PaymentLineItem, error) {
	err := r.scd.Insert(p)
	return p, err
}

func (r *repo) FindByUID(uid string) (PaymentLineItem, error) {
	return r.scd.FindByUID(uid)
}

func (r *repo) Update(uid string, updated PaymentLineItem) (PaymentLineItem, error) {
	old, err := r.scd.FindByUID(uid)
	if err != nil {
		return PaymentLineItem{}, err
	}
	newVer := old.CopyForNewVersion()
	newVer.Amount = updated.Amount
	newVer.IssuedAt = updated.IssuedAt
	newVer.ContractorID = updated.ContractorID
	err = r.scd.Insert(newVer)
	return newVer, err
}

func (r *repo) SoftDelete(uid string) error {
	old, err := r.scd.FindByUID(uid)
	if err != nil {
		return err
	}
	newVer := old.CopyForNewVersion()
	newVer.Amount = 0
	return r.scd.Insert(newVer)
}

func (r *repo) FindLatestByContractor(contractorID uuid.UUID) ([]PaymentLineItem, error) {
	var list []PaymentLineItem
	err := r.scd.GetLatest().Where("contractor_id = ?", contractorID).Find(&list).Error
	return list, err
}
