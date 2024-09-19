package logical

import (
	"encoding/json"
	"errors"
	"exchanger/internal/entity"
	"fmt"
	"io/ioutil"
	"os"
)

// Function for parsing JSON-file with currency rates
func loadRates() ([]entity.CurrencyRate, error) {
	
	file, err := os.Open("../storage/filtered_currency_rates.json")
	if err != nil {
		return nil, errors.New("unable to open a file with exchange rates")
	}
	defer file.Close()

	
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.New("error reading a file with exchange rates")
	}

	
	var rates []entity.CurrencyRate
	err = json.Unmarshal(byteValue, &rates)
	if err != nil {
		return nil, errors.New("error parsing a file with exchange rates")
	}

	return rates, nil
}

// Function for converting currencies
func ConvertCurrency(fromCurrency, toCurrency string, amount float64) (float64, error) {
	// Download currency rates
	rates, err := loadRates()
	if err != nil {
		return 0, err
	}

	// Find the hryvnia (UAH) exchange rate to all currencies
	var uahToCurrencyRates = make(map[string]float64)
	for _, rate := range rates {
		if rate.CurrencyCodeL != "UAH" {
			uahToCurrencyRates[rate.CurrencyCodeL] = 1 / rate.Amount
		}
	}

	// If the hryvnia is not specified in the exchange rates, set it as the base currency
	uahToCurrencyRates["UAH"] = 1.0

	// Currency exchange rate to hryvnia
	fromRate, fromExists := uahToCurrencyRates[fromCurrency]
	toRate, toExists := uahToCurrencyRates[toCurrency]

	if !fromExists {
		return 0, fmt.Errorf("currency %s not found", fromCurrency)
	}
	if !toExists {
		return 0, fmt.Errorf("currency %s not found", toCurrency)
	}

	// Conversion from source currency to hryvnia and then to the target currency
	convertedAmount := (amount / fromRate) * toRate

	return convertedAmount, nil
}
