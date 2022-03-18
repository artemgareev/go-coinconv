package protocol

import (
	"time"
)

// PriceConversionQueryParameters specifies the query parameters to PriceConversion
type PriceConversionQueryParameters struct {
	Amount    float64    `url:"amount"`
	ID        *int64     `url:"id,omitempty"`
	Symbol    *string    `url:"symbol,omitempty"`
	Time      *time.Time `url:"time,omitempty"`
	Convert   *string    `url:"convert,omitempty"`
	ConvertID *int64     `url:"convert_id,omitempty"`
}

// PriceConversionResponse describes PriceConversion API response
type PriceConversionResponse struct {
	Data   []Listing `json:"data"`
	Status Status    `json:"status"`
}

type Listing struct {
	ID          int              `json:"id"`
	Symbol      string           `json:"symbol"`
	Name        string           `json:"name"`
	Amount      float64          `json:"amount"`
	LastUpdated string           `json:"last_updated"`
	Quote       map[string]Quote `json:"quote"`
}

type Quote struct {
	Price       float64   `json:"price"`
	LastUpdated time.Time `json:"last_updated"`
}

type Status struct {
	Timestamp    string  `json:"timestamp"`
	ErrorCode    int     `json:"error_code"`
	ErrorMessage *string `json:"error_message"`
	Elapsed      int     `json:"elapsed"`
	CreditCount  int     `json:"credit_count"`
}
