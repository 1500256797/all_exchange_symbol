package exchanges

import (
	"all_exchange_symbol/models"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type OKX struct {
	Name string
}

type OKXInstrument struct {
	InstType string `json:"instType"`
	InstId   string `json:"instId"`
	State    string `json:"state"`
}

func NewOKX() *OKX {
	return &OKX{Name: "okx"}
}

func (o *OKX) GetName() string {
	return o.Name
}

func (o *OKX) FetchSpotSymbols() ([]models.Symbol, error) {
	resp, err := http.Get("https://www.okx.com/api/v5/public/instruments?instType=SPOT")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code string          `json:"code"`
		Data []OKXInstrument `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var symbols []models.Symbol
	for _, s := range result.Data {
		symbols = append(symbols, models.Symbol{
			Exchange:  o.Name,
			Type:      "spot",
			Symbol:    s.InstId,
			CreatedAt: time.Now(),
		})
	}

	return symbols, nil
}

func (o *OKX) FetchFuturesSymbols() ([]models.Symbol, error) {
	resp, err := http.Get("https://www.okx.com/api/v5/public/instruments?instType=SWAP")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Code string          `json:"code"`
		Data []OKXInstrument `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var symbols []models.Symbol
	for _, s := range result.Data {
		symbols = append(symbols, models.Symbol{
			Exchange:  o.Name,
			Type:      "futures",
			Symbol:    s.InstId,
			CreatedAt: time.Now(),
		})
	}

	return symbols, nil
}
