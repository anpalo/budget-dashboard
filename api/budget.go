package api

import (
    "encoding/json"
    "net/http"
    "budget-dashboard/budget"
)

func MonthlyTotalsHandler(monthlyTotals map[string]budget.MonthCatTotals) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(monthlyTotals); err != nil {
            http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
        }
    }
}

func TotalSavingsHandler(rows []budget.BudgetRow, headers []string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        // Get total savings from the first row
        totalSavings := budget.GetTotalSavings(rows, headers)

        // Wrap in JSON object
        resp := map[string]float64{
            "krw": totalSavings,
        }

        if err := json.NewEncoder(w).Encode(resp); err != nil {
            http.Error(w, "Failed to encode JSON: "+err.Error(), http.StatusInternalServerError)
            return
        }
    }
}