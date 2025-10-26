
package main

import (
	"log"
	"fmt"
	"budget-dashboard/budget"
	"budget-dashboard/currencies"
	// "budget-dashboard/stocks"
	// "budget-dashboard/api"
	"net/http"
)

func main() {

	// http.HandleFunc("/api/upload-csv", api.UploadCSVHandler)

	db, err := budget.ConnectDB()
	if err != nil {
    	fmt.Println("DB connection error:", err)
    }
	defer db.Close() 
	fmt.Println("DB connected!")

	dailyMap, err := budget.FetchDailyTotalsFromDB(db)
	if err != nil {
		log.Fatal(err)
	}
	monthlyTotals, err:= budget.FetchMonthlyTotalsFromDB(dailyMap)
	if err != nil {
		log.Fatal(err)
	}


	budget.PrintMonthlyTotals(monthlyTotals)
	budget.PrintHighestSpendingCategory(monthlyTotals)


	rates, err := currencies.GetRates()
	if err != nil {
	    fmt.Println("Error fetching rates:", err)
	    return
	}
	fmt.Println("Rates:", rates)

	
    // http.HandleFunc("/api/monthly-totals", api.MonthlyTotalsHandler(monthlyTotals))
    // http.HandleFunc("/api/currency", api.CurrencyConversionHandler())

	// // http.HandleFunc("/api/total-savings", api.TotalSavingsHandler(rows, headers))
	// // http.HandleFunc("/api/stocks", api.StocksHandler())
	// http.HandleFunc("/api/symbol-search", api.SymbolSearchHandler())


fs := http.FileServer(http.Dir("./frontend"))
http.Handle("/", fs)


	fmt.Println("Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}



