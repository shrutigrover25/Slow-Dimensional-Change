package payment

import (
  "time"
  "github.com/google/uuid"
)

type PaymentLineItem struct {
  ID           uuid.UUID `gorm:"type:uuid"`
  UID          uuid.UUID `gorm:"type:uuid;primaryKey"`
  Version      int
  ContractorID uuid.UUID
  Amount       float64
  IssuedAt     time.Time
  CreatedAt    time.Time
  UpdatedAt    time.Time
}

func (PaymentLineItem) TableName() string { return "payment_line_items" }
func (p PaymentLineItem) GetID() string { return p.ID.String() }
func (p PaymentLineItem) GetUID() string { return p.UID.String() }
func (p PaymentLineItem) GetVersion() int { return p.Version }
func (p PaymentLineItem) CopyForNewVersion() PaymentLineItem{
  return PaymentLineItem{
    ID:           p.ID,
    ContractorID: p.ContractorID,
    Amount:       p.Amount,
    IssuedAt:     p.IssuedAt,
    Version:      p.Version + 1,
  }
}
