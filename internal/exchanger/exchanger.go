package exchanger

import (
	"encoding/json"
	"exchanger/internal/entity"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func GetCurrencyRates() {
	// URL for writing to the NBU API
	apiURL := "https://bank.gov.ua/NBU_Exchange/exchange?json"

	resp, err := http.Get(apiURL)
	if err != nil {
		fmt.Println("Помилка запиту:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Помилка читання відповіді:", err)
		return
	}

	var allRates []entity.CurrencyRate
	err = json.Unmarshal(body, &allRates)
	if err != nil {
		fmt.Println("Помилка парсингу JSON:", err)
		return
	}

	// Filter the currencies you need
	var filteredRates []entity.CurrencyRate
	for _, rate := range allRates {
		if rate.CurrencyCodeL == "UAH" || rate.CurrencyCodeL == "USD" || rate.CurrencyCodeL == "EUR" || rate.CurrencyCodeL == "GBP" {
			filteredRates = append(filteredRates, rate)
		}
	}

	// Verify filtered data
	for _, rate := range filteredRates {
		fmt.Printf("Валюта: %s, Курс: %.6f, Одиниці: %d\n", rate.CurrencyCodeL, rate.Amount, rate.Units)
	}

	file, err := os.Create("../storage/filtered_currency_rates.json")
	if err != nil {
		fmt.Println("Помилка створення файлу:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") 
	err = encoder.Encode(filteredRates)
	if err != nil {
		fmt.Println("Помилка збереження JSON:", err)
		return
	}

	fmt.Println("Дані успішно збережено в filtered_currency_rates.json")
}
