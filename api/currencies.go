package api

import (
"encoding/json"
"net/http"
"budget-dashboard/currencies"
)


func CurrencyConversionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rates, err := currencies.GetRates()
		if err != nil {
			http.Error(w, "Failed to fetch exchange rates", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rates)
	}
}