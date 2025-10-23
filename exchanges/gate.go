package exchanges

import (
	"all_exchange_symbol/models"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type Gate struct {
	Name string
}

type GateSymbol struct {
	Id              string `json:"id"`
	Base            string `json:"base"`
	Quote           string `json:"quote"`
	Fee             string `json:"fee"`
	MinBaseAmount   string `json:"min_base_amount"`
	MinQuoteAmount  string `json:"min_quote_amount"`
	AmountPrecision int    `json:"amount_precision"`
	Precision       int    `json:"precision"`
	TradeStatus     string `json:"trade_status"`
	SellStart       int64  `json:"sell_start"`
	BuyStart        int64  `json:"buy_start"`
}

type GateFuturesContract struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Quanto      bool   `json:"quanto"`
	Leverage    string `json:"leverage"`
	InDelisting bool   `json:"in_delisting"`
	TradeStatus string `json:"trade_status"`
}

func NewGate() *Gate {
	return &Gate{Name: "gate"}
}

func (g *Gate) GetName() string {
	return g.Name
}

func (g *Gate) FetchSpotSymbols() ([]models.Symbol, error) {
	resp, err := http.Get("https://api.gateio.ws/api/v4/spot/currency_pairs")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []GateSymbol

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var symbols []models.Symbol
	for _, s := range result {
		symbols = append(symbols, models.Symbol{
			Exchange:  g.Name,
			Type:      "spot",
			Symbol:    s.Id,
			CreatedAt: time.Now(),
		})
	}

	return symbols, nil
}

func (g *Gate) FetchFuturesSymbols() ([]models.Symbol, error) {
	resp, err := http.Get("https://api.gateio.ws/api/v4/futures/usdt/contracts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result []GateFuturesContract

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var symbols []models.Symbol
	for _, s := range result {
		symbols = append(symbols, models.Symbol{
			Exchange:  g.Name,
			Type:      "futures",
			Symbol:    s.Name,
			CreatedAt: time.Now(),
		})
	}

	return symbols, nil
}
