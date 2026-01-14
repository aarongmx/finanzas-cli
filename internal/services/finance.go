package services

import (
	"time"

	"github.com/aarongmx/finanzas-personales/internal/models"

	"gorm.io/gorm"
)

type FinanceService struct {
	db *gorm.DB
}

func NewFinanceService(db *gorm.DB) *FinanceService {
	return &FinanceService{db: db}
}

type Summary struct {
	TotalIncome  float64
	TotalExpense float64
	Balance      float64
}

func normalizeDate(t *time.Time) time.Time {
	if t == nil {
		return time.Now()
	}
	return *t
}

func (s *FinanceService) AddExpense(
	amount float64,
	category, note string,
	occurredAt *time.Time,
) error {
	tx := models.Transaction{
		Type:       models.Expense,
		Amount:     amount,
		Category:   models.Category{},
		Note:       note,
		OccurredAt: normalizeDate(occurredAt),
	}

	return s.db.Create(&tx).Error
}

func (s *FinanceService) AddIncome(
	amount float64,
	category string,
	occurredAt *time.Time,
) error {
	tx := models.Transaction{
		Type:       models.Income,
		Amount:     amount,
		Category:   models.Category{},
		OccurredAt: normalizeDate(occurredAt),
	}

	return s.db.Create(&tx).Error
}

func (s *FinanceService) GetSummary() (Summary, error) {
	var income float64
	var expense float64

	if err := s.db.
		Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("type = ?", models.Income).
		Scan(&income).Error; err != nil {
		return Summary{}, err
	}

	if err := s.db.
		Model(&models.Transaction{}).
		Select("COALESCE(SUM(amount), 0)").
		Where("type = ?", models.Expense).
		Scan(&expense).Error; err != nil {
		return Summary{}, err
	}

	return Summary{
		TotalIncome:  income,
		TotalExpense: expense,
		Balance:      income - expense,
	}, nil
}
