package budget 

import (
	"fmt"
	// "strings"
	// "sort"
	"golang.org/x/text/language"
    "golang.org/x/text/message"
    "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "log"

)


type BudgetRow struct {
	Date   string
	Values map[string]float64
	DailyTotal float64
}

type CatTotal struct {
	Category string
	Total float64
}

type MonthCatTotals struct {
	Month  string
	Totals []CatTotal
}

type YearlyCatTotals struct {
	Category string
	Total float64
}


func ConnectDB() (*sql.DB, error){
	db, err := sql.Open("sqlite3", "./budget.db")
	if err != nil {
    	return nil, err
	}
	return db, nil

}


func FetchDailyTotalsFromDB(db *sql.DB) (map[string]map[string]float64, error) {
	rows, err := db.Query("SELECT category, amount, expense_date FROM csv_data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	budgetMap := make(map[string]map[string]float64)
	for rows.Next() {
		var category string
		var amount float64
	    var date string

		err := rows.Scan(&category, &amount, &date)
		if err != nil {
			log.Fatal(err)
		}
		if _, ok := budgetMap[date]; !ok {
			budgetMap[date] = make(map[string]float64)
		}
		if _, ok := budgetMap[date][category]; !ok {
			budgetMap[date][category] = 0 
		}
		budgetMap[date][category] += amount
	}


	if err := rows.Err(); err != nil {
	    // log.Fatal(err)
	    return nil, err
	}
	return budgetMap, nil
}


func FetchMonthlyTotalsFromDB(dailyMap map[string]map[string]float64) (map[string]MonthCatTotals, error) {
	monthlyTotals := make(map[string]MonthCatTotals)

	for dateStr, catMap := range dailyMap {
		month := dateStr[len(dateStr)-3:]

		monthData, exists := monthlyTotals[month]
		if !exists {
			monthData = MonthCatTotals{
				Month:  month,
				Totals: []CatTotal{},
			}
		}

		for cat, val := range catMap {
			found := false
			for i := range monthData.Totals {
				if monthData.Totals[i].Category == cat {
					monthData.Totals[i].Total += val
					found = true
					break
				}
			}
			if !found {
				monthData.Totals = append(monthData.Totals, CatTotal{Category: cat, Total: val})
			}
		}

		monthlyTotals[month] = monthData
	}

	return monthlyTotals, nil
}



func GetTotalSavings(db *sql.DB) (float64, error) { 
	var total float64
	row := db.QueryRow("SELECT total_savings FROM budget_summary ORDER BY month DESC LIMIT 1;")
	
	err := row.Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}




func PrintMonthlyTotals(monthlyTotals map[string]MonthCatTotals){
	months := []string{"Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"}
	p := message.NewPrinter(language.Korean)

	for _, month := range months {
		m, exists := monthlyTotals[month]
		if !exists {
			continue
		}
	
    	monthTotal := 0.0
    	for _, ct := range m.Totals {
        	monthTotal += ct.Total
    	}

	    fmt.Printf("Month: %s\n", m.Month)
	    p.Printf(" Total Monthly Spending: %d원\n", int64(monthTotal))

	    for _, ct := range m.Totals {
	        pct := (ct.Total / monthTotal) * 100
	        p.Printf(" %s: %d원 (%.2f%%)\n", ct.Category, int64(ct.Total), pct)
	    }
    }
}



func PrintHighestSpendingCategory(monthlyTotals map[string]MonthCatTotals){
	months := []string{"Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"}
	p := message.NewPrinter(language.Korean)
	
	for _, month := range months {
		m, exists := monthlyTotals[month]
		if !exists {
			continue
		}
		maxTotal := 0.0
		maxCat := ""
    	for _, ct := range m.Totals {
        	if ct.Total > maxTotal {
        		maxTotal = ct.Total
        		maxCat = ct.Category
        	}
    	}
		p.Printf("In %s, you spent the most on %s (%d원)\n", month, maxCat, int64(maxTotal))
	}

}





