package timelog

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mercor/internal/scd"
)

type Repository interface {
	Insert(t Timelog) (Timelog, error)
	FindByUID(uid string) (Timelog, error)
	Update(uid string, updated Timelog) (Timelog, error)
	SoftDelete(uid string) error
	FindLatestByContractor(contractorID uuid.UUID) ([]Timelog, error)
}

type repo struct {
	scd *scd.SCDManager[Timelog]
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{scd: scd.NewManager[Timelog](db)}
}

func (r *repo) Insert(t Timelog) (Timelog, error) {
	err := r.scd.Insert(t)
	return t, err
}

func (r *repo) FindByUID(uid string) (Timelog, error) {
	return r.scd.FindByUID(uid)
}

func (r *repo) Update(uid string, updated Timelog) (Timelog, error) {
	old, err := r.scd.FindByUID(uid)
	if err != nil {
		return Timelog{}, err
	}
	newVer := old.CopyForNewVersion()
	newVer.StartTime = updated.StartTime
	newVer.EndTime = updated.EndTime
	newVer.ContractorID = updated.ContractorID
	err = r.scd.Insert(newVer)
	return newVer, err
}

func (r *repo) SoftDelete(uid string) error {
	old, err := r.FindByUID(uid)
	if err != nil {
		return err
	}
	newVer := old.CopyForNewVersion()
	newVer.EndTime = newVer.StartTime // Mark invalid
	return r.scd.Insert(newVer)
}

func (r *repo) FindLatestByContractor(contractorID uuid.UUID) ([]Timelog, error) {
	var list []Timelog
	err := r.scd.GetLatest().Where("contractor_id = ?", contractorID).Find(&list).Error
	return list, err
}
