package stocks

import (
	"encoding/json"
	"net/http"
	"net/url"
	"io/ioutil"
)


type SymbolMatch struct {
	Symbol string `json:"1. symbol"`
	Name string `json:"2. name"`
	Type string `json:"3. type"`
	Region string `json:"4. region"`
	MarketOpen string `json:"5. marketOpen"`
	MarketClose string `json:"6. marketClose"`
	Timezone string `json:"7. timezone"`
	Currency string `json:"8. currency"`
	MatchScore string `json:"9. matchScore"`
}

type SymbolSearchResponse struct {
	BestMatches []SymbolMatch `json:"bestMatches"`
}

func SearchSymbol(searchString string) ([]SymbolMatch, error) {

	baseURL := "https://www.alphavantage.co/query"
	u, _ := url.Parse(baseURL)
	q := u.Query()
	q.Set("apikey", ALPHA_VANTAGE_API_KEY)
	q.Set("function","SYMBOL_SEARCH")
	q.Set("keywords", searchString)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
    	return nil, err
	}
	
	var data SymbolSearchResponse

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err

	}

	return data.BestMatches, nil

}