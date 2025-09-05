/*
TODO:
- Average daily spending per month.
- Highest and lowest spending categories per month.
*/

package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sort"
	"golang.org/x/text/language"
    "golang.org/x/text/message"
)

type BudgetRow struct {
	Date   string
	Values map[string]float64
	DailyTotal float64
}

type MonthCatTotals struct {
	Month  string
	Totals map[string]float64
}

type YearlyCatTotals struct {
	Category string
	Total float64
}

var skipCategories = map[string]bool{
	"Daily Total":            true,
	"Income":               true,
	"Monthly Sum":          true,
	"Adjusted Monthly Sum": true,
	"Net Month Income":     true,
	"Yearly Savings Total": true,
	"Daily Sum": true,
}

func parseCSV(filename string) ([]BudgetRow, []string, error){
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
	}
	headers := records[0]
	var budget []BudgetRow

	for _, row := range records[1:] {
		if row[0] == "" || !isDate(row[0]) {
			continue
		}

		rowData := BudgetRow{
			Date:   row[0],
			Values: make(map[string]float64),
		}

		for i, h := range headers[1:] {
			rowData.Values[h] = parseToFloat(row[i+1])
		}
		rowData.DailyTotal = computeDailyTotal(rowData.Values)
		budget = append(budget, rowData)
	}
	return budget, headers, nil
}

func isDate(s string) bool {
	return len(s) > 0 && s[0] >= '0' && s[0] <= '9'
}

func parseToFloat(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	if s == "" {
		return 0
	}
	num, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return num
}

func computeDailyTotal(values map[string]float64) float64 {
	total := 0.0
	for k, v := range values {
		if !skipCategories[k] {
			total += v
		}
	}
	return total
}

func printBudgetRow(row BudgetRow, headers []string) {
	fmt.Printf("Date: %s Values: {", row.Date)
	first := true
	for _, h := range headers[1:] {
		if skipCategories[h] {
			continue
		}
		if !first {
			fmt.Print(", ")
		}
		fmt.Printf("%s: %.0f", h, row.Values[h])
		first = false
	}
	fmt.Println("}")
}

func computeMonthlyTotals(budget []BudgetRow) map[string]MonthCatTotals{
	monthlyTotals := make(map[string]MonthCatTotals)
	for _, row := range budget {
		month := strings.Split(row.Date, "-")[1] // to get month
		m, exists := monthlyTotals[month]
		if !exists {
			m = MonthCatTotals{
				Month: month, 
				Totals: make(map[string]float64),
			}
		}

		for cat, val := range row.Values {
			if skipCategories[cat] {
				continue
			}
			m.Totals[cat] += val
		}
		monthlyTotals[month] = m
	}
	return monthlyTotals
}

func printMonthlyTotals(monthlyTotals map[string]MonthCatTotals){
	for _, m := range monthlyTotals {
    	monthTotal := 0.0
    	for cat, val := range m.Totals {
        	if skipCategories[cat] {
            	continue
        	}
        	monthTotal += val
    	}
    	p := message.NewPrinter(language.Korean)
	    fmt.Printf("Month: %s\n", m.Month)
	    p.Printf(" Total Monthly Spending: %d원\n", int64(monthTotal))

	    for cat, total := range m.Totals {
	        if skipCategories[cat] {
	            continue
	        }
	        pct := (total / monthTotal) * 100
	        p.Printf(" %s: %d원 (%.2f%%)\n", cat, int64(total), pct)
	    }
    }
}

func computeMonthlyAverages(budget []BudgetRow, monthlyTotals map[string]MonthCatTotals) map[string]float64 {
	dayCounts := make(map[string]int)
	for _, row := range budget {
		month := strings.Split(row.Date, "-")[1] // to get month
		dayCounts[month] ++
	}
	averages := make(map[string]float64)
	for month, totals := range monthlyTotals {
		monthTotal := 0.0
		for cat, val := range totals.Totals {
			if skipCategories[cat] {
				continue
			}
			monthTotal += val
		}
		averages[month] = monthTotal / float64(dayCounts[month])
	}
	return averages
}

func printMonthAverage(month string, averages map[string]float64) {
    p := message.NewPrinter(language.Korean)
    p.Printf("Avg. Daily Spending in %s: %d원\n", month, int64(averages[month]))
}

func printAllMonthAverages(averages map[string]float64) {
    months := []string{"Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"}
    p := message.NewPrinter(language.Korean)

    fmt.Println("Average Daily gSpending by Month:")
    for _, m := range months {
        p.Printf("%s: %d원\n", m, int64(averages[m]))
    }
}

func computeYearlyTotals(budget []BudgetRow, headers []string) (map[string]float64, float64){ 
	yearlyTotal := 0.0
	totals := make(map[string]float64)
	for _, h := range headers[1:] {
		if skipCategories[h] {
			continue
		}
		sum := 0.0
		for _, row := range budget {
			if val, ok := row.Values[h]; ok {
				sum += val
			}
		}
		totals[h] = sum
		yearlyTotal += sum
	}
	return totals, yearlyTotal
}

func printYearlyTotals(categoryTotals map[string]float64, yearlyTotal float64){
	var totals []YearlyCatTotals
	for cat, total := range categoryTotals{
		totals = append(totals, YearlyCatTotals{cat, total})
	}
	sort.Slice(totals, func(i, j int) bool {
    return totals[i].Total > totals[j].Total
	})

	p := message.NewPrinter(language.Korean)

	p.Printf("Yearly Total Spending: %d원\n", int64(yearlyTotal))
	for _, t := range totals {
		pct := (t.Total / yearlyTotal) * 100
        p.Printf("%s: %d원 (%.2f%%)\n", t.Category, int64(t.Total), pct)
    }
}

func main() {
	budget, headers, err := parseCSV("./ExampleBudget.csv")
	if err != nil {
		fmt.Println("Failed to parse CSV:", err)
		return
	}

	// for _, row := range budget {
	// 	printBudgetRow(row, headers)
	// }

	monthlyTotals := computeMonthlyTotals(budget)
	printMonthlyTotals(monthlyTotals)
	yearlyCategoryTotals, yearlyTotal := computeYearlyTotals(budget, headers)
	printYearlyTotals(yearlyCategoryTotals, yearlyTotal)
	printAllMonthAverages(computeMonthlyAverages(budget, monthlyTotals))

}