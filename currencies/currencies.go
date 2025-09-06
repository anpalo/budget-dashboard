package currencies

import (
  "encoding/json"	
  "fmt"
  "net/http"
  "net/url"
  "io/ioutil"
)

type FXResponse struct {
	Success bool `json:"success"`
	Base string `json:"base"`
	Rates map[string]float64 `json:"rates"`
}

func GetRates() (map[string]float64, error) {

	baseURL := "https://api.fxratesapi.com/latest"

	u, _ := url.Parse(baseURL)
	q := u.Query()
	q.Set("api_key", FX_API_KEY)
	q.Set("base", "KRW")
	q.Set("symbols", "USD,EUR,GBP")
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var data FXResponse

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err

	}

	return data.Rates, nil

  // fmt.Println("KRW → Currency: ")
  // for key, value := range data.Rates{
  // 	fmt.Printf("1 KRW = %.6f %s\n", value, key)
  // }

}

func main() {

	rates, err := GetRates()
	if err != nil {
	    fmt.Println("Error fetching rates:", err)
	    return
	}
	fmt.Println("Rates:", rates)

}


