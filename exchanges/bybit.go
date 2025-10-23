package exchanges

import (
	"all_exchange_symbol/models"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Bybit struct {
	Name string
}

type BybitSymbol struct {
	Symbol    string `json:"symbol"`
	BaseCoin  string `json:"baseCoin"`
	QuoteCoin string `json:"quoteCoin"`
	Status    string `json:"status"`
}

type BybitFuturesSymbol struct {
	Symbol       string `json:"symbol"`
	ContractType string `json:"contractType"`
	Status       string `json:"status"`
	BaseCoin     string `json:"baseCoin"`
	QuoteCoin    string `json:"quoteCoin"`
	LaunchTime   string `json:"launchTime"`
	DeliveryTime string `json:"deliveryTime"`
}

func NewBybit() *Bybit {
	return &Bybit{Name: "bybit"}
}

func (b *Bybit) GetName() string {
	return b.Name
}

func (b *Bybit) FetchSpotSymbols() ([]models.Symbol, error) {
	resp, err := http.Get("https://api.bybit.com/v5/market/instruments-info?category=spot")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			Category string        `json:"category"`
			List     []BybitSymbol `json:"list"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var symbols []models.Symbol
	for _, s := range result.Result.List {
		symbols = append(symbols, models.Symbol{
			Exchange:  b.Name,
			Type:      "spot",
			Symbol:    s.Symbol,
			CreatedAt: time.Now(),
		})
	}

	return symbols, nil
}

func (b *Bybit) FetchFuturesSymbols() ([]models.Symbol, error) {
	resp, err := http.Get("https://api.bybit.com/v5/market/instruments-info?category=linear")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		RetCode int    `json:"retCode"`
		RetMsg  string `json:"retMsg"`
		Result  struct {
			Category string               `json:"category"`
			List     []BybitFuturesSymbol `json:"list"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var symbols []models.Symbol
	for _, s := range result.Result.List {
		symbols = append(symbols, models.Symbol{
			Exchange:  b.Name,
			Type:      "futures",
			Symbol:    s.Symbol,
			CreatedAt: time.Now(),
		})
	}

	return symbols, nil
}
