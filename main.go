package currencyGoConverter

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

const (
	baseUrl = "https://api.exchangeratesapi.io/latest?base="
)

type ConvertHandler struct {
}

func main() {
	http.Handle("/convert-all", &ConvertHandler{})
	log.Fatal(http.ListenAndServe(":8080", nil))

}

/// Server with endpoint /convert-all/?base=PLN&amount=2.34 that returns:
// {"base": "PLN", "amount": 2.34, "results": {"USD": 0.56, ...}}
// converting: call to external API to get actual rates returns {"currency": amount} struct to channel.
// endpoint response is read from channel
type JsonCurrencyAmount struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type ResponseStruct struct {
	Date         string             `json:"date"`
	Rates        map[string]float64 `json:"rates"`
	BaseCurrency string             `json:"base"`
}

func convertCurrency(currency string, amount float64, rate float64) JsonCurrencyAmount {
	return JsonCurrencyAmount{currency, amount * rate}
}

func ProcessResponse(resp *http.Response) ResponseStruct {
	var dr ResponseStruct
	if err := json.NewDecoder(resp.Body).Decode(&dr); err != nil {
		log.Fatal(err)
	}
	return dr
}

func GetExchangeRate(baseCurrency string, resultCurrency string) float64 {
	url := baseUrl + baseCurrency
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	dr := ProcessResponse(response)
	result := dr.Rates[resultCurrency]
	return result
}

func (h *ConvertHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var data bytes.Buffer
	_, err := data.ReadFrom(r.Body)
	if err != nil {
		panic(err)
	}
	queryParams := r.URL.Query()
	baseCurrency := queryParams.Get("base")
	amount, _ := strconv.ParseFloat(queryParams.Get("amount"), 64)
	resultCurrency := queryParams.Get("result")
	exchangeRate := GetExchangeRate(baseCurrency, resultCurrency)
	result := convertCurrency(resultCurrency, amount, exchangeRate)
	json.NewEncoder(w).Encode(result)

}
