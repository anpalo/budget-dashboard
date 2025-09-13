package api

import (
"encoding/json"
"net/http"
"budget-dashboard/stocks"
)


func StocksHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := r.URL.Query().Get("symbol")
        if symbol == "" {
            http.Error(w, "Please provide a symbol using the 'symbol' query parameter, e.g., ?symbol=AAPL", http.StatusBadRequest)
            return
        }
        
		timeSeries, err := stocks.GetStocks(symbol)
		if err != nil {
			http.Error(w, "Failed to fetch stock data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(timeSeries)
	}
}

func SymbolSearchHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		search := r.URL.Query().Get("q")
		if search == "" {
			http.Error(w, "Please provide a search term using the 'q' query parameter, e.g., ?q=Apple", http.StatusBadRequest)
			return
		}

		matches, err := stocks.SearchSymbol(search)
		if err != nil {
			http.Error(w, "Failed to search for stock", http.StatusInternalServerError)
			return
		}

		if len(matches) == 0{
			http.Error(w, "Sorry, we couldn’t find any stock matching your search.", http.StatusNotFound)
			return 
		}



		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(matches)
	}
}