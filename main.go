/*
TODO:
- Add Currency Conversion Capability using API
- Add stock viewing / update option using API
*/

package main

import (
	"fmt"
	"budget-dashboard/budget"
)

func main() {
	rows, headers, err := budget.ParseCSV("./ExampleBudget.csv")
	if err != nil {
		fmt.Println("Failed to parse CSV:", err)
		return
	}

	// for _, row := range budget {
	// 	budget.PrintBudgetRow(row, headers)
	// }

	monthlyTotals := budget.ComputeMonthlyTotals(rows)
	budget.PrintMonthlyTotals(monthlyTotals)

	yearlyCategoryTotals, yearlyTotal := budget.ComputeYearlyTotals(rows, headers)
	budget.PrintYearlyTotals(yearlyCategoryTotals, yearlyTotal)

	averages := budget.ComputeMonthlyAverages(rows, monthlyTotals)
	budget.PrintAllMonthAverages(averages)
}