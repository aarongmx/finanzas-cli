package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aarongmx/finanzas-personales/internal/db"
	"github.com/aarongmx/finanzas-personales/internal/services"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: finanzas [add-expense|add-income]")
		return
	}

	dbConn, err := db.Connect("finanzas.db")
	if err != nil {
		panic(err)
	}

	service := services.NewFinanceService(dbConn)

	switch os.Args[1] {

	case "add-expense":
		if len(os.Args) < 5 {
			fmt.Println("Uso: finanzas add-expense MONTO CATEGORIA [NOTA]")
			return
		}

		amount, _ := strconv.ParseFloat(os.Args[2], 64)
		category := os.Args[3]
		note := os.Args[4]

		if err := service.AddExpense(amount, category, note); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("✅ Gasto registrado")

	case "add-income":
		if len(os.Args) < 4 {
			fmt.Println("Uso: finanzas add-income MONTO CATEGORIA")
			return
		}

		amount, _ := strconv.ParseFloat(os.Args[2], 64)
		category := os.Args[3]

		if err := service.AddIncome(amount, category); err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("✅ Ingreso registrado")

	default:
		fmt.Println("Comando no reconocido")
	}
}
