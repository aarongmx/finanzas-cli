package services

import (
	"github.com/aarongmx/finanzas-personales/internal/models"

	"gorm.io/gorm"
)

type FinanceService struct {
	db *gorm.DB
}

func NewFinanceService(db *gorm.DB) *FinanceService {
	return &FinanceService{db: db}
}

func (s *FinanceService) AddExpense(
	amount float64,
	category, note string,
) error {
	tx := models.Transaction{
		Type:     models.Expense,
		Amount:   amount,
		Category: category,
		Note:     note,
	}

	return s.db.Create(&tx).Error
}

func (s *FinanceService) AddIncome(
	amount float64,
	category string,
) error {
	tx := models.Transaction{
		Type:     models.Income,
		Amount:   amount,
		Category: category,
	}

	return s.db.Create(&tx).Error
}
