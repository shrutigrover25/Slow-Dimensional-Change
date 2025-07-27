package timelog

import (
	"time"

	"github.com/google/uuid"
)

type Timelog struct {
  ID           uuid.UUID `gorm:"type:uuid"`
  UID          uuid.UUID `gorm:"type:uuid;primaryKey"`
  Version      int
  ContractorID uuid.UUID
  StartTime    time.Time
  EndTime      time.Time
  CreatedAt    time.Time
  UpdatedAt    time.Time
}

func (Timelog) TableName() string { return "timelogs" }
func (t Timelog) GetID() string { return t.ID.String() }
func (t Timelog) GetUID() string { return t.UID.String() }
func (t Timelog) GetVersion() int { return t.Version }
func (t Timelog) CopyForNewVersion() Timelog {
  return Timelog{
    ID:           t.ID,
    ContractorID: t.ContractorID,
    StartTime:    t.StartTime,
    EndTime:      t.EndTime,
    UID:        uuid.New(),
		Version:    t.Version + 1,
  }
}
