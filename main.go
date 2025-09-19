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
	// "budget-dashboard/stocks"
	"budget-dashboard/api"
	"net/http"
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

	totalSavings := budget.GetTotalSavings(rows, headers)
	fmt.Printf("Total Savings (KRW): %d원\n", int64(totalSavings))



	rates, err := currencies.GetRates()
	if err != nil {
	    fmt.Println("Error fetching rates:", err)
	    return
	}
	fmt.Println("Rates:", rates)

	// searchString := "LULU"
	// matches, err := stocks.SearchSymbol(searchString)
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }

	// if len(matches) == 0 {
    // 	fmt.Println("No matching symbols found")
    // 	return
	// }

	// chosenSymbol := matches[0].Symbol
	// chosenName := matches[0].Name
	// timeSeries, err := stocks.GetStocks(chosenSymbol)
	// if err != nil {
    // 	fmt.Println("Error fetching stock data:", err)
    // 	return
	// }

	// latestPrice := ""
    // for _, values := range timeSeries {
    //     latestPrice = values["4. close"]
    //     break 
    // }

    // fmt.Printf("%s (%s) | Latest Price: %s %s\n",
    // chosenName, chosenSymbol, latestPrice, matches[0].Currency)


    http.HandleFunc("/api/monthly-totals", api.MonthlyTotalsHandler(monthlyTotals))
	http.HandleFunc("/api/total-savings", api.TotalSavingsHandler(rows, headers))
    http.HandleFunc("/api/currency", api.CurrencyConversionHandler())
	// http.HandleFunc("/api/stocks", api.StocksHandler())
	// http.HandleFunc("/api/symbol-search", api.SymbolSearchHandler())


	fs := http.FileServer(http.Dir("./frontend"))
	http.Handle("/", fs)

	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}



