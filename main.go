package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aarongmx/finanzas-personales/internal/db"
	"github.com/aarongmx/finanzas-personales/internal/services"
)

func parseDate(args []string) (*time.Time, error) {
	for i, arg := range args {
		if arg == "--date" && i+1 < len(args) {
			t, err := time.Parse("2006-01-02", args[i+1])
			if err != nil {
				return nil, err
			}
			return &t, nil
		}
	}
	return nil, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso:")
		fmt.Println("  finanzas add-income MONTO CATEGORIA [--date YYYY-MM-DD]")
		fmt.Println("  finanzas add-expense MONTO CATEGORIA NOTA [--date YYYY-MM-DD]")
		fmt.Println("  finanzas summary")
		return
	}

	dbConn, err := db.Connect("finanzas.db")
	if err != nil {
		panic(err)
	}

	service := services.NewFinanceService(dbConn)
	command := os.Args[1]

	switch command {

	case "add-income":
		if len(os.Args) < 4 {
			fmt.Println("Uso: finanzas add-income MONTO CATEGORIA [--date YYYY-MM-DD]")
			return
		}

		amount, err := strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			fmt.Println("Monto inválido")
			return
		}

		category := os.Args[3]

		occurredAt, err := parseDate(os.Args)
		if err != nil {
			fmt.Println("Fecha inválida. Usa YYYY-MM-DD")
			return
		}

		if err := service.AddIncome(amount, category, occurredAt); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("✅ Ingreso registrado")

	case "add-expense":
		if len(os.Args) < 5 {
			fmt.Println("Uso: finanzas add-expense MONTO CATEGORIA NOTA [--date YYYY-MM-DD]")
			return
		}

		amount, err := strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			fmt.Println("Monto inválido")
			return
		}

		category := os.Args[3]
		note := os.Args[4]

		occurredAt, err := parseDate(os.Args)
		if err != nil {
			fmt.Println("Fecha inválida. Usa YYYY-MM-DD")
			return
		}

		if err := service.AddExpense(amount, category, note, occurredAt); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("✅ Gasto registrado")

	case "summary":
		summary, err := service.GetSummary()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("📊 Resumen financiero")
		fmt.Printf("Ingresos: $%.2f\n", summary.TotalIncome)
		fmt.Printf("Gastos:   $%.2f\n", summary.TotalExpense)
		fmt.Printf("Balance:  $%.2f\n", summary.Balance)

	default:
		fmt.Println("Comando no reconocido")
	}
}
