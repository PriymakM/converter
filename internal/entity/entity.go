package entity

type CurrencyRate struct {
	StartDate    string  `json:"StartDate"`
	TimeSign     string  `json:"TimeSign"`
	CurrencyCode string  `json:"CurrencyCode"`
	CurrencyCodeL string `json:"CurrencyCodeL"`
	Units        int     `json:"Units"`
	Amount       float64 `json:"Amount"`
}

type RateResponse struct {
	CurrencyCode   string  `json:"cc"`
	ExchangeDate   string  `json:"exchangedate"`
	Rate           float64 `json:"rate"`
	CurrencyCodeID int     `json:"r030"`
	Description    string  `json:"txt"`
}
