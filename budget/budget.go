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
	/* Next() and Scan() from database/sql package */
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

// func ComputeMonthlyAverages(budget []BudgetRow, monthlyTotals map[string]MonthCatTotals) map[string]float64 {
// 	dayCounts := make(map[string]int)
// 	for _, row := range budget {
// 		month := strings.Split(row.Date, "-")[1] // to get month
// 		dayCounts[month] ++
// 	}
// 	averages := make(map[string]float64)
// 	for month, totals := range monthlyTotals {
// 		monthTotal := 0.0
// 		for _, ct := range totals.Totals {
// 			if skipCategories[ct.Category] {
// 				continue
// 			}
// 			monthTotal += ct.Total
// 		}
// 		averages[month] = monthTotal / float64(dayCounts[month])
// 	}
// 	return averages
// }

// func PrintMonthAverage(month string, averages map[string]float64) {
//     p := message.NewPrinter(language.Korean)
//     p.Printf("Avg. Daily Spending in %s: %d원\n", month, int64(averages[month]))
// }

// func PrintAllMonthAverages(averages map[string]float64) {
//     months := []string{"Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"}
//     p := message.NewPrinter(language.Korean)

//     fmt.Println("Average Daily Spending by Month:")
//     for _, m := range months {
//         p.Printf("%s: %d원\n", m, int64(averages[m]))
//     }
// }

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

// func ComputeYearlyTotals(budget []BudgetRow, headers []string) (map[string]float64, float64){ 
// 	yearlyTotal := 0.0
// 	totals := make(map[string]float64)
// 	for _, h := range headers[1:] {
// 		if skipCategories[h] {
// 			continue
// 		}
// 		sum := 0.0
// 		for _, row := range budget {
// 			if val, ok := row.Values[h]; ok {
// 				sum += val
// 			}
// 		}
// 		totals[h] = sum
// 		yearlyTotal += sum
// 	}
// 	return totals, yearlyTotal
// }

// func PrintYearlyTotals(categoryTotals map[string]float64, yearlyTotal float64){
// 	var totals []YearlyCatTotals
// 	for cat, total := range categoryTotals{
// 		totals = append(totals, YearlyCatTotals{cat, total})
// 	}
// 	sort.Slice(totals, func(i, j int) bool {
//     return totals[i].Total > totals[j].Total
// 	})

// 	p := message.NewPrinter(language.Korean)

// 	p.Printf("Yearly Total Spending: %d원\n", int64(yearlyTotal))
// 	for _, t := range totals {
// 		pct := (t.Total / yearlyTotal) * 100
//         p.Printf("%s: %d원 (%.2f%%)\n", t.Category, int64(t.Total), pct)
//     }
// }


func GetTotalSavings(db *sql.DB) {
	var total float64
	query := `
		SELECT category, SUM(amount)
		FROM csv_data`
}
// func GetTotalSavings(budget []BudgetRow, headers []string) float64 {
//     if len(budget) == 0 || len(headers) == 0 {
//         return 0
//     }

//     firstRow := budget[0]
//     lastHeader := headers[len(headers)-1] // last column = yearly savings total

//     if val, ok := firstRow.Values[lastHeader]; ok {
//         return val
//     }
//     return 0
// }



