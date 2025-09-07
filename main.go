/*
TODO:
- Add Currency Conversion Capability using API
- Add stock viewing / update option using API
*/

package main

import (
	"fmt"
	"budget-dashboard/budget"
	"budget-dashboard/currencies"
	"budget-dashboard/stocks"
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
	budget.PrintHighestSpendingCategory(monthlyTotals)

	yearlyCategoryTotals, yearlyTotal := budget.ComputeYearlyTotals(rows, headers)
	budget.PrintYearlyTotals(yearlyCategoryTotals, yearlyTotal)

	averages := budget.ComputeMonthlyAverages(rows, monthlyTotals)
	budget.PrintAllMonthAverages(averages)

	rates, err := currencies.GetRates()
	if err != nil {
	    fmt.Println("Error fetching rates:", err)
	    return
	}
	fmt.Println("Rates:", rates)

	matches, err := stocks.SearchSymbol("KO")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, match := range matches {
        fmt.Printf("%s (%s) - %s, %s\n", match.Symbol, match.Name, match.Region, match.Currency)
	
	}

}