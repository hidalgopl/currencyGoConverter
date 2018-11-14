package main

import (
	"net/http"
	"testing"
)

func TestConvertCurrency(t *testing.T) {
	amount := 1.5
	currency := "USD"
	rate := 2.0
	result := convertCurrency(currency, amount, rate)
	if result.Currency != currency || result.Amount != rate*amount {
		t.Fatal("Bleh")
	}
}

func TestProcessResponse(t *testing.T) {
	response, _ := http.Get("https://api.exchangeratesapi.io/latest?base=USD")
	result := ProcessResponse(response)
	if result.BaseCurrency != "USD" {
		t.Fatal("Something is in the air")
	}
}
