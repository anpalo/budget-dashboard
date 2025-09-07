package stocks

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AlphaVantageResponse struct {
	MetaData   map[string]string              `json:"Meta Data"`
	TimeSeries map[string]map[string]string `json:"Time Series (5min)"`
}

func GetStocks(symbol string) (map[string]map[string]string, error) {
	baseURL := "https://www.alphavantage.co/query"
	u, _ := url.Parse(baseURL)

	q := u.Query()
	q.Set("apikey", ALPHA_VANTAGE_API_KEY)
	q.Set("function", "TIME_SERIES_INTRADAY")
	q.Set("interval", "5min")
	q.Set("symbol", symbol)
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

	var data AlphaVantageResponse
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data.TimeSeries, nil
}
