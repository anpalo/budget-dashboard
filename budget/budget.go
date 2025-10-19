package budget

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"sort"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"budget-dashboard/utils"
)

type BudgetRow struct {
	Date       string
	Values     map[string]float64
	DailyTotal float64
}

type CatTotal struct {
	Category string
	Total    float64
}

type MonthCatTotals struct {
	Month  string
	Totals []CatTotal
}

type YearlyCatTotals struct {
	Category string
	Total    float64
}

var skipCategories = map[string]bool{
	"Daily Total":          true,
	"Income":               true,
	"Monthly Sum":          true,
	"Adjusted Monthly Sum": true,
	"Net Month Income":     true,
	"Yearly Savings Total": true,
	"Daily Sum":            true,
}

func ParseCSV(filename string) ([]BudgetRow, []string, error) {
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
		if row[0] == "" || !utils.IsDate(row[0]) {
			continue
		}

		rowData := BudgetRow{
			Date:   row[0],
			Values: make(map[string]float64),
		}

		for i, h := range headers[1:] {
			rowData.Values[h] = utils.ParseToFloat(row[i+1])
		}

		rowData.DailyTotal = ComputeDailyTotal(rowData.Values)
		budget = append(budget, rowData)
	}

	return budget, headers, nil
}

func PrintBudgetRow(row BudgetRow, headers []string) {
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

func ComputeDailyTotal(values map[string]float64) float64 {
	total := 0.0
	for k, v := range values {
		if !skipCategories[k] {
			total += v
		}
	}
	return total
}

func ComputeMonthlyTotals(budget []BudgetRow) map[string]MonthCatTotals {
	monthlyTotals := make(map[string]MonthCatTotals)

	for _, row := range budget {
		month := strings.Split(row.Date, "-")[1] // to get month
		monthData, exists := monthlyTotals[month]
		if !exists {
			monthData = MonthCatTotals{
				Month:  month,
				Totals: []CatTotal{},
			}
		}

		for cat, val := range row.Values {
			if skipCategories[cat] {
				continue
			}

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

		sort.Slice(monthData.Totals, func(i, j int) bool {
			return monthData.Totals[i].Total > monthData.Totals[j].Total
		})

		monthlyTotals[month] = monthData
	}

	return monthlyTotals
}

func PrintMonthlyTotals(monthlyTotals map[string]MonthCatTotals) {
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

func ComputeMonthlyAverages(budget []BudgetRow, monthlyTotals map[string]MonthCatTotals) map[string]float64 {
	dayCounts := make(map[string]int)
	for _, row := range budget {
		month := strings.Split(row.Date, "-")[1]
		dayCounts[month]++
	}

	averages := make(map[string]float64)
	for month, totals := range monthlyTotals {
		monthTotal := 0.0
		for _, ct := range totals.Totals {
			if skipCategories[ct.Category] {
				continue
			}
			monthTotal += ct.Total
		}
		averages[month] = monthTotal / float64(dayCounts[month])
	}
	return averages
}

func PrintMonthAverage(month string, averages map[string]float64) {
	p := message.NewPrinter(language.Korean)
	p.Printf("Avg. Daily Spending in %s: %d원\n", month, int64(averages[month]))
}

func PrintAllMonthAverages(averages map[string]float64) {
	months := []string{"Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct","Nov","Dec"}
	p := message.NewPrinter(language.Korean)
	fmt.Println("Average Daily Spending by Month:")
	for _, m := range months {
		p.Printf("%s: %d원\n", m, int64(averages[m]))
	}
}

func PrintHighestSpendingCategory(monthlyTotals map[string]MonthCatTotals) {
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

func ComputeYearlyTotals(budget []BudgetRow, headers []string) (map[string]float64, float64) {
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

func PrintYearlyTotals(categoryTotals map[string]float64, yearlyTotal float64) {
	var totals []YearlyCatTotals
	for cat, total := range categoryTotals {
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

func GetTotalSavings(budget []BudgetRow, headers []string) float64 {
	if len(budget) == 0 || len(headers) == 0 {
		return 0
	}
	firstRow := budget[0]
	lastHeader := headers[len(headers)-1] // last column = yearly savings total
	if val, ok := firstRow.Values[lastHeader]; ok {
		return val
	}
	return 0
}

// func PrintTotalSavingsKRW(savings float64) {
// 	p := message.NewPrinter(language.Korean)
// 	p.Printf("Total Savings: %d원\n", int64(savings))
// }
