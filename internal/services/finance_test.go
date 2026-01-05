package services_test

import (
	"testing"
	"time"

	"github.com/aarongmx/finanzas-personales/internal/models"
	"github.com/aarongmx/finanzas-personales/internal/services"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("no se pudo abrir db: %v", err)
	}

	if err := db.AutoMigrate(&models.Transaction{}); err != nil {
		t.Fatalf("no se pudo migrar: %v", err)
	}

	return db
}

func TestAddIncome(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewFinanceService(db)

	fixedDate := time.Date(2024, 12, 31, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		amount     float64
		category   string
		date       *time.Time
		expectDate bool
	}{
		{
			name:       "income with explicit date",
			amount:     1000,
			category:   "salary",
			date:       &fixedDate,
			expectDate: true,
		},
		{
			name:       "income without date uses now",
			amount:     500,
			category:   "freelance",
			date:       nil,
			expectDate: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.AddIncome(tt.amount, tt.category, tt.date)
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}

			var tx models.Transaction
			if err := db.Last(&tx).Error; err != nil {
				t.Fatalf("no se pudo leer transacción: %v", err)
			}

			if tx.Amount != tt.amount {
				t.Errorf("amount esperado %.2f, obtenido %.2f", tt.amount, tx.Amount)
			}

			if tx.Category != tt.category {
				t.Errorf("category esperada %s, obtenida %s", tt.category, tx.Category)
			}

			if tx.Type != models.Income {
				t.Errorf("type esperado income, obtenido %s", tx.Type)
			}

			if tt.expectDate && !tx.OccurredAt.Equal(*tt.date) {
				t.Errorf("fecha incorrecta")
			}
		})
	}
}

func TestAddExpense(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewFinanceService(db)

	fixedDate := time.Date(2025, 1, 2, 8, 30, 0, 0, time.UTC)

	err := service.AddExpense(
		250,
		"food",
		"groceries",
		&fixedDate,
	)

	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	var tx models.Transaction
	if err := db.First(&tx).Error; err != nil {
		t.Fatalf("no se pudo leer gasto: %v", err)
	}

	if tx.Type != models.Expense {
		t.Errorf("type esperado expense, obtenido %s", tx.Type)
	}

	if tx.Note != "groceries" {
		t.Errorf("nota incorrecta")
	}

	if !tx.OccurredAt.Equal(fixedDate) {
		t.Errorf("fecha incorrecta")
	}
}

func TestGetSummary(t *testing.T) {
	db := setupTestDB(t)
	service := services.NewFinanceService(db)

	tests := []struct {
		name            string
		transactions    []models.Transaction
		expectedIncome  float64
		expectedExpense float64
	}{
		{
			name: "mixed income and expenses",
			transactions: []models.Transaction{
				{Type: models.Income, Amount: 2000},
				{Type: models.Expense, Amount: 500},
				{Type: models.Expense, Amount: 300},
			},
			expectedIncome:  2000,
			expectedExpense: 800,
		},
		{
			name: "only income",
			transactions: []models.Transaction{
				{Type: models.Income, Amount: 1500},
			},
			expectedIncome:  1500,
			expectedExpense: 0,
		},
		{
			name:            "empty database",
			transactions:    nil,
			expectedIncome:  0,
			expectedExpense: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			db.Exec("DELETE FROM transactions")

			for _, tx := range tt.transactions {
				tx.OccurredAt = time.Now()
				if err := db.Create(&tx).Error; err != nil {
					t.Fatalf("error insertando tx: %v", err)
				}
			}

			summary, err := service.GetSummary()
			if err != nil {
				t.Fatalf("error inesperado: %v", err)
			}

			if summary.TotalIncome != tt.expectedIncome {
				t.Errorf("income esperado %.2f, obtenido %.2f",
					tt.expectedIncome, summary.TotalIncome)
			}

			if summary.TotalExpense != tt.expectedExpense {
				t.Errorf("expense esperado %.2f, obtenido %.2f",
					tt.expectedExpense, summary.TotalExpense)
			}

			expectedBalance := tt.expectedIncome - tt.expectedExpense
			if summary.Balance != expectedBalance {
				t.Errorf("balance esperado %.2f, obtenido %.2f",
					expectedBalance, summary.Balance)
			}
		})
	}
}
