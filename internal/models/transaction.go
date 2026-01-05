package models

import "time"

type TransactionType string

const (
	Income  TransactionType = "income"
	Expense TransactionType = "expense"
)

type Transaction struct {
	ID         uint            `gorm:"primaryKey"`
	Type       TransactionType `gorm:"size:10;not null"`
	Amount     float64         `gorm:"not null"`
	CategoryID uint            `gorm:"not null;index"`
	Category   Category        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
	Note       string          `gorm:"size:255"`
	OccurredAt time.Time       `gorm:"not null;index"`
	CreatedAt  time.Time
}
