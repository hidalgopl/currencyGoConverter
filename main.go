package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

const (
	baseUrl = "https://api.exchangeratesapi.io/latest?base="
)

type ConvertHandler struct {
}

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

func GetExchangeRate(baseCurrency string, resultCurrency string, amount float64, resultChannel chan JsonCurrencyAmount) float64 {
	exchangeUrl := baseUrl + baseCurrency
	response, err := http.Get(exchangeUrl)
	if err != nil {
		log.Fatal(err)
	}
	dr := ProcessResponse(response)
	result := dr.Rates[resultCurrency]
	defer response.Body.Close()
	resultChannel <- convertCurrency(resultCurrency, amount, result)
	return 0
}

func (h *ConvertHandler) HandleQueryParams(qParams url.Values) (string, float64, string) {
	baseCurrency := qParams.Get("base")
	if baseCurrency == "" {
		log.Fatal("No base param")
	}
	amount, err := strconv.ParseFloat(qParams.Get("amount"), 64)
	if err != nil {
		log.Fatal(err)
	}
	resultCurrency := qParams.Get("result")
	if resultCurrency == "" {
		log.Fatal("No result param")
	}
	return baseCurrency, amount, resultCurrency
}

func (h *ConvertHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data bytes.Buffer
	_, err := data.ReadFrom(r.Body)
	if err != nil {
		panic(err)
	}
	queryParams := r.URL.Query()
	baseCurrency, amount, resultCurrency := h.HandleQueryParams(queryParams)
	resultChannel := make(chan JsonCurrencyAmount)
	go GetExchangeRate(baseCurrency, resultCurrency, amount, resultChannel)
	json.NewEncoder(w).Encode(<-resultChannel)

}

func main() {
	/// Server with endpoint /convert-all/?base=PLN&amount=2.34 that returns:
	// {"base": "PLN", "amount": 2.34, "results": {"USD": 0.56, ...}}
	// converting: call to external API to get actual rates returns {"currency": amount} struct to channel.
	// endpoint response is read from channel
	http.Handle("/convert-all", &ConvertHandler{})
	log.Fatal(http.ListenAndServe(":8080", nil))

}
