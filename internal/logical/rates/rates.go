package rates

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// getFirstDayOfLastMonth calculates the first day of the previous month
func getFirstDayOfLastMonth() time.Time {
	now := time.Now()
	year, month, _ := now.Date()
	location := time.UTC

	if month == time.January {
		month = time.December
		year--
	} else {
		month--
	}

	return time.Date(year, month, 1, 0, 0, 0, 0, location)
}

// getFirstDayOfLastMonth calculates the first day of the current month
func getFirstDayOfCurrentMonth() time.Time {
	now := time.Now()
	year, month, _ := now.Date()
	location := time.UTC

	return time.Date(year, month, 1, 0, 0, 0, 0, location)
}

// getExchangeRate gets the exchange rate for the specified day
func getExchangeRate(currency, dateStr string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?json&valcode=%s&date=%s", currency, dateStr)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error in obtaining exchange rates: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("incorrect response status for %s: %d", currency, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error when reading the answer: %v", err)
	}

	var rate []map[string]interface{}
	if err := json.Unmarshal(body, &rate); err != nil {
		return nil, fmt.Errorf("error when decoding data for %s: %v", currency, err)
	}

	if len(rate) > 0 {
		return rate[0], nil
	}

	return nil, fmt.Errorf("Data for %s not found", currency)
}

// GetLastMonthRates exchange rate for the last month or the current month and saves the data to a JSON file
func getLastMonthRates(month string) error {
	firstDayLastMonth := getFirstDayOfLastMonth()
	firstDayCurrentMonth := getFirstDayOfCurrentMonth()
	dateStrLast := firstDayLastMonth.Format("20060102")
	dateStrCurr := firstDayCurrentMonth.Format("20060102")

	currencies := []string{"USD", "EUR", "GBP"}
	var results []map[string]interface{}

	for _, currency := range currencies {
		var rate map[string]interface{}
		var err error

		if month == "current" {
			rate, err = getExchangeRate(currency, dateStrCurr)
		} else if month == "last" {
			rate, err = getExchangeRate(currency, dateStrLast)
		}

		if err != nil {
			fmt.Println(err)
			continue
		}

		results = append(results, rate)
	}

	var file *os.File
	var err error

	if month == "current" {
		file, err = os.Create("../storage/currency_rates_current_month.json")

		fmt.Println("Currency rates for the current month have been successfully saved to a file currency_rates_current_month.json")

	} else if month == "last" {
		file, err = os.Create("../storage/currency_rates_last_month.json")

		fmt.Println("currency rates for the previous month have been successfully saved to a file currency_rates_last_month.json")
	}

	if err != nil {
		return fmt.Errorf("error when creating a file: %v", err)
	}
	defer file.Close()

	encodedData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("error when decoding data to a file: %v", err)
	}

	if _, err := file.Write(encodedData); err != nil {
		return fmt.Errorf("an error when writing data to a file: %v", err)
	}

	return nil
}

// LoadRatesFromFile Loads data from a JSON file and returns a currency exchange rate mappa
func LoadRatesFromFile(month string) (map[string]float64, error) {
	getLastMonthRates(month)

	var file *os.File
	var err error
	if month == "current" {
		file, err = os.Open("../storage/currency_rates_current_month.json")
	} else if month == "last" {
		file, err = os.Open("../storage/currency_rates_last_month.json")
	}

	if err != nil {
		return nil, fmt.Errorf("error when opening a file: %v", err)
	}
	defer file.Close()

	var results []map[string]interface{}
	if err := json.NewDecoder(file).Decode(&results); err != nil {
		return nil, fmt.Errorf("error when decoding data from a file: %v", err)
	}

	// Creating a mappa for storing currency rates
	ratesMap := make(map[string]float64)

	// Go through the results and fill in the mappa
	for _, result := range results {
		cc, ok := result["cc"].(string)
		if !ok {
			return nil, fmt.Errorf("error in the data format: unable to retrieve the currency code")
		}

		rate, ok := result["rate"].(float64)
		if !ok {
			return nil, fmt.Errorf("error in the data format: failed to retrieve the course")
		}

		ratesMap[cc] = rate
	}

	return ratesMap, nil
}

func LoadCurrentRatesFromFile() (map[string]float64, error) {

	file, err := os.Open("../storage/filtered_currency_rates.json")
	if err != nil {
		return nil, fmt.Errorf("error when opening a file: %v", err)
	}
	defer file.Close()

	var results []map[string]interface{}
	if err := json.NewDecoder(file).Decode(&results); err != nil {
		return nil, fmt.Errorf("error when decoding data from a file: %v", err)
	}

	// create a mappa for storing currency rates
	ratesMap := make(map[string]float64)

	for _, result := range results {
		cc, ok := result["CurrencyCodeL"].(string)
		if !ok {
			return nil, fmt.Errorf("error in the data format: unable to retrieve the currency code")
		}

		rate, ok := result["Amount"].(float64)
		if !ok {
			return nil, fmt.Errorf("error in the data format: failed to retrieve the course")
		}

		ratesMap[cc] = rate
	}

	return ratesMap, nil
}
