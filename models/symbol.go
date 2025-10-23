package models

import (
	"time"

	"gorm.io/gorm"
)

type Symbol struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Exchange    string    `gorm:"not null;index" json:"exchange"`
	Type        string    `gorm:"not null;index" json:"type"` // "spot" or "futures"
	Symbol      string    `gorm:"not null;index" json:"symbol"`
	Combination string    `gorm:"not null;unique" json:"combination"` // exchange-type-symbol
	CreatedAt   time.Time `json:"created_at"`
}

func (s *Symbol) BeforeCreate(tx *gorm.DB) error {
	s.Combination = s.Exchange + "-" + s.Type + "-" + s.Symbol
	return nil
}
