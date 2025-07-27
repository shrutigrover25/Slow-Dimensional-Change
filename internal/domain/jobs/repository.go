package jobs

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mercor/internal/scd"
)

type Repository interface {
	Create(job Job) error
	FindByUID(uid string) (Job, error)
	Update(uid string, newJob Job) (Job, error)
	UpdateStatus(uid string, newStatus string) (Job, error)
	FindLatestByCompany(companyID uuid.UUID) ([]Job, error)
}

type repo struct {
	scd *scd.SCDManager[Job]
}

func NewRepository(db *gorm.DB) Repository {
	return &repo{scd: scd.NewManager[Job](db)}
}

func (r *repo) Create(j Job) error {
	return r.scd.Insert(j)
}

func (r *repo) FindByUID(uid string) (Job, error) {
	return r.scd.FindByUID(uid)
}

func (r *repo) Update(uid string, newJob Job) (Job, error) {
	old, err := r.FindByUID(uid)
	if err != nil {
		return Job{}, err
	}
	updated := old.CopyForNewVersion()
	updated.Title = newJob.Title
	updated.Rate = newJob.Rate
	updated.Status = newJob.Status
	updated.CompanyID = newJob.CompanyID
	updated.ContractorID = newJob.ContractorID
	err = r.scd.Insert(updated)
	return updated, err
}

func (r *repo) UpdateStatus(uid string, newStatus string) (Job, error) {
	old, err := r.FindByUID(uid)
	if err != nil {
		return Job{}, err
	}
	newItem := old.CopyForNewVersion()
	newItem.Status = newStatus
	err = r.scd.Insert(newItem)
	return newItem, err
}

func (r *repo) FindLatestByCompany(companyID uuid.UUID) ([]Job, error) {
	var jobs []Job
	err := r.scd.GetLatest().
		Where("company_id = ?", companyID).
		Where("status = ?", "active").
		Find(&jobs).Error
	return jobs, err
}
