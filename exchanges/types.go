package exchanges

import "all_exchange_symbol/models"

type ExchangeInterface interface {
	GetName() string
	FetchSpotSymbols() ([]models.Symbol, error)
	FetchFuturesSymbols() ([]models.Symbol, error)
}

type BaseSymbol struct {
	Symbol     string `json:"symbol"`
	Status     string `json:"status"`
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
}
