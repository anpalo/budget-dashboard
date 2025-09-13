package api

import (
"encoding/json"
"net/http"
"budget-dashboard/budget"
)

func MonthlyTotalsHandlers(monthlyTotals map[string]budget.MonthCatTotals) http.Handler{ 
    return func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(monthlyTotals)
    }
}